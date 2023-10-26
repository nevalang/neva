import { Component } from "../generated/src";

export function ComponentView(props: { name: string; entity: Component }) {
  return (
    <h3 style={{ marginBottom: "10px", display: "flex", alignItems: "center" }}>
      {" "}
      <i
        className="codicon codicon-symbol-class"
        style={{ marginRight: "5px" }}
      />
      {props.name}
    </h3>
  );
}
