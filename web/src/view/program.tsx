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
import { Module, Program } from "../types/program"
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
  const [drawerOpen, setDrawerOpen] = useState(false)

  const addNewNode = (componentName: string, workerName: string) => {
    setModule(prev => ({
      ...prev,
      workers: {
        ...prev.workers,
        [workerName]: componentName,
      },
    }))
  }

  return (
    <>
      <NetworkEditor
        module={module}
        onNodeClick={(componentName: string) => {
          setModule(program.scope[componentName] as Module)
        }}
        onAddNewNode={() => setDrawerOpen(true)}
      />
      <Drawer
        position={Position.LEFT}
        isOpen={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        title="Scope"
        style={{ overflow: "scroll" }}
      >
        <Scope
          scope={program.scope}
          onRemove={onRemoveFromScope}
          onAdd={onAddToScope}
          onDragEnd={console.log}
          onClick={addNewNode}
        />
      </Drawer>
    </>
  )
}

export { ProgramEditor }
