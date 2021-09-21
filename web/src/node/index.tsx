import classNames from "classnames"
import * as React from "react"
import * as rf from "reaflow"
import { hasLink, NodeData, removeNode, EdgeData } from "reaflow"
import { Port } from "~port"

interface NodeProps {
  name: string
  ports: {
    in: { name: string }[]
    out: { name: string }[]
  }

  onClick(any)
  onPortClick(portId: string)
}

export function Node(props: NodeProps) {
  return (
    <rf.Node
      className={classNames("node", {
        activeNode: node.id == selectedIds[0],
      })}
      style={{ transition: "none" }}
      dragType="port"
      // onEnter={(_, port) => console.log(port)}
      // onLeave={(_, port) => console.log(port)}
      onClick={(_, node) => props.onClick(node.id)}
      onRemove={(_event, node) => {
        props.onRemove(node.id)
      }}
      port={
        <Port />
      }
      remove={<rf.Remove />}
    />
  )
}
