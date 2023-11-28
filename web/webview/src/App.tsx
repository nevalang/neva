import { useMemo } from "react";
import {
  ImportsView,
  TypesView,
  ConstantView,
  InterfacesView,
  ComponentsView,
} from "./components";
import { useResolveFile } from "./core/vscode_state";
import { getFileViewState } from "./core/file_view_state";

export default function App() {
  const resolveFileResp = useResolveFile();
  const fileViewState = useMemo(
    () => getFileViewState(resolveFileResp),
    [resolveFileResp]
  );

  return (
    <div className="app">
      {fileViewState.imports.length > 0 && (
        <div style={{ marginBottom: "20px" }}>
          <h2>Imports</h2>
          <ImportsView imports={fileViewState.imports} />
        </div>
      )}

      {fileViewState.entities.types.length > 0 && (
        <>
          <h2>Types</h2>
          <div style={{ marginBottom: "20px" }}>
            <TypesView types={fileViewState.entities.types} />
          </div>
        </>
      )}

      {fileViewState.entities.constants.length > 0 && (
        <div style={{ marginBottom: "20px" }}>
          <h2>Const</h2>
          <ConstantView constants={fileViewState.entities.constants} />
        </div>
      )}

      {fileViewState.entities.interfaces.length > 0 && (
        <div style={{ marginBottom: "20px" }}>
          <h2>Interfaces</h2>
          <InterfacesView interfaces={fileViewState.entities.interfaces} />
        </div>
      )}

      {fileViewState.entities.components.length > 0 && (
        <div>
          <h2>Components</h2>
          <ComponentsView components={fileViewState.entities.components} />
        </div>
      )}
    </div>
  );
}
