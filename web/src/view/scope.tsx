import * as React from "react"
import { Component, IO } from "../types/program"

interface ScopeProps {
  scope: {
    [key: string]: Component
  }
}

function Scope(props: ScopeProps) {
  return <ul className="ul">{props.scope.toString()}</ul>
}

export { Scope }
