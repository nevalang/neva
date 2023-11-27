import { useState, useEffect, createContext } from "react";
import { TextDocument, Uri } from "vscode";
import { Module } from "../generated/sourcecode";

export interface VSCodeState {
  workspaceUri: Uri;
  openedDocument: TextDocument;
  indexedModule?: Module;
}

type VSCodeEvent = {
  data: VSCodeState & {
    type: string;
  };
};

const vscodeApi = acquireVsCodeApi<VSCodeState>();

export function useVSCodeState(): VSCodeState | undefined {
  const persistedState = vscodeApi.getState();
  const [state, setState] = useState<VSCodeState | undefined>(persistedState);

  useEffect(() => {
    const listener = (event: VSCodeEvent) => {
      let newState: VSCodeState;

      if (event.data.type === "index") {
        newState = event.data;
      } else if (event.data.type === "tab_change") {
        newState = {
          ...(state || {}),
          ...event.data,
        };
      }

      setState(newState!);
      vscodeApi.setState(newState!);
    };

    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, [state]);

  return state;
}

export const vscodeStateContext = createContext<VSCodeState | undefined>(
  undefined
);
