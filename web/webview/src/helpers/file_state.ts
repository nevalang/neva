import {
  ConstEntity,
  InterfaceEntity,
  ComponentEntity,
  Component,
  TypeEntity,
  File,
  Const,
  Interface,
} from "../generated/sourcecode";
import * as ts from "../generated/typesystem";
import { VSCodeState } from "./vscode_state";

// File state is grouped and sorted render-friendly object
interface FileState {
  imports: Array<{ alias: string; path: string }>;
  entities: GroupedEntities;
}

// we use arrays instead of objects because it's faster to render
export interface GroupedEntities {
  types: Array<{ name: string; entity: ts.Def }>;
  interfaces: Array<{ name: string; entity: Interface }>;
  constants: Array<{ name: string; entity: Const }>;
  components: Array<{ name: string; entity: Component }>;
}

export function getFileState(state: VSCodeState | undefined): FileState {
  const result: FileState = {
    imports: [],
    entities: { types: [], interfaces: [], constants: [], components: [] },
  };

  // if tab opened first time and there were no updates from vscode yet
  if (!state || !state.indexedModule || !state.indexedModule.packages) {
    return result;
  }

  const { currentFileName, currentPackageName } = getCurrentPackageAndFileName(
    state.openedDocument.fileName,
    state.workspaceUri.path
  );

  const currentFile: File =
    state.indexedModule.packages![currentPackageName][currentFileName];

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

  return result;
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
