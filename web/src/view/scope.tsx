import * as React from "react"
import { Component, IO } from "../types/program"
import {
  Button,
  Card,
  Dialog,
  Elevation,
  InputGroup,
  Tag,
} from "@blueprintjs/core"
import { Draggable } from "./shared"
import { useState } from "react"

interface ScopeProps {
  scope: { [key: string]: Component }
  onAdd(name: string, path: string)
  onRemove(name: string): void
  onDragEnd(name: string): void
  onClick(component: string, worker: string): void
}

function Scope(props: ScopeProps) {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [componentName, setComponentName] = useState("")
  const [workerName, setWorkerName] = useState("")

  const handleClick = () => {
    props.onClick(componentName, workerName)
    setIsDialogOpen(false)
    setWorkerName("")
    setComponentName("")
  }

  return (
    <div className="scope">
      <Dialog
        isOpen={isDialogOpen}
        onClose={() => {
          setIsDialogOpen(false)
          setWorkerName("")
          setComponentName("")
        }}
      >
        <InputGroup
          placeholder="worker name"
          value={workerName}
          onChange={e => setWorkerName(e.target.value)}
        />
        <Button text="submit" onClick={handleClick} />
      </Dialog>
      {Object.entries(props.scope).map(([name, component]) => (
        <Card
          className="scope__card"
          onClick={() => {
            setIsDialogOpen(true)
            setComponentName(name)
          }}
          interactive
          key={name}
        >
          <h3>{name}</h3>
          {/* <Tag fill minimal onRemove={() => props.onRemove(name)}>
          </Tag> */}
          <ComponentIO io={component.io} />
        </Card>
      ))}
    </div>
  )
}

interface IOPreviewProps {
  io: IO
}

function ComponentIO(props: IOPreviewProps) {
  const inports = Object.entries(props.io.in)
  const outports = Object.entries(props.io.out)

  return (
    <table className="bp3-html-table .modifier">
      <tbody>
        <tr>
          <td>IN:</td>
        </tr>
        {inports.map(([name, typ]) => (
          <tr key={name}>
            <td>{name}</td>
            <td>{typ}</td>
          </tr>
        ))}
        <tr>
          <td>OUT:</td>
        </tr>
        {outports.map(([name, typ]) => (
          <tr key={name}>
            <td>{name}</td>
            <td>{typ}</td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export { Scope }
