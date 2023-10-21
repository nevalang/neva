import path from "path";
import { ExtensionContext, workspace } from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  SocketTransport,
  TransportKind,
} from "vscode-languageclient/node";

export const clientId = "nevaLSPClient";
export const clientName = "Neva LSP Client";

export function setupLsp(context: ExtensionContext): LanguageClient {
  let command = context.asAbsolutePath(path.join("out", "lsp"));

  console.log(command);

  const transport: SocketTransport = {
    kind: TransportKind.socket,
    port: 8080,
  };

  let serverOptions: ServerOptions = {
    run: {
      command,
      transport,
    },
    debug: {
      command,
      args: ["-debug"],
      transport,
    },
  };

  let clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "neva" }],
  };

  return new LanguageClient(clientId, clientName, serverOptions, clientOptions);
}
