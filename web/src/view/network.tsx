import * as React from "react"
import { EdgeData, hasLink, NodeData, PortData } from "reaflow"
import * as rf from "reaflow"
import { Connection, IO, Module } from "../types/program"

function netGraph<N, E>(
  module: Module
): {
  nodes: NodeData<N>[]
  edges: EdgeData<E>[]
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

  const constNodeOut = {}
  for (const name in module.constants) {
    constNodeOut[name] = module.constants[name].typ
  }
  const constNode = node("const", { in: {}, out: constNodeOut })

  return createWorkerNodes(module).concat(inportsNode, outportsNode, constNode)
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
      text: id,
      from: from.node,
      fromPort: from.port, // TODO: array ports
      to: to.node,
      toPort: to.port,
    }
  })
}

interface NetworkProps {
  module: Module
  onNodeClick(string): void
}

function Network(props: NetworkProps) {
  const [selections, setSelections] = React.useState<string[]>([])
  const { nodes, edges } = netGraph(props.module)

  console.log({ nodes, edges })

  console.log(
    edges.every(edge => {
      const f = nodes
        .find(node => node.id == edge.from)
        .ports.some(port => port.id == edge.fromPort)

      const t = nodes
        .find(node => node.id == edge.to)
        .ports.some(port => port.id == edge.toPort)

      console.log({ edgeID: edge.id, f, t })

      return f && t
    })
  )

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
        edges={[]}
        selections={selections}
        onNodeLinkCheck={(_, from, to) => !hasLink(edges, from, to)}
        onCanvasClick={() => setSelections([])}
        node={
          <rf.Node
            className="node"
            dragType="port"
            onClick={(_, node) => props.onNodeClick(node.id)}
          />
        }
      />
    </div>
  )
}

export { Network }
