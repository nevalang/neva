import {
  ExtensionContext,
  window,
  commands,
  ViewColumn,
  Uri,
  WebviewPanel,
} from "vscode";
import {
  GenericNotificationHandler,
  LanguageClient,
} from "vscode-languageclient/node";
import { setupLsp } from "./lsp";
import { getWebviewContent, sendMsgToWebview } from "./webview";

let lspClient: LanguageClient;

export async function activate(context: ExtensionContext) {
  console.info("neva module detected, extension activated");

  // Run language server, initialize client and establish connection
  lspClient = setupLsp(context);
  lspClient.start();
  lspClient.onNotification("neva/analyzer_message", (message: string) => {
    window.showWarningMessage(message);
  });

  console.info("language-server started, client connection established");

  // Listen to language server events and update current indexed module state
  let initialIndex: unknown;
  lspClient.onNotification("neva/workdir_indexed", (newIndex: unknown) => {
    console.info(
      "language-server notification - workdir has been indexed",
      newIndex
    );
    initialIndex = newIndex;
  });

  // Register preview command that opens webview
  context.subscriptions.push(
    commands.registerCommand(
      "neva.openPreview",
      getPreviewCommand(
        context,
        () => initialIndex, // note that we must use closure so function call gets actual value and not undefined
        (handler: GenericNotificationHandler) => {
          lspClient.onNotification("neva/workdir_indexed", handler); // register webview update-function
        }
      )
    )
  );

  console.info("preview command registered");
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}

function getPreviewCommand(
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
      ViewColumn.Beside,
      {
        enableScripts: true,
        localResourceRoots: [
          (Uri as any).joinPath(context.extensionUri, "out"),
          (Uri as any).joinPath(context.extensionUri, "webview/dist"),
        ],
      }
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
