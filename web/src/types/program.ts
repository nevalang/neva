export interface Program {
  scope: Components
  root: string
}

export type Components = { [key: string]: Component }

export type Component = Operator | Module

export interface Operator {
  type: ComponentTypes.OPERATOR
  io: IO
}

export interface Module {
  type: ComponentTypes.MODULE
  io: IO
  deps: { [key: string]: IO }
  workers: { [key: string]: string }
  constants: { [key: string]: Const }
  net: Connection[]
}

export enum ComponentTypes {
  OPERATOR,
  MODULE,
}

export interface IO {
  in: { [key: string]: string }
  out: { [key: string]: string }
}

export interface Const {
  typ: string
  value: any
}

export interface Connection {
  from: PortAddr
  to: PortAddr
}

export interface PortAddr {
  node: string
  port: string
  idx?: number
}
