import { NodeProps } from "reactflow";
import * as src from "../../../../generated/sourcecode";
import { InterfaceNode } from "../../../interface_node";

export interface IInterfaceNodeProps {
  title: string;
  component: src.Component;
  isDimmed: boolean;
  isRelated: boolean;
  entityName: string;
}

export function ComponentNode(props: NodeProps<IInterfaceNodeProps>) {
  const interfaceProps = {
    ...props,
    data: {
      ...props.data,
      interface: props.data.component.interface!,
    },
  };
  return <InterfaceNode {...interfaceProps} />;
}
