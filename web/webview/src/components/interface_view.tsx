import {
  VSCodeTextField,
  VSCodeDataGrid,
  VSCodeDataGridRow,
  VSCodeDataGridCell,
} from "@vscode/webview-ui-toolkit/react";
import { Interface, Port } from "../generated/sourcecode";
import { useMemo } from "react";

interface PortEntries {
  inports: PortEntry[];
  outports: PortEntry[];
}
interface PortEntry {
  name: string;
  port: Port;
}
export function InterfaceView(props: { name: string; entity: Interface }) {
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
      {props.name && (
        <h3
          style={{
            marginBottom: "10px",
            display: "flex",
            alignItems: "center",
          }}
        >
          {props.name}
        </h3>
      )}
      <div style={{ marginBottom: "20px" }}>
        <h4>Inports</h4>
        <VSCodeDataGrid generateHeader="sticky">
          <VSCodeDataGridRow>
            {inports.map((port) => (
              <>
                <VSCodeDataGridCell grid-column={1}>
                  {port.name}
                </VSCodeDataGridCell>
                <VSCodeDataGridCell grid-column={2}>
                  {port.port.typeExpr?.inst?.ref?.name}
                </VSCodeDataGridCell>
              </>
            ))}
          </VSCodeDataGridRow>
        </VSCodeDataGrid>
      </div>
      <div style={{ marginBottom: "20px" }}>
        <h4>Outports</h4>
        <VSCodeDataGrid generateHeader="sticky">
          <VSCodeDataGridRow>
            {outports.map((port) => (
              <>
                <VSCodeDataGridCell grid-column={1}>
                  {port.name}
                </VSCodeDataGridCell>
                <VSCodeDataGridCell grid-column={2}>
                  {port.port.typeExpr?.inst?.ref?.name}
                </VSCodeDataGridCell>
              </>
            ))}
          </VSCodeDataGridRow>
        </VSCodeDataGrid>
      </div>
    </>
  );
}
