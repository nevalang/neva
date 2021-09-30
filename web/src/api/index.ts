import { Program } from "../types/program"

interface Api {
  GetProgram(str: string): Promise<Program>
  UpdateProgram(Program)
}

export { Api }
