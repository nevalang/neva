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
  window.showInformationMessage(
    "Neva module detected, start indexing workdir."
  );

  // Run language server, initialize client and establish connection
  lspClient = setupLsp(context);
  lspClient.start();
  lspClient.onNotification("neva/analyzer_message", (message: string) => {
    window.showWarningMessage(message);
  });

  // Listen to language server events and update current indexed module state
  let initialIndex: unknown;
  lspClient.onNotification("neva/workdir_indexed", (newIndex: unknown) => {
    window.showInformationMessage("Workdir indexed.");
    initialIndex = newIndex;
  });

  // Register preview command that opens webview
  context.subscriptions.push(
    commands.registerCommand(
      "neva.openPreview",
      getPreviewCommand(
        context,
        initialIndex, // note that initial index could be undefined in case langauge server hasn't indexed workdir yet
        (f: GenericNotificationHandler) => {
          lspClient.onNotification("neva/workdir_indexed", f);
        }
      )
    )
  );
}

export function deactivate(): Thenable<void> | undefined {
  return lspClient && lspClient.stop();
}

function getPreviewCommand(
  context: ExtensionContext,
  initialIndex: unknown,
  onWebviewCreated: (f: GenericNotificationHandler) => void
): () => void {
  let panel: WebviewPanel | undefined;

  return () => {
    if (!window.activeTextEditor) {
      window.showWarningMessage(
        "You need to open neva file before calling this command."
      );
      return;
    }

    if (!initialIndex) {
      window.showWarningMessage(
        "Working directory is not indexed yet. Please wait for a little bit."
      );
    }

    const column = window.activeTextEditor
      ? window.activeTextEditor.viewColumn
      : undefined;

    if (panel) {
      panel.reveal(column);
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

    // send initial index to webview
    sendMsgToWebview(panel, window.activeTextEditor!.document, initialIndex);

    // subscribe to further language server updates
    onWebviewCreated((update: unknown) => {
      sendMsgToWebview(panel!, window.activeTextEditor!.document, update);
    });

    panel.onDidDispose(
      () => {
        panel = undefined;
      },
      null,
      context.subscriptions
    );
  };
}
