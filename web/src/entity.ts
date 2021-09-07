interface Store {
  components: {
    [key: string]: Component
  }
}

type Component = Operator | Module

enum ComponentTypes {
  OPERATOR,
  MODULE,
}

interface Operator {
  type: ComponentTypes.OPERATOR
  io: IO
}

interface Module {   
  type: ComponentTypes.MODULE
  io: IO
  net: Net
}

interface IO {
  in: Ports
  out: Ports
}

interface Ports {
  [key: string]: PortType
}

enum PortType {
  INT,
  STR,
  BOOL,
}

export {
  Store,
  Component,
  ComponentTypes,
  Operator,
  Module,
  IO,
  Ports,
  PortType,
}

interface Net {
  nodes: Node[]
  connections: Connection[]
}

interface Node {
  name: string
  component: string
  ports: Port[]
}

interface Port {
  name: string
  direction: Direction
}

enum Direction {
  IN,
  OUT,
}

interface Connection {
  from: PortAddr
  to: PortAddr
}

interface PortAddr {
  node: string
  port: string
  idx: number
}

export { Net, Node, Port, Direction, Connection, PortAddr }
