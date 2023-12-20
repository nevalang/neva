import { BaseEdge, EdgeProps, getBezierPath } from "reactflow";

interface INormalEdgeProps {
  isHighlighted: boolean;
}

export function NormalEdge(props: EdgeProps<INormalEdgeProps>) {
  console.log("===", props.data);

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
