import { useCallback, useMemo, MouseEvent, useEffect } from "react";
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
import { buildReactFlowGraph } from "./helpers";
import { InterfaceNode } from "./nodes/interface_node";

interface ICanvasProps {
  fileViewState: FileViewState;
}

const edgeTypes = { normal: NormalEdge };
const nodeTypes = {
  // type: NormalNode,
  interface: InterfaceNode,
  // const: NormalNode,
  component: InterfaceNode, // component and interface nodes are the same at presentation level
};

export function Canvas(props: ICanvasProps) {
  const { nodes, edges } = useMemo(
    () => buildReactFlowGraph(props.fileViewState),
    [props.fileViewState]
  );

  const [nodesState, setNodesState, onNodesChange] = useNodesState(nodes);
  const [edgesState, setEdgesState, onEdgesChange] = useEdgesState(edges);
  useEffect(() => {
    setNodesState(nodes);
    setEdgesState(edges);
  }, [nodes, edges, setNodesState, setEdgesState]);

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

function handleNodeMouseEnter(
  hoveredNode: Node,
  edgesState: Edge[],
  nodesState: Node[]
) {
  const newEdges: Edge[] = [];
  const relatedNodeIds: Set<string> = new Set();

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
