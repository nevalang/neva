import * as React from "react"
import { Component, IO } from "../types/program"
import { Button, Card, Dialog, InputGroup } from "@blueprintjs/core"
import { useState } from "react"

interface ScopeProps {
  scope: { [key: string]: Component }
  onAdd(): void
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
      <Card className="scope__card" interactive onClick={props.onAdd}>
        <h3>New</h3>
      </Card>
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
    <table className="scope__io bp3-html-table .modifier">
      <tbody>
        <tr>
          <td>In:</td>
        </tr>
        {inports.map(([name, typ]) => (
          <tr key={name}>
            <td>{name}</td>
            <td>{typ}</td>
          </tr>
        ))}
        <tr>
          <td>Out:</td>
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
