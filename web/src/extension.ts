import { Disposable, ExtensionContext, window } from "vscode";
import { NevaEditor } from "./editor";

const viewType = "neva.editNeva";

export async function activate(context: ExtensionContext) {
  console.log("vscode-neva activated");

  const editor = new NevaEditor(context);

  const disposable: Disposable = window.registerCustomEditorProvider(
    viewType,
    editor,
    { supportsMultipleEditorsPerDocument: true }
  );

  context.subscriptions.push(disposable);
}
