package threejs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	return nil
}

type Encoder struct{}

func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	tmpl, err := template.New("threejs").Parse(templateHTML)
	if err != nil {
		return err
	}

	data, err := prepareData(prog)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

type nodeData struct {
	X     int                 `json:"x"`
	Y     int                 `json:"y"`
	Z     int                 `json:"z"`
	Ref   string              `json:"ref"`
	Ports map[string]portList `json:"ports"`
}

type portList map[string]portData

type portData struct {
	Type   string  `json:"type"`
	Pos    string  `json:"pos"`
	Offset float64 `json:"offset"`
}

type connectionData struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Color int    `json:"color"`
}

type templateData struct {
	NodesJSON       template.JS
	ConnectionsJSON template.JS
}

func prepareData(prog *ir.Program) (templateData, error) {
	nodes := make(map[string]nodeData)
	pathMap := make(map[string]string) // Maps a port's path to the unified node name

	// Helper to normalize node names
	getNodeName := func(path string) string {
		for _, suffix := range []string{"/in", "/out", ".in", ".out"} {
			if strings.HasSuffix(path, suffix) {
				return strings.TrimSuffix(path, suffix)
			}
		}
		return path
	}
	
	// 1. Collect all nodes defined in Funcs
	for _, f := range prog.Funcs {
		// Determine the node name for this function call
		name := ""
		if len(f.IO.In) > 0 {
			name = getNodeName(f.IO.In[0].Path)
		} else if len(f.IO.Out) > 0 {
			name = getNodeName(f.IO.Out[0].Path)
		}
		
		if name == "" {
			continue
		}

		if _, exists := nodes[name]; !exists {
			nodes[name] = nodeData{
				Ref:   f.Ref,
				Ports: make(map[string]portList),
			}
		}
		
		if nodes[name].Ports["in"] == nil {
			nodes[name].Ports["in"] = make(portList)
		}
		if nodes[name].Ports["out"] == nil {
			nodes[name].Ports["out"] = make(portList)
		}

		// Map all paths involved in this func to the unified name
		// and populate ports
		for j, port := range f.IO.In {
			pathMap[port.Path] = name
			
			portName := port.Port
			if port.IsArray {
				portName = fmt.Sprintf("%s[%d]", port.Port, port.Idx)
			}
			nodes[name].Ports["in"][portName] = portData{
				Type:   "in",
				Pos:    "top",
				Offset: float64(j) - float64(len(f.IO.In)-1)/2.0, 
			}
		}
		
		for j, port := range f.IO.Out {
			pathMap[port.Path] = name

			portName := port.Port
			if port.IsArray {
				portName = fmt.Sprintf("%s[%d]", port.Port, port.Idx)
			}
			nodes[name].Ports["out"][portName] = portData{
				Type:   "out",
				Pos:    "bottom",
				Offset: float64(j) - float64(len(f.IO.Out)-1)/2.0,
			}
		}
	}

	// 2. Process connections
	// We use pathMap to resolve node names. If a path is missing (implicit node),
	// we derive the name and create the node on the fly.
	
	// Build adjacency list for layout
	adj := make(map[string][]string)
	inDegree := make(map[string]int)
	
	connections := make([]connectionData, 0, len(prog.Connections))
	for from, to := range prog.Connections {
		// Resolve From Node
		fromNode, ok := pathMap[from.Path]
		if !ok {
			fromNode = getNodeName(from.Path)
			pathMap[from.Path] = fromNode
		}
		
		// Resolve To Node
		toNode, ok := pathMap[to.Path]
		if !ok {
			toNode = getNodeName(to.Path)
			pathMap[to.Path] = toNode
		}

		color := 0x4dd0e1 // Default light blue
		if from.Port == "err" {
			color = 0xffa726 // Orange
		}
		
		fromPortName := from.Port
		if from.IsArray {
			fromPortName = fmt.Sprintf("%s[%d]", from.Port, from.Idx)
		}
		
		toPortName := to.Port
		if to.IsArray {
			toPortName = fmt.Sprintf("%s[%d]", to.Port, to.Idx)
		}

		connections = append(connections, connectionData{
			From:  fmt.Sprintf("%s:%s", fromNode, fromPortName),
			To:    fmt.Sprintf("%s:%s", toNode, toPortName),
			Color: color,
		})
		
		// Ensure From Node exists (handle implicit)
		if _, ok := nodes[fromNode]; !ok {
			nodes[fromNode] = nodeData{
				Ref: fromNode,
				Ports: make(map[string]portList),
			}
		}
		// Ensure From Port exists
		node := nodes[fromNode]
		if node.Ports == nil { node.Ports = make(map[string]portList) }
		if node.Ports["out"] == nil { node.Ports["out"] = make(portList) }
		if _, ok := node.Ports["out"][fromPortName]; !ok {
			node.Ports["out"][fromPortName] = portData{Type: "out", Pos: "bottom", Offset: 0}
		}
		nodes[fromNode] = node // Write back updated copy

		// Ensure To Node exists (handle implicit)
		if _, ok := nodes[toNode]; !ok {
			nodes[toNode] = nodeData{
				Ref: toNode,
				Ports: make(map[string]portList),
			}
		}
		// Ensure To Port exists
		node = nodes[toNode]
		if node.Ports == nil { node.Ports = make(map[string]portList) }
		if node.Ports["in"] == nil { node.Ports["in"] = make(portList) }
		if _, ok := node.Ports["in"][toPortName]; !ok {
			node.Ports["in"][toPortName] = portData{Type: "in", Pos: "top", Offset: 0}
		}
		nodes[toNode] = node // Write back updated copy

		adj[fromNode] = append(adj[fromNode], toNode)
		inDegree[toNode]++
		if _, exists := inDegree[fromNode]; !exists {
			inDegree[fromNode] = 0
		}
	}

	// Calculate ranks (BFS)
	ranks := make(map[string]int)
	queue := []string{}
	
	// Find roots
	for node := range nodes {
		if inDegree[node] == 0 {
			ranks[node] = 0
			queue = append(queue, node)
		}
	}
	// If no roots (cycles), pick one arbitrarily
	if len(queue) == 0 && len(nodes) > 0 {
		for node := range nodes {
			ranks[node] = 0
			queue = append(queue, node)
			break
		}
	}

	visited := make(map[string]bool)
	for _, n := range queue {
		visited[n] = true
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		
		currentRank := ranks[curr]
		
		for _, neighbor := range adj[curr] {
			if r, seen := ranks[neighbor]; !seen || r < currentRank+1 {
				ranks[neighbor] = currentRank + 1
				if !visited[neighbor] {
					queue = append(queue, neighbor)
					visited[neighbor] = true
				}
			}
		}
	}

	// Group by rank
	nodesByRank := make(map[int][]string)
	maxRank := 0
	for node, rank := range ranks {
		nodesByRank[rank] = append(nodesByRank[rank], node)
		if rank > maxRank {
			maxRank = rank
		}
	}

	// Assign coordinates
	// Y based on rank (Top-Down), X based on position in rank (Centered)
	const xSpacing = 12
	const ySpacing = 8

	for r := 0; r <= maxRank; r++ {
		rankNodes := nodesByRank[r]
		for i, nodeName := range rankNodes {
			node := nodes[nodeName]
			// Invert Y so rank 0 is at top
			// Center X around 0
			node.Y = (maxRank/2 - r) * ySpacing 
			node.X = (i - (len(rankNodes)-1)/2.0) * xSpacing
			node.Z = 0
			nodes[nodeName] = node
		}
	}
	
	// Handle disconnected nodes (rank not assigned)
	for nodeName, node := range nodes {
		if _, ok := ranks[nodeName]; !ok {
			node.X = -10 
			node.Y = 0
			nodes[nodeName] = node
		}
	}

	nodesJSON, err := json.Marshal(nodes)
	if err != nil {
		return templateData{}, err
	}
	
	connsJSON, err := json.Marshal(connections)
	if err != nil {
		return templateData{}, err
	}

	return templateData{
		NodesJSON:       template.JS(nodesJSON),
		ConnectionsJSON: template.JS(connsJSON),
	}, nil
}

const templateHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>3D Dataflow Program</title>
    <style>
        body { margin: 0; overflow: hidden; }
        canvas { display: block; }
        #info {
            position: absolute;
            top: 10px;
            left: 10px;
            color: white;
            font-family: monospace;
            background: rgba(0,0,0,0.5);
            padding: 5px;
            border-radius: 3px;
        }
    </style>
</head>
<body>
    <div id="info">Drag to rotate, Scroll to zoom</div>
    <script type="importmap">
        {
          "imports": {
            "three": "https://cdn.jsdelivr.net/npm/three@0.181.2/build/three.module.js",
            "three/addons/": "https://cdn.jsdelivr.net/npm/three@0.181.2/examples/jsm/"
          }
        }
    </script>
    <script type="module">
        import * as THREE from 'three';
        import { OrbitControls } from 'three/addons/controls/OrbitControls.js';

        // --- Scene Setup ---
        const scene = new THREE.Scene();
        scene.background = new THREE.Color(0x1a1a2e);

        const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        camera.position.set(0, 0, 40); 

        const renderer = new THREE.WebGLRenderer({ antialias: true });
        renderer.setSize(window.innerWidth, window.innerHeight);
        document.body.appendChild(renderer.domElement);

        // --- Lighting ---
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.6);
        scene.add(ambientLight);

        const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight.position.set(5, 10, 7);
        scene.add(directionalLight);

        // --- Orbit Controls ---
        const controls = new OrbitControls(camera, renderer.domElement);
        controls.enableDamping = true;
        controls.dampingFactor = 0.05;

        // --- Dataflow Program Definition ---
        const nodes = {{.NodesJSON}};
        const connections = {{.ConnectionsJSON}};

        // --- Node & Port Dimensions ---
        const NODE_WIDTH = 5;
        const NODE_HEIGHT = 3;
        const NODE_DEPTH = 1;
        const PORT_RADIUS = 0.3;
        const PORT_HEIGHT = 0.1;

        const portMap = {};

        // --- Create Nodes ---
        function createNode(name, data) {
            const nodeGroup = new THREE.Group();
            nodeGroup.position.set(data.x, data.y, data.z);
            nodeGroup.name = name;

            const boxGeometry = new THREE.BoxGeometry(NODE_WIDTH, NODE_HEIGHT, NODE_DEPTH);
            const boxMaterial = new THREE.MeshStandardMaterial({
                color: 0xcccccc,
                roughness: 0.6,
                metalness: 0.3,
                emissive: 0x1a1a2e,
                emissiveIntensity: 0.1
            });
            const nodeMesh = new THREE.Mesh(boxGeometry, boxMaterial);
            nodeGroup.add(nodeMesh);

            // Node Label
            const canvas = document.createElement('canvas');
            const context = canvas.getContext('2d');
            const fontSize = 64;
            context.font = ` + "`" + `${fontSize}px Arial` + "`" + `;
            const textWidth = context.measureText(name).width;
            const textHeight = fontSize;
            canvas.width = textWidth + 20;
            canvas.height = textHeight + 20;
            context.font = ` + "`" + `${fontSize}px Arial` + "`" + `;
            context.fillStyle = '#333333';
            context.textAlign = 'center';
            context.textBaseline = 'middle';
            context.fillText(name, canvas.width / 2, canvas.height / 2);

            const texture = new THREE.CanvasTexture(canvas);
            const labelMaterial = new THREE.MeshBasicMaterial({ map: texture, transparent: true });
            const labelPlane = new THREE.Mesh(new THREE.PlaneGeometry(NODE_WIDTH * 0.9, NODE_HEIGHT * 0.5), labelMaterial);
            labelPlane.position.z = NODE_DEPTH / 2 + 0.1;
            nodeGroup.add(labelPlane);

            // Create Ports
            if (data.ports) {
                for (const typeKey in data.ports) {
                    for (const portName in data.ports[typeKey]) {
                        const portData = data.ports[typeKey][portName];
                        const portAddr = ` + "`" + `${name}:${portName}` + "`" + `;

                        const portGeometry = new THREE.CylinderGeometry(PORT_RADIUS, PORT_RADIUS, PORT_HEIGHT, 16);
                        const portMaterial = new THREE.MeshStandardMaterial({
                            color: portData.type === 'in' ? 0x66bb6a : 0xef5350,
                            emissive: portData.type === 'in' ? 0x388e3c : 0xc62828,
                            emissiveIntensity: 0.5,
                            roughness: 0.5
                        });
                        const portMesh = new THREE.Mesh(portGeometry, portMaterial);

                        if (portData.pos === 'top') {
                            portMesh.position.y = NODE_HEIGHT / 2;
                            portMesh.position.x = portData.offset * (NODE_WIDTH / 2 - PORT_RADIUS);
                            portMesh.rotation.x = Math.PI / 2;
                        } else { // 'bottom'
                            portMesh.position.y = -NODE_HEIGHT / 2;
                            portMesh.position.x = portData.offset * (NODE_WIDTH / 2 - PORT_RADIUS);
                            portMesh.rotation.x = -Math.PI / 2;
                        }
                        portMesh.position.y += (portData.type === 'in' ? -PORT_HEIGHT / 2 : PORT_HEIGHT / 2);

                        nodeGroup.add(portMesh);
                        portMap[portAddr] = new THREE.Vector3();
                        portMesh.userData.portAddr = portAddr;

                        // Port Label
                        const portLabelCanvas = document.createElement('canvas');
                        const portLabelContext = portLabelCanvas.getContext('2d');
                        const portFontSize = 32;
                        portLabelContext.font = ` + "`" + `${portFontSize}px Arial` + "`" + `;
                        const portTextWidth = portLabelContext.measureText(portName).width;
                        const portTextHeight = portFontSize;
                        portLabelCanvas.width = portTextWidth + 10;
                        portLabelCanvas.height = portTextHeight + 10;
                        portLabelContext.font = ` + "`" + `${portFontSize}px Arial` + "`" + `;
                        portLabelContext.fillStyle = '#ffffff';
                        portLabelContext.textAlign = 'center';
                        portLabelContext.textBaseline = 'middle';
                        portLabelContext.fillText(portName, portLabelCanvas.width / 2, portLabelCanvas.height / 2);

                        const portTexture = new THREE.CanvasTexture(portLabelCanvas);
                        const portLabelMaterial = new THREE.MeshBasicMaterial({ map: portTexture, transparent: true });
                        const portLabelPlane = new THREE.Mesh(new THREE.PlaneGeometry(PORT_RADIUS * 4, PORT_RADIUS * 2), portLabelMaterial);

                        if (portData.pos === 'top') {
                            portLabelPlane.position.copy(portMesh.position);
                            portLabelPlane.position.y += PORT_HEIGHT + PORT_RADIUS;
                        } else {
                            portLabelPlane.position.copy(portMesh.position);
                            portLabelPlane.position.y -= PORT_HEIGHT + PORT_RADIUS;
                        }
                        nodeGroup.add(portLabelPlane);
                    }
                }
            }

            scene.add(nodeGroup);
            return nodeGroup;
        }

        for (const name in nodes) {
            createNode(name, nodes[name]);
        }

        scene.updateMatrixWorld(true);

        scene.traverse(obj => {
            if (obj.userData.portAddr) {
                portMap[obj.userData.portAddr] = obj.getWorldPosition(new THREE.Vector3());
            }
        });

        connections.forEach(conn => {
            const fromPos = portMap[conn.from];
            const toPos = portMap[conn.to];

            if (fromPos && toPos) {
                const curve = new THREE.CatmullRomCurve3([
                    fromPos,
                    new THREE.Vector3(
                        (fromPos.x + toPos.x) / 2,
                        (fromPos.y + toPos.y) / 2,
                        (fromPos.z + toPos.z) / 2 + 2
                    ),
                    toPos
                ]);

                const geometry = new THREE.TubeGeometry(curve, 64, 0.1, 8, false);
                const material = new THREE.MeshBasicMaterial({
                    color: conn.color,
                    emissive: conn.color,
                    emissiveIntensity: 1.5
                });
                const line = new THREE.Mesh(geometry, material);
                scene.add(line);
            } else {
                console.warn(` + "`" + `Missing port for connection: ${conn.from} -> ${conn.to}` + "`" + `);
            }
        });

        function animate() {
            requestAnimationFrame(animate);
            controls.update();
            renderer.render(scene, camera);
        }
        animate();

        window.addEventListener('resize', () => {
            camera.aspect = window.innerWidth / window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.setSize(window.innerWidth, window.innerHeight);
        });
    </script>
</body>
</html>
`
