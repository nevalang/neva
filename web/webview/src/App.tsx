import { useMemo } from "react";
import {
  ImportsView,
  TypesView,
  ConstantView,
  InterfacesView,
  ComponentsView,
} from "./components";
import { useVSCodeState, vscodeStateContext } from "./helpers/vscode_state";
import { getFileState } from "./helpers/file_state";

export default function App() {
  const index = useVSCodeState();
  const fileState = useMemo(() => getFileState(index), [index]);

  return (
    <vscodeStateContext.Provider value={index}>
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
    </vscodeStateContext.Provider>
  );
}
