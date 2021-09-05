import * as React from "react";
import * as ReactDOM from "react-dom";
import { Canvas, Port, Node } from "reaflow";

const nodes = [
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
];

const edges = [
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
];

function App() {
  return (
    <Canvas
      nodes={nodes}
      edges={edges}
      node={
        <Node
          port={
            <Port
              onClick={(_, port) => {
                console.log("onClick port: ", port);
              }}
              onEnter={(_, port) => {
                console.log("onEnter port: ", port);
              }}
              onLeave={(_, port) => {
                console.log("onLeave port: ", port);
              }}
              style={{ fill: "black", stroke: "white", strokeWidth: "1px" }}
              rx={10}
              ry={10}
            />
          }
        />
      }
    />
  );
}

ReactDOM.render(<App />, document.getElementById("root"));
