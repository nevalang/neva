import { Node, Edge, MarkerType, XYPosition, Position } from "reactflow";
import * as src from "../generated/sourcecode";
import { ComponentViewState, FileViewState } from "../core/file_view_state";
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

  // TODO
  // fileViewState.entities.types
  // fileViewState.entities.constants
  // fileViewState.entities.interfaces

  const nodes: Node[] = [];
  const edges: Edge[] = [];

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
  componentViewState: ComponentViewState,
  reactflowNodes: Node[],
  reactflowEdges: Edge[],
  dagreGraph: dagre.graphlib.Graph
) {
  const { nodes, interface: iface, net } = componentViewState;

  for (const nodeView of nodes) {
    const reactflowNode = {
      id: `${entityName}-${nodeView.name}`,
      type: "normal",
      position: defaultPosition,
      data: {
        ports: nodeView.interface,
        label: nodeView.name,
      },
    };
    reactflowNodes.push(reactflowNode);
    dagreGraph.setNode(reactflowNode.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  }

  if (iface) {
    const ioNodes = getIONodes(entityName, iface, defaultPosition);

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

  return { nodes: reactflowNodes, edges: reactflowEdges };
}

function getIONodes(name: string, iface: src.Interface, position: XYPosition) {
  const defaultData = {
    type: "normal",
    position: position,
  };

  const inportsNode = {
    ...defaultData,
    id: `${name}-in`,
    data: {
      ports: {
        io: { out: {} },
      } as src.Interface,
      label: "in",
    },
  };
  for (const portName in iface!.io?.in) {
    const inport = iface!.io?.in[portName];
    inportsNode.data.ports.io!.out![portName] = inport;
  }

  const outportsNode = {
    ...defaultData,
    id: `${name}-out`,
    data: {
      ports: {
        io: { in: {} },
      } as src.Interface,
      label: "out",
    },
  };
  for (const portName in iface!.io?.out) {
    const outport = iface!.io?.out[portName];
    outportsNode.data.ports.io!.in![portName] = outport;
  }

  return {
    in: inportsNode,
    out: outportsNode,
  };
}
