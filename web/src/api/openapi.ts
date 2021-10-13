import { AxiosPromise } from "axios"
import {
  DefaultApiFactory,
  Module as SDKModule,
  Component as SDKComponent,
  Operator as SDKOperator,
  Program as SDKProgram,
} from "../generated_web/index"
import {
  Component,
  ComponentTypes,
  Program,
  IO,
  Components,
} from "../types/program"
import { Api } from "./"

export class OpenApi implements Api {
  client: {
    programGet(options?: any): AxiosPromise<SDKProgram>
  }

  constructor(backendURL: string) {
    this.client = DefaultApiFactory(undefined, backendURL)
  }

  async getProgram(): Promise<Program> {
    try {
      const resp = await this.client.programGet()
      return this.castProgram(resp.data)
    } catch (err) {
      throw err
    }
  }

  castProgram(from: SDKProgram): Program {
    return {
      root: from.root,
      scope: this.castScope(from.scope),
    }
  }

  castScope(scope: { [key: string]: SDKComponent }): Components {
    const result = {}

    for (const name in scope) {
      const el = scope[name]
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

      result[name] = res
    }

    return result
  }
}
