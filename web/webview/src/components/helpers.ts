import { Node, Edge, MarkerType, XYPosition, Position } from "reactflow";
import * as src from "../generated/sourcecode";
import * as ts from "../generated/typesystem";
import {
  ComponentViewState,
  FileViewState,
  NodesViewState,
} from "../core/file_view_state";
import dagre from "dagre";

const defaultPosition = { x: 0, y: 0 };
const nodeWidth = 342.5;
const nodeHeight = 70;
const direction = "TB";
const isHorizontal = false;

export function buildReactFlowGraph(fileViewState: FileViewState) {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));
  dagreGraph.setGraph({ rankdir: direction });

  const nodes: Node[] = [];
  const edges: Edge[] = [];

  for (const typeDef of fileViewState.entities.types) {
    buildAndInsertTypeDefNode(typeDef.name, typeDef.entity, nodes, dagreGraph);
  }

  for (const constant of fileViewState.entities.constants) {
    buildAndInsertConstNode(constant.name, constant.entity, nodes, dagreGraph);
  }

  for (const iface of fileViewState.entities.interfaces) {
    buildAndInsertInterfaceNode(iface.name, iface.entity, nodes, dagreGraph);
  }

  for (const component of fileViewState.entities.components) {
    buildAndInsertComponentSubgraph(
      component.name,
      component.entity,
      nodes,
      edges,
      dagreGraph
    );
  }

  dagre.layout(dagreGraph);

  nodes.forEach((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.targetPosition = (isHorizontal ? "left" : "top") as Position;
    node.sourcePosition = (isHorizontal ? "right" : "bottom") as Position;

    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };

    return node;
  });

  return {
    nodes: nodes,
    edges: edges,
  };
}

function buildAndInsertComponentSubgraph(
  entityName: string,
  component: ComponentViewState,
  reactflowNodes: Node[],
  reactflowEdges: Edge[],
  dagreGraph: dagre.graphlib.Graph
) {
  if (component.interface) {
    buildAndInsertInterfaceNodes(
      component.interface,
      entityName,
      reactflowNodes,
      dagreGraph
    );
  }
  buildAndInsertComponentNodes(
    entityName,
    component.nodes,
    reactflowNodes,
    dagreGraph
  );
  buildAndInsertNetEdges(component.net, entityName, reactflowEdges, dagreGraph);
}

function buildAndInsertNetEdges(
  net: src.Connection[],
  entityName: string,
  reactflowEdges: Edge[],
  dagreGraph: dagre.graphlib.Graph
) {
  for (const connection of net) {
    const { senderSide, receiverSide } = connection;
    if (!senderSide || !receiverSide) {
      continue;
    }

    const senderNode = senderSide.portAddr
      ? senderSide.portAddr.node
      : `${senderSide.constRef?.pkg}.${senderSide.constRef?.name}`;

    const senderOutport = senderSide.portAddr
      ? senderSide.portAddr.port
      : "out";

    for (const receiver of receiverSide) {
      const senderPart = senderSide.portAddr
        ? senderSide.portAddr.meta?.text
        : senderSide.constRef?.meta?.text;

      const reactflowEdge = {
        id: `${entityName}-${senderPart} -> ${receiver.portAddr?.meta?.text}`,
        source: `${entityName}-${senderNode}`,
        sourceHandle: senderOutport,
        target: `${entityName}-${receiver.portAddr!.node!}`,
        targetHandle: receiver.portAddr?.port,
        markerEnd: { type: MarkerType.Arrow },
        type: "normal",
        data: {
          isHighlighted: false,
        },
      };

      reactflowEdges.push(reactflowEdge);
      dagreGraph.setEdge(reactflowEdge.source, reactflowEdge.target);
    }
  }
}

