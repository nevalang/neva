import { ResolveFileResponce } from "../generated/lsp_api";
import * as src from "../generated/sourcecode";
import * as ts from "../generated/typesystem";

// FileView is an object optimized for fast and easy rendering of the UI
interface FileView {
  imports: { alias: string; path: string }[];
  entities: FileViewEntities;
}

// we use arrays instead of objects because it's faster to render than maps
export interface FileViewEntities {
  types: { name: string; entity: ts.Def }[];
  interfaces: { name: string; entity: src.Interface }[];
  constants: { name: string; entity: src.Const }[];
  components: {
    name: string;
    entity: src.Component;
    nodesPorts: { [keyof: string]: src.Interface };
  }[];
}

export function getFileView(state: ResolveFileResponce | undefined): FileView {
  const result: FileView = {
    imports: [],
    entities: { types: [], interfaces: [], constants: [], components: [] },
  };

  if (!state) {
    return result;
  }

  for (const alias in state.file.imports) {
    result.imports.push({
      alias,
      path: state.file.imports[alias],
    });
  }

  result.imports.sort();

  // create object entries once so we can render them fast later
  for (const entityName in state.file.entities) {
    const entity = state.file.entities[entityName];

    switch (entity.kind) {
      case src.TypeEntity:
        if (entity.type === undefined) {
          continue;
        }
        result.entities.types.push({
          name: entityName,
          entity: entity.type as ts.Def,
        });
        break;
      case src.ConstEntity:
        if (entity.const === undefined) {
          break;
        }
        result.entities.constants.push({
          name: entityName,
          entity: entity.const,
        });
        break;
      case src.InterfaceEntity:
        if (entity.interface === undefined) {
          break;
        }
        result.entities.interfaces.push({
          name: entityName,
          entity: entity.interface,
        });
        break;
      case src.ComponentEntity: {
        if (entity.component === undefined) {
          break;
        }

        result.entities.components.push({
          name: entityName,
          entity: entity.component,
          nodesPorts: state.extra.nodesPorts[entityName],
        });

        break;
      }
    }
  }

  return result;
}
