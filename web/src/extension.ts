import {
  ExtensionContext,
  window,
  commands,
  ViewColumn,
  Uri,
  WebviewPanel,
} from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { NevaEditor } from "./editor";
import { setupLsp } from "./lsp";
import { getWebviewContent, sendMsgToWebview } from "./webview";

let lspClient: LanguageClient;
const viewType = "neva.editNeva";

export async function activate(context: ExtensionContext) {
  console.log("vscode-neva: activated");

  // Run language server, initialize client and establish connection
  lspClient = setupLsp(context);
  lspClient.start();
  lspClient.onNotification("neva/analyzer_message", (message: string) => {
    window.showWarningMessage(message);
  });

  // Track the current panel with a webview
  let currentPanel: WebviewPanel | undefined = undefined;

  // Register preview command
  context.subscriptions.push(
    commands.registerCommand("neva.openPreview", () => {
      const columnToShowIn = window.activeTextEditor
        ? window.activeTextEditor.viewColumn
        : undefined;

      // If we already have a panel, show it in the target column
      if (currentPanel) {
        currentPanel.reveal(columnToShowIn);
      } else {
        // Otherwise, create a new panel
        currentPanel = window.createWebviewPanel(
          "neva",
          "Neva: Preview",
          ViewColumn.Beside,
          {
            enableScripts: true,
            localResourceRoots: [
              (Uri as any).joinPath(context.extensionUri, "out"),
              (Uri as any).joinPath(context.extensionUri, "webview/dist"),
            ],
          }
        );

        // Set content
        currentPanel.webview.html = getWebviewContent(
          currentPanel.webview,
          context.extensionUri
        );

        lspClient.onNotification("neva/workdir_indexed", (indexedModule) => {
          sendMsgToWebview(
            currentPanel!,
            window.activeTextEditor?.document!,
            indexedModule
          );
        });

        // Reset when the current panel is closed
        currentPanel.onDidDispose(
          () => {
            currentPanel = undefined;
          },
          null,
          context.subscriptions
        );
      }
    })
  );
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}
