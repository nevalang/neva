import { Component } from "../generated/sourcecode";
import { ComponentView } from "./component_view";

interface IComponentViewProps {
  components: Array<{ name: string; entity: Component }>;
}

export function ComponentsView(props: IComponentViewProps) {
  return (
    <>
      {props.components.map((entry) => (
        <ComponentView
          name={entry.name}
          entity={entry.entity}
          style={{ marginBottom: "20px" }}
        />
      ))}
    </>
  );
}
