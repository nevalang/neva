import {
  VSCodeDataGrid,
  VSCodeDataGridRow,
  VSCodeDataGridCell,
  VSCodeTextField,
} from "@vscode/webview-ui-toolkit/react";
import { Const, Msg } from "../generated/src";

export function ConstantView(props: {
  constants: Array<{ name: string; entity: Const }>;
}) {
  return (
    <VSCodeDataGrid>
      {props.constants.map((constant) => (
        <VSCodeDataGridRow>
          <VSCodeDataGridCell grid-column="1">
            <VSCodeTextField style={{ width: "100%" }} value={constant.name} />
          </VSCodeDataGridCell>
          <VSCodeDataGridCell grid-column="2">
            <VSCodeTextField
              style={{ width: "100%" }}
              value={formatConstValue(constant.entity.value!)}
            />
          </VSCodeDataGridCell>
        </VSCodeDataGridRow>
      ))}
    </VSCodeDataGrid>
  );
}

const formatConstValue = (msg: Msg): string => {
  switch (true) {
    case msg.bool !== undefined:
      return String(msg.bool);
    case msg.int !== undefined:
      return String(msg.int);
    case msg.float !== undefined:
      return String(msg.float);
    case msg.str !== undefined:
      return String(msg.str);
  }

  return "unknown value";
};
