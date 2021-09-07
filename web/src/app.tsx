import * as React from "react"
import * as rf from "reaflow"

export function App() {
  const [selections, setSelections] = React.useState([])

  const [nodes, setNodes] = React.useState([
    {
      id: "in",
      text: "in",
      ports: [
        {
          id: "x",
          height: 10,
          width: 10,
          side: "SOUTH",
        },
      ],
    },
    {
      id: "multi",
      text: "multi",
      ports: [
        {
          id: "nums[0]",
          height: 10,
          width: 10,
          side: "NORTH",
        },
        {
          id: "nums[1]",
          height: 10,
          width: 10,
          side: "NORTH",
        },
        {
          id: "mul",
          height: 10,
          width: 10,
          side: "SOUTH",
        },
      ],
    },
    {
      id: "out",
      text: "out",
      ports: [
        {
          id: "y",
          height: 10,
          width: 10,
          side: "NORTH",
        },
      ],
    },
  ])

  const [edges, setEdges] = React.useState([
    {
      id: "in.x-multi.nums[0]",
      from: "in",
      to: "multi",
      fromPort: "x",
      toPort: "nums[0]",
    },
    {
      id: "in.x-multi.nums[1]",
      from: "in",
      to: "multi",
      fromPort: "x",
      toPort: "nums[1]",
    },
    {
      id: "multi.mul-out.y",
      from: "multi",
      to: "out",
      fromPort: "mul",
      toPort: "y",
    },
  ])

  return (
    <div
      style={{
        position: "absolute",
        left: 0,
        right: 0,
        top: 0,
        bottom: 0,
        background: "#171010",
      }}
    >
      <rf.Canvas
        nodes={nodes}
        edges={edges}
        selections={selections}
        onCanvasClick={() => setSelections([])}
        onNodeLinkCheck={() => true}
        onNodeLink={(event, from, to) => {
          console.log(event, from, to)
        }}
        edge={
          <rf.Edge
            onClick={(_, edge) => setSelections([edge.id])}
            onRemove={(_, e) =>
              setEdges(edges.filter(edge => edge.id !== e.id))
            }
          />
        }
        node={
          <rf.Node
            className="node"
            drugtype="all"
            onClick={(_, node) => setSelections([node.id])}
            onRemove={(_, node) => {
              console.log("onRemove node: ", node)
            }}
            port={
              <rf.Port
                onClick={(_, port) => {}}
                onEnter={(_, port) => {}}
                onLeave={(_, port) => {}}
                style={{ fill: "black", stroke: "white", strokeWidth: "1px" }}
                rx={10}
                ry={10}
              />
            }
          />
        }
      />
    </div>
  )
}
