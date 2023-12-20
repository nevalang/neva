import { Edge, MarkerType, Node, XYPosition } from "reactflow";
import * as src from "../../../generated/sourcecode";
import {
  ComponentViewState,
  NodesViewState,
} from "../../../core/file_view_state";
import { IInterfaceNodeProps } from "../../flow/nodes/interface_node";

const defaultPosition = { x: 0, y: 0 };

export function buildComponentNetwork(
  entityName: string,
  component: ComponentViewState
) {
  const nodes: Node[] = [];
  const edges: Edge[] = [];

  buildAndInsertInterfaceNodes(component.interface!, entityName, nodes);
  buildAndInsertComponentNodes(entityName, component.nodes, nodes);
  buildAndInsertNetEdges(component.net, entityName, edges);

  return { nodes, edges };
}

function buildAndInsertNetEdges(
  net: src.Connection[],
  entityName: string,
  reactflowEdges: Edge[]
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
    }
  }
}

function buildAndInsertInterfaceNodes(
  iface: src.Interface,
  entityName: string,
  reactflowNodes: Node[]
) {
  const { inports, outports } = getComponentIONodes(
    entityName,
    iface,
    defaultPosition
  );
  reactflowNodes.push(inports);
  reactflowNodes.push(outports);
}

function buildAndInsertComponentNodes(
  entityName: string,
  nodes: NodesViewState[],
  reactflowNodes: Node[]
) {
  for (const nodeView of nodes) {
    buildAndInsertComponentNode(entityName, nodeView, reactflowNodes);
  }
}

function buildAndInsertComponentNode(
  entityName: string,
  nodeView: NodesViewState,
  reactflowNodes: Node[]
) {
  const reactflowNode = {
    id: `${entityName}-${nodeView.name}`,
    type: "component",
    position: defaultPosition,
    data: {
      title: nodeView.name,
      interface: nodeView.interface,
      isDimmed: false,
      isRelated: false,
      entityName: entityName,
    } as IInterfaceNodeProps,
  };
  reactflowNodes.push(reactflowNode);
}

function getComponentIONodes(
  entityName: string,
  iface: src.Interface,
  position: XYPosition
) {
  const defaultData = {
    type: "interface",
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
      isDimmed: false,
      isRelated: false,
      entityName: entityName,
    } as IInterfaceNodeProps,
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
      isDimmed: false,
      isRelated: false,
      entityName: entityName,
    } as IInterfaceNodeProps,
  };
  for (const portName in iface!.io?.out) {
    const outport = iface!.io?.out[portName];
    outportsNode.data.interface.io!.in![portName] = outport;
  }

  return {
    inports: inportsNode,
    outports: outportsNode,
  };
}
