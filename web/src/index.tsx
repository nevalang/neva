// import "regenerator-runtime/runtime"
import * as React from "react"
import * as ReactDOM from "react-dom"

import { App } from "./view/app"
import { Program } from "~types/program"
import { OpenApi } from "./api"

export interface Api {
  getProgram(): Promise<Program>
}

const api = new OpenApi()

ReactDOM.render(<App api={api} />, document.getElementById("root"))
