# Filter

In classical FBP filter works differently.
Is uses sort of optional IP to specify field in the IP with bool.


# Sequencizer component

From fbp book:

> Во всех существующих системах FBP есть очень полезный компонент, который просто принимает и отдаёт все IP из своего первого входного порта, за которыми следуют все IP из его второго входного порта, и так далее, пока все пакеты в слотах не закончатся. В DFDM это называлось Sequencizer (некоторые мои друзья любят играть с английским языком). Этот компонент часто используется для принудительного создания последовательности данных, которые генерируются случайным образом из различных источников. Одним из примеров могут быть контрольные итоги, генерируемые различными процессами, которые затем вы хотите распечатать в фиксированном порядке в отчете. Вы знаете последовательность, в которой хотите, чтобы они отображались, но не знаете время, в которое они будут созданы.

# Simplify syntax

1. Allow string in network for cases with only one receiver - `node.port: node.inport`
2. Think about some special syntax? `node.port -> node.inport`

# General Purpose Router?

```yml
io:
  arg: [x,y,z]
  in:
    a[]: x


```

# SubStreams to array outports

Introduce component that allows to turn substream values into array-outport firings.

## Problem

Slots are compiled-timed

# Graphical notation for network

- Triangles for in-out
- Squares for constants/memory
- Circles for components (ops/mods)

# Context as a program primitive?

Introduce outport to root module?

# Gradual typing

If component `A` needs struct with `X` field
It can take struct with `X, Y` fields

--- Only if we need structs at all ---

# The flatter the faster

Optimize runtime structures including proto?

# Website in a rpg-manner

Tutorials like chapters in game
They require knowledge of specific topics
You go into them in a right order

# Minimize flow

Instead of require user to use `Filter` make `More` (and stuff like that) more friendly.
Instead of having `More.Out.Result` (true or false) you can have `Out.True` and `Out.False` which passes
given element next.

But how to solve task "for every user that is older than 30 years send a message `yo`"?

# Data editor

Data editor is a mind-map-like GUI
that allowes create graph
where one can leads to another

a way to visualize message interface creation

# BLACK ADAPTERS MAGIC!!!

Модуль, который динамически создаёт другие модули.
Кейс - динамическое создание адаптера между компонентом, который принимает список
и компонентом, который имеет аррай-портс интерфейс.

При старте такой модуль создаёт аррай портс с кол-ом ячеек соотв. длине списка.
При получении значения он пишет в этот порт.

# Mock autogen

Every module depends on components via interfaces.
So it should be possible to generate mock modules,
that would allow to program behaviour.

## Motivation

Test is simply a program that uses (e.g. `std/testing`) test utils

## Mock API (go:generate-ish?)

???

```
prog1.yml
prog2.yml
common/
  mod1.yml
prog1/
  mod2.yml
prog2/
  mod3.yml

respect run prog1.yml
```

# Debugger

## Undo/Redo

Instead of keeping log of all sends/receives, keep only previous values.

## Editing of messages values

## Live changing of networks

- Обмен байткодом между рантаймами
- DEBUGGER (Обёртка над компилятором и рантаймом, в рантайме, вероятно, мидлварь)
- FBP SHELL
- mocker?
- type system (типы должны быть максимально просты и совместимы с `gRPC`, `graphQL` и `json schema`)
- Close all the ports when there are no senders to receive a message from.

# REPL

...

# WEB (old)

```tsx
import * as React from "react";
import * as rf from "reaflow";
import { hasLink, NodeData, removeNode, getEdgesByNode } from "reaflow";

export function App() {
  const [selections, setSelections] = React.useState([]);

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
  ]);

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
  ]);

  const [draggingPort, setDraggingPort] = React.useState("");

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
      selections={selections}
    >
      <rf.Canvas
        // https://www.eclipse.org/elk/reference/options.html
        layoutOptions={
          {
            // "elk.nodeLabels.placement": "INSIDE V_CENTER H_RIGHT",
            // "elk.algorithm": "org.eclipse.elk.layered",
            // "elk.direction": "DOWN",
            // nodeLayering: "INTERACTIVE",
            // "org.eclipse.elk.edgeRouting": "ORTHOGONAL",
            // "elk.layered.unnecessaryBendpoints": "true",
            // "elk.layered.spacing.edgeNodeBetweenLayers": "20",
            // "org.eclipse.elk.layered.nodePlacement.bk.fixedAlignment": "BALANCED",
            // "org.eclipse.elk.layered.cycleBreaking.strategy": "DEPTH_FIRST",
            // "org.eclipse.elk.insideSelfLoops.activate": "true",
            // separateConnectedComponents: "false",
            // "spacing.componentComponent": "20",
            // spacing: "25",
            // "spacing.nodeNodeBetweenLayers": "20",
          }
        }
        nodes={nodes}
        edges={edges}
        selections={selections}
        onCanvasClick={() => setSelections([])}
        onNodeLinkCheck={(_, from, to) => !hasLink(edges, from, to)}
        onNodeLink={(_, fromNode, toNode) => {
          // TODO link ports, not nodes!
          setEdges([
            ...edges,
            {
              id: `${fromNode.id}-${toNode.id}`,
              from: fromNode.id,
              to: toNode.id,
            },
          ]);
        }}
        edge={
          <rf.Edge
            onClick={(_, edge) => setSelections([edge.id])}
            // onEnter={console.log}
            // onLeave={console.log}
            onRemove={(_, e) => {
              setEdges(edges.filter((edge) => edge.id !== e.id));
            }}
            onAdd={console.log}
          />
        }
        node={
          <rf.Node
            className="node"
            dragType="all"
            // onEnter={(_, port) => console.log(port)}
            // onLeave={(_, port) => console.log(port)}
            onClick={(_, node) => setSelections([node.id])}
            onRemove={(_event, node) => {
              const results = removeNode(nodes, edges, [node.id]);
              setNodes(results.nodes);
              setEdges(results.edges);
            }}
            port={
              <rf.Port
                onClick={(_, port) => setSelections([port.id])}
                // onEnter={(_, port) => console.log(port)}
                // onLeave={(_, port) => console.log(port)}
                onDragStart={(...a) => console.log("start", ...a)}
                onDragEnd={(...a) => console.log("end", ...a)}
                style={{ fill: "black", stroke: "white", strokeWidth: "1px" }}
                rx={10}
                ry={10}
              />
            }
          />
        }
      />
    </div>
  );
}
```
