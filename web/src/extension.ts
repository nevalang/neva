import { ExtensionContext, window } from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { NevaEditor } from "./editor";
import { setupLsp } from "./lsp";

let lspClient: LanguageClient;
const viewType = "neva.editNeva";

export async function activate(context: ExtensionContext) {
  console.log("vscode-neva: activated");

  // Language Server
  lspClient = setupLsp(context);
  lspClient.start();

  // Custom Editor
  const editor = new NevaEditor(context, lspClient);
  context.subscriptions.push(
    window.registerCustomEditorProvider(viewType, editor, {
      supportsMultipleEditorsPerDocument: true,
    })
  );
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}
