import {
  Module as SDKModule,
  Component as SDKComponent,
  Operator as SDKOperator,
  Program as SDKProgram,
} from "../../sdk"

import {
  Component,
  ComponentTypes,
  Program,
  IO,
  Components,
  Const,
  Operator,
} from "../../types/program"

import { Caster } from "."

export class CasterImpl implements Caster {
  castProgram(from: SDKProgram): Program {
    return {
      root: from.root,
      scope: this.castScope(from.scope),
    }
  }

  castScope(scope: { [key: string]: SDKComponent }): Components {
    const result = {}

    for (const name in scope) {
      const sdkComponent = scope[name]
      let component: Component

      if ((sdkComponent as SDKModule).net == undefined) {
        component = this.castOperator(sdkComponent as SDKOperator)
      } else {
        const mod = sdkComponent as SDKModule

        const deps: { [key: string]: IO } = {}
        for (const k in mod.deps) {
          const dep = mod.deps[k]
          deps[k] = { in: dep.in, out: dep.out }
        }

        const constants: { [key: string]: Const } = {}
        for (const k in mod.const) {
          const c = mod.const[k]
          constants[k] = {
            typ: c.type,
            value: c.value,
          }
        }

        component = {
          type: ComponentTypes.MODULE,
          io: { in: sdkComponent.io.in, out: sdkComponent.io.out },
          deps: deps,
          workers: { ...mod.workers },
          constants,
          net: mod.net.map(v => ({
            from: { node: v.from.node, port: v.from.port, idx: v.from.idx },
            to: { node: v.to.node, port: v.to.port, idx: v.to.idx },
          })),
        }
      }

      result[name] = component
    }

    return result
  }

  castOperator(sdkComponent: SDKOperator): Operator {
    return {
      type: ComponentTypes.OPERATOR,
      io: {
        in: sdkComponent.io.in,
        out: sdkComponent.io.out,
      },
    }
  }
}
