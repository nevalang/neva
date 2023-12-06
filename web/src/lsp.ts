import path from "path";
import net from "net";
import cp from "child_process";
import { window, ExtensionContext, workspace } from "vscode";
import { Trace } from "vscode-jsonrpc";
import { LanguageClient, ServerOptions } from "vscode-languageclient/node";

export const clientId = "nevaLSPClient";
export const clientName = "Neva LSP Client";

export function setupLsp(context: ExtensionContext): LanguageClient {
  console.info(
    "initializing lsp-client, extension mode: ",
    context.extensionMode
  );

  const serverOptions: ServerOptions = () =>
    new Promise((resolve) => {
      const serverProcess = cp.spawn(context.asAbsolutePath(path.join("lsp")));
      serverProcess.stdout.on("data", (data) => {
        console.info(data.toString());
      });
      serverProcess.stderr.on("data", (data) => {
        console.error(data.toString());
      });
      serverProcess.on("exit", (code, signal) => {
        console.warn(`server exited with code ${code} and signal ${signal}`);
      });
      const outputChannel = window.createOutputChannel(
        "Neva Language Server Logs"
      );
      serverProcess.stdout.on("data", (data) => {
        outputChannel.appendLine(data.toString());
      });
      resolve({ reader: serverProcess.stdout, writer: serverProcess.stdin });
    });

  const client = new LanguageClient(clientId, clientName, serverOptions, {
    documentSelector: [{ scheme: "file", language: "neva" }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("**/*.*"),
    },
  });

  client.setTrace(Trace.Verbose);

  client
    .start()
    .then(() =>
      console.info("language-server started, client connection established")
    )
    .catch(console.error);

  return client;
}
