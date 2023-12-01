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
  useNodesState,
  useEdgesState,
  XYPosition,
  useStore,
  BaseEdge,
  EdgeProps,
  getStraightPath,
} from "reactflow";
import "reactflow/dist/style.css";
// import SmartBezierEdge from "@tisoap/react-flow-smart-edge";
import dagre from "dagre";
import * as src from "../generated/sourcecode";
import { ComponentViewState } from "../core/file_view_state";

const nodeTypes = { normal: NormalNode };
const edgeTypes = { smart: NormalEdge };

interface INetViewProps {
  name: string;
  componentViewState: ComponentViewState;
}

export default function NetView(props: INetViewProps) {
  const { nodes, edges } = useMemo(() => {
    const dagreGraph = new dagre.graphlib.Graph();
    dagreGraph.setDefaultEdgeLabel(() => ({}));
    return getReactFlowElements(
      props.name,
      props.componentViewState,
      dagreGraph
    );
  }, [props.name, props.componentViewState]);

  const [nodesState, , onNodesChange] = useNodesState(nodes);
  const [edgesState, setEdgesState, onEdgesChange] = useEdgesState(edges);

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ReactFlow
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onInit={(instance) => instance.fitView()}
        nodes={nodesState}
        edges={edgesState}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        fitView
        nodesConnectable={false}
        onNodeMouseEnter={(_, node: Node) => {
          setEdgesState(
            edgesState.map((edge) =>
              edge.source === node.id || edge.target === node.id
                ? {
                    ...edge,
                    data: { isHighlighted: true },
                  }
                : edge
            )
          );
        }}
        onNodeMouseLeave={() => {
          setEdgesState(
            edgesState.map((edge) =>
              edge.data.isHighlighted
                ? {
                    ...edge,
                    data: { isHighlighted: false },
                  }
                : edge
            )
          );
        }}
      >
        <Controls />
        <MiniMap />
        <Background variant={BackgroundVariant.Dots} gap={10} size={0.5} />
      </ReactFlow>
    </div>
  );
}

function NormalEdge(props: EdgeProps<{ isHighlighted: boolean }>) {
  const [edgePath] = getStraightPath({
    sourceX: props.sourceX,
    sourceY: props.sourceY,
    targetX: props.targetX,
    targetY: props.targetY,
  });

  const style = props.data?.isHighlighted ? { stroke: "white" } : {};

  return <BaseEdge {...props} path={edgePath} style={style} />;
}

function NormalNode(props: NodeProps<{ ports: src.Interface; label: string }>) {
  const { io } = props.data.ports;
  const isPortsVisible = useStore((s) => s.transform[2] >= 0.51);
  const isTypesVisible = useStore((s) => s.transform[2] >= 1);

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
          {inports.map(([inportName, inportType]) => (
            <Handle
              content="asd"
              type="target"
              id={inportName}
              key={inportName}
              position={Position.Top}
              isConnectable={true}
            >
              {/* {inportName} */}
              {isPortsVisible && inportName}
              {isTypesVisible &&
                inportType.typeExpr &&
                inportType.typeExpr.meta && (
                  <span className="portType">
                    {" "}
                    {(inportType.typeExpr.meta as src.Meta).Text}
                  </span>
                )}
            </Handle>
          ))}
        </div>
      )}
      <div className="nodeName">{props.data.label}</div>
      {outports.length > 0 && (
        <div className="outports">
          {outports.map(([outportName, outportType]) => (
            <Handle
              type="source"
              id={outportName}
              key={outportName}
              position={Position.Bottom}
              isConnectable={true}
            >
              {isPortsVisible && outportName}
              {isTypesVisible &&
                outportType.typeExpr &&
                outportType.typeExpr.meta && (
                  <span className="portType">
                    {" "}
                    {(outportType.typeExpr.meta as src.Meta).Text}{" "}
                  </span>
                )}
            </Handle>
          ))}
        </div>
      )}
    </div>
  );
}

