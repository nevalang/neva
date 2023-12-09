import { NodeProps } from "reactflow";
import * as ts from "../../generated/typesystem";

export interface ITypeNodeProps {
  title: string;
  type: ts.Def;
}

export function TypeNode(props: NodeProps<ITypeNodeProps>) {
  return (
    <div className={"react-flow__node-default"}>
      <div className="nodeBody">
        <div className="nodeName">{props.data.title}</div>
        <div className="nodeType">{props.data.type.meta.text}</div>
      </div>
    </div>
  );
}
