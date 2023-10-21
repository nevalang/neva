import path from "path";
import { ExtensionContext, workspace } from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

export const clientId = "nevaLSPClient";
export const clientName = "Neva LSP Client";

export function getClient(context: ExtensionContext): LanguageClient {
  let serverModule = context.asAbsolutePath(path.join("out", "lsp"));

  console.log(serverModule);

  let serverOptions: ServerOptions = {
    run: {
      command: serverModule,
      transport: TransportKind.stdio,
    },
    debug: {
      command: serverModule,
      transport: TransportKind.stdio,
      options: {},
    },
  };

  let clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "neva" }],
  };

  return new LanguageClient(clientId, clientName, serverOptions, clientOptions);
}
