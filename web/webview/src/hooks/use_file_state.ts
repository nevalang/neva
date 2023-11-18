import { TextDocument, Uri } from "vscode";
import { useState, useEffect, useMemo } from "react";
import {
  ConstEntity,
  InterfaceEntity,
  ComponentEntity,
  Component,
  TypeEntity,
  File,
  Const,
  Interface,
  Module,
} from "../generated/sourcecode";
import * as ts from "../generated/typesystem";

interface VSCodeMessageData {
  workspaceUri: Uri;
  openedDocument: TextDocument;
  programState: Module;
  isDarkTheme: boolean;
}

// we use arrays instead of objects because it's faster to render
interface GroupedEntities {
  types: Array<{ name: string; entity: ts.Def }>;
  interfaces: Array<{ name: string; entity: Interface }>;
  constants: Array<{ name: string; entity: Const }>;
  components: Array<{ name: string; entity: Component }>;
}

const vscodeApi = acquireVsCodeApi<VSCodeMessageData>();

// File state is grouped and sorted render-friendly object
interface FileState {
  imports: Array<{ alias: string; path: string }>;
  entities: GroupedEntities;
}

// UseFileState returns state that is easy to render. It also does memorization to avoid re-rendering
export function useFileState(): FileState {
  const persistedState = vscodeApi.getState(); // load persistent state

  // copy persistent state from disk it to memory
  const [state, setState] = useState<VSCodeMessageData | undefined>(
    persistedState
  );

  // subscribe to updates from vscode extension
  useEffect(() => {
    const listener = (event: { data: VSCodeMessageData }) => {
      setState(event.data); // update both memory
      vscodeApi.setState(event.data); // and persistent state (to reuse when tab is reopened)
    };
    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  // use state to compute render-friendly state and memoize the result
  const fileState: FileState = useMemo(() => {
    const result: FileState = {
      imports: [],
      entities: { types: [], interfaces: [], constants: [], components: [] },
    };

    // if tab opened first time and there were no updates from vscode yet
    if (state === undefined) {
      return result;
    }

    try {
      const { currentFileName, currentPackageName } =
        getCurrentPackageAndFileName(
          state.openedDocument.fileName,
          state.workspaceUri.path
        );

      const currentFile: File =
        state.programState.packages![currentPackageName][currentFileName];

      for (const alias in currentFile.imports) {
        result.imports.push({
          alias,
          path: currentFile.imports[alias],
        });
      }

      result.imports.sort();

      // object to array for faster rendering
      for (const name in currentFile.entities) {
        const entity = currentFile.entities[name];

        switch (entity.kind) {
          case TypeEntity:
            if (entity.type === undefined) {
              continue;
            }
            result.entities.types.push({
              name: name,
              entity: entity.type as ts.Def,
            });
            break;
          case ConstEntity:
            if (entity.const === undefined) {
              break;
            }
            result.entities.constants.push({
              name: name,
              entity: entity.const,
            });
            break;
          case InterfaceEntity:
            if (entity.interface === undefined) {
              break;
            }
            result.entities.interfaces.push({
              name: name,
              entity: entity.interface,
            });
            break;
          case ComponentEntity:
            if (entity.component === undefined) {
              break;
            }
            result.entities.components.push({
              name: name,
              entity: entity.component,
            });
            break;
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      return result; // eslint-disable-line no-unsafe-finally
    }
  }, [state]);

  return fileState;
}

const getCurrentPackageAndFileName = (
  openedFileName: string,
  workspacePath: string
) => {
  const relativePath = openedFileName.replace(workspacePath + "/", "");
  const pathParts = relativePath.split("/");

  const currentPackageName = pathParts.slice(0, -1).join("/"); // all but the last segment (filename)
  const currentFileNameWithExtension = pathParts[pathParts.length - 1]; // last semgent (filename)
  const currentFileName = currentFileNameWithExtension.split(".")[0]; // filename without .neva

  return { currentPackageName, currentFileName };
};
