import { useEffect, useState, createContext } from "react";
import { Editor } from "./editor/editor";
import { FileViewState, getFileViewState } from "../core/file_view_state";
import { ResolveFileResponce } from "../generated/lsp_api";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { Entity } from "../components/entity";

const vscodeApi = acquireVsCodeApi<ResolveFileResponce>();

export interface IFileContext {
  resp: ResolveFileResponce;
  state: FileViewState;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const FileContext = createContext<IFileContext>(undefined as any);

export default function App() {
  const persistentState = vscodeApi.getState();
  const [resp, setResp] = useState(persistentState);

  useEffect(() => {
    const listener = (event: { data: ResolveFileResponce }) => {
      vscodeApi.setState(event.data!);
      setResp(event.data!);
    };
    window.addEventListener("message", listener);
    vscodeApi.postMessage("ready");
    return () => window.removeEventListener("message", listener);
  }, []);

  if (resp === undefined) {
    return null;
  }

  const contextValue = {
    resp: resp,
    state: getFileViewState(resp),
  };

  return (
    <div className="app">
      <FileContext.Provider value={contextValue}>
        <MemoryRouter initialEntries={["/"]}>
          <Routes>
            <Route path="/" element={<Editor />} />
            <Route path=":entityName" element={<Entity />} />
          </Routes>
        </MemoryRouter>
      </FileContext.Provider>
    </div>
  );
}
