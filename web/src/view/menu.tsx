import * as React from "react"
import { Link } from "react-router-dom"
import { RouterProps } from "react-router"

interface Props extends RouterProps {
  onOpen(): void
}

function Menu(props: Props) {
  return (
    <ul className="menu">
      <li>
        <Link to="program">New</Link>
      </li>
      <li>
        <Link to="program" onClick={props.onOpen}>
          Open
        </Link>
      </li>
    </ul>
  )
}

export { Menu }
