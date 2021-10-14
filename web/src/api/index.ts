import { Program } from "~types/program"

export interface Api {
  createProgram(path: string, program: Program): Promise<Program>
  editProgram(path: string, program: Program): Promise<Program>
  getProgram(path: string): Promise<Program>
}
