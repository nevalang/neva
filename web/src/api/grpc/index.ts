import {} from "../../sdk/devserver_pb"
import { DevPromiseClient } from "../../sdk/devserver_grpc_web_pb"

import { Operator, Program } from "~types/program"

export class GrpcClient {
  client: DevPromiseClient

  constructor() {
    this.client = new DevPromiseClient("http://localhost:8090")
  }

  async getPaths(): Promise<string[]> {
    return []
  }

  async getProgram(path: string): Promise<Program> {
    return
  }

  async createProgram(path: string, program: Program): Promise<Program> {
    return
  }

  async editProgram(path: string, program: Program): Promise<Program> {
    return
  }

  async getOperators(): Promise<{ [key: string]: Operator }> {
    return
  }
}
