import { DefaultApi, Module, SDKOperator, Program as SDKProgram } from "../../sdk"
import {
  Component,
  ComponentTypes,
  Program,
  Operator,
  IO,
} from "../types/program"

export interface Api {
  getProgram(): Promise<Program>
}

export class OpenApi {
  client: DefaultApi

  constructor(client: DefaultApi) {
    this.client = client
  }

  async getProgram(): Promise<Program> {
    try {
      console.log('here')
      const resp = await this.client.programGet()
      const prog = castProgram(resp)
      return prog
    } catch (err) {
      throw err
    }
  }
}

function castProgram(from: SDKProgram) {
  return {
    root: from.root,
    scope: castScope(from.scope),
  }
}

function castScope(scope: { [key: string]: SDKOperator | Module }): {
  [key: string]: Component
} {
  const r = {}

  for (const k in scope) {
    const el = scope[k]
    let res: Component

    if (el instanceof SDKOperator) {
      res = {
        type: ComponentTypes.OPERATOR,
        io: { in: el.io._in, out: el.io.out },
      } as Operator
    } else {
      const deps: { [key: string]: IO } = {}
      for (const k in el.deps) {
        const dep = el.deps[k]
        deps[k] = { in: dep._in, out: dep.out }
      }

      res = {
        type: ComponentTypes.MODULE,
        io: { in: el.io._in, out: el.io.out },
        deps: deps,
        workers: el.workers,
        net: el.het.map(v => ({
          from: {
            node: v.from.node,
            port: v.from.port,
            idx: v.from.idx,
          },
          to: {
            node: v.to.node,
            port: v.to.port,
            idx: v.to.idx,
          },
        })),
      }
    }

    r[k] = res
  }

  return r
}
