import { Trace } from "vscode-jsonrpc";
import path from "path";
import net from "net";
import { ExtensionContext, ExtensionMode, workspace } from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  SocketTransport,
  StreamInfo,
  TransportKind,
} from "vscode-languageclient/node";

export const clientId = "nevaLSPClient";
export const clientName = "Neva LSP Client";

export function setupLsp(context: ExtensionContext): LanguageClient {
  let serverOptions: ServerOptions;
  if (context.extensionMode === ExtensionMode.Production) {
    let command = context.asAbsolutePath(path.join("out", "lsp"));
    console.log({ command });

    const transport: SocketTransport = {
      kind: TransportKind.socket,
      port: 9000,
    };

    serverOptions = { command, transport };
  } else {
    serverOptions = () => {
      let socket = net.connect({ port: 9000 });
      let result: StreamInfo = {
        writer: socket,
        reader: socket,
      };
      return Promise.resolve(result);
    };
  }

  let clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "neva" }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("**/*.*"),
    },
  };

  const client = new LanguageClient(
    clientId,
    clientName,
    serverOptions,
    clientOptions
  );
  client.setTrace(Trace.Verbose);

  return client;
}
