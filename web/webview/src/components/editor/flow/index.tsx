import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  BackgroundVariant,
  Edge,
  Node,
  FitViewOptions,
} from "reactflow";
import "reactflow/dist/style.css";
import { InterfaceNode } from "../../interface_node";
import { TypeNode } from "./nodes/type_node";
import { ConstNode } from "./nodes/const_node";
import { useNavigate } from "react-router-dom";

const nodeTypes = {
  type: TypeNode,
  const: ConstNode,
  component: InterfaceNode, // component and interface nodes are the same at presentation level
  interface: InterfaceNode,
};

interface IFlowProps {
  nodes: Node[];
  edges: Edge[];
}

const fitViewOptions: FitViewOptions = {
  duration: 500,
  padding: 20,
  minZoom: 0.5,
  maxZoom: 1,
};

export function Flow(props: IFlowProps) {
  const navigate = useNavigate();

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      <ReactFlow
        nodeTypes={nodeTypes}
        nodes={props.nodes}
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
          navigate(`/${node.data.entityName}`);
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
