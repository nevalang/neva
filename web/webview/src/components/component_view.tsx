import {
  VSCodeDataGrid,
  VSCodeDataGridRow,
  VSCodeDataGridCell,
  VSCodeTextField,
} from "@vscode/webview-ui-toolkit/react";
import { Component, EntityRef, Node } from "../generated/sourcecode";
import { InterfaceView } from "./interface_view";
import { useMemo } from "react";
// import NetView from "./network_view";

export function ComponentView(props: { name: string; entity: Component }) {
  const { name, entity } = props;
  return (
    <>
      <h3
        style={{ marginBottom: "10px", display: "flex", alignItems: "center" }}
      >
        {name}
      </h3>
      {entity.interface && <InterfaceView name="" entity={entity.interface} />}
      {entity.nodes && <NodesView nodes={entity.nodes} />}
      {/* {entity.nodes && entity.net && (
        <NetView nodes={entity.nodes} net={entity.net} />
      )} */}
    </>
  );
}

function NodesView(props: { nodes: { [key: string]: Node } }) {
  const nodes = useMemo(() => {
    const result = [];
    for (const name in props.nodes) {
      result.push({
        name: name,
        entity: props.nodes[name],
      });
    }
    return result;
  }, [props.nodes]);

  return (
    <>
      <h4>Nodes</h4>
      <VSCodeDataGrid>
        {nodes.map(({ name, entity }) => (
          <VSCodeDataGridRow style={{ marginBottom: "20px" }}>
            <VSCodeDataGridCell grid-column="1">{name}</VSCodeDataGridCell>
            <VSCodeDataGridCell grid-column="2">
              {formatEntityRef(entity.entityRef)}
            </VSCodeDataGridCell>
          </VSCodeDataGridRow>
        ))}
      </VSCodeDataGrid>
    </>
  );
}

function formatEntityRef(ref?: EntityRef): string {
  if (!ref) {
    return "";
  }
  if (!ref.pkg) {
    return String(ref.name);
  }
  return String(ref.pkg) + "." + String(ref.name);
}
