import {
  TextEditor,
  window,
  ExtensionContext,
  WebviewPanel,
  ViewColumn,
  Uri,
} from "vscode";
import { getWebviewContent, sendMessageToWebView } from "./webview";
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

    console.info("new panel has been created");

    let resp: unknown;
    try {
      resp = await client.sendRequest("foobar", {
        document: window.activeTextEditor.document,
      });
      sendMessageToWebView(panel, resp);
    } catch (e) {
      console.error(e);
      return;
    }

    console.info("message sent to webview", resp);
  };
}
