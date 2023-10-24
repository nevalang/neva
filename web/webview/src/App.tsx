/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */

import { useEffect, useState } from "react";
import {
  VSCodeProgressRing,
  VSCodeTextField,
} from "@vscode/webview-ui-toolkit/react";
import { File, Component } from "./generated/types";

interface State {
  originalFileContent: string;
  parsedFile: File;
  uri: string;
  isDarkTheme: boolean;
}

const vscodeApi = acquireVsCodeApi<State>();

export default function App() {
  const defaultState = vscodeApi.getState();
  const [state, setState] = useState<State | undefined>(defaultState);

  useEffect(() => {
    const listener = (event: any) => {
      const msg = event.data;
      const newState = {
        originalFileContent: msg.document,
        uri: msg.uri,
        isDarkTheme: msg.isDarkTheme,
        parsedFile: msg.file,
      };
      setState(newState);
      vscodeApi.setState(newState);
    };

    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  if (state === undefined || state.parsedFile === undefined) {
    return (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          justifyContent: "center",
          alignContent: "center",
          alignItems: "center",
        }}
      >
        <VSCodeProgressRing />
      </div>
    );
  }

  const { parsedFile: file } = state;

  return (
    <div className="app">
      {file.imports && (
        <Imports imports={file.imports} style={{ marginBottom: "20px" }} />
      )}
      {file.entities &&
        Object.entries(file.entities)
          .filter((entry) => {
            const [, entity] = entry;
            return (
              entity.kind === "component_entity" &&
              entity.component !== undefined
            );
          })
          .map((entry) => {
            const [name, entity] = entry;
            return <Component name={name} entity={entity.component!} />;
          })}
    </div>
  );
}

function Component(props: { name: string; entity: Component }) {
  return <h3 style={{ marginBottom: "10px" }}>{props.name}</h3>;
}

function Imports(props: {
  imports: {
    [key: string]: string;
  };
  style?: object;
}) {
  return (
    <div {...props.style}>
      <h2 style={{ marginBottom: "10px" }}>
        {/* <span style={{ marginRight: "5px" }} className="codicon codicon-plug" /> */}
        Use
      </h2>
      {Object.entries(props.imports).map((entry, idx, imports) => {
        const [, path] = entry;
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
