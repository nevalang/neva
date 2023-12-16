import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  BackgroundVariant,
  Edge,
  Node,
  FitViewOptions,
  NodeTypes,
} from "reactflow";
import "reactflow/dist/style.css";
import { getLayoutedNodes } from "../editor/get_layouted_nodes";

interface IFlowProps {
  nodes: Node[];
  edges: Edge[];
  nodeTypes: NodeTypes;
}

const fitViewOptions: FitViewOptions = {
  duration: 0,
  padding: 20,
  minZoom: 0.5,
  maxZoom: 1,
};

export function Flow(props: IFlowProps) {
  const navigate = useNavigate();

  const [layoutedNodes, setLayoutedNodes] = useState<Node[]>([]);
  useEffect(() => {
    (async () => {
      setLayoutedNodes(await getLayoutedNodes(props.nodes, props.edges));
    })();
  }, [props.nodes, props.edges]);

  if (layoutedNodes.length === 0) {
    return null;
  }

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ReactFlow
        nodeTypes={props.nodeTypes}
        nodes={layoutedNodes}
        edges={props.edges}
        fitView
        fitViewOptions={fitViewOptions}
        nodesFocusable
        panOnScroll
        zoomOnScroll={false}
        elementsSelectable={false}
        nodesDraggable={false}
        nodesConnectable={false}
        onNodeClick={(_, node: Node) => {
          if (node.type !== "parent") {
            navigate(`/${node.data.entityName}`);
          }
        }}
        minZoom={0.3}
        maxZoom={2}
      >
        <Controls fitViewOptions={fitViewOptions} />
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
    interface: "var(--foreground)",
    component: "var(--component)",
    parent: "var(--background)",
  }[nodeType.type!]!);
