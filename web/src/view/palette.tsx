import * as React from "react"
import { Component } from "../types/program"

interface PaletteProps {
  scope: {
    [key: string]: Component
  }
}

function Palette(props: PaletteProps) {
  return <ul className="ul"></ul>
}

export { Palette }
