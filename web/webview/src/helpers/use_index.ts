import { TextDocument, Uri } from "vscode";
import { useState, useEffect } from "react";
import { Module } from "../generated/sourcecode";

export interface State {
  uri: Uri;
  document: TextDocument;
  module?: Module;
}

type VSCodeEvent = {
  data: State & {
    type: string;
  };
};

const vscodeApi = acquireVsCodeApi<State>();

export function useIndex(): State | undefined {
  const persistedState = vscodeApi.getState();
  const [state, setState] = useState<State | undefined>(persistedState);

  useEffect(() => {
    const listener = (event: VSCodeEvent) => {
      let newState: State;

      if (event.data.type === "index") {
        newState = event.data; // index event contains the whole state
      } else if (event.data.type === "tab_change") {
        // tab change event doesn't contain indexed module
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
