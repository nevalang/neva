import { useContext, useMemo } from "react";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  BackgroundVariant,
  Handle,
  Position,
  Edge,
  MarkerType,
  NodeProps,
  Node,
} from "reactflow";
import "reactflow/dist/style.css";
import dagre from "dagre";
import { VSCodeState, vscodeStateContext } from "../helpers/vscode_state";
import * as src from "../generated/sourcecode";

const nodeTypes = { normal_node: NormalNode };

export default function NetView(props: {
  nodes: {
    name: string;
    entity: src.Node;
  }[];
  net: src.Connection[];
}) {
  const vscodeState = useContext(vscodeStateContext) as VSCodeState; // component shouldn't be rendered at all without this
  const { nodes, edges } = useMemo(
    () => getLayoutedElements(props.nodes, props.net, vscodeState),
    [props.nodes, props.net, vscodeState]
  );

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ReactFlow
        nodeTypes={nodeTypes}
        onInit={(instance) => instance.fitView()}
        nodes={nodes}
        edges={edges}
      >
        <Controls />
        <MiniMap />
        <Background variant={BackgroundVariant.Dots} gap={10} size={0.5} />
      </ReactFlow>
    </div>
  );
}

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
// const initialNodes = [
//   {
//     type: "normal_node",
//     id: "in",
//     position: defaultPosition,
//     isHidden: false,
//     data: {
//       ports: {
//         in: [],
//         out: ["enter"],
//       },
//     },
//   },
//   {
//     type: "normal_node",
//     id: "out",
//     position: defaultPosition,
//     data: {
//       ports: {
//         in: ["exit"],
//         out: [],
//       },
//     },
//   },
//   {
//     type: "normal_node",
//     id: "readFirstInt",
//     position: defaultPosition,
//     data: {
//       ports: {
//         in: ["sig"],
//         out: ["v", "err"],
//       },
//     },
//   },
//   {
//     type: "normal_node",
//     id: "readSecondInt",
//     position: defaultPosition,
//     data: {
//       ports: {
//         in: ["sig"],
//         out: ["v", "err"],
//       },
//     },
//   },
//   {
//     type: "normal_node",
//     id: "add",
//     position: defaultPosition,
//     data: {
//       ports: {
//         in: ["a", "b"],
//         out: ["v"],
//       },
//     },
//   },
//   {
//     type: "normal_node",
//     id: "print",
//     position: defaultPosition,
//     data: {
//       ports: {
//         in: ["v"],
//         out: ["v"],
//       },
//     },
//   },
// ];
// const initialEdges: Edge[] = [
//   {
//     id: "in.enter -> readFirstInt.sig",
//     source: "in",
//     sourceHandle: "enter",
//     target: "readFirstInt",
//     targetHandle: "sig",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "readFirstInt.err -> print.v",
//     source: "readFirstInt",
//     sourceHandle: "err",
//     target: "print",
//     targetHandle: "v",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "readFirstInt.v -> add.a",
//     source: "readFirstInt",
//     sourceHandle: "v",
//     target: "add",
//     targetHandle: "a",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "readFirstInt.v -> readSecondInt.sig",
//     source: "readFirstInt",
//     sourceHandle: "v",
//     target: "readSecondInt",
//     targetHandle: "sig",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "readSecondInt.err -> print.v",
//     source: "readSecondInt",
//     sourceHandle: "err",
//     target: "print",
//     targetHandle: "v",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "readSecondInt.v -> add.b",
//     source: "readSecondInt",
//     sourceHandle: "v",
//     target: "add",
//     targetHandle: "b",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "add.v -> print.v",
//     source: "add",
//     sourceHandle: "v",
//     target: "print",
//     targetHandle: "v",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
//   {
//     id: "print.v -> out.exit",
//     source: "print",
//     sourceHandle: "v",
//     target: "out",
//     targetHandle: "exit",
//     markerEnd: {
//       type: MarkerType.Arrow,
//     },
//   },
// ];

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 342.5;
const nodeHeight = 70;

const getLayoutedElements = (
  nodes: { name: string; entity: src.Node }[],
  net: src.Connection[],
  vscodeState: VSCodeState, // TODO use vscode state for ports
  direction = "TB"
) => {
  const isHorizontal = direction === "LR";
  dagreGraph.setGraph({ rankdir: direction });

  const reactflowNodes: Node[] = [];
  for (const node of nodes) {
    const reactflowNode = {
      id: node.name,
      type: "normal_node",
      position: defaultPosition,
      data: {
        // TODO we need more than current parsed file to render this
        ports: {
          in: {},
          out: {},
        },
      },
    };
    reactflowNodes.push(reactflowNode);
    dagreGraph.setNode(node.name, { width: nodeWidth, height: nodeHeight });
  }

  const reactflowEdges: Edge[] = [];
  for (const connection of net) {
    const { senderSide, receiverSide } = connection;
    if (!senderSide || !receiverSide) {
      continue;
    }

    for (const receiver of receiverSide) {
      const source = senderSide.portAddr
        ? senderSide.portAddr.node
        : `${senderSide.constRef?.pkg}.${senderSide.constRef?.name}`;

      const sourceHandle = senderSide.portAddr
        ? senderSide.portAddr.port
        : "out"; // TODO check that this is how constant works

      const reactflowEdge = {
        id: `${senderSide.portAddr || senderSide.constRef} -> ${
          receiver.portAddr
        }`,
        source: source || "unknown",
        sourceHandle: sourceHandle,
        target: receiver.portAddr?.node || "unknown",
        targetHandle: receiver.portAddr?.port || "unknown",
        markerEnd: {
          type: MarkerType.Arrow,
        },
      };
      reactflowEdges.push(reactflowEdge);
    }
  }

  reactflowEdges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  reactflowNodes.forEach((node) => {
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

  return { nodes: reactflowNodes, edges: reactflowEdges };
};
