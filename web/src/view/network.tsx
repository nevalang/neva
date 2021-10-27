import * as React from "react"
import { Edge, EdgeData, NodeData, Port, PortData } from "reaflow"
import * as rf from "reaflow"
import { Connection, IO, Module } from "../types/program"
import { ContextMenu, Menu, MenuItem } from "@blueprintjs/core"
import { MouseEvent, useState } from "react"

interface NetworkProps {
  module: Module
  onNodeClick(nodeName: string): void
  onAddNode(event: MouseEvent): void
  onAddConnection(connection: Connection): void
  onRemoveConnection(connection: Connection): void
  onRemoveNode(nodeName: string)
}

function NetworkEditor(props: NetworkProps) {
  const [selections, setSelections] = useState<string[]>([])

  // dragging state
  const [srcNodeId, setSrcNodeId] = useState<string | null>(null)
  const [srcPortId, setSrcPortId] = useState<string | null>(null)
  const [dstNodeId, setDstNodeId] = useState<string | null>(null)
  const [dstPortId, setDstPortId] = useState<string | null>(null)

  const renderContextMenu = (e: MouseEvent) => {
    e.preventDefault()
    ContextMenu.show(
      <Menu>
        <MenuItem text="Add node" onClick={props.onAddNode} />
      </Menu>,
      { left: e.clientX, top: e.clientY },
      () => {}
    )
  }

  const { nodes, edges } = rfGraph(props.module)

  return (
    <div
      style={{
        position: "relative",
        width: "100%",
        height: "100%",
      }}
      onContextMenu={event => renderContextMenu(event.nativeEvent as any)}
    >
      <div
        style={{
          position: "absolute",
          left: 0,
          right: 0,
          top: 0,
          bottom: 0,
          background: "black",
        }}
      >
        <rf.Canvas
          nodes={nodes}
          edges={edges}
          selections={selections}
          onCanvasClick={() => setSelections([])}
          node={nodeProps => (
            <rf.Node
              dragType="port"
              onClick={(_event, node) => setSelections([node.id])}
              onRemove={(_event, node) => props.onRemoveNode(node.id)}
              port={
                <Port
                  onEnter={(_event, port) => {
                    setDstNodeId(nodeProps.id)
                    setDstPortId(port.id)
                  }}
                  onLeave={(_event, _port) => {}}
                  onDragStart={(_event, _pos, port) => {
                    setSrcNodeId(nodeProps.id)
                    setSrcPortId(port.id)
                  }}
                  onDragEnd={(_event, _pos, _port) => {
                    props.onAddConnection({
                      from: {
                        node: srcNodeId,
                        port: removePortPrefix(srcPortId),
                      },
                      to: {
                        node: dstNodeId,
                        port: removePortPrefix(dstPortId),
                      },
                    })
                    setSrcPortId(null)
                    setDstPortId(null)
                  }}
                />
              }
              onEnter={(_event, node) => {}}
            />
          )}
          edge={edge => (
            <Edge
              onClick={(_event, edge) => setSelections([edge.id])}
              onRemove={(_event, edge) => {
                props.onRemoveConnection({
                  from: {
                    node: edge.from,
                    port: removePortPrefix(edge.fromPort),
                  },
                  to: {
                    node: edge.to,
                    port: removePortPrefix(edge.toPort),
                  },
                })
              }}
              // onAdd={}
            />
          )}
        />
      </div>
    </div>
  )
}

function removePortPrefix(withPrefix: string): string {
  return withPrefix.substring(withPrefix.indexOf("_") + 1)
}

function rfGraph<N, E>(
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
  let nodes: NodeData[] = []

  const { in: inports, out: outports } = module.io
  const inportsNode = node("in", { in: {}, out: inports })
  const outportsNode = node("out", { in: outports, out: {} })
  nodes = nodes.concat(inportsNode, outportsNode)

  if (Object.keys(module.constants).length > 0) {
    const ports = {}
    for (const name in module.constants) {
      ports[name] = module.constants[name].typ
    }
    nodes = nodes.concat(node("const", { in: {}, out: ports }))
  }

  if (Object.keys(module.workers).length > 0) {
    nodes = nodes.concat(...workerNodes(module))
  }

  return nodes
}

function workerNodes(module: Module): NodeData[] {
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
    ports: ports(name, io),
  }
}

function ports(nodeName: string, io: IO): PortData[] {
  const ports: PortData[] = []

  for (const inportName in io.in) {
    ports.push({
      id: nodeName + "_" + inportName,
      side: "NORTH",
      height: 10,
      width: 10,
    })
  }
  for (const outportName in io.out) {
    ports.push({
      id: nodeName + "_" + outportName,
      side: "SOUTH",
      height: 10,
      width: 10,
    })
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
      // text: id,
      from: from.node,
      fromPort: from.node + "_" + from.port, // TODO: array ports
      to: to.node,
      toPort: to.node + "_" + to.port,
    }
  })
}

export { NetworkEditor }