function buildAndInsertInterfaceNodes(
  iface: src.Interface,
  entityName: string,
  reactflowNodes: Node[],
  dagreGraph: dagre.graphlib.Graph
) {
  const ioNodes = getComponentIONodes(entityName, iface, defaultPosition);

  reactflowNodes.push(ioNodes.in);
  dagreGraph.setNode(ioNodes.in.id, {
    width: nodeWidth,
    height: nodeHeight,
  });

  reactflowNodes.push(ioNodes.out);
  dagreGraph.setNode(ioNodes.out.id, {
    width: nodeWidth,
    height: nodeHeight,
  });
}

function buildAndInsertComponentNodes(
  entityName: string,
  nodes: NodesViewState[],
  reactflowNodes: Node[],
  dagreGraph: dagre.graphlib.Graph
) {
  for (const nodeView of nodes) {
    buildAndInsertComponentNode(
      entityName,
      nodeView,
      reactflowNodes,
      dagreGraph
    );
  }
}

function buildAndInsertTypeDefNode(
  entityName: string,
  typeDef: ts.Def,
  reactflowNodes: Node[],
  dagreGraph: dagre.graphlib.Graph
) {
  const reactflowNode = {
    id: entityName,
    type: "type",
    position: defaultPosition,
    data: {
      kind: src.TypeEntity,
      title: entityName,
      typeDef: typeDef,
    },
  };
  reactflowNodes.push(reactflowNode);
  dagreGraph.setNode(reactflowNode.id, {
    width: nodeWidth,
    height: nodeHeight,
  });
}

function buildAndInsertConstNode(
  entityName: string,
  constant: src.Const,
  reactflowNodes: Node[],
  dagreGraph: dagre.graphlib.Graph
) {
  const reactflowNode = {
    id: entityName,
    type: "const",
    position: defaultPosition,
    data: {
      kind: src.ConstEntity,
      title: entityName,
      constant: constant,
    },
  };
  reactflowNodes.push(reactflowNode);
  dagreGraph.setNode(reactflowNode.id, {
    width: nodeWidth,
    height: nodeHeight,
  });
}

function buildAndInsertInterfaceNode(
  entityName: string,
  iface: src.Interface,
  reactflowNodes: Node[],
  dagreGraph: dagre.graphlib.Graph
) {
  const reactflowNode = {
    id: entityName,
    type: "interface",
    position: defaultPosition,
    data: {
      kind: src.InterfaceEntity,
      title: entityName,
      interface: iface,
    },
  };
  reactflowNodes.push(reactflowNode);
  dagreGraph.setNode(reactflowNode.id, {
    width: nodeWidth,
    height: nodeHeight,
  });
}

function buildAndInsertComponentNode(
  entityName: string,
  nodeView: NodesViewState,
  reactflowNodes: Node[],
  dagreGraph: dagre.graphlib.Graph
) {
  const reactflowNode = {
    id: `${entityName}-${nodeView.name}`,
    type: "component",
    position: defaultPosition,
    data: {
      kind: src.ComponentEntity,
      title: nodeView.name,
      interface: nodeView.interface,
    },
  };
  reactflowNodes.push(reactflowNode);
  dagreGraph.setNode(reactflowNode.id, {
    width: nodeWidth,
    height: nodeHeight,
  });
}

function getComponentIONodes(
  entityName: string,
  iface: src.Interface,
  position: XYPosition
) {
  const defaultData = {
    type: "component",
    position: position,
  };

  const inportsNode = {
    ...defaultData,
    id: `${entityName}-in`,
    data: {
      interface: {
        io: { out: {} },
      } as src.Interface,
      title: "in",
      kind: src.ComponentEntity,
    },
  };
  for (const portName in iface!.io?.in) {
    const inport = iface!.io?.in[portName];
    inportsNode.data.interface.io!.out![portName] = inport;
  }

  const outportsNode = {
    ...defaultData,
    id: `${entityName}-out`,
    data: {
      interface: {
        io: { in: {} },
      } as src.Interface,
      title: "out",
      kind: src.ComponentEntity,
    },
  };
  for (const portName in iface!.io?.out) {
    const outport = iface!.io?.out[portName];
    outportsNode.data.interface.io!.in![portName] = outport;
  }

  return {
    in: inportsNode,
    out: outportsNode,
  };
}
