import axios from "axios"
import { Program } from "~types/program"

export class RestApi {
  apiAddress: string

  constructor(address: string) {
    this.apiAddress = address
  }

  async getProgram(path: string): Promise<Program> {
    try {
      const resp = await axios.get(`${this.apiAddress}/program?path=${path}`)
      return resp.data
    } catch (err) {
      throw new Error("axios get" + err)
    }
  }

  async UpdateProgram(program: Program): Promise<Program> {
    try {
      const resp = await axios.post(`program`, program)
      return resp.data
    } catch (err) {
      throw new Error("axios post" + err)
    }
  }
}
