import * as React from "react"
import * as ReactDOM from "react-dom"
import { createConfiguration, DefaultApi } from "../sdk"
import { OpenApi } from "./api"
import { App } from "./view/app"

const config = createConfiguration({})
const d = new DefaultApi(config)
const api = new OpenApi(d)

ReactDOM.render(<App api={api} />, document.getElementById("root"))
