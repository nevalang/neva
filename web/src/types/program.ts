export interface Program {
  // descriptor?: ProgramDescriptor
  scope: { [key: string]: Component }
  root: string
}

export interface ProgramDescriptor {
  path: string
  content: string
}

export enum ComponentTypes {
  OPERATOR,
  MODULE,
}

export type Component = Operator | Module

export interface Operator {
  type: ComponentTypes.OPERATOR
  io: IO
  // name: string
}

export interface Module {
  type: ComponentTypes.MODULE
  io: IO
  net: Connection[]
  deps: { [key: string]: IO }
  workers: { [key: string]: string }
}

export interface IO {
  in: { [key: string]: TypeDescriptor }
  out: { [key: string]: TypeDescriptor }
}

export interface TypeDescriptor {
  typeName: string
  genericArguments: TypeDescriptor[]
}

export interface Connection {
  from: PortAddr
  to: PortAddr
}

export interface PortAddr {
  node: string
  port: string
  idx: number
}
