import { VSCodeTextField } from "@vscode/webview-ui-toolkit/react";

export function ImportsView(props: {
  imports: Array<{ alias: string; path: string }>;
  style?: object;
}) {
  return (
    <div style={props.style}>
      <h2>Use</h2>
      {props.imports.map((entry, idx, imports) => {
        const { path } = entry;
        return (
          <section
            style={{
              width: "500px",
              marginBottom: idx === imports.length - 1 ? 0 : "10px",
            }}
          >
            <VSCodeTextField readOnly style={{ width: "100%" }} value={path}>
              <span
                slot="end"
                className="codicon codicon-go-to-file import_goto_icon"
                onClick={() => console.log("navigation not implemented")}
              />
            </VSCodeTextField>
          </section>
        );
      })}
    </div>
  );
}
