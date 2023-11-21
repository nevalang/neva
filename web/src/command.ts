import {
  TextEditor,
  window,
  ExtensionContext,
  WebviewPanel,
  ViewColumn,
  Uri,
} from "vscode";
import { GenericNotificationHandler } from "vscode-languageclient";
import {
  getWebviewContent,
  sendIndexMsgToWebView,
  sendTabChangeMsgToWebView,
} from "./webview";

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

    onWebviewCreated((indexedModule: unknown) => {
      sendIndexMsgToWebView(
        panel!,
        window.activeTextEditor!.document,
        indexedModule
      );
      console.info("upd message sent to webview", initialIndex);
    });

    window.onDidChangeActiveTextEditor((editor: TextEditor | undefined) => {
      console.info("active text editor changed", editor);
      if (!editor || !editor.document.fileName.endsWith(".neva")) {
        return;
      }
      sendTabChangeMsgToWebView(panel!, editor.document);
      console.info("tab changed message was sent to webview");
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

    sendIndexMsgToWebView(
      panel,
      window.activeTextEditor!.document,
      initialIndex
    );
    console.info("initial message to webview", initialIndex);
  };
}
