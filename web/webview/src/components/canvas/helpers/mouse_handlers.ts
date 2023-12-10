import { Node, Edge } from "reactflow";

export function handleNodeMouseEnter(
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

export function handleNodeMouseLeave(edgesState: Edge[], nodesState: Node[]) {
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
