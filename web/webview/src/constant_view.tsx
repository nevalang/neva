import { Const } from "./generated/src";

export function ConstantView(props: { name: string; entity: Const }) {
  return <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>;
}
