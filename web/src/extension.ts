import { ExtensionContext, window, commands } from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { setupLsp } from "./lsp";
import { getPreviewCommand } from "./command";

let lspClient: LanguageClient;

export async function activate(context: ExtensionContext) {
  console.info("neva module detected, extension activated");

  // Run language server, initialize client and establish connection
  lspClient = setupLsp(context);
  lspClient.onNotification("neva/analyzer_message", (message: string) => {
    window.showWarningMessage(message);
  });

  // Register preview command that opens webview
  context.subscriptions.push(
    commands.registerCommand(
      "neva.openPreview",
      getPreviewCommand(context, lspClient)
    )
  );
  console.info("preview command registered");
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}
