import {
  DefaultApiFactory,
  Program as SDKProgram,
  Operator as SDKOperator,
} from "../../sdk"
import { Api } from "../"
import { Operator, Program } from "../../types/program"
import { CasterImpl } from "./caster"

export interface Caster {
  castProgram(from: SDKProgram): Program
  castOperator(from: SDKOperator): Operator
}

class OpenApi implements Api {
  client: ReturnType<typeof DefaultApiFactory>
  caster: Caster

  constructor(backendURL: string) {
    this.client = DefaultApiFactory(undefined, backendURL)
    this.caster = new CasterImpl()
  }

  async getPaths(): Promise<string[]> {
    try {
      const resp = await this.client.pathsGet()
      return resp.data
    } catch (err) {
      throw new Error("getPaths: " + err.message)
    }
  }

  async getOperators(): Promise<{ [key: string]: Operator }> {
    try {
      const resp = await this.client.operatorsGet()
      const res = {}
      for (const k in resp.data) {
        res[k] = this.caster.castOperator(resp.data[k])
      }
      return res
    } catch (err) {
      throw new Error("getPaths: " + err.message)
    }
  }

  async getProgram(path: string): Promise<Program> {
    try {
      const sdkProg = await this.client.programGet(path)
      return this.caster.castProgram(sdkProg.data)
    } catch (err) {
      throw new Error("programGet: " + err.message)
    }
  }

  async createProgram(path: string, program: Program): Promise<Program> {
    try {
      const resp = await this.client.programPost(path, program)
      return this.caster.castProgram(resp.data)
    } catch (err) {
      throw new Error("programPost: " + err.message)
    }
  }

  async editProgram(path: string, program: Program): Promise<Program> {
    try {
      const resp = await this.client.programPatch(path, program)
      return this.caster.castProgram(resp.data)
    } catch (err) {
      throw new Error("programPatch: " + err.message)
    }
  }
}

export { OpenApi }
