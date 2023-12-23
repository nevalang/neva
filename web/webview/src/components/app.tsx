import { useEffect, useState, createContext } from "react";
import { Editor } from "./editor/editor";
import { FileViewState, getFileViewState } from "../core/file_view_state";
import { GetFileViewResponce } from "../generated/lsp_api";
import { MemoryRouter, Route, Routes, useLocation } from "react-router-dom";
import { Entity } from "../components/entity";

const vscodeApi = acquireVsCodeApi<VSCodePersistentState>();

export interface IFileContext {
  resp: GetFileViewResponce;
  state: FileViewState;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const FileContext = createContext<IFileContext>(undefined as any);

export interface VSCodePersistentState {
  resp: GetFileViewResponce;
  pathname: string;
}

export default function App() {
  const persistentState = vscodeApi.getState();
  const [resp, setResp] = useState(persistentState?.resp);

  useEffect(() => {
    const listener = (event: { data: GetFileViewResponce }) => {
      vscodeApi.setState({
        resp: event.data!,
        pathname: persistentState?.pathname || "/",
      });
      setResp(event.data!);
    };
    window.addEventListener("message", listener);
    vscodeApi.postMessage("ready");
    return () => window.removeEventListener("message", listener);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // this is only must be done once

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
        <MemoryRouter initialEntries={[persistentState?.pathname || "/"]}>
          <LocationSaver>
            <Routes>
              <Route path="/" element={<Editor />} />
              <Route path=":entityName" element={<Entity />} />
            </Routes>
          </LocationSaver>
        </MemoryRouter>
      </FileContext.Provider>
    </div>
  );
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function LocationSaver(props: any) {
  const location = useLocation();

  useEffect(() => {
    const presistentState = vscodeApi.getState();
    if (!presistentState) {
      return;
    }
    vscodeApi.setState({
      pathname: location.pathname,
      resp: presistentState.resp,
    });
  }, [location.pathname]);

  return props.children;
}
