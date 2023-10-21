import { ExtensionContext, commands } from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { NevaEditor, registerEditor } from "./editor";
import { clientId, getClient } from "./client";

let client: LanguageClient;
const viewType = "neva.editNeva";

export async function activate(context: ExtensionContext) {
  console.log("vscode-neva: activated");

  // Custom Editor
  const editor = new NevaEditor(context);
  context.subscriptions.push(registerEditor(viewType, new NevaEditor(context)));

  // Language Server
  client = getClient(context);
  await client.start();
}

export function deactivate(): Thenable<void> | undefined {
  return client && client.stop();
}
