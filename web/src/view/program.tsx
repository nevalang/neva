import * as React from "react"
import { useState } from "react"
import { Drawer, Position } from "@blueprintjs/core"
import { RouterProps } from "react-router"
import omit from "lodash.omit"

import { Connection, Module, Program } from "../types/program"
import { NetworkEditor } from "./network"
import { Scope } from "./scope"

interface ProgramEditorProps extends RouterProps {
  program: Program
  onRemoveFromScope(name: string)
  onAddToScope()
}

function ProgramEditor({
  program,
  history,
  onRemoveFromScope,
  onAddToScope,
}: ProgramEditorProps) {
  const [module, setModule] = useState(program.scope[program.root] as Module)
  const [isOpen, setIsOpen] = useState(false)

  const addWorker = (componentName: string, workerName: string) => {
    if (module.workers[workerName]) {
      console.log(workerName, "exist already")
      return
    }
    setModule(prev => ({
      ...prev,
      deps: { ...prev.deps, [componentName]: program.scope[componentName].io },
      workers: { ...prev.workers, [workerName]: componentName },
    }))
    setIsOpen(false)
  }

  const addConnection = (connection: Connection) => {
    if (
      !module.net.find(
        c =>
          c.from.node == connection.from.node &&
          c.from.port == connection.from.port &&
          c.from.idx == connection.from.idx &&
          c.to.node == connection.to.node &&
          c.to.port == connection.to.port &&
          c.to.idx == connection.to.idx
      )
    ) {
      setModule(prev => ({
        ...prev,
        net: [...prev.net, connection],
      }))
    }
  }

  const removeConnection = (toRemove: Connection) => {
    setModule(prev => ({
      ...prev,
      net: prev.net.filter(
        c =>
          c.from.node !== toRemove.from.node &&
          c.from.port !== toRemove.from.port &&
          c.to.node !== toRemove.to.node &&
          c.to.port !== toRemove.to.port
      ),
    }))
  }

  const removeNode = (nodeName: string) => {
    if (["in", "out", "const"].includes(nodeName)) {
      console.log(nodeName, "cannot be removed")
      return
    }
    setModule(prev => ({
      ...prev,
      workers: omit(prev.workers, nodeName),
      net: prev.net.filter(
        c => c.from.node !== nodeName && c.to.node !== nodeName
      ),
    }))
  }

  return (
    <>
      <Drawer
        position={Position.LEFT}
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Scope"
        style={{ overflow: "scroll" }}
      >
        <Scope
          scope={program.scope}
          onRemove={onRemoveFromScope}
          onNew={console.log}
          onDragEnd={console.log}
          onSelect={addWorker}
        />
      </Drawer>
      <NetworkEditor
        module={module}
        onNodeClick={(componentName: string) =>
          setModule(program.scope[componentName] as Module)
        }
        onAddNode={() => setIsOpen(true)}
        onAddConnection={addConnection}
        onRemoveConnection={removeConnection}
        onRemoveNode={removeNode}
      />
    </>
  )
}

export { ProgramEditor }
