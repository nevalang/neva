import * as React from "react"
import * as ReactDOM from "react-dom"

import { GrpcClient } from "./api/grpc"
import { App } from "./view/app"

ReactDOM.render(<App api={new GrpcClient()} />, document.getElementById("root"))
