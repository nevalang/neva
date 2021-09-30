import { Program } from "~types/program"

export enum LoadingStatus {
  PENDING,
  FINISHED,
  FAILED,
}

export interface State {
  program:
    | {
        data: Program
        status: LoadingStatus.PENDING
        error?: Error
      }
    | {
        data: Program
        status: LoadingStatus.FINISHED
        error: null
      }
    | {
        data: Program
        status: LoadingStatus.FAILED
        error: Error
      }
}

export function createStore(data: Program): {
  data: Program
  status: LoadingStatus.FINISHED
  error: null
} {
  return {
    data,
    status: LoadingStatus.FINISHED,
    error: null,
  }
}
