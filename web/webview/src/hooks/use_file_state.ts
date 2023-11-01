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
  Program,
} from "../generated/src";
import * as ts from "../generated/typesystem";

interface VSCodeMessage {
  workspaceUri: Uri;
  openedDocument: TextDocument;
  programState: Program;
  isDarkTheme: boolean;
}

// we use arrays instead of objects because it's faster to render
interface GroupedEntities {
  types: Array<{ name: string; entity: ts.Def }>;
  interfaces: Array<{ name: string; entity: Interface }>;
  constants: Array<{ name: string; entity: Const }>;
  components: Array<{ name: string; entity: Component }>;
}

const vscodeApi = acquireVsCodeApi<VSCodeMessage>();

// File state is grouped and sorted render-friendly object
interface FileState {
  imports: Array<{ alias: string; path: string }>;
  entities: GroupedEntities;
}

// UseFileState returns state that is easy to render. It also does memorization to avoid re-rendering
export function useFileState(): FileState {
  const persistedState = vscodeApi.getState(); // load persistent state
  const [state, setState] = useState<VSCodeMessage | undefined>(persistedState); // copy it to memory

  console.log({ state });

  // subscribe to state updates from vscode
  useEffect(() => {
    const listener = (event: { data: VSCodeMessage }) => {
      setState(event.data); // update both local state
      vscodeApi.setState(event.data); // and persistent state to use when tab is reopened
    };
    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

  // given program state find file corresponding to current opened document and transform it to a render-friendly form, memoize the result
  const fileState: FileState = useMemo(() => {
    const result: FileState = {
      imports: [],
      entities: {
        types: [],
        interfaces: [],
        constants: [],
        components: [],
      },
    };

    // if tab opened first time and there were no updates from vscode yet
    if (state === undefined) {
      return result;
    }

    try {
      const workspacePath = state.workspaceUri.path;

      // TODO take current file from the program and use it down here
      const currentFile: File = {};

      if (currentFile.imports === undefined) {
        return result;
      }

      for (const alias in state.programState.imports) {
        result.imports.push({
          alias,
          path: currentFile.imports[alias],
        });
      }

      result.imports.sort();

      if (currentFile.entities === undefined) {
        return result;
      }

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
