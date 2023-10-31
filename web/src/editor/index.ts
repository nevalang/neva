import {
  window,
  WebviewPanel,
  Uri,
  CustomTextEditorProvider,
  CancellationToken,
  TextDocument,
  ExtensionContext,
} from "vscode";
import { LanguageClient } from "vscode-languageclient/node";
import { getWebviewContent, sendMsgToWebview } from "./helpers";

export class NevaEditor implements CustomTextEditorProvider {
  private readonly context: ExtensionContext;
  private readonly client: LanguageClient;

  constructor(context: ExtensionContext, client: LanguageClient) {
    this.context = context;
    this.client = client;

    this.client.onNotification("neva/workdir_indexed", (parsedProgram) => {
      console.log({ parsedProgram });
    });

    this.client.onNotification("neva/analyzer_message", (message: string) => {
      window.showWarningMessage(message);
    });
  }

  resolveCustomTextEditor(
    document: TextDocument,
    webviewPanel: WebviewPanel,
    _: CancellationToken
  ): void | Thenable<void> {
    const extensionUri = this.context.extensionUri;

    webviewPanel.webview.options = {
      enableScripts: true,
      localResourceRoots: [
        (Uri as any).joinPath(extensionUri, "out"),
        (Uri as any).joinPath(extensionUri, "webview/dist"),
      ],
    };

    webviewPanel.webview.html = getWebviewContent(
      webviewPanel.webview,
      extensionUri
    );

    this.client.onNotification("neva/workdir_indexed", (parsedProgram) => {
      console.log({ parsedProgram });
      sendMsgToWebview(webviewPanel, document, parsedProgram);
    });

    this.client.onNotification("neva/analyzer_message", (message: string) => {
      console.log({ message });
      window.showWarningMessage(message);
    });
  }
}
