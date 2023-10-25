import {
  Webview,
  WebviewPanel,
  window,
  Uri,
  ColorThemeKind,
  CustomTextEditorProvider,
  CancellationToken,
  TextDocument,
  ExtensionContext,
} from "vscode";
import { LanguageClient } from "vscode-languageclient/node";

export const registerEditor = (viewType: string, editor: NevaEditor) =>
  window.registerCustomEditorProvider(viewType, editor, {
    supportsMultipleEditorsPerDocument: true,
  });

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
    token: CancellationToken
  ): void | Thenable<void> {
    const extensionUri = this.context.extensionUri;

    webviewPanel.webview.options = {
      enableScripts: true,
      localResourceRoots: [
        (Uri as any).joinPath(extensionUri, "out"),
        (Uri as any).joinPath(extensionUri, "webview/dist"),
      ],
    };

    // load react-app
    webviewPanel.webview.html = getWebviewContent(
      webviewPanel.webview,
      extensionUri
    );

    // subscribe to server's "your file parsed, now render it" events
    this.client.onNotification("neva/renderFile", (parsedFile) => {
      sendUpdWebviewMsg(webviewPanel, document, parsedFile);
    });

    // const disposables: Disposable[] = [];

    // let isUpdating = {
    //   creationTime: Date.now(),
    //   current: false,
    //   editTime: undefined,
    // };

    // workspace.onDidChangeTextDocument(
    //   (event) => {
    //     if (event.document.uri.toString() === document.uri.toString()) {
    //       if (!isUpdating.current) {
    //         console.log("update window", isUpdating);
    //         updateWebview(webviewPanel, event.document);
    //       } else {
    //         isUpdating.current = false;
    //       }
    //     }
    //   },
    //   undefined,
    //   disposables
    // );

    // window.onDidChangeActiveColorTheme(
    //   (event: ColorTheme) => {
    //     let isDarkTheme = event.kind === ColorThemeKind.Dark;
    //     webviewPanel.webview.postMessage({
    //       type: "theme",
    //       isDarkTheme: isDarkTheme,
    //     });
    //   },
    //   undefined,
    //   disposables
    // );

    // webviewPanel.webview.onDidReceiveMessage(
    //   (message: any) => {
    //     const command = message.command;
    //     const text = message.text;

    //     switch (command) {
    //       case "update":
    //         const edit = new WorkspaceEdit();
    //         isUpdating.current = true;
    //         isUpdating.editTime = Date.now() as any;
    //         console.log("update", message.uri, text);
    //         edit.replace(
    //           message.uri,
    //           new Range(0, 0, document.lineCount, 0),
    //           text
    //         );
    //         workspace.applyEdit(edit);
    //     }
    //   },
    //   undefined,
    //   disposables
    // );

    // webviewPanel.onDidDispose(() => {
    //   while (disposables.length) {
    //     const disposable = disposables.pop();
    //     if (disposable) {
    //       disposable.dispose();
    //     }
    //   }
    // });
  }
}

function getWebviewContent(webview: Webview, extensionUri: Uri) {
  const stylesUri = getUri(webview, extensionUri, [
    "webview",
    "dist",
    "assets",
    "index.css",
  ]);

  const scriptUri = getUri(webview, extensionUri, [
    "webview",
    "dist",
    "assets",
    "index.js",
  ]);

  const codiconsUri = getUri(webview, extensionUri, [
    "webview",
    "dist",
    "codicons",
    "codicon.css",
  ]);

  return /*html*/ `
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="stylesheet" type="text/css" href="${stylesUri}">
        <link href="${codiconsUri}" rel="stylesheet" />
        <title>Neva Editor</title>
      </head>
      <body>
        <div id="root"></div>
        <script type="module" nonce="${getNonce()}" src="${scriptUri}"></script>
      </body>
    </html>
  `;
}

function sendUpdWebviewMsg(
  panel: WebviewPanel,
  document: TextDocument,
  file: any
) {
  panel.webview.postMessage({
    type: "update",
    original: document,
    parsed: file,
    isDarkTheme: window.activeColorTheme.kind === ColorThemeKind.Dark,
  });
}

function getNonce() {
  let text = "";
  const possible =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  for (let i = 0; i < 32; i++) {
    text += possible.charAt(Math.floor(Math.random() * possible.length));
  }
  return text;
}

function getUri(webview: Webview, extensionUri: Uri, pathList: string[]) {
  return webview.asWebviewUri((Uri as any).joinPath(extensionUri, ...pathList));
}
