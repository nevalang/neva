import { BaseEdge, EdgeProps, getBezierPath } from "reactflow";

export function NormalEdge(props: EdgeProps<{ isHighlighted: boolean }>) {
  const [edgePath] = getBezierPath({
    sourceX: props.sourceX,
    sourceY: props.sourceY,
    targetX: props.targetX,
    targetY: props.targetY,
  });

  const style = props.data?.isHighlighted
    ? { strokeOpacity: 1, stroke: "white" }
    : { strokeOpacity: 0.75 };

  return <BaseEdge {...props} path={edgePath} style={style} />;
}
