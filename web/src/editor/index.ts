import {
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
  }

  // every time new tab with .neva file (re)opened
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

    console.log(getWebviewContent);

    webviewPanel.webview.html = getWebviewContent(
      webviewPanel.webview,
      extensionUri
    );

    this.client.onNotification("neva/file_parsed", (parsedFile) => {
      sendMsgToWebview(webviewPanel, document, parsedFile);
    });

    // const disposables: Disposable[] = [];

    // workspace.onDidChangeTextDocument(console.log, this, disposables);

    // window.onDidChangeActiveColorTheme(
    //   (event: ColorTheme) =>
    //     webviewPanel.webview.postMessage({
    //       type: "theme",
    //       isDarkTheme: event.kind === ColorThemeKind.Dark,
    //     }),
    //   this,
    //   disposables
    // );

    // webviewPanel.webview.onDidReceiveMessage(console.log, this, disposables);

    // webviewPanel.onDidDispose(console.log);
  }
}
