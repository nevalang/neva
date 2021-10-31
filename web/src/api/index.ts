import { Operator, Program } from "~types/program"

export interface Api {
  getPaths(): Promise<string[]>
  getProgram(path: string): Promise<Program>
  createProgram(path: string, program: Program): Promise<Program>
  editProgram(path: string, program: Program): Promise<Program>
  getOperators(): Promise<{ [key: string]: Operator }>
}
