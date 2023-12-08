import { NodeProps } from "reactflow";
import * as src from "../../generated/sourcecode";

export interface IConstNodeProps {
  title: string;
  const: src.Const;
}

export function ConstNode(props: NodeProps<IConstNodeProps>) {
  return (
    <div className={"react-flow__node-default"}>
      <div className="nodeBody">
        <div className="nodeName">{props.data.title}</div>
        <div className="nodeType">{formatConstType(props.data.const)}</div>
        <div className="nodeType">{formatConstValue(props.data.const)}</div>
      </div>
    </div>
  );
}

function formatConstValue(constant: src.Const): string {
  if (constant.ref) {
    return constant.ref.meta!.text!;
  }
  return constant.value?.meta?.text || "";
}

function formatConstType(constant: src.Const): string {
  if (constant.ref) {
    return "";
  }
  return constant.value?.typeExpr?.meta?.text || "";
}
