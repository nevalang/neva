// import classNames from "classnames"
import * as React from "react"
import { IO } from "../types/program"

interface Props {
  io: IO
  deps: { [key: string]: IO }
  workers: { [key: string]: string }
}

function Menu(props: Props) {
  return (
    <ul className="ul">
      {/* <li onClick={props.onNew}>new program</li>
      <li onClick={props.onEdit}>edit program</li> */}
    </ul>
  )
}

export { Menu }
