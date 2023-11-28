import { ComponentViewState } from "../core/file_view_state";
import { ComponentView } from "./component_view";

interface IComponentViewProps {
  components: Array<{ name: string; entity: ComponentViewState }>;
}

export function ComponentsView(props: IComponentViewProps) {
  return (
    <>
      {props.components.map((entry) => (
        <ComponentView
          name={entry.name}
          viewState={entry.entity}
          style={{ marginBottom: "20px" }}
        />
      ))}
    </>
  );
}
