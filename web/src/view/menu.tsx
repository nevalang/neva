import * as React from "react"
import { Link } from "react-router-dom"
import { RouterProps } from "react-router"
import { IO } from "../types/program"

interface Props extends RouterProps {}

function Menu(props: Props) {
  return (
    <ul className="menu">
      <Link to="program">Program</Link>
    </ul>
  )
}

export { Menu }
