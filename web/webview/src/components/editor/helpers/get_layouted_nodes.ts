import { Edge, Node } from "reactflow";
import ELK, { ElkNode } from "elkjs";

const nodeWidth = 150;
const nodeHeight = 100;

const elk = new ELK();

const nodeTypes = ["type", "const", "interface", "component"];

type NodeType = (typeof nodeTypes)[number];

const layoutOptions: { [key: NodeType]: object } = {
  type: { "elk.algorithm": "rectpacking" },
  const: { "elk.algorithm": "rectpacking" },
  interface: { "elk.algorithm": "rectpacking" },
  component: {
    "elk.algorithm": "mrtree",
    "elk.spacing.nodeNode": "50",
  },
};

export default async function getLayoutedNodes(
  nodes: Node[],
  edges: Edge[]
): Promise<Node[]> {
  const graph: ElkNode = {
    id: "root",
    layoutOptions: {
      "elk.algorithm": "box",
      "elk.direction": "DOWN",
      "elk.spacing.nodeNode": "50",
    },
    children: nodeTypes
      .map((nodeType) => ({
        id: nodeType,
        type: "parent",
        width: nodeWidth,
        height: nodeHeight,
        layoutOptions: {
          "elk.direction": "DOWN",
          "elk.spacing.nodeNode": "20",
          ...layoutOptions[nodeType],
        },
        children: nodes
          .filter((node) => node.type === nodeType)
          .map((node) => ({
            ...node,
            width: nodeWidth,
            height: nodeHeight,
          })),
      }))
      .filter((node) => node.children.length > 0),
    edges: edges.map((edge) => ({
      id: edge.id,
      sources: [edge.source],
      targets: [edge.target],
    })),
  };

  const layoutedGraph = await elk.layout(graph);

  const layoutedNodes: Node[] = [];
  for (const groupNode of layoutedGraph.children!) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const data = (groupNode as any).data;
    layoutedNodes.push({
      ...groupNode,
      id: groupNode.id,
      position: { x: groupNode.x!, y: groupNode.y! },
      data: { ...data, label: groupNode.id },
      style: { width: groupNode.width, height: groupNode.height },
    });
    for (const childNode of groupNode.children!) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = (childNode as any).data;
      layoutedNodes.push({
        ...childNode,
        id: childNode.id,
        position: { x: childNode.x!, y: childNode.y! },
        data: { ...data, label: childNode.id },
        style: { width: childNode.width, height: childNode.height },
        parentNode: groupNode.id,
      });
    }
  }

  return layoutedNodes;
}
