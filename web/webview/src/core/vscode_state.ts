import { useState, useEffect, createContext } from "react";
import { ResolveFileResponce } from "../generated/lsp_api";

const vscodeApi = acquireVsCodeApi<ResolveFileResponce>();

export function useResolveFile(): ResolveFileResponce | undefined {
  const persistedState = vscodeApi.getState();
  const [state, setState] = useState<ResolveFileResponce | undefined>(
    persistedState
  );

  useEffect(() => {
    const listener = (event: { data: ResolveFileResponce }) => {
      // console.log("HERE!!!", event);
      setState(event.data!);
      vscodeApi.setState(event.data!);
    };
    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, [state]);

  return state;
}

export const vscodeStateContext = createContext<
  ResolveFileResponce | undefined
>(undefined);
