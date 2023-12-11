import { useCallback, MouseEvent } from "react";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  BackgroundVariant,
  Edge,
  Node,
  useNodesState,
  useEdgesState,
  FitViewOptions,
  ReactFlowInstance,
} from "reactflow";
import "reactflow/dist/style.css";
import { NormalEdge } from "./edge";
import { InterfaceNode } from "./nodes/interface_node";
import { TypeNode } from "./nodes/type_node";
import { ConstNode } from "./nodes/const_node";
import {
  handleNodeMouseEnter,
  handleNodeMouseLeave,
} from "./helpers/mouse_handlers";

const edgeTypes = { normal: NormalEdge };
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
  duration: 300,
  padding: 20,
  minZoom: 0.5,
  maxZoom: 1,
};

export function Flow(props: IFlowProps) {
  console.log("=== props", props.nodes, props.edges);

  const [nodesState, setNodesState, onNodesChange] = useNodesState(props.nodes);
  const [edgesState, setEdgesState, onEdgesChange] = useEdgesState(props.edges);

  console.log("=== state", nodesState, edgesState);

  const onNodeMouseEnter = useCallback(
    (_: MouseEvent, hoveredNode: Node) => {
      if (hoveredNode.type !== "component") {
        return;
      }
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
        nodes={nodesState}
        edges={edgesState}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeMouseEnter={onNodeMouseEnter}
        onNodeMouseLeave={onNodeMouseLeave}
        fitView
        // onInit={fitView}
        fitViewOptions={fitViewOptions}
        nodesFocusable
        panOnScroll
        zoomOnScroll={false}
        elementsSelectable={false}
        nodesDraggable={false}
        nodesConnectable={false}
        minZoom={0.1}
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
