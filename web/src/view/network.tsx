import * as React from "react"
import { EdgeData, hasLink, NodeData, Port, PortData } from "reaflow"
import * as rf from "reaflow"
import { Connection, IO, Module } from "../types/program"
import { ContextMenu, Menu, MenuItem } from "@blueprintjs/core"
import { MouseEvent, useState } from "react"

interface NetworkProps {
  module: Module
  onNodeClick(componentName: string): void
  onAddNewNode(event: MouseEvent): void
}

const componentByNode = (nodeName: string, module: Module): string => {
  return module.workers[nodeName]
}

function NetworkEditor(props: NetworkProps) {
  const [selections, setSelections] = useState<string[]>([])

  // dragging state
  const [srcNodeId, setSrcNodeId] = useState(null)
  const [srcPortId, setSrcPortId] = useState(null)
  const [targetNodeId, setTargetNodeId] = useState(null)
  const [targetPortId, setTargetPortId] = useState(null)

  const renderContextMenu = (e: MouseEvent) => {
    e.preventDefault()
    ContextMenu.show(
      <Menu>
        <MenuItem text="Add node" onClick={props.onAddNewNode} />
      </Menu>,
      { left: e.clientX, top: e.clientY },
      () => {}
    )
  }

  const { nodes: n, edges: e } = netGraph(props.module)
  const [nodes, setNodes] = useState(n)
  const [edges, setEdges] = useState(e)

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
              onClick={(_, nodeData) =>
                props.onNodeClick(componentByNode(nodeData.id, props.module))
              }
              port={
                <Port
                  onEnter={(_event, port) => {
                    setTargetPortId(port.id)
                  }}
                  onLeave={(_event, _port) => {
                    setTargetPortId(null)
                  }}
                  onDragStart={(_event, _pos, port) => {
                    setSrcNodeId(nodeProps.id)
                    setSrcPortId(port.id)
                  }}
                  onDragEnd={(_event, _pos, _port) => {
                    setSrcPortId(null)
                    setTargetPortId(null)
                    setEdges([
                      ...edges,
                      {
                        id: `${srcNodeId}.${srcPortId}-${targetNodeId}.${targetPortId}`,
                        from: srcNodeId,
                        fromPort: srcPortId,
                        to: targetNodeId,
                        toPort: targetPortId,
                      },
                    ])
                  }}
                />
              }
              onEnter={(_event, node) => {
                setTargetNodeId(node.id)
              }}
            />
          )}
        />
      </div>
    </div>
  )
}

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
