import * as React from "react"
import {
  Classes,
  Icon,
  Intent,
  TreeNodeInfo,
  Tree as BPTree,
} from "@blueprintjs/core"

interface TreeProps {}

function Tree(props: TreeProps) {
  const contents: TreeNodeInfo[] = [
    // {
    // },
  ]
  return <BPTree contents={[]} />
}

export { Tree }
