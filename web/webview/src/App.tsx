// import NetworkEditor from "./network_editor";
import { useEffect, useState } from "react";

const vscodeApi = acquireVsCodeApi();

export default function App() {
  //   const defaultState = vscodeApi.getState();
  const [file, setFile] = useState();

  useEffect(() => {
    const listener = (event: any) => {
      const message = event.data;

      //   if (message.type != "update") {
      //     return
      //   }

      setFile(message.file);

      vscodeApi.setState({
        content: message.document,
        uri: message.uri,
        isDarkTheme: message.isDarkTheme,
      });
    };

    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  return JSON.stringify(file);
  //   return <NetworkEditor />;
}
