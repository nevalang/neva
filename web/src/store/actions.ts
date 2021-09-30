import { Api } from "~api"
import { Program } from "~types/program"

export type PendingAction<T> = {
  type: T
}

export type SuccessAction<T, S> = {
  type: T
  payload?: S
}

export type FailAction<T, F> = {
  type: T
  error: F
}

export type AsyncActions<T1, T2, T3, P, S, F> =
  | PendingAction<T1>
  | SuccessAction<T2, S>
  | FailAction<T3, F>

export enum ActionTypes {
  // LOAD
  START_LOAD_PROGRAM,
  SUCCESS_LOAD_PROGRAM,
  FAIL_LOAD_PROGRAM,
  // UPDATE
  START_UPD_PROGRAM,
  SUCCESS_UPD_PROGRAM,
  FAIL_UPD_PROGRAM,
}

type LoadProgramActions = AsyncActions<
  ActionTypes.START_LOAD_PROGRAM,
  ActionTypes.SUCCESS_LOAD_PROGRAM,
  ActionTypes.FAIL_LOAD_PROGRAM,
  undefined,
  Program,
  Error
>

type UpdateProgramActions = AsyncActions<
  ActionTypes.START_UPD_PROGRAM,
  ActionTypes.SUCCESS_UPD_PROGRAM,
  ActionTypes.FAIL_UPD_PROGRAM,
  undefined,
  Program,
  Error
>

export type Action = LoadProgramActions | UpdateProgramActions

export type Dispatch = (action: Action) => void

export class Dispatcher {
  api: Api

  constructor(api: Api) {
    this.api = api
  }

  async loadProgram(dispatch: Dispatch, path: string) {
    dispatch({ type: ActionTypes.START_LOAD_PROGRAM })

    try {
      const program = await this.api.GetProgram(path)
      dispatch({ type: ActionTypes.SUCCESS_LOAD_PROGRAM })
    } catch (error) {
      dispatch({ type: ActionTypes.FAIL_LOAD_PROGRAM, error: error })
    }
  }

  async updateProgram(dispatch: Dispatch, program: Program, path: string) {
    dispatch({ type: ActionTypes.START_UPD_PROGRAM })

    try {
      await this.api.UpdateProgram(path)
      dispatch({ type: ActionTypes.SUCCESS_UPD_PROGRAM, payload: program })
    } catch (error) {
      dispatch({ type: ActionTypes.FAIL_UPD_PROGRAM, error: error })
    }
  }
}
