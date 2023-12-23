import { GetFileViewResponce } from "../generated/lsp_api";
import * as src from "../generated/sourcecode";
import * as ts from "../generated/typesystem";

// Object optimized for fast and easy rendering of the UI
export interface FileViewState {
  imports: { alias: string; path: string }[];
  entities: FileViewEntities;
}

export interface FileViewEntities {
  types: { name: string; entity: ts.Def }[];
  interfaces: { name: string; entity: src.Interface }[];
  constants: { name: string; entity: src.Const }[];
  components: { name: string; entity: ComponentViewState }[];
}

export interface ComponentViewState {
  interface?: src.Interface;
  net: src.Connection[];
  nodes: NodesViewState[];
}

export interface NodesViewState {
  name: string;
  node: src.Node;
  interface: src.Interface;
}

export function getFileViewState(
  state: GetFileViewResponce | undefined
): FileViewState {
  const result: FileViewState = {
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

        const componentViewState = {
          name: entityName,
          entity: {
            ...entity.component,
            nodes: [],
          } as ComponentViewState,
        };

        if (!entity.component.nodes) {
          result.entities.components.push(componentViewState);
          break;
        }

        const componentNodesExtra = state.extra.nodesPorts[entityName];
        const nodesViewState: NodesViewState[] = [];

        for (const nodeName in entity.component.nodes) {
          const node = entity.component.nodes;
          nodesViewState.push({
            name: nodeName,
            node: node,
            interface: componentNodesExtra[nodeName],
          });
        }

        componentViewState.entity.nodes = nodesViewState;
        result.entities.components.push(componentViewState);

        break;
      }
    }
  }

  return result;
}
