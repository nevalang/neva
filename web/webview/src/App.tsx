/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */

import { VSCodeTextField } from "@vscode/webview-ui-toolkit/react";
import { Component, Const, Interface } from "./generated/types";
import { UseFileState } from "./hooks";

export default function App() {
  const { imports, entities } = UseFileState();
  const { types, constants, interfaces, components } = entities;

  return (
    <div className="app">
      <Imports imports={imports} style={{ marginBottom: "20px" }} />

      <h2 style={{ marginBottom: "10px" }}>Types</h2>
      {types.map((entry) => {
        return JSON.stringify(entry);
      })}

      <h2 style={{ marginBottom: "10px" }}>Constants</h2>
      {constants.map((entry) => {
        const { name, entity } = entry;
        return <ConstantView name={name} entity={entity} />;
      })}

      <h2 style={{ marginBottom: "10px" }}>Interfaces</h2>
      {interfaces.map((entry) => {
        const { name, entity } = entry;
        return <InterfaceView name={name} entity={entity} />;
      })}

      <h2 style={{ marginBottom: "10px" }}>Components</h2>
      {components.map((entry) => {
        const { name, entity } = entry;
        return <ComponentView name={name} entity={entity} />;
      })}
    </div>
  );
}

function Imports(props: {
  imports: Array<{ alias: string; path: string }>;
  style?: object;
}) {
  return (
    <div {...props.style}>
      <h2 style={{ marginBottom: "10px" }}>Use</h2>
      {props.imports.map((entry, idx, imports) => {
        const { path } = entry;
        return (
          <section
            style={{
              width: "500px",
              marginBottom: idx === imports.length - 1 ? 0 : "10px",
            }}
          >
            <VSCodeTextField readOnly style={{ width: "100%" }} value={path}>
              <span
                slot="end"
                className="codicon codicon-go-to-file import_goto_icon"
                onClick={() => console.log("navigation not implemented")}
              />
            </VSCodeTextField>
          </section>
        );
      })}
    </div>
  );
}

function ConstantView(props: { name: string; entity: Const }) {
  return <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>;
}

function InterfaceView(props: { name: string; entity: Interface }) {
  return <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>;
}

function ComponentView(props: { name: string; entity: Component }) {
  return <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>;
}
