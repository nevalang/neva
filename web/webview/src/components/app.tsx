import { useEffect, useMemo, useState } from "react";
import { Editor } from "./editor/editor";
import { getFileViewState } from "../core/file_view_state";
import { ResolveFileResponce } from "../generated/lsp_api";

const vscodeApi = acquireVsCodeApi<ResolveFileResponce>();

export default function App() {
  const persistentState = vscodeApi.getState();
  const [state, setState] = useState(persistentState);

  useEffect(() => {
    const listener = (event: { data: ResolveFileResponce }) => {
      vscodeApi.setState(event.data!);
      setState(event.data!);
    };
    window.addEventListener("message", listener);
    vscodeApi.postMessage("ready");
    return () => window.removeEventListener("message", listener);
  }, []);

  if (state === undefined) {
    return null;
  }

  return (
    <div className="app">
      <Editor fileViewState={getFileViewState(state)} />
    </div>
  );
}
