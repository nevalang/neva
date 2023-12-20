import { Node } from "reactflow";
import ELK, { ElkNode } from "elkjs";

const nodeWidth = 150;
const nodeHeight = 100;

const elk = new ELK();

const parents = ["component", "type", "const", "interface"]; // node-types

export async function getLayoutedNodes(nodes: Node[]): Promise<Node[]> {
  const graph: ElkNode = {
    id: "root",
    layoutOptions: {
      "elk.algorithm": "rectpacking",
      "elk.direction": "DOWN",
      "elk.spacing.nodeNode": "100",
    },
    children: parents
      .map((parentNodeId) => ({
        id: parentNodeId,
        type: "parent",
        width: nodeWidth,
        height: nodeHeight,
        layoutOptions: {
          "elk.direction": "DOWN",
          "elk.spacing.nodeNode": "20",
          "elk.algorithm": "rectpacking",
          "elk.padding": "[top=40,left=20,bottom=20,right=20]",
        },
        children: nodes
          .filter((node) => node.type === parentNodeId)
          .map((node) => ({
            ...node,
            width: nodeWidth,
            height: nodeHeight,
          })),
      }))
      .filter((node) => node.children.length > 0),
  };

  const layoutedGraph = await elk.layout(graph);

  const layoutedNodes: Node[] = [];
  for (const parentNode of layoutedGraph.children!) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const data = (parentNode as any).data;
    layoutedNodes.push({
      ...parentNode,
      id: parentNode.id,
      position: { x: parentNode.x!, y: parentNode.y! },
      data: { ...data, label: labels[parentNode.id] },
      style: { width: parentNode.width, height: parentNode.height },
    });
    for (const childNode of parentNode.children!) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = (childNode as any).data;
      layoutedNodes.push({
        ...childNode,
        id: childNode.id,
        position: { x: childNode.x!, y: childNode.y! },
        data: { ...data, label: childNode.id },
        style: { width: childNode.width, height: childNode.height },
        parentNode: parentNode.id,
      });
    }
  }

  return layoutedNodes;
}

const labels: { [key: string]: string } = {
  type: "TYPES",
  const: "CONSTANTS",
  interface: "INTERFACES",
  component: "COMPONENTS",
};
