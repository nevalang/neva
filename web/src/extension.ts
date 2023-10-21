import { ExtensionContext } from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { NevaEditor, registerEditor } from "./editor";
import { setupLsp } from "./client";

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
    registerEditor(viewType, new NevaEditor(context, lspClient))
  );
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}
