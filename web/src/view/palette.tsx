import * as React from "react"
import { Component } from "../types/program"

interface PaletteProps {
  scope: {
    [key: string]: Component
  }
  onClick(component: Component): void
}

function Palette(props: PaletteProps) {
  return (
    <div className="palette">
      {Object.entries(props.scope).map(([name, component]) => (
        <div
          className="palette-card"
          onClick={() => props.onClick(component)}
          key={name}
        >
          {name}
        </div>
      ))}
    </div>
  )
}

interface PaletteCardProps {
  name: string
  component: Component
}

function PaletteCard(props: PaletteCardProps) {
  return <div className="palette-card">{props.name}</div>
}

export { Palette }
