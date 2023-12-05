import { useEffect, useMemo, useState } from "react";
import { Canvas } from "./canvas/canvas";
import { getFileViewState } from "./core/file_view_state";
import { ResolveFileResponce } from "./generated/lsp_api";

const vscodeApi = acquireVsCodeApi<ResolveFileResponce>();

export default function App() {
  const persistentState = vscodeApi.getState();
  const [state, setState] = useState(persistentState);

  useEffect(() => {
    window.addEventListener(
      "message",
      (event: { data: ResolveFileResponce }) => {
        vscodeApi.setState(event.data!);
        setState(event.data!);
      }
    );
  }, []);

  const fileViewState = useMemo(() => getFileViewState(state), [state]);

  return (
    <div className="app">
      <Canvas fileViewState={fileViewState} />
    </div>
  );
}
