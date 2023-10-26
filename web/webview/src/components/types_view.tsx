import {
  VSCodeTextField,
  VSCodeDataGrid,
  VSCodeDataGridRow,
  VSCodeDataGridCell,
} from "@vscode/webview-ui-toolkit/react";
import * as ts from "../generated/typesystem";

export function TypesView(props: {
  types: Array<{ name: string; entity: ts.Def }>;
}) {
  return (
    <VSCodeDataGrid>
      {props.types.map((typeDef) => (
        <VSCodeDataGridRow>
          <VSCodeDataGridCell grid-column="1">
            <VSCodeTextField value={typeDef.name} />
          </VSCodeDataGridCell>
          <VSCodeDataGridCell grid-column="2">
            <VSCodeTextField
              value={typeDef?.entity?.bodyExpr?.inst?.ref?.name}
            />
          </VSCodeDataGridCell>
        </VSCodeDataGridRow>
      ))}
    </VSCodeDataGrid>
  );
}
