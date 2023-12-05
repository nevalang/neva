import { useMemo } from "react";
import { Canvas } from "./canvas/canvas";
import { useResolveFile } from "./core/vscode_state";
import { getFileViewState } from "./core/file_view_state";

export default function App() {
  const resolveFileResp = useResolveFile();

  const fileViewState = useMemo(
    () => getFileViewState(resolveFileResp),
    [resolveFileResp]
  );

  return (
    <div className="app">
      <Canvas fileViewState={fileViewState} />
    </div>
  );
}
