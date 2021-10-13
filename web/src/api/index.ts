import { Program } from "~types/program"

export interface Api {
  getProgram(): Promise<Program>
}
