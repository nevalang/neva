/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */

import { VSCodeTextField } from "@vscode/webview-ui-toolkit/react";
import { Component, Const, Interface, Port } from "./generated/types";
import { useFileState } from "./hooks";
import { useMemo } from "react";

export default function App() {
  const { imports, entities } = useFileState();
  const { types, constants, interfaces, components } = entities;

  console.log({ imports, entities });

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
        console.log({ name, entity });
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

interface PortEntries {
  inports: PortEntry[];
  outports: PortEntry[];
}

interface PortEntry {
  name: string;
  port: Port;
}

function InterfaceView(props: { name: string; entity: Interface }) {
  const portEntries: PortEntries = useMemo(() => {
    const result: PortEntries = { inports: [], outports: [] };
    if (props.entity.io === undefined) {
      return result;
    }
    const { in: inports, out: outports } = props.entity.io;
    for (const name in inports) {
      result.inports.push({
        name: name,
        port: inports[name],
      });
    }
    for (const name in outports) {
      result.outports.push({
        name: name,
        port: outports[name],
      });
    }
    return result;
  }, [props.entity]);

  const { inports, outports } = portEntries;

  return (
    <>
      <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>
      <div>
        {inports.map((inport) => (
          <>{inport.name}</>
        ))}
      </div>
      <div>
        {outports.map((outport) => (
          <>{outport.name}</>
        ))}
      </div>
    </>
  );
}

function ComponentView(props: { name: string; entity: Component }) {
  return <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>;
}
