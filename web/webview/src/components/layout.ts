import { Node, Edge } from "reactflow";
import ELK from "elkjs";
import * as src from "../generated/sourcecode";
import { ITypeNodeProps } from "./nodes/type_node";
import { IConstNodeProps } from "./nodes/const_node";
import { IInterfaceNodeProps } from "./nodes/interface_node";

const nodeWidth = 350;
const nodeHeight = 100;

const initialGroups = [
  {
    id: "types",
    width: nodeWidth,
    height: nodeHeight,
  },
  {
    id: "const",
    width: nodeWidth,
    height: nodeHeight,
  },
  {
    id: "interfaces",
    width: nodeWidth,
    height: nodeHeight,
  },
  {
    id: "components",
    width: nodeWidth,
    height: nodeHeight,
  },
];

const initialNodes = [
  // types group
  {
    id: "t1",
    group: "types",
    type: "type",
    data: {
      kind: src.TypeEntity,
      title: "T1",
      type: {
        meta: {
          text: "[]int",
        },
      },
    } as ITypeNodeProps,
  },
  // const group
  {
    id: "c1",
    group: "const",
    type: "const",
    data: {
      kind: src.ConstEntity,
      title: "C1",
      const: {
        value: {
          meta: {
            text: "42",
          },
          typeExpr: {
            meta: {
              text: "int",
            },
          },
        },
      },
    } as IConstNodeProps,
  },
  // interface group
  {
    id: "i1",
    group: "interface",
    type: "interface",
    data: {
      kind: src.InterfaceEntity,
      title: "I1",
      interface: {
        io: {
          in: {},
          out: {},
        },
      },
    } as IInterfaceNodeProps,
  },
  // component group
  {
    id: "cmp1",
    group: "component",
    type: "component",
    data: {
      kind: src.InterfaceEntity,
      title: "CMP1",
      interface: {
        io: {
          in: {},
          out: {},
        },
      },
    } as IInterfaceNodeProps,
  },
];

const initialEdges: any[] = [];

const elk = new ELK();

const graph = {
  id: "root",
  layoutOptions: {
    "elk.algorithm": "mrtree",
    "elk.direction": "DOWN",
  },
  children: initialGroups.map((group) => ({
    id: group.id,
    width: group.width,
    height: group.height,
    layoutOptions: {
      "elk.direction": "DOWN",
    },
    children: initialNodes
      .filter((node) => node.group === group.id)
      .map((node) => ({
        id: node.id,
        width: 100,
        height: 50,
        layoutOptions: {
          "elk.direction": "DOWN",
        },
      })),
  })),
  edges: initialEdges.map((edge) => ({
    id: edge.id,
    sources: [edge.source],
    targets: [edge.target],
  })),
};

export default async function createLayout(): Promise<{
  nodes: Node[];
  edges: Edge[];
}> {
  const layout = await elk.layout(graph);

  const nodes = layout.children!.reduce((result: any[], current) => {
    result.push({
      id: current.id,
      position: { x: current.x, y: current.y },
      data: { label: current.id },
      style: { width: current.width, height: current.height },
    });

    current.children!.forEach((child) =>
      result.push({
        id: child.id,
        position: { x: child.x, y: child.y },
        data: { label: child.id },
        style: { width: child.width, height: child.height },
        parentNode: current.id,
      })
    );

    return result;
  }, []);

  return {
    nodes,
    edges: initialEdges,
  };
}
