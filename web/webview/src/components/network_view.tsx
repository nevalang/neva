import { useCallback, useMemo, MouseEvent } from "react";
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
  getBezierPath,
} from "reactflow";
import "reactflow/dist/style.css";
import dagre from "dagre";
import * as src from "../generated/sourcecode";
import { ComponentViewState } from "../core/file_view_state";
import classnames from "classnames";

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

  const [nodesState, setNodesState, onNodesChange] = useNodesState(nodes);
  const [edgesState, setEdgesState, onEdgesChange] = useEdgesState(edges);

  const onNodeMouseEnter = useCallback(
    (_: MouseEvent, hoveredNode: Node) => {
      const { newEdges, newNodes } = handleNodeMouseEnter(
        hoveredNode,
        edgesState,
        nodesState
      );
      setEdgesState(newEdges);
      setNodesState(newNodes);
    },
    [edgesState, nodesState, setEdgesState, setNodesState]
  );

  const onNodeMouseLeave = useCallback(() => {
    const { newEdges, newNodes } = handleNodeMouseLeave(edgesState, nodesState);
    setEdgesState(newEdges);
    setNodesState(newNodes);
  }, [edgesState, nodesState, setEdgesState, setNodesState]);

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
        onNodeMouseEnter={onNodeMouseEnter}
        onNodeMouseLeave={onNodeMouseLeave}
      >
        <Controls />
        <MiniMap />
        <Background variant={BackgroundVariant.Dots} gap={10} size={0.5} />
      </ReactFlow>
    </div>
  );
}

function handleNodeMouseLeave(edgesState: Edge[], nodesState: Node[]) {
  const newEdges = edgesState.map((edge) =>
    edge.data.isHighlighted ? { ...edge, data: { isHighlighted: false } } : edge
  );
  const newNodes = nodesState.map((node) => ({
    ...node,
    data: {
      ...node.data,
      isDimmed: false,
      isHighlighted: false,
    },
  }));
  return { newEdges, newNodes };
}

function handleNodeMouseEnter(
  hoveredNode: Node,
  edgesState: Edge[],
  nodesState: Node[]
) {
  const newEdges: Edge[] = [];
  const relatedNodeIds: Set<string> = new Set();
  relatedNodeIds.add(hoveredNode.id);

  edgesState.forEach((edge) => {
    const isEdgeRelated =
      edge.source === hoveredNode.id || edge.target === hoveredNode.id;
    const newEdge = isEdgeRelated
      ? { ...edge, data: { isHighlighted: true } }
      : edge;
    newEdges.push(newEdge);
    if (isEdgeRelated) {
      const isIncoming = edge.source === hoveredNode.id;
      const relatedNodeId = isIncoming ? edge.target : edge.source;
      relatedNodeIds.add(relatedNodeId);
    }
  });

  const newNodes = nodesState.map((node) =>
    relatedNodeIds.has(node.id)
      ? {
          ...node,
          data: {
            ...node.data,
            isHighlighted: true,
          },
        }
      : { ...node, data: { ...node.data, isDimmed: true } }
  );

  return { newEdges, newNodes };
}

function NormalEdge(props: EdgeProps<{ isHighlighted: boolean }>) {
  const [edgePath] = getBezierPath({
    sourceX: props.sourceX,
    sourceY: props.sourceY,
    targetX: props.targetX,
    targetY: props.targetY,
  });

  const style = props.data?.isHighlighted
    ? { strokeOpacity: 1, stroke: "white" }
    : { strokeOpacity: 0.75 };

  return <BaseEdge {...props} path={edgePath} style={style} />;
}

function NormalNode(
  props: NodeProps<{
    ports: src.Interface;
    label: string;
    isHighlighted: boolean;
    isDimmed: boolean;
  }>
) {
  const { io } = props.data.ports;
  const arePortsVisible = useStore((s) => s.transform[2] >= 0.6);
  const areTypesVisible = useStore((s) => s.transform[2] >= 1);

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

  const cn = classnames("react-flow__node-default", {
    highlighted: props.data.isHighlighted,
    dimmed: props.data.isDimmed,
  });

  return (
    <div className={cn}>
      {inports.length > 0 && (
        <div
          className={classnames("ports", "in", { hidden: !arePortsVisible })}
        >
          {inports.map(([inportName, inportType]) => (
            <Handle
              content="asd"
              type="target"
              id={inportName}
              key={inportName}
              position={Position.Top}
              isConnectable={true}
            >
              {inportName}
              {areTypesVisible &&
                inportType.typeExpr &&
                inportType.typeExpr.meta && (
                  <span className="portType">
                    {" "}
                    {(inportType.typeExpr.meta as src.Meta).text}
                  </span>
                )}
            </Handle>
          ))}
        </div>
      )}
      <div className="nodeName">{props.data.label}</div>
      {outports.length > 0 && (
        <div
          className={classnames("ports", "out", { hidden: !arePortsVisible })}
        >
          {outports.map(([outportName, outportType]) => (
            <Handle
              type="source"
              id={outportName}
              key={outportName}
              position={Position.Bottom}
              isConnectable={true}
            >
              {outportName}
              {areTypesVisible &&
                outportType.typeExpr &&
                outportType.typeExpr.meta && (
                  <span className="portType">
                    {" "}
                    {(outportType.typeExpr.meta as src.Meta).text}{" "}
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
  entityName: string,
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

  const reactflowNodes: Node[] = [];

  for (const nodeView of nodes) {
    const reactflowNode = {
      id: `${entityName}-${nodeView.name}`,
      type: "normal",
      position: defaultPosition,
      data: {
        ports: nodeView.interface,
        label: nodeView.name,
      },
    };
    reactflowNodes.push(reactflowNode);
    dagreGraph.setNode(reactflowNode.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  }

  if (iface) {
    const ioNodes = getIONodes(entityName, iface, defaultPosition);

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
      const senderPart = senderSide.portAddr
        ? senderSide.portAddr.meta?.text
        : senderSide.constRef?.meta?.text;

      const reactflowEdge = {
        id: `${entityName}-${senderPart} -> ${receiver.portAddr?.meta?.text}`,
        source: `${entityName}-${senderNode}`,
        sourceHandle: senderOutport,
        target: `${entityName}-${receiver.portAddr?.node!}`, // eslint-disable-line @typescript-eslint/no-non-null-asserted-optional-chain
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
