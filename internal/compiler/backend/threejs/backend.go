package threejs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"

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
	
	// Simple grid layout
	const gridSize = 10
	x, y := 0, 0
	
	// Add functions as nodes
	for i, f := range prog.Funcs {
		name := ""
		// Try to find name from Input ports
		if len(f.IO.In) > 0 {
			name = f.IO.In[0].Path
		} else if len(f.IO.Out) > 0 {
			name = f.IO.Out[0].Path
		}
		
		if name == "" {
			continue
		}

		nodes[name] = nodeData{
			X:     (i % gridSize) * 10,
			Y:     (i / gridSize) * 10,
			Z:     0,
			Ref:   f.Ref,
			Ports: make(map[string]portList),
		}
		
		if nodes[name].Ports["in"] == nil {
			nodes[name].Ports["in"] = make(portList)
		}
		if nodes[name].Ports["out"] == nil {
			nodes[name].Ports["out"] = make(portList)
		}

		for j, port := range f.IO.In {
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
	
	connections := make([]connectionData, 0, len(prog.Connections))
	for from, to := range prog.Connections {
		connections = append(connections, connectionData{
			From:  from.String(),
			To:    to.String(),
			Color: 0x4dd0e1, // Default color
		})
		
		// Check if nodes exist, if not create dummy ones
		if _, ok := nodes[from.Path]; !ok {
			portName := from.Port
			if from.IsArray {
				portName = fmt.Sprintf("%s[%d]", from.Port, from.Idx)
			}
			nodes[from.Path] = nodeData{
				X: x * 10, Y: y * 10, Z: 0, 
				Ref: from.Path,
				Ports: map[string]portList{"out": {portName: {Type: "out", Pos: "bottom", Offset: 0}}},
			}
			x++
		}
		// Update existing node ports if missing
		if node, ok := nodes[from.Path]; ok {
			if node.Ports == nil { node.Ports = make(map[string]portList) }
			if node.Ports["out"] == nil { node.Ports["out"] = make(portList) }
			
			portName := from.Port
			if from.IsArray {
				portName = fmt.Sprintf("%s[%d]", from.Port, from.Idx)
			}

			if _, ok := node.Ports["out"][portName]; !ok {
				node.Ports["out"][portName] = portData{Type: "out", Pos: "bottom", Offset: 0}
			}
		}

		if _, ok := nodes[to.Path]; !ok {
			portName := to.Port
			if to.IsArray {
				portName = fmt.Sprintf("%s[%d]", to.Port, to.Idx)
			}
			nodes[to.Path] = nodeData{
				X: x * 10, Y: y * 10, Z: 0,
				Ref: to.Path,
				Ports: map[string]portList{"in": {portName: {Type: "in", Pos: "top", Offset: 0}}},
			}
			x++
		}
		if node, ok := nodes[to.Path]; ok {
			if node.Ports == nil { node.Ports = make(map[string]portList) }
			if node.Ports["in"] == nil { node.Ports["in"] = make(portList) }
			
			portName := to.Port
			if to.IsArray {
				portName = fmt.Sprintf("%s[%d]", to.Port, to.Idx)
			}
			
			if _, ok := node.Ports["in"][portName]; !ok {
				node.Ports["in"][portName] = portData{Type: "in", Pos: "top", Offset: 0}
			}
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
        camera.position.set(0, 10, 15);

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
            const fontSize = 80;
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
            const labelPlane = new THREE.Mesh(new THREE.PlaneGeometry(NODE_WIDTH * 0.8, NODE_HEIGHT * 0.4), labelMaterial);
            labelPlane.position.z = NODE_DEPTH / 2 + 0.01;
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
                        const portFontSize = 40;
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
