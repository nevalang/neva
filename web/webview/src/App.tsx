import { VSCodeDivider } from "@vscode/webview-ui-toolkit/react";
import {
  ImportsView,
  TypesView,
  ConstantView,
  InterfaceView,
  ComponentView,
} from "./components";
import { useIndex } from "./helpers/use_index";
import { useMemo } from "react";
import { getFileState } from "./helpers/get_file_state";

export default function App() {
  const index = useIndex();
  const fileState = useMemo(() => getFileState(index), [index]);

  return (
    <div className="app">
      <ImportsView
        imports={fileState.imports}
        style={{ marginBottom: "20px" }}
      />

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Types</h2>
      <div style={{ marginBottom: "20px" }}>
        <TypesView types={fileState.entities.types} />
      </div>

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Const</h2>
      <ConstantView constants={fileState.entities.constants} />

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Interfaces</h2>
      {fileState.entities.interfaces.map((entry) => {
        const { name, entity } = entry;
        return <InterfaceView name={name} entity={entity} />;
      })}

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Components</h2>
      {fileState.entities.components.map((entry) => {
        const { name, entity } = entry;
        return <ComponentView name={name} entity={entity} />;
      })}
    </div>
  );
}
