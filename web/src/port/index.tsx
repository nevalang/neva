import classNames from "classnames"
import * as React from "react"
import * as rf from "reaflow"

interface PortProps {
  id: string
  name: string
  isActive: boolean

  onClick({ id: string })
  onEnter({ id: string })
  onLeave({ id: string })
  onDragStart({ id: string })
  onDragEnd({ id: string })
}

export function Port(props: PortProps) {
  return (
    <rf.Port
      className={classNames({ activePort: props.isActive })}
      onClick={(_, port) => props.onClick(port)}
      onEnter={(_, port) => props.onEnter(port)}
      onLeave={(_, port) => props.onLeave(port)}
      onDragStart={props.onDragStart}
      onDragEnd={props.onDragEnd}
      style={{
        fill: "#5c3f9b",
        stroke: "#000000",
        strokeWidth: "1px",
      }}
      rx={10}
      ry={10}
    />
  )
}
