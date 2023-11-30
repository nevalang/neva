import NetView from "./network_view";
import { ComponentViewState } from "../core/file_view_state";

export function ComponentView(props: {
  name: string;
  viewState: ComponentViewState;
  style?: object;
}) {
  return (
    <div style={props.style}>
      <h3
        style={{ marginBottom: "10px", display: "flex", alignItems: "center" }}
      >
        {props.name}
      </h3>
      {props.viewState.nodes.length > 0 &&
        props.viewState.interface &&
        props.viewState.net && (
          <NetView name={props.name} componentViewState={props.viewState} />
        )}
    </div>
  );
}
