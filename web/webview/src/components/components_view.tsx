import { Component } from "../generated/sourcecode";
import { ComponentView } from "./component_view";

interface IComponentViewProps {
  components: Array<{ name: string; entity: Component }>;
}

export function ComponentsView(props: IComponentViewProps) {
  return (
    <>
      {props.components.map((entry) => {
        const { name, entity } = entry;
        return <ComponentView name={name} entity={entity} />;
      })}
    </>
  );
}
