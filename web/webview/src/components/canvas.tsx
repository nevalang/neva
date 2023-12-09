import { useCallback, MouseEvent, useEffect, useState } from "react";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  BackgroundVariant,
  Edge,
  Node,
  useNodesState,
  useEdgesState,
} from "reactflow";
import "reactflow/dist/style.css";
import * as src from "../generated/sourcecode";
import { FileViewState } from "../core/file_view_state";
import { NormalEdge } from "./edge";
import { InterfaceNode } from "./nodes/interface_node";
import { TypeNode } from "./nodes/type_node";
import { ConstNode } from "./nodes/const_node";
import getLayoutedNodes from "./get_layouted_nodes";
import { buildGraph } from "./build_graph";
import { handleNodeMouseEnter, handleNodeMouseLeave } from "./mouse_handlers";

const edgeTypes = { normal: NormalEdge };
const nodeTypes = {
  type: TypeNode,
  const: ConstNode,
  component: InterfaceNode, // component and interface nodes are the same at presentation level
  interface: InterfaceNode,
};

interface ICanvasProps {
  fileViewState: FileViewState;
}

export function Canvas(props: ICanvasProps) {
  const [graph, setGraph] = useState<{ nodes: Node[]; edges: Edge[] }>({
    nodes: [],
    edges: [],
  });

  useEffect(() => {
    (async () => {
      const graph = buildGraph(props.fileViewState);
      console.log({ graph });
      const layoutedNodes = await getLayoutedNodes(graph.nodes, graph.edges);
      console.log({ layoutedNodes });
      setGraph({ nodes: layoutedNodes, edges: graph.edges });
    })();
  }, [props.fileViewState]);

  const [nodesState, setNodesState, onNodesChange] = useNodesState(graph.nodes);
  const [edgesState, setEdgesState, onEdgesChange] = useEdgesState(graph.edges);
  useEffect(() => {
    setNodesState(graph.nodes);
    setEdgesState(graph.edges);
  }, [graph, setNodesState, setEdgesState]);

  const onNodeMouseEnter = useCallback(
    (_: MouseEvent, hoveredNode: Node) => {
      if (hoveredNode.data.kind != src.ComponentEntity) {
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
