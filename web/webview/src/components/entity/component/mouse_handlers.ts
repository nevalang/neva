import { Node, Edge } from "reactflow";

export function handleNodeMouseEnter(
  hoveredNode: Node,
  nodesState: Node[],
  edgesState: Edge[]
): { nodes: Node[]; edges: Edge[] } {
  const newEdges: Edge[] = [];
  const relatedNodeIds: Set<string> = new Set();

  console.log(edgesState);

  edgesState.forEach((edge) => {
    const isEdgeRelated =
      edge.source === hoveredNode.id || edge.target === hoveredNode.id;
    const newEdge = isEdgeRelated
      ? {
          ...edge,
          data: {
            ...edge.data,
            isHighlighted: true,
          },
        }
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
            isRelated: true,
          },
        }
      : {
          ...node,
          data: {
            ...node.data,
            isDimmed: node.id !== hoveredNode.id,
          },
        }
  );

  return { edges: newEdges, nodes: newNodes };
}

export function handleNodeMouseLeave(
  nodesState: Node[],
  edgesState: Edge[]
): { nodes: Node[]; edges: Edge[] } {
  const newEdges = edgesState.map((edge) =>
    edge.data.isHighlighted
      ? {
          ...edge,
          data: {
            ...edge.data,
            isHighlighted: false,
          },
        }
      : edge
  );

  const newNodes = nodesState.map((node) => ({
    ...node,
    data: {
      ...node.data,
      isDimmed: false,
      isRelated: false,
    },
  }));

  return { edges: newEdges, nodes: newNodes };
}
