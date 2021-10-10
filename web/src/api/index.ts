import {
  DefaultApiFactory,
  Module as SDKModule,
  Operator as SDKOperator,
  Program as SDKProgram,
} from "../generated_web/index"
import { Component, ComponentTypes, Program, IO } from "../types/program"

const client = DefaultApiFactory(undefined, "http://localhost:8090")

export class OpenApi {
  async getProgram(): Promise<Program> {
    try {
      const resp = await client.programGet()
      console.log(resp)
      const prog = castProgram(resp.data)
      return prog
    } catch (err) {
      console.log("so", err)
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

function castScope(scope: { [key: string]: SDKOperator | SDKModule }): {
  [key: string]: Component
} {
  const r = {}

  for (const k in scope) {
    const el = scope[k]
    let res: Component

    if ((el as any).deps == undefined) {
      res = {
        type: ComponentTypes.OPERATOR,
        io: {
          in: (el as SDKOperator).io.in,
          out: (el as SDKOperator).io.out,
        },
      }
    } else {
      const deps: { [key: string]: IO } = {}

      for (const k in (el as SDKModule).deps) {
        const dep = (el as SDKModule).deps[k]
        deps[k] = { in: dep.in, out: dep.out }
      }

      res = {
        type: ComponentTypes.MODULE,
        io: { in: el.io.in, out: el.io.out },
        deps: deps,
        workers: (el as SDKModule).workers,
        net: (el as SDKModule).het.map(v => ({
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
