interface Program {
  components: {
    [key: string]: Component
  }
}

interface Component {
  io: IO
  net?: Network
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

interface Network {
  nodes: Node[]
  connections: Connection[]
}

interface Node {
  name: string
  componentName: string
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

export {
  Program,
  Component,
  IO,
  Ports,
  PortType,
  Network,
  Node,
  Port,
  Direction,
  Connection,
  PortAddr,
}
