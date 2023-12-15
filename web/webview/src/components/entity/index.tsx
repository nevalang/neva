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

  if (entity.kind === src.ComponentEntity) {
    return <Component />;
  }

  return location.pathname;
}
