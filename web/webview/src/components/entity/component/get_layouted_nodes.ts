import dagre from "dagre";
import { Node, Edge, Position } from "reactflow";

export type Graph = { nodes: Node[]; edges: Edge[] };

export function getLayoutedNodes(graph: Graph): Graph {
  const dagreGraph = new dagre.graphlib.Graph({ directed: true });
  dagreGraph.setDefaultEdgeLabel(() => ({}));

  const nodeWidth = 200;
  const nodeHeight = 50;

  dagreGraph.setGraph({
    rankdir: "TB",
    nodesep: 500,
    edgesep: 200,
    ranksep: 100,
    ranker: "longest-path",
  });

  graph.nodes.forEach((node) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });
  graph.edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  graph.nodes.forEach((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.targetPosition = "top" as Position;
    node.sourcePosition = "bottom" as Position;

    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };
  });

  return { nodes: graph.nodes, edges: graph.edges };
}
