import { useMemo } from "react";
import { ComponentsView } from "./components/components_view";
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
      {fileViewState.entities.components.length > 0 && (
        <ComponentsView components={fileViewState.entities.components} />
      )}
    </div>
  );
}
