import { NodeProps } from "reactflow";
import * as ts from "../../../generated/typesystem";
import { useZoom } from "./use_zoom";
import classnames from "classnames";

export interface ITypeNodeProps {
  title: string;
  type: ts.Def;
}

export function TypeNode(props: NodeProps<ITypeNodeProps>) {
  const { isZoomMiddle } = useZoom();
  return (
    <div className={classnames("react-flow__node-default", props.type)}>
      <div className="nodeBody">
        <div className="nodeName">{props.data.title}</div>
        <div className="nodeType" style={{ opacity: isZoomMiddle ? 1 : 0 }}>
          {props.data.type.meta.text}
        </div>
      </div>
    </div>
  );
}
