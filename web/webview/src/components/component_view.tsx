import { useMemo } from "react";
import {
  VSCodeDataGrid,
  VSCodeDataGridRow,
  VSCodeDataGridCell,
} from "@vscode/webview-ui-toolkit/react";
import { Component, EntityRef, Node } from "../generated/sourcecode";
import { InterfaceView } from "./interface_view";
import NetView from "./network_view";

export function ComponentView(props: {
  name: string;
  entity: Component;
  style?: object;
}) {
  const { name, entity } = props;

  const nodes = useMemo(() => {
    const result = [];
    for (const name in props.entity.nodes) {
      result.push({
        name: name,
        entity: props.entity.nodes[name],
      });
    }
    return result;
  }, [props.entity.nodes]);

  return (
    <div style={props.style}>
      <h3
        style={{ marginBottom: "10px", display: "flex", alignItems: "center" }}
      >
        {name}
      </h3>
      {entity.interface && <InterfaceView name="" entity={entity.interface} />}
      {nodes && <NodesView nodes={nodes} />}
      {nodes && entity.net && <NetView nodes={nodes} net={entity.net} />}
    </div>
  );
}

function NodesView(props: {
  nodes: {
    name: string;
    entity: Node;
  }[];
}) {
  return (
    <>
      <h4>Nodes</h4>
      <VSCodeDataGrid>
        {props.nodes.map(({ name, entity }) => (
          <VSCodeDataGridRow>
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
