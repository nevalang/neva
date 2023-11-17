import {
  ExtensionContext,
  window,
  commands,
  ViewColumn,
  Uri,
  WebviewPanel,
} from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { setupLsp } from "./lsp";
import { getWebviewContent, sendMsgToWebview } from "./webview";

let lspClient: LanguageClient; // module-scope var for deactivate
const viewType = "neva.editNeva";

export async function activate(context: ExtensionContext) {
  window.showInformationMessage(
    "vscode-neva: neva module detected, start indexing workdir..."
  );

  // Run language server, initialize client and establish connection
  lspClient = setupLsp(context);
  lspClient.start();
  lspClient.onNotification("neva/analyzer_message", (message: string) => {
    window.showWarningMessage(message);
  });

  // Register event listener that catches updates from language-server to send them to webview later
  let indexedModule: unknown;
  lspClient.onNotification("neva/workdir_indexed", (upd: unknown) => {
    window.showInformationMessage("vscode-neva: workdir indexed");
    indexedModule = upd;
  });

  // Track the current panel with a webview
  let currentPanel: WebviewPanel | undefined = undefined;

  // Register preview command
  context.subscriptions.push(
    commands.registerCommand(
      "neva.openPreview",
      previewCommand(currentPanel, context, indexedModule)
    )
  );
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}

function previewCommand(
  currentPanel: WebviewPanel | undefined,
  context: ExtensionContext,
  indexedModule: unknown
): () => void {
  return () => {
    const columnToShowIn = window.activeTextEditor
      ? window.activeTextEditor.viewColumn
      : undefined;

    if (currentPanel) {
      currentPanel.reveal(columnToShowIn);
    } else {
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

      currentPanel.webview.html = getWebviewContent(
        currentPanel.webview,
        context.extensionUri
      );

      if (window.activeTextEditor) {
        if (indexedModule) {
          sendMsgToWebview(
            currentPanel,
            window.activeTextEditor.document,
            indexedModule
          );
        } else {
          window.showWarningMessage(
            "vscode-neva: workdir not indexed yet, please wait a little bit and again later"
          );
        }
      } else {
        window.showWarningMessage(
          "vscode-neva: you need to open neva file before calling this command"
        );
      }

      currentPanel.onDidDispose(
        () => {
          currentPanel = undefined;
        },
        null,
        context.subscriptions
      );
    }
  };
}
