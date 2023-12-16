import { NodeProps } from "reactflow";
import * as src from "../../../generated/sourcecode";
import { useZoom } from "./use_zoom";
import classnames from "classnames";

export interface IConstNodeProps {
  title: string;
  const: src.Const;
}

export function ConstNode(props: NodeProps<IConstNodeProps>) {
  const { isZoomMiddle, isZoomClose } = useZoom();

  return (
    <div className={classnames("react-flow__node-default", props.type)}>
      <div className="nodeBody">
        <div className="nodeName">{props.data.title}</div>
        <div className="nodeType" style={{ opacity: isZoomClose ? 1 : 0 }}>
          {formatConstType(props.data.const)}
        </div>
        <div className="nodeType" style={{ opacity: isZoomMiddle ? 1 : 0 }}>
          {formatConstValue(props.data.const)}
        </div>
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
