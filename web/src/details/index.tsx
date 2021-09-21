import classNames from "classnames"
import * as React from "react"

interface PortDetailProps {
  name: string
  type: string
}

export function PortDetails(props: PortDetailProps) {
  return (
    <div className="details">
      <span className="detailsName">{props.name}</span>
      <span className="detailsType">{props.type}</span>
    </div>
  )
}
