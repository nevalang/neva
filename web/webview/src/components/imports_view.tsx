import {
  VSCodeDataGrid,
  VSCodeDataGridCell,
  VSCodeDataGridRow,
} from "@vscode/webview-ui-toolkit/react";

export function ImportsView(props: {
  imports: Array<{ alias: string; path: string }>;
  style?: object;
}) {
  return (
    <VSCodeDataGrid>
      {props.imports.map((importDef) => (
        <VSCodeDataGridRow>
          <VSCodeDataGridCell grid-column="1">
            {importDef.alias}
          </VSCodeDataGridCell>
          <VSCodeDataGridCell grid-column="2">
            {importDef.path}
          </VSCodeDataGridCell>
        </VSCodeDataGridRow>
      ))}
    </VSCodeDataGrid>
  );
}
