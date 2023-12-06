import { Handle, NodeProps, Position, useStore } from "reactflow";
import * as src from "../generated/sourcecode";
import { useMemo } from "react";
import classnames from "classnames";

export function NormalNode(
  props: NodeProps<{
    ports: src.Interface;
    label: string;
    isHighlighted: boolean;
    isDimmed: boolean;
  }>
) {
  const { io } = props.data.ports;
  const arePortsVisible = useStore((s) => s.transform[2] >= 0.6);
  const areTypesVisible = useStore((s) => s.transform[2] >= 1);

  const { inports, outports } = useMemo(() => {
    const result = { inports: [], outports: [] };
    if (!io) {
      return result;
    }
    return {
      inports: Object.entries(io.in || {}),
      outports: Object.entries(io.out || {}),
    };
  }, [io]);

  const cn = classnames("react-flow__node-default", {
    highlighted: props.data.isHighlighted,
    dimmed: props.data.isDimmed,
  });

  return (
    <div className={cn}>
      {inports.length > 0 && (
        <div
          className={classnames("ports", "in", { hidden: !arePortsVisible })}
        >
          {inports.map(([inportName, inportType]) => (
            <Handle
              content="asd"
              type="target"
              id={inportName}
              key={inportName}
              position={Position.Top}
              isConnectable={true}
            >
              {inportName}
              {areTypesVisible &&
                inportType.typeExpr &&
                inportType.typeExpr.meta && (
                  <span className="portType">
                    {" "}
                    {(inportType.typeExpr.meta as src.Meta).text}
                  </span>
                )}
            </Handle>
          ))}
        </div>
      )}
      <div className="nodeName">{props.data.label}</div>
      {outports.length > 0 && (
        <div
          className={classnames("ports", "out", { hidden: !arePortsVisible })}
        >
          {outports.map(([outportName, outportType]) => (
            <Handle
              type="source"
              id={outportName}
              key={outportName}
              position={Position.Bottom}
              isConnectable={true}
            >
              {outportName}
              {areTypesVisible &&
                outportType.typeExpr &&
                outportType.typeExpr.meta && (
                  <span className="portType">
                    {" "}
                    {(outportType.typeExpr.meta as src.Meta).text}{" "}
                  </span>
                )}
            </Handle>
          ))}
        </div>
      )}
    </div>
  );
}
