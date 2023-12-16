import { useMemo } from "react";
import classnames from "classnames";
import { Handle, NodeProps, HandleType, Position } from "reactflow";
import * as src from "../../../generated/sourcecode";
import { useZoom } from "./use_zoom";

export interface IInterfaceNodeProps {
  title: string;
  interface: src.Interface;
  isDimmed: boolean;
  isRelated: boolean;
  entityName: string;
}

export function InterfaceNode(props: NodeProps<IInterfaceNodeProps>) {
  const { inports, outports } = useMemo(() => {
    const result = { inports: [], outports: [] };
    if (!props.data.interface.io) {
      return result;
    }
    return {
      inports: Object.entries(props.data.interface.io.in || {}),
      outports: Object.entries(props.data.interface.io.out || {}),
    };
  }, [props.data.interface.io]);

  const { isZoomMiddle, isZoomClose } = useZoom();

  return (
    <div
      className={classnames("react-flow__node-default", props.type, {
        related: props.data.isRelated,
        dimmed: props.data.isDimmed,
      })}
    >
      <Ports
        direction="in"
        ports={inports}
        position={Position.Top}
        type="target"
        isVisible={isZoomMiddle}
        areTypesVisible={isZoomClose}
      />
      <div className="nodeBody">
        <div className="nodeName">{props.data.title}</div>
      </div>
      <Ports
        direction="out"
        ports={outports}
        position={Position.Bottom}
        type="source"
        isVisible={isZoomMiddle}
        areTypesVisible={isZoomClose}
      />
    </div>
  );
}

function Ports(props: {
  ports: [string, src.Port][];
  position: Position;
  type: HandleType;
  isVisible: boolean;
  areTypesVisible: boolean;
  direction: "in" | "out";
}) {
  if (!props.ports) {
    return null;
  }

  return (
    <div
      className={classnames("ports", props.direction, {
        hidden: !props.isVisible,
      })}
    >
      {props.ports.map(([portName, portType]) => (
        <Handle
          id={portName}
          type={props.type}
          position={props.position}
          isConnectable={false}
          key={portName}
        >
          {portName}
          {props.areTypesVisible &&
            portType.typeExpr &&
            portType.typeExpr.meta && (
              <span className="portType">
                {" "}
                {(portType.typeExpr.meta as src.Meta).text}
              </span>
            )}
        </Handle>
      ))}
    </div>
  );
}
