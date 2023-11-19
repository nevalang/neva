import path from "path";
import net from "net";
import { ExtensionContext, ExtensionMode, workspace } from "vscode";
import { Trace } from "vscode-jsonrpc";
import {
  LanguageClient,
  ServerOptions,
  StreamInfo,
  TransportKind,
} from "vscode-languageclient/node";

export const clientId = "nevaLSPClient";
export const clientName = "Neva LSP Client";

export function setupLsp(context: ExtensionContext): LanguageClient {
  let serverOptions: ServerOptions;

  console.info(
    "initializing lsp-client, extension mode: ",
    context.extensionMode
  );

  serverOptions =
    context.extensionMode === ExtensionMode.Production
      ? {
          command: context.asAbsolutePath(path.join("out", "lsp")),
          transport: {
            kind: TransportKind.socket,
            port: 9000,
          },
        }
      : () => {
          let socket = net.connect({ port: 9000 });
          let result: StreamInfo = {
            writer: socket,
            reader: socket,
          };
          return Promise.resolve(result);
        };

  const client = new LanguageClient(clientId, clientName, serverOptions, {
    documentSelector: [{ scheme: "file", language: "neva" }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("**/*.*"),
    },
  });

  client.setTrace(Trace.Verbose);

  return client;
}
