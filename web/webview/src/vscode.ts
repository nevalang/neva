import type { WebviewApi } from "vscode-webview";

export class VSCodeAPIWrapper {
  private readonly api: WebviewApi<unknown>;

  constructor(api: WebviewApi<unknown>) {
    this.api = api;
  }

  public postMessage(message: unknown) {
    this.api.postMessage(message);
  }

  public getState(): unknown {
    return this.api.getState();
  }

  public setState<T>(newState: T): T {
    return this.api.setState(newState);
  }
}
