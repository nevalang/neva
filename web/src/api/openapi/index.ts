import { DefaultApiFactory, Program as SDKProgram } from "../../sdk"
import { Api } from "../"
import { Program } from "../../types/program"
import { CasterImpl } from "./caster"

export interface Caster {
  castProgram(from: SDKProgram): Program
}

class OpenApi implements Api {
  client: ReturnType<typeof DefaultApiFactory>
  caster: Caster

  constructor(backendURL: string) {
    this.client = DefaultApiFactory(undefined, backendURL)
    this.caster = new CasterImpl()
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
