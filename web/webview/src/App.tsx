import { VSCodeDivider } from "@vscode/webview-ui-toolkit/react";
import { useFileState } from "./hooks";
import { TypesView } from "./types_view";
import { ImportsView } from "./imports_view";
import { ConstantView } from "./constant_view";
import { InterfaceView } from "./interface_view";
import { ComponentView } from "./component_view";

export default function App() {
  const { imports, entities } = useFileState();
  const { types, constants, interfaces, components } = entities;

  console.log({ imports, entities });

  return (
    <div className="app">
      <ImportsView imports={imports} style={{ marginBottom: "20px" }} />

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Types</h2>
      <div style={{ marginBottom: "20px" }}>
        <TypesView types={types} />
      </div>

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Constants</h2>
      {constants.map((entry) => {
        const { name, entity } = entry;
        return <ConstantView name={name} entity={entity} />;
      })}

      <VSCodeDivider style={{ marginBottom: "20px" }} />

      <h2>Interfaces</h2>
      {interfaces.map((entry) => {
        const { name, entity } = entry;
        console.log({ name, entity });
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
