import { useMemo } from "react";
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
import * as src from "../generated/sourcecode";
import { ComponentViewState } from "../core/file_view_state";

const nodeTypes = { normal_node: NormalNode };

interface INetViewProps {
  componentViewState: ComponentViewState;
}

export default function NetView(props: INetViewProps) {
  const { nodes, edges } = useMemo(
    () => getReactFlowElements(props.componentViewState),
    [props.componentViewState]
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

function NormalNode(props: NodeProps<{ ports: src.Interface }>) {
  const { io } = props.data.ports;

  const { inports, outports } = useMemo(() => {
    const result = { inports: [], outports: [] };
    if (!io) {
      return result;
    }
    return {
      inports: Object.entries(io.in || {}),
      outports: Object.entries(io.out || {}),
    };
  }, [io]);

  return (
    <div className="react-flow__node-default">
      {inports.length > 0 && (
        <div className="inports">
          {inports.map(([inportName]) => (
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
      {outports.length > 0 && (
        <div className="outports">
          {outports.map(([outportName]) => (
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

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 342.5;
const nodeHeight = 70;

const getReactFlowElements = (
  componentViewState: ComponentViewState,
  direction = "TB"
) => {
  const { nodes, interface: iface, net } = componentViewState;

  const isHorizontal = direction === "LR";
  dagreGraph.setGraph({ rankdir: direction });

  const reactflowNodes: Node[] = [];
  for (const nodeView of nodes) {
    const reactflowNode = {
      id: nodeView.name,
      type: "normal_node",
      position: defaultPosition,
      data: { ports: nodeView.interface },
    };
    reactflowNodes.push(reactflowNode);
    dagreGraph.setNode(nodeView.name, { width: nodeWidth, height: nodeHeight });
  }

  if (iface) {
    const ioNodes = getIONodes(iface);

    reactflowNodes.push(ioNodes.in);
    dagreGraph.setNode(ioNodes.in.id, {
      width: nodeWidth,
      height: nodeHeight,
    });

    reactflowNodes.push(ioNodes.out);
    dagreGraph.setNode(ioNodes.out.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  }

  console.log(reactflowNodes);

  const reactflowEdges: Edge[] = [];
  for (const connection of net!) {
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
        : "out";

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

  console.log({ reactflowNodes, iface });

  return { nodes: reactflowNodes, edges: reactflowEdges };
};

function getIONodes(iface: src.Interface) {
  const inportsNode = {
    id: "in",
    type: "normal_node",
    position: defaultPosition,
    data: {
      ports: {
        io: {
          in: {},
          out: {},
        },
      } as src.Interface,
    },
  };

  for (const portName in iface!.io?.in) {
    const port = iface!.io?.in[portName];
    inportsNode.data.ports.io!.out![portName] = port; // inport for component is outport for inport-node in network
  }

  const outportsNode = {
    id: "out",
    type: "normal_node",
    position: defaultPosition,
    data: {
      ports: {
        io: {
          in: {},
          out: {},
        },
      } as src.Interface,
    },
  };

  for (const portName in iface!.io?.out) {
    const port = iface!.io?.out[portName];
    outportsNode.data.ports.io!.in![portName] = port; // outport for component is inport for outport-node in network
  }

  return {
    in: inportsNode,
    out: outportsNode,
  };
}
