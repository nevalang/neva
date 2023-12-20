import { ReactNode, useCallback, useEffect } from "react";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  BackgroundVariant,
  Edge,
  Node,
  FitViewOptions,
  NodeTypes,
  useNodesState,
  useEdgesState,
  Panel,
  EdgeTypes,
} from "reactflow";
import "reactflow/dist/style.css";

interface IFlowProps {
  nodes: Node[];
  edges: Edge[];
  nodeTypes: NodeTypes;
  edgeTypes?: EdgeTypes;
  title: string;
  onNodeClick?: (node: Node) => void;
  nodesDraggable?: boolean;
  leftTopPanel?: ReactNode;
  onNodeMouseEnter?: (
    node: Node,
    nodes: Node[],
    edges: Edge[]
  ) => { nodes: Node[]; edges: Edge[] };
  onNodeMouseLeave?: (
    nodes: Node[],
    edges: Edge[]
  ) => { nodes: Node[]; edges: Edge[] };
}

const defaultFitViewOptions: FitViewOptions = {
  duration: 0,
  padding: 20,
  minZoom: 0.5,
  maxZoom: 1,
};

const fitViewControlOptions: FitViewOptions = {
  ...defaultFitViewOptions,
  duration: 300,
};

export function Flow(props: IFlowProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState(props.nodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(props.edges);

  useEffect(() => {
    setNodes(props.nodes);
    setEdges(props.edges);
  }, [props.edges, props.nodes, setEdges, setNodes]);

  const handleOnMouseEnter = useCallback(
    (_: unknown, node: Node): void => {
      if (!props.onNodeMouseEnter) {
        return;
      }
      const { nodes: newNodes, edges: newEdges } = props.onNodeMouseEnter(
        node,
        nodes,
        edges
      );
      setNodes(newNodes);
      setEdges(newEdges);
    },
    [edges, nodes, props, setEdges, setNodes]
  );

  const handleOnMouseLeave = useCallback(() => {
    if (!props.onNodeMouseLeave) {
      return;
    }
    const { nodes: newNodes, edges: newEdges } = props.onNodeMouseLeave(
      nodes,
      edges
    );
    setNodes(newNodes);
    setEdges(newEdges);
  }, [edges, nodes, props, setEdges, setNodes]);

  if (nodes.length === 0) {
    return null;
  }

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeMouseEnter={handleOnMouseEnter}
        onNodeMouseLeave={handleOnMouseLeave}
        nodeTypes={props.nodeTypes}
        edgeTypes={props.edgeTypes}
        fitView
        fitViewOptions={defaultFitViewOptions}
        nodesFocusable
        panOnScroll
        zoomOnScroll={false}
        elementsSelectable={false}
        nodesDraggable={Boolean(props.nodesDraggable)}
        nodesConnectable={false}
        onNodeClick={(_, node: Node) =>
          props.onNodeClick && props.onNodeClick(node)
        }
        minZoom={0.3}
        maxZoom={2}
      >
        <Panel position="top-center">{props.title}</Panel>
        <Panel position="top-left">{props.leftTopPanel}</Panel>
        <Controls fitViewOptions={fitViewControlOptions} />
        <MiniMap
          position="top-right"
          zoomable
          pannable
          nodeStrokeWidth={3}
          nodeColor={nodeColor}
          nodeBorderRadius={10}
          maskColor="rgba(255, 255, 255, 0.1)"
          maskStrokeColor="var(--text)"
          nodeStrokeColor="var(--light)"
        />
        <Background variant={BackgroundVariant.Dots} gap={10} size={0.5} />
      </ReactFlow>
    </div>
  );
}

const nodeColor = (nodeType: Node): string =>
  ({
    type: "var(--type)",
    const: "var(--const)",
    interface: "var(--interface)",
    component: "var(--component)",
    parent: "var(--dark)",
  }[nodeType.type!]!);
