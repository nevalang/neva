import { useCallback, useMemo } from "react";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  BackgroundVariant,
  Connection,
  Handle,
  Position,
  Edge,
  Node,
  MarkerType,
  NodeProps,
} from "reactflow";
import dagre from "dagre";
import "reactflow/dist/style.css";

interface IPorts {
  in: string[];
  out: string[];
}

function NormalNode(props: NodeProps<{ ports: IPorts }>) {
  return (
    <div className="react-flow__node-default">
      {props.data.ports.in.length > 0 && (
        <div className="inports">
          {props.data.ports.in.map((inportName) => (
            <Handle
              content="asd"
              type="target"
              id={inportName}
              key={inportName}
              position={Position.Top}
              isConnectable={true}
            >
              {inportName}
            </Handle>
          ))}
        </div>
      )}
      <div className="nodeName">{props.id}</div>
      {props.data.ports.out.length > 0 && (
        <div className="outports">
          {props.data.ports.out.map((outportName) => (
            <Handle
              type="source"
              id={outportName}
              key={outportName}
              position={Position.Bottom}
              isConnectable={true}
            >
              {outportName}
            </Handle>
          ))}
        </div>
      )}
    </div>
  );
}

const defaultPosition = { x: 0, y: 0 };
const initialNodes = [
  {
    type: "normalNode",
    id: "in",
    position: defaultPosition,
    isHidden: false,
    data: {
      ports: {
        in: [],
        out: ["enter"],
      },
    },
  },
  {
    type: "normalNode",
    id: "out",
    position: defaultPosition,
    data: {
      ports: {
        in: ["exit"],
        out: [],
      },
    },
  },
  {
    type: "normalNode",
    id: "readFirstInt",
    position: defaultPosition,
    data: {
      ports: {
        in: ["sig"],
        out: ["v", "err"],
      },
    },
  },
  {
    type: "normalNode",
    id: "readSecondInt",
    position: defaultPosition,
    data: {
      ports: {
        in: ["sig"],
        out: ["v", "err"],
      },
    },
  },
  {
    type: "normalNode",
    id: "add",
    position: defaultPosition,
    data: {
      ports: {
        in: ["a", "b"],
        out: ["v"],
      },
    },
  },
  {
    type: "normalNode",
    id: "print",
    position: defaultPosition,
    data: {
      ports: {
        in: ["v"],
        out: ["v"],
      },
    },
  },
];
const initialEdges: Edge[] = [
  {
    id: "in.enter -> readFirstInt.sig",
    source: "in",
    sourceHandle: "enter",
    target: "readFirstInt",
    targetHandle: "sig",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "readFirstInt.err -> print.v",
    source: "readFirstInt",
    sourceHandle: "err",
    target: "print",
    targetHandle: "v",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "readFirstInt.v -> add.a",
    source: "readFirstInt",
    sourceHandle: "v",
    target: "add",
    targetHandle: "a",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "readFirstInt.v -> readSecondInt.sig",
    source: "readFirstInt",
    sourceHandle: "v",
    target: "readSecondInt",
    targetHandle: "sig",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "readSecondInt.err -> print.v",
    source: "readSecondInt",
    sourceHandle: "err",
    target: "print",
    targetHandle: "v",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "readSecondInt.v -> add.b",
    source: "readSecondInt",
    sourceHandle: "v",
    target: "add",
    targetHandle: "b",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "add.v -> print.v",
    source: "add",
    sourceHandle: "v",
    target: "print",
    targetHandle: "v",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
  {
    id: "print.v -> out.exit",
    source: "print",
    sourceHandle: "v",
    target: "out",
    targetHandle: "exit",
    markerEnd: {
      type: MarkerType.Arrow,
    },
  },
];

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 342.5;
const nodeHeight = 70;

const getLayoutedElements = (
  nodes: Node[],
  edges: Edge[],
  direction = "TB"
) => {
  const isHorizontal = direction === "LR";
  dagreGraph.setGraph({ rankdir: direction });

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  nodes.forEach((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.targetPosition = (isHorizontal ? "left" : "top") as Position;
    node.sourcePosition = (isHorizontal ? "right" : "bottom") as Position;

    // We are shifting the dagre node position (anchor=center center) to the top left
    // so it matches the React Flow node anchor point (top left).
    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };

    return node;
  });

  return { nodes, edges };
};

const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
  initialNodes,
  initialEdges
);

export default function App() {
  const [nodes, _, onNodesChange] = useNodesState(layoutedNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(layoutedEdges);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  const nodeTypes = useMemo(
    () => ({
      normalNode: NormalNode,
    }),
    []
  );

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodeTypes={nodeTypes}
        onInit={(instance) => instance.fitView()}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
      >
        <Controls />
        <MiniMap />
        <Background variant={BackgroundVariant.Dots} gap={5} size={0.5} />
      </ReactFlow>
    </div>
  );
}
