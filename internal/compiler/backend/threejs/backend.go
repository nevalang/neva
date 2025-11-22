package threejs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"sort"
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

// To fix crossings, we need to know which port connects to which node.
type outgoingConn struct {
	targetNode string
	portOffset float64
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
	parentToChildren := make(map[string][]outgoingConn)

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
		// Ensure From Port exists and get offset
		node := nodes[fromNode]
		if node.Ports == nil { node.Ports = make(map[string]portList) }
		if node.Ports["out"] == nil { node.Ports["out"] = make(portList) }
		
		// If implicit, default offset 0. If explicit (from Funcs loop), it's already set.
		if _, ok := node.Ports["out"][fromPortName]; !ok {
			node.Ports["out"][fromPortName] = portData{Type: "out", Pos: "bottom", Offset: 0}
		}
		
		fromOffset := node.Ports["out"][fromPortName].Offset
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
		
		parentToChildren[fromNode] = append(parentToChildren[fromNode], outgoingConn{
			targetNode: toNode,
			portOffset: fromOffset,
		})
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
	const xSpacing = 15
	const ySpacing = 10

	// We will store computed X for nodes to use for next rank
	nodeX := make(map[string]float64)

	for r := 0; r <= maxRank; r++ {
		rankNodes := nodesByRank[r]
		
		// Sort rankNodes to minimize crossings
		if r > 0 {
			sort.Slice(rankNodes, func(i, j int) bool {
				nodeA := rankNodes[i]
				nodeB := rankNodes[j]
				
				weightA := getParentWeight(nodeA, nodeX, parentToChildren)
				weightB := getParentWeight(nodeB, nodeX, parentToChildren)
				
				if weightA != weightB {
					return weightA < weightB
				}
				return nodeA < nodeB // Stable tie-break
			})
		} else {
			// Rank 0: just sort by name for stability
			sort.Strings(rankNodes)
		}

		for i, nodeName := range rankNodes {
			node := nodes[nodeName]
			// Invert Y so rank 0 is at top
			// Center X around 0
			node.Y = (maxRank/2 - r) * ySpacing 
			
			x := (float64(i) - float64(len(rankNodes)-1)/2.0) * float64(xSpacing)
			
			node.X = int(x)
			node.Z = 0
			nodes[nodeName] = node
			nodeX[nodeName] = x
		}
	}
	
	// Handle disconnected nodes (rank not assigned)
	for nodeName, node := range nodes {
		if _, ok := ranks[nodeName]; !ok {
			node.X = -20 
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

// getParentWeight calculates the "ideal" X position for a node based on its parents' positions and port offsets.
func getParentWeight(childNode string, nodeX map[string]float64, parentToChildren map[string][]outgoingConn) float64 {
	sumX := 0.0
	count := 0.0
	
	for parent, conns := range parentToChildren {
		if pX, ok := nodeX[parent]; ok {
			for _, conn := range conns {
				if conn.targetNode == childNode {
					// The ideal X for child is ParentX + PortOffset
					// This tends to align the child with the port it connects to.
					// Multiply offset to give it weight in screen space (5 units approx)
					sumX += pX + conn.portOffset*5 
					count++
				}
			}
		}
	}
	
	if count == 0 {
		return 0 // No parents with assigned positions
	}
	return sumX / count
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
            pointer-events: none;
        }
    </style>
</head>
<body>
    <div id="info">
        Scroll to zoom, Drag to rotate
    </div>
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
        camera.position.set(0, 0, 60); 

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
        const orbitControls = new OrbitControls(camera, renderer.domElement);
        orbitControls.enableDamping = true;
        orbitControls.dampingFactor = 0.05;

        // --- Dataflow Program Definition ---
        const nodesData = {{.NodesJSON}};
        const connectionsData = {{.ConnectionsJSON}};

        // --- Constants ---
        const NODE_WIDTH = 5;
        const NODE_HEIGHT = 3;
        const NODE_DEPTH = 1;
        const PORT_RADIUS = 0.3;
        const PORT_HEIGHT = 0.1;

        // --- State ---
        const nodeMeshes = []; 
        const nodeMap = {}; // name -> { mesh, ports: { name -> localPos } }
        const connectionLines = []; 

        // --- Create Nodes ---
        function createNode(name, data) {
            const nodeGroup = new THREE.Group();
            nodeGroup.position.set(data.x, data.y, data.z);
            nodeGroup.name = name;
            // Store velocity for force layout
            nodeGroup.userData.velocity = new THREE.Vector3();
            nodeGroup.userData.force = new THREE.Vector3();

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

            // Ports
            const ports = {};
            if (data.ports) {
                for (const typeKey in data.ports) {
                    for (const portName in data.ports[typeKey]) {
                        const portData = data.ports[typeKey][portName];
                        
                        const portGeometry = new THREE.CylinderGeometry(PORT_RADIUS, PORT_RADIUS, PORT_HEIGHT, 16);
                        const portMaterial = new THREE.MeshStandardMaterial({
                            color: portData.type === 'in' ? 0x66bb6a : 0xef5350,
                            emissive: portData.type === 'in' ? 0x388e3c : 0xc62828,
                            emissiveIntensity: 0.5,
                            roughness: 0.5
                        });
                        const portMesh = new THREE.Mesh(portGeometry, portMaterial);

                        // Local position calculation
                        const localPos = new THREE.Vector3();
                        if (portData.pos === 'top') {
                            localPos.set(
                                portData.offset * (NODE_WIDTH / 2 - PORT_RADIUS),
                                NODE_HEIGHT / 2 + (portData.type === 'in' ? -PORT_HEIGHT/2 : PORT_HEIGHT/2),
                                0
                            );
                            portMesh.rotation.x = Math.PI / 2;
                        } else { // bottom
                            localPos.set(
                                portData.offset * (NODE_WIDTH / 2 - PORT_RADIUS),
                                -NODE_HEIGHT / 2 + (portData.type === 'in' ? -PORT_HEIGHT/2 : PORT_HEIGHT/2),
                                0
                            );
                            portMesh.rotation.x = -Math.PI / 2;
                        }
                        
                        portMesh.position.copy(localPos);
                        nodeGroup.add(portMesh);
                        
                        ports[portName] = localPos; // Store local offset

                        // Port Label
                        const pCanvas = document.createElement('canvas');
                        const pCtx = pCanvas.getContext('2d');
                        const pFontSize = 32;
                        pCtx.font = ` + "`" + `${pFontSize}px Arial` + "`" + `;
                        const pTextWidth = pCtx.measureText(portName).width;
                        const pTextHeight = pFontSize;
                        pCanvas.width = pTextWidth + 10;
                        pCanvas.height = pTextHeight + 10;
                        pCtx.font = ` + "`" + `${pFontSize}px Arial` + "`" + `;
                        pCtx.fillStyle = '#ffffff';
                        pCtx.textAlign = 'center';
                        pCtx.textBaseline = 'middle';
                        pCtx.fillText(portName, pCanvas.width / 2, pCanvas.height / 2);

                        const pTexture = new THREE.CanvasTexture(pCanvas);
                        const pLabelMat = new THREE.MeshBasicMaterial({ map: pTexture, transparent: true });
                        const pLabelPlane = new THREE.Mesh(new THREE.PlaneGeometry(PORT_RADIUS * 4, PORT_RADIUS * 2), pLabelMat);

                        if (portData.pos === 'top') {
                            pLabelPlane.position.copy(localPos);
                            pLabelPlane.position.y += PORT_HEIGHT + PORT_RADIUS;
                        } else {
                            pLabelPlane.position.copy(localPos);
                            pLabelPlane.position.y -= PORT_HEIGHT + PORT_RADIUS;
                        }
                        nodeGroup.add(pLabelPlane);
                    }
                }
            }

            scene.add(nodeGroup);
            nodeMeshes.push(nodeGroup); 
            nodeMap[name] = { mesh: nodeGroup, ports: ports };
        }

        // Initialize Nodes
        for (const name in nodesData) {
            createNode(name, nodesData[name]);
        }

        // Initialize Connections
        const materialCache = {};
        function getLineMaterial(color) {
            if (!materialCache[color]) {
                materialCache[color] = new THREE.LineBasicMaterial({ 
                    color: color,
                    linewidth: 2 
                });
            }
            return materialCache[color];
        }

        connectionsData.forEach(conn => {
            const fromParts = conn.from.split(':');
            const toParts = conn.to.split(':');
            const fromNodeName = fromParts[0];
            const fromPortName = fromParts[1];
            const toNodeName = toParts[0];
            const toPortName = toParts[1];

            const fromNode = nodeMap[fromNodeName];
            const toNode = nodeMap[toNodeName];

            if (fromNode && toNode) {
                const geometry = new THREE.BufferGeometry();
                const material = getLineMaterial(conn.color);
                const line = new THREE.Line(geometry, material);
                
                scene.add(line);
                
                connectionLines.push({
                    mesh: line,
                    fromNode: fromNode,
                    fromPort: fromPortName,
                    toNode: toNode,
                    toPort: toPortName
                });
            }
        });

        // --- Force Layout Simulation ---
        const REPULSION_STRENGTH = 500;
        const SPRING_LENGTH = 15;
        const SPRING_STRENGTH = 0.05;
        const DAMPING = 0.90;
        const CENTER_STRENGTH = 0.01;

        function updatePhysics() {
            // 1. Repulsion (All pairs)
            for (let i = 0; i < nodeMeshes.length; i++) {
                const nodeA = nodeMeshes[i];
                for (let j = i + 1; j < nodeMeshes.length; j++) {
                    const nodeB = nodeMeshes[j];
                    
                    const diff = new THREE.Vector3().subVectors(nodeA.position, nodeB.position);
                    const distSq = diff.lengthSq();
                    
                    if (distSq > 0.1 && distSq < 5000) { // Limit range
                        const force = diff.normalize().multiplyScalar(REPULSION_STRENGTH / distSq);
                        nodeA.userData.force.add(force);
                        nodeB.userData.force.sub(force);
                    }
                }
                
                // 2. Center Gravity (Weak)
                const centerDir = new THREE.Vector3().subVectors(new THREE.Vector3(0,0,0), nodeA.position);
                nodeA.userData.force.add(centerDir.multiplyScalar(CENTER_STRENGTH));
            }

            // 3. Spring Force (Connections)
            connectionLines.forEach(conn => {
                const nodeA = conn.fromNode.mesh;
                const nodeB = conn.toNode.mesh;
                
                const diff = new THREE.Vector3().subVectors(nodeA.position, nodeB.position);
                const dist = diff.length();
                const displacement = dist - SPRING_LENGTH;
                
                const force = diff.normalize().multiplyScalar(-SPRING_STRENGTH * displacement);
                
                nodeA.userData.force.add(force);
                nodeB.userData.force.sub(force);
            });

            // 4. Apply Force to Velocity & Position
            nodeMeshes.forEach(node => {
                const vel = node.userData.velocity;
                vel.add(node.userData.force);
                vel.multiplyScalar(DAMPING);
                
                // Limit speed
                if (vel.lengthSq() > 100) vel.setLength(10);
                
                node.position.add(vel);
                node.userData.force.set(0, 0, 0); // Reset force
            });
        }

        function updateConnections() {
            connectionLines.forEach(conn => {
                const start = new THREE.Vector3();
                const end = new THREE.Vector3();
                
                // Get world positions of ports
                const startLocal = conn.fromNode.ports[conn.fromPort] || new THREE.Vector3(0,0,0);
                const endLocal = conn.toNode.ports[conn.toPort] || new THREE.Vector3(0,0,0);
                
                start.copy(startLocal).applyMatrix4(conn.fromNode.mesh.matrixWorld);
                end.copy(endLocal).applyMatrix4(conn.toNode.mesh.matrixWorld);

                // Create curved line points
                // Simple bezier control points
                const dist = start.distanceTo(end);
                const control1 = start.clone().add(new THREE.Vector3(0, dist * 0.25, dist * 0.1));
                const control2 = end.clone().add(new THREE.Vector3(0, -dist * 0.25, dist * 0.1));
                
                const curve = new THREE.CubicBezierCurve3(start, control1, control2, end);
                const points = curve.getPoints(20);
                
                conn.mesh.geometry.setFromPoints(points);
            });
        }

        // Pre-calculate layout to avoid shaking
        for (let i = 0; i < 300; i++) {
            updatePhysics();
        }
        
        // Initial connection update
        scene.updateMatrixWorld();
        updateConnections();

        // --- Animation Loop ---
        function animate() {
            requestAnimationFrame(animate);
            orbitControls.update();
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
