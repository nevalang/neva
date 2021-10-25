import * as React from "react"
import * as ReactDOM from "react-dom"

import { OpenApi } from "./api/openapi"
import { App } from "./view/app"

const api = new OpenApi("http://localhost:8090")

ReactDOM.render(<App api={api} />, document.getElementById("root"))
