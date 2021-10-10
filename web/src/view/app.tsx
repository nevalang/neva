import classNames from "classnames"
import * as React from "react"
import { EdgeData, NodeData, PortData } from "reaflow"
import * as rf from "reaflow"
import {
  ComponentTypes,
  Connection,
  IO,
  Module,
  Program,
} from "../types/program"
import { OpenApi } from "../api"

function moduleNodesAndEdges<T>(module: Module): {
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
  return net.map<EdgeData>(({ from, to }) => ({
    id: `${from.node}.${from.port}[${from.idx}]-${to.node}.${to.port}[${to.idx}]`,
    from: from.node,
    fromPort: `${from.port}[${from.idx}]`,
    to: to.node,
    toPort: `${to.port}[${to.idx}]`,
  }))
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
  api: OpenApi // FIXME
}

function App(props: AppProps) {
  const [program, setProgram] = React.useState<Program>(defaultProgram)

  React.useEffect(() => {
    async function wrap() {
      try {
        const program = await props.api.getProgram()
        setProgram(program)
      } catch (err) {
        console.error(err)
      }
    }
    wrap()
  }, [])

  const root = program.scope[program.root] as Module
  const { nodes, edges } = moduleNodesAndEdges(root)

  return (
    <div
      style={{
        position: "absolute",
        left: 0,
        right: 0,
        top: 0,
        bottom: 0,
        background: "#171010",
      }}
    >
      <rf.Canvas
        nodes={nodes}
        edges={edges}
        onNodeLink={(_, fromNode, toNode) => {}}
        edge={<rf.Edge />}
        node={node => (
          <rf.Node
            className={classNames("node", {})}
            style={{ transition: "none" }}
            dragType="port"
            port={
              <rf.Port
                onDragStart={(...a) => console.log("start", ...a)}
                onDragEnd={(...a) => console.log("end", ...a)}
                style={{
                  fill: "#5c3f9b",
                  stroke: "#000000",
                  strokeWidth: "1px",
                }}
                rx={10}
                ry={10}
              />
            }
            remove={<rf.Remove />}
          />
        )}
      />
    </div>
  )
}

export { App }
