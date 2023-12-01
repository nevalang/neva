import { useMemo } from "react";
import {} from "@vscode/webview-ui-toolkit/react";
import { Interface } from "../generated/sourcecode";
import NetView from "./network_view";
import { ComponentViewState } from "../core/file_view_state";

export function InterfaceView(props: { name: string; entity: Interface }) {
  const virtualComponentState: ComponentViewState = useMemo(
    () => ({
      nodes: [
        {
          name: props.name,
          interface: props.entity,
          node: {},
        },
      ],
      net: [],
    }),
    [props.entity, props.name]
  );

  return (
    <NetView name={props.name} componentViewState={virtualComponentState} />
  );
}
