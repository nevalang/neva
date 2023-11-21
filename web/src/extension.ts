import { ExtensionContext, window, commands } from "vscode";
import {
  GenericNotificationHandler,
  LanguageClient,
} from "vscode-languageclient/node";
import { setupLsp } from "./lsp";
import { getPreviewCommand } from "./command";

let lspClient: LanguageClient;

export async function activate(context: ExtensionContext) {
  console.info("neva module detected, extension activated");

  // Run language server, initialize client and establish connection
  lspClient = setupLsp(context);
  lspClient.start().then(console.info).catch(console.error);
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
