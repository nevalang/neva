import * as React from "react"
import * as ReactDOM from "react-dom"
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom"

import { OpenApi } from "./api/openapi"
import { App } from "./view/app"

const api = new OpenApi("http://localhost:8090") // TODO use env

ReactDOM.render(<App api={api} />, document.getElementById("root"))
