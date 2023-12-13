import {
  window,
  ExtensionContext,
  WebviewPanel,
  ViewColumn,
  Uri,
  workspace,
} from "vscode";
import { getWebviewContent } from "./webview";
import { LanguageClient } from "vscode-languageclient/node";

export function getPreviewCommand(
  context: ExtensionContext,
  client: LanguageClient
): () => Promise<void> {
  let panel: WebviewPanel | undefined;

  return async () => {
    console.info("webview triggered");

    if (panel) {
      panel.reveal();
      console.info("existing panel revealed");
      return;
    }

    if (!window.activeTextEditor) {
      window.showWarningMessage("You need to open neva file to open preview.");
      return;
    }

    console.info("existing panel not found, trying to create new one");

    // Render empty webview
    panel = window.createWebviewPanel(
      "neva",
      "Neva: Preview",
      ViewColumn.Active,
      {
        enableScripts: true,
        localResourceRoots: [
          (Uri as any).joinPath(context.extensionUri, "out"),
          (Uri as any).joinPath(context.extensionUri, "webview/dist"),
        ],
      }
    );
    panel.webview.html = getWebviewContent(panel.webview, context.extensionUri);
    panel.onDidDispose(() => (panel = undefined), null, context.subscriptions);
    // panel.iconPath = "";
    console.info("new panel has been created");

    // Request index object from LSP server
    let resp: unknown;
    try {
      resp = await client.sendRequest("resolve_file", {
        document: window.activeTextEditor.document,
        workspaceUri: workspace.workspaceFolders![0].uri,
      });
    } catch (e) {
      console.error(e);
      return;
    }
    console.info("got response from LSP server: ", resp);

    panel.webview.onDidReceiveMessage((msg) => {
      if (msg === "ready") {
        panel!.webview.postMessage(resp);
        console.info("webview ready, message sent: ", msg);
        return;
      }
      console.info("unknown message from webview: ", msg);
    });
  };
}
