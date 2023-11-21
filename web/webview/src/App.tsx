// import { VSCodeDivider } from "@vscode/webview-ui-toolkit/react";
import {
  ImportsView,
  TypesView,
  ConstantView,
  InterfacesView,
  ComponentsView,
} from "./components";
import { useIndex } from "./helpers/use_index";
import { useMemo } from "react";
import { getFileState } from "./helpers/get_file_state";

export default function App() {
  const index = useIndex();
  const fileState = useMemo(() => getFileState(index), [index]);

  return (
    <div className="app">
      {fileState.imports.length > 0 && (
        <div style={{ marginBottom: "20px" }}>
          <h2>Imports</h2>
          <ImportsView imports={fileState.imports} />
        </div>
      )}

      {fileState.entities.types.length > 0 && (
        <>
          <h2>Types</h2>
          <div style={{ marginBottom: "20px" }}>
            <TypesView types={fileState.entities.types} />
          </div>
        </>
      )}

      {fileState.entities.constants.length > 0 && (
        <div style={{ marginBottom: "20px" }}>
          <h2>Const</h2>
          <ConstantView constants={fileState.entities.constants} />
        </div>
      )}

      {fileState.entities.interfaces.length > 0 && (
        <div style={{ marginBottom: "20px" }}>
          <h2>Interfaces</h2>
          <InterfacesView interfaces={fileState.entities.interfaces} />
        </div>
      )}

      {fileState.entities.components.length > 0 && (
        <div>
          <h2>Components</h2>
          <ComponentsView components={fileState.entities.components} />
        </div>
      )}
    </div>
  );
}
