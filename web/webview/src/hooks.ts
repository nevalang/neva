import { useState, useEffect, useMemo } from "react";
import {
  ConstEntity,
  InterfaceEntity,
  ComponentEntity,
  Component,
  TypeEntity,
  Const,
  Interface,
  File,
} from "./generated/types";
import { TextDocument } from "vscode";

interface VSCodeUpdMsg {
  original: TextDocument;
  parsed: File;
  isDarkTheme: boolean;
}

// we use arrays instead of objects because it's faster to render
interface GroupedEntities {
  types: Array<unknown>;
  interfaces: Array<{ name: string; entity: Interface }>;
  constants: Array<{ name: string; entity: Const }>;
  components: Array<{ name: string; entity: Component }>;
}

const vscodeApi = acquireVsCodeApi<VSCodeUpdMsg>();

// File state is grouped and sorted render-friendly object
interface FileState {
  imports: Array<{ alias: string; path: string }>;
  entities: GroupedEntities;
}

// UseFileState returns state that is easy to render. It also does memorization to avoid re-rendering
export function useFileState(): FileState {
  const persistedState = vscodeApi.getState();
  const [state, setState] = useState<VSCodeUpdMsg | undefined>(persistedState);

  useEffect(() => {
    const listener = (event: { data: VSCodeUpdMsg }) => {
      setState(event.data);
      vscodeApi.setState(event.data);
    };
    window.addEventListener("message", listener);
    return () => window.removeEventListener("message", listener);
  }, []);

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

    if (state === undefined) {
      return result;
    }

    try {
      if (state.parsed.imports === undefined) {
        return result;
      }

      for (const alias in state.parsed.imports) {
        result.imports.push({
          alias,
          path: state.parsed.imports[alias],
        });
      }

      result.imports.sort();

      if (state.parsed.entities === undefined) {
        return result;
      }

      for (const name in state.parsed.entities) {
        const entity = state.parsed.entities[name];

        switch (entity.kind) {
          case TypeEntity:
            if (entity.type === undefined) {
              continue;
            }
            result.entities.types.push({ name: name, entity: entity.type });
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
