import { Node, Edge } from "reactflow";
import * as src from "../../generated/sourcecode";
import * as ts from "../../generated/typesystem";
import { FileViewState } from "../../core/file_view_state";
import { ITypeNodeProps } from "../flow/nodes/type_node";
import { IConstNodeProps } from "../flow/nodes/const_node";
import { IInterfaceNodeProps } from "../flow/nodes/interface_node";

const defaultPosition = { x: 0, y: 0 };

export function buildFileNodes(fileViewState: FileViewState): Node[] {
  const nodes: Node[] = [];

  for (const typeDef of fileViewState.entities.types) {
    handleTypeNode(typeDef.name, typeDef.entity, nodes);
  }

  for (const constant of fileViewState.entities.constants) {
    handleConstNode(constant.name, constant.entity, nodes);
  }

  for (const iface of fileViewState.entities.interfaces) {
    handleInterfaceNode(iface.name, iface.entity, nodes, "interface");
  }

  for (const component of fileViewState.entities.components) {
    handleInterfaceNode(
      component.name,
      component.entity.interface!,
      nodes,
      "component"
    );
  }

  return nodes;
}

function handleTypeNode(
  entityName: string,
  typeDef: ts.Def,
  reactflowNodes: Node[]
) {
  const reactflowNode = {
    id: entityName,
    type: "type",
    position: defaultPosition,
    data: {
      title: entityName,
      type: typeDef,
      entityName: entityName,
    } as ITypeNodeProps,
  };
  reactflowNodes.push(reactflowNode);
}

function handleConstNode(
  entityName: string,
  constant: src.Const,
  reactflowNodes: Node[]
) {
  const reactflowNode = {
    id: entityName,
    type: "const",
    position: defaultPosition,
    data: {
      title: entityName,
      const: constant,
      entityName: entityName,
    } as IConstNodeProps,
  };
  reactflowNodes.push(reactflowNode);
}

function handleInterfaceNode(
  entityName: string,
  iface: src.Interface,
  reactflowNodes: Node[],
  type: "component" | "interface"
) {
  const reactflowNode = {
    id: entityName,
    type: type,
    position: defaultPosition,
    data: {
      title: entityName,
      interface: iface,
      isDimmed: false,
      isRelated: false,
      entityName: entityName,
    } as IInterfaceNodeProps,
  };
  reactflowNodes.push(reactflowNode);
}
