import {
  ContextMenu,
  Drawer,
  Menu,
  MenuItem,
  Position,
} from "@blueprintjs/core"
import * as React from "react"
import { useState } from "react"
import { RouterProps } from "react-router"
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
  const [isScopeVisible, setIsScopeVisible] = useState(false)

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
  }

  const addConnection = (connection: Connection) => {
    setModule(prev => ({
      ...prev,
      net: [...prev.net, connection],
    }))
  }

  return (
    <>
      <NetworkEditor
        module={module}
        onNodeClick={(componentName: string) => {
          setModule(program.scope[componentName] as Module)
        }}
        onAddNode={() => setIsScopeVisible(true)}
        onNewConnection={addConnection}
      />
      <Drawer
        position={Position.LEFT}
        isOpen={isScopeVisible}
        onClose={() => setIsScopeVisible(false)}
        title="Scope"
        style={{ overflow: "scroll" }}
      >
        <Scope
          scope={program.scope}
          onRemove={onRemoveFromScope}
          onAdd={onAddToScope}
          onDragEnd={console.log}
          onClick={addWorker}
        />
      </Drawer>
    </>
  )
}

export { ProgramEditor }
