import {
  window,
  ExtensionContext,
  WebviewPanel,
  ViewColumn,
  Uri,
} from "vscode";
import { GenericNotificationHandler } from "vscode-languageclient";
import { getWebviewContent, sendMsgToWebview } from "./webview";

export function getPreviewCommand(
  context: ExtensionContext,
  getInitialIndex: () => unknown,
  onWebviewCreated: (f: GenericNotificationHandler) => void
): () => void {
  let panel: WebviewPanel | undefined;

  return () => {
    const initialIndex = getInitialIndex();
    console.info("webview triggered: ", { initialIndex });

    if (!window.activeTextEditor) {
      window.showWarningMessage("You need to open neva file to open preview.");
      return;
    }

    const column = window.activeTextEditor
      ? window.activeTextEditor.viewColumn
      : undefined;

    if (panel) {
      panel.reveal(column);
      console.info("existing panel revealed");
      return;
    }

    panel = window.createWebviewPanel(
      "neva",
      "Neva: Preview 👀",
      ViewColumn.Beside,
      {
        enableScripts: true,
        localResourceRoots: [
          (Uri as any).joinPath(context.extensionUri, "out"),
          (Uri as any).joinPath(context.extensionUri, "webview/dist"),
        ],
      }
    );

    panel.iconPath = (Uri as any).joinPath(
      context.extensionUri,
      "webview/dist/logo.svg"
    );
    panel.webview.html = getWebviewContent(panel.webview, context.extensionUri);

    onWebviewCreated((update: unknown) => {
      sendMsgToWebview(panel!, window.activeTextEditor!.document, update);
      console.info("upd message sent to webview", initialIndex);
    });

    panel.onDidDispose(
      () => {
        panel = undefined;
      },
      null,
      context.subscriptions
    );

    console.info("existing panel not found, new panel has been created");

    if (!initialIndex) {
      window.showWarningMessage(
        "Working directory is not indexed yet. Just wait for a little bit."
      );
      return;
    }

    sendMsgToWebview(panel, window.activeTextEditor!.document, initialIndex);
    console.info("initial message to webview", initialIndex);
  };
}
