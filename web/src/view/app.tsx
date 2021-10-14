// import classNames from "classnames"
import * as React from "react"
import { EdgeData, hasLink, NodeData, PortData } from "reaflow"
import * as rf from "reaflow"
import {
  ComponentTypes,
  Connection,
  IO,
  Module,
  Program,
} from "../types/program"
import { Api } from "../api"

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
  api: Api
}

function App(props: AppProps) {
  const [program, setProgram] = React.useState<Program>(defaultProgram)
  const [selections, setSelections] = React.useState<string[]>([])
  const [draggingPort, setDraggingPort] = React.useState("")

  React.useEffect(() => {
    async function aux() {
      try {
        const program = await props.api.getProgram("examples/program/pkg")
        setProgram(program)
      } catch (err) {
        console.error(err)
      }
    }
    aux()
  }, [])

  const root = program.scope[program.root] as Module
  const { nodes, edges } = netGraph(root)

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
        maxWidth={window.innerWidth}
        maxHeight={window.innerHeight}
        nodes={nodes}
        edges={edges}
        selections={selections}
        onNodeLinkCheck={(_, from, to) => !hasLink(edges, from, to)}
        onCanvasClick={() => setSelections([])}
        onNodeLink={(_, fromNode, toNode) => {
          // TODO link ports, not nodes!
          // setEdges([
          //   ...edges,
          //   {
          //     id: `${fromNode.id}-${toNode.id}`,
          //     from: fromNode.id,
          //     to: toNode.id,
          //   },
          // ]);
        }}
        node={
          <rf.Node
            className="node"
            dragType="port"
            // onEnter={(_, port) => console.log(port)}
            // onLeave={(_, port) => console.log(port)}
            onClick={(_, node) => setSelections([node.id])}
            onRemove={(_event, node) => {
              // const results = removeNode(nodes, edges, [node.id]);
              // setNodes(results.nodes);
              // setEdges(results.edges);
            }}
            port={
              <rf.Port
                onClick={(_, port) => setSelections([port.id])}
                // onEnter={(_, port) => console.log(port)}
                // onLeave={(_, port) => console.log(port)}
                onDragStart={(...a) => console.log("start", ...a)}
                onDragEnd={(...a) => console.log("end", ...a)}
                style={{ fill: "black", stroke: "white", strokeWidth: "1px" }}
                rx={10}
                ry={10}
              />
            }
          />
        }
        edge={edge => (
          <rf.Edge
            onClick={(_, edge) => setSelections([edge.id])}
            // onEnter={console.log}
            // onLeave={console.log}
            onRemove={(_, e) => {
              // setEdges(edges.filter(edge => edge.id !== e.id))
            }}
            onAdd={console.log}
          />
        )}
        animated={false}
      />
    </div>
  )
}

export { App }
