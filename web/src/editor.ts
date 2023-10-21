import {
  Disposable,
  Webview,
  WebviewPanel,
  window,
  Uri,
  ColorThemeKind,
  ColorTheme,
  CustomTextEditorProvider,
  CancellationToken,
  TextDocument,
  ExtensionContext,
  WorkspaceEdit,
  Range,
  workspace,
} from "vscode";

export const registerEditor = (viewType: string, editor: NevaEditor) =>
  window.registerCustomEditorProvider(viewType, editor, {
    supportsMultipleEditorsPerDocument: true,
  });

export class NevaEditor implements CustomTextEditorProvider {
  private readonly context: ExtensionContext;

  constructor(context: ExtensionContext) {
    this.context = context;
  }

  resolveCustomTextEditor(
    document: TextDocument,
    webviewPanel: WebviewPanel,
    token: CancellationToken
  ): void | Thenable<void> {
    // const file = File.fromJSON(protoString);

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

    const disposables: Disposable[] = [];

    let isUpdating = {
      creationTime: Date.now(),
      current: false,
      editTime: undefined,
    };

    workspace.onDidChangeTextDocument(
      (event) => {
        if (event.document.uri.toString() === document.uri.toString()) {
          if (!isUpdating.current) {
            console.log("update window", isUpdating);
            updateWindow(webviewPanel, event.document);
          } else {
            isUpdating.current = false;
          }
        }
      },
      undefined,
      disposables
    );

    window.onDidChangeActiveColorTheme(
      (event: ColorTheme) => {
        let isDarkTheme = event.kind === ColorThemeKind.Dark;
        webviewPanel.webview.postMessage({
          type: "theme",
          isDarkTheme: isDarkTheme,
        });
      },
      undefined,
      disposables
    );

    webviewPanel.webview.onDidReceiveMessage(
      (message: any) => {
        const command = message.command;
        const text = message.text;

        switch (command) {
          case "update":
            const edit = new WorkspaceEdit();
            isUpdating.current = true;
            isUpdating.editTime = Date.now() as any;
            console.log("update", message.uri, text);
            edit.replace(
              message.uri,
              new Range(0, 0, document.lineCount, 0),
              text
            );
            workspace.applyEdit(edit);
        }
      },
      undefined,
      disposables
    );

    webviewPanel.onDidDispose(() => {
      while (disposables.length) {
        const disposable = disposables.pop();
        if (disposable) {
          disposable.dispose();
        }
      }
    });

    updateWindow(webviewPanel, document);
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

  return /*html*/ `
    <!DOCTYPE html>
    <html lang="en">
      <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="stylesheet" type="text/css" href="${stylesUri}">
        <title>Neva Editor</title>
      </head>
      <body>
        <div id="root"></div>
        <script type="module" nonce="${getNonce()}" src="${scriptUri}"></script>
      </body>
    </html>
  `;
}

function updateWindow(panel: WebviewPanel, document: TextDocument) {
  const currentTheme = window.activeColorTheme;
  const isDarkTheme = currentTheme.kind === ColorThemeKind.Dark;

  panel.webview.postMessage({
    type: "revive",
    value: document.getText(),
    uri: document.uri,
    isDarkTheme: isDarkTheme,
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
