import {
  VSCodeDivider,
  VSCodePanelTab,
  VSCodePanels,
} from "@vscode/webview-ui-toolkit/react";
import {
  ImportsView,
  TypesView,
  ConstantView,
  InterfaceView,
  ComponentView,
} from "./components";
import { useFileState } from "./hooks/use_file_state";

export default function App() {
  const { imports, entities } = useFileState();
  const { types, constants, interfaces, components } = entities;

  return (
    <div className="app">
      <ImportsView imports={imports} style={{ marginBottom: "20px" }} />

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Types</h2>
      <div style={{ marginBottom: "20px" }}>
        <TypesView types={types} />
      </div>

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Const</h2>
      <ConstantView constants={constants} />

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Interfaces</h2>
      {interfaces.map((entry) => {
        const { name, entity } = entry;
        return <InterfaceView name={name} entity={entity} />;
      })}

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Components</h2>
      {components.map((entry) => {
        const { name, entity } = entry;
        return <ComponentView name={name} entity={entity} />;
      })}
    </div>
  );
}
