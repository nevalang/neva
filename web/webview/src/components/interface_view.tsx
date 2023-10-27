import {
  VSCodeTextField,
  VSCodeDataGrid,
  VSCodeDataGridRow,
  VSCodeDataGridCell,
} from "@vscode/webview-ui-toolkit/react";
import { Interface, Port } from "../generated/src";
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
          <VSCodeDataGridRow rowType="header">
            {inports.map((port, idx) => (
              <VSCodeDataGridCell grid-column={idx + 1}>
                <VSCodeTextField style={{ width: "100%" }} value={port.name} />
              </VSCodeDataGridCell>
            ))}
          </VSCodeDataGridRow>
          <VSCodeDataGridRow rowType="default">
            {inports.map((port, idx) => (
              <VSCodeDataGridCell grid-column={idx + 1}>
                <VSCodeTextField
                  style={{ width: "100%" }}
                  value={port.port.typeExpr?.inst?.ref?.name}
                />
              </VSCodeDataGridCell>
            ))}
          </VSCodeDataGridRow>
        </VSCodeDataGrid>
      </div>
      <div style={{ marginBottom: "20px" }}>
        <h4>Outports</h4>
        <VSCodeDataGrid generateHeader="sticky">
          <VSCodeDataGridRow rowType="header">
            {outports.map((port, idx) => (
              <VSCodeDataGridCell grid-column={idx + 1}>
                <VSCodeTextField style={{ width: "100%" }} value={port.name} />
              </VSCodeDataGridCell>
            ))}
          </VSCodeDataGridRow>
          <VSCodeDataGridRow rowType="default">
            {outports.map((port, idx) => (
              <VSCodeDataGridCell grid-column={idx + 1}>
                <VSCodeTextField
                  style={{ width: "100%" }}
                  value={port.port.typeExpr?.inst?.ref?.name}
                />
              </VSCodeDataGridCell>
            ))}
          </VSCodeDataGridRow>
        </VSCodeDataGrid>
      </div>
    </>
  );
}
