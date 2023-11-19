import { TextDocument, Uri } from "vscode";
import { useState, useEffect } from "react";
import { Module } from "../generated/sourcecode";

export interface VSCodeMessageData {
  workspaceUri: Uri;
  openedDocument: TextDocument;
  programState: Module;
  isDarkTheme: boolean;
}

const vscodeApi = acquireVsCodeApi<VSCodeMessageData>();

export function useIndex(): VSCodeMessageData | undefined {
  const persistedState = vscodeApi.getState();
  const [state, setState] = useState<VSCodeMessageData | undefined>(
    persistedState
  );

  useEffect(() => {
    const listener = (event: { data: VSCodeMessageData }) => {
      setState(event.data);
      vscodeApi.setState(event.data);
    };
    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  return state;
}
