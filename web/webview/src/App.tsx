/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-namespace */

import { useEffect, useState } from "react";
import {
  VSCodeProgressRing,
  VSCodeTextField,
} from "@vscode/webview-ui-toolkit/react";
import { File } from "./generated/types";

const vscodeApi = acquireVsCodeApi();

export default function App() {
  //   const defaultState = vscodeApi.getState();
  const [file, setFile] = useState<File>();

  useEffect(() => {
    const listener = (event: any) => {
      const message = event.data;

      setFile(message.file);

      vscodeApi.setState({
        content: message.document,
        uri: message.uri,
        isDarkTheme: message.isDarkTheme,
      });
    };

    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  if (file === undefined) {
    return (
      <div style={{ width: "100%", height: "100%" }}>
        <VSCodeProgressRing />
      </div>
    );
  }

  return (
    <div className="app">
      {file.imports && <Imports imports={file.imports} />}
    </div>
  );
}

function Imports(props: {
  imports: {
    [key: string]: string;
  };
}) {
  return (
    <>
      <h2 style={{ marginBottom: "10px" }}>
        {/* <span style={{ marginRight: "5px" }} className="codicon codicon-plug" /> */}
        Use
      </h2>
      {Object.entries(props.imports).map((entry, idx, imports) => {
        const [alias, path] = entry;
        return (
          <section
            style={{
              width: "500px",
              marginBottom: idx === imports.length - 1 ? 0 : "10px",
            }}
          >
            <VSCodeTextField style={{ width: "100%" }} value={path}>
              <span
                slot="end"
                className="codicon codicon-go-to-file import_goto_icon"
                onClick={() => console.log("navigation not implemented")}
              />
            </VSCodeTextField>
          </section>
        );
      })}
    </>
  );
}
