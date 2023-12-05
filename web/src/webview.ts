import { Webview, Uri } from "vscode";

export function getWebviewContent(webview: Webview, extensionUri: Uri) {
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

// deliverMessageToWebview tries to deliver a given message to the webview.
// On each try it increases the timeout to give webview more time to respond.
export async function deliverMessageToWebview<T>(
  webview: Webview,
  msg: T,
  tries: number
) {
  const inner = async (step: number): Promise<boolean> => {
    if (step > tries) {
      return false;
    }
    if (await webview.postMessage(msg)) {
      return true;
    }
    await sleep(step * 1000);
    return inner(step + 1);
  };

  await inner(1); // it's important to count from 1 and not 1 since we use step as a multiplier
}

const sleep = async (ms: number): Promise<any> =>
  new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