const getReactFlowElements = (
  name: string,
  componentViewState: ComponentViewState,
  dagreGraph: dagre.graphlib.Graph
) => {
  const { nodes, interface: iface, net } = componentViewState;

  const direction = "TB";
  const isHorizontal = false;
  dagreGraph.setGraph({ rankdir: direction });

  const defaultPosition = { x: 0, y: 0 };
  const nodeWidth = 342.5;
  const nodeHeight = 70;

  // const containerNode = {
  //   id: `${name}-container`,
  //   type: "group",
  //   data: { label: name },
  //   position: defaultPosition,
  // };

  const reactflowNodes: Node[] = [];
  // dagreGraph.setNode(containerNode.id, {
  //   width: 1000,
  //   height: 1000,
  // });

  for (const nodeView of nodes) {
    const reactflowNode = {
      id: `${name}-${nodeView.name}`,
      type: "normal",
      position: defaultPosition,
      data: {
        ports: nodeView.interface,
        label: nodeView.name,
      },
      // parentNode: `${name}-container`,
      // extent: "parent" as const,
    };
    reactflowNodes.push(reactflowNode);
    dagreGraph.setNode(reactflowNode.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  }

  if (iface) {
    const ioNodes = getIONodes(name, iface, defaultPosition);

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

  const reactflowEdges: Edge[] = [];
  for (const connection of net!) {
    const { senderSide, receiverSide } = connection;
    if (!senderSide || !receiverSide) {
      continue;
    }

    const senderNode = senderSide.portAddr
      ? senderSide.portAddr.node
      : `${senderSide.constRef?.pkg}.${senderSide.constRef?.name}`;

    const senderOutport = senderSide.portAddr
      ? senderSide.portAddr.port
      : "out";

    for (const receiver of receiverSide) {
      const reactflowEdge = {
        // FIXME senderSide.portAddr and senderSide.constRef are formated like [object Object] (maybe do the same trick)
        id: `${name}-${senderSide.portAddr || senderSide.constRef} -> ${
          receiver.portAddr
        }`,
        source: `${name}-${senderNode}`!,
        sourceHandle: senderOutport,
        target: `${name}-${receiver.portAddr?.node!}`, // eslint-disable-line @typescript-eslint/no-non-null-asserted-optional-chain
        targetHandle: receiver.portAddr?.port,
        markerEnd: { type: MarkerType.Arrow },
        type: "smart",
        data: {
          isHighlighted: false,
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

    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };

    // if (node.parentNode) {
    //   const parentNodeWithPosition = dagreGraph.node(node.parentNode);
    //   node.position = {
    //     x: nodeWithPosition.x - parentNodeWithPosition.x,
    //     y: nodeWithPosition.y - parentNodeWithPosition.y,
    //   };
    // } else {
    //   node.position = {
    //     x: nodeWithPosition.x - nodeWidth / 2,
    //     y: nodeWithPosition.y - nodeHeight / 2,
    //   };
    // }

    return node;
  });

  return { nodes: reactflowNodes, edges: reactflowEdges };
};

function getIONodes(name: string, iface: src.Interface, position: XYPosition) {
  const defaultData = {
    type: "normal",
    position: position,
  };

  const inportsNode = {
    ...defaultData,
    id: `${name}-in`,
    data: {
      ports: {
        io: { out: {} },
      } as src.Interface,
      label: "in",
    },
  };
  for (const portName in iface!.io?.in) {
    const inport = iface!.io?.in[portName];
    inportsNode.data.ports.io!.out![portName] = inport;
  }

  const outportsNode = {
    ...defaultData,
    id: `${name}-out`,
    data: {
      ports: {
        io: { in: {} },
      } as src.Interface,
      label: "out",
    },
  };
  for (const portName in iface!.io?.out) {
    const outport = iface!.io?.out[portName];
    outportsNode.data.ports.io!.in![portName] = outport;
  }

  return {
    in: inportsNode,
    out: outportsNode,
  };
}
