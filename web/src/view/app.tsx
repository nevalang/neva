// import classNames from "classnames"
import * as React from "react"
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom"
import { EdgeData, hasLink, NodeData, PortData } from "reaflow"

import {
  ComponentTypes,
  Connection,
  IO,
  Module,
  Program,
} from "../types/program"
import { Api } from "../api"
import { Network } from "./network"

function netGraph<T>(module: Module): {
  nodes: NodeData<T>[]
  edges: EdgeData<T>[]
} {
  return {
    nodes: moduleNodes(module),
    edges: netEdges(module.net),
  }
}

function moduleNodes(module: Module): NodeData[] {
  const { in: inports, out: outports } = module.io
  const inportsNode = node("in", { in: {}, out: inports })
  const outportsNode = node("out", { in: outports, out: {} })

  return createWorkerNodes(module).concat(inportsNode, outportsNode)
}

function createWorkerNodes(module: Module): NodeData[] {
  const nodes: NodeData[] = []
  for (const workerName in module.workers) {
    const depName = module.workers[workerName]
    const depIO = module.deps[depName]
    nodes.push(node(workerName, depIO))
  }
  return nodes
}

function node(name: string, io: IO): NodeData {
  return {
    id: name,
    text: name,
    ports: ports(io),
  }
}

function ports(io: IO): PortData[] {
  const ports: PortData[] = []

  for (const inportName in io.in) {
    ports.push({ id: inportName, side: "NORTH", height: 10, width: 10 })
  }
  for (const outportName in io.out) {
    ports.push({ id: outportName, side: "SOUTH", height: 10, width: 10 })
  }

  return ports
}

function netEdges(net: Connection[]): EdgeData[] {
  return net.map<EdgeData>(({ from, to }) => {
    let fromStr = `${from.node}.${from.port}`
    if (from.idx !== undefined) {
      fromStr += `[${from.idx}]`
    }

    let toStr = `${to.node}.${to.port}`
    if (to.idx !== undefined) {
      toStr += `[${to.idx}]`
    }

    const id = `${fromStr}-${toStr}`
    return {
      id,
      from: from.node,
      fromPort: from.port, // TODO: array ports
      to: to.node,
      toPort: to.port,
      text: id,
    }
  })
}

const defaultProgram: Program = {
  root: "",
  scope: {
    "": {
      type: ComponentTypes.MODULE,
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
  api: Api // TODO
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
        <Route path="/menu">
          {/* <Menu /> */}
        </Route>
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
