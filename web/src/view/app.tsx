import * as React from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"

import { ComponentTypes, Module, Program } from "../types/program"
import { Api } from "../api"
import { Network } from "./network"

const defaultProgram: Program = {
  root: "",
  scope: {
    "": {
      type: ComponentTypes.MODULE,
      constants: {},
      io: {
        in: {},
        out: {},
      },
      deps: {},
      net: [],
      workers: {},
    },
  },
}

interface AppProps {
  api: Api
}

function App(props: AppProps) {
  const [state, setState] = React.useState({
    program: defaultProgram,
    activeModuleName: defaultProgram.root,
  })

  const moduleName = (nodeName: string, program: Program): string => {
    if (nodeName == "in" || nodeName == "out") {
      return state.activeModuleName
    }

    const root = program.scope[state.activeModuleName] as Module
    const dep = root.workers[nodeName]
    if (dep === undefined) {
      return state.activeModuleName
    }

    return dep
  }

  React.useEffect(() => {
    async function aux() {
      try {
        const program = await props.api.getProgram("examples/program/pkg")
        setState({ program: program, activeModuleName: program.root })
      } catch (err) {
        console.error(err)
      }
    }
    aux()
  }, [])

  return (
    <Router>
      <Switch>
        <Route path="/menu">{/* <Menu /> */}</Route>
        <Route path="/">
          <Network
            module={state.program.scope[state.activeModuleName] as Module}
            onNodeClick={(nodeName: string) => {
              setState({
                program: state.program,
                activeModuleName: moduleName(nodeName, state.program),
              })
            }}
          />
        </Route>
      </Switch>
    </Router>
  )
}

export { App }
