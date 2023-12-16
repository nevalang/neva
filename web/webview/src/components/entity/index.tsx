import { useContext } from "react";
import { useLocation } from "react-router-dom";
import * as src from "../../generated/sourcecode";
import { FileContext } from "../app";
import { Component } from "./component";

export function Entity() {
  const fileContext = useContext(FileContext);
  const location = useLocation();
  const entityName = location.pathname.substr(1);
  const entity = (fileContext.resp.file.entities || {})[entityName];

  if (entity.kind === src.TypeEntity) {
    const typ = fileContext.state.entities.types.find(
      (entity) => entity.name === entityName
    );
    return JSON.stringify(typ);
  }

  if (entity.kind === src.InterfaceEntity) {
    const constant = fileContext.state.entities.interfaces.find(
      (entity) => entity.name === entityName
    );
    return JSON.stringify(constant);
  }

  if (entity.kind === src.ConstEntity) {
    const constant = fileContext.state.entities.constants.find(
      (entity) => entity.name === entityName
    );
    return JSON.stringify(constant);
  }

  if (entity.kind === src.ComponentEntity) {
    const viewState = fileContext.state.entities.components.find(
      (entity) => entity.name === entityName
    );
    return <Component viewState={viewState!.entity} />;
  }

  return location.pathname;
}
