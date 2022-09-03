import * as grpcWeb from 'grpc-web';

import {
  DebugRequest,
  DebugResponse,
  GetProgramRequest,
  GetProgramResponse,
  ListProgramsRequest,
  ListProgramsResponse,
  StartDebugRequest,
  StartDebugResponse,
  UpdateProgramRequest,
  UpdateProgramResponse} from './devserver_pb';

export class DevClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  listPrograms(
    request: ListProgramsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ListProgramsResponse) => void
  ): grpcWeb.ClientReadableStream<ListProgramsResponse>;

  getProgram(
    request: GetProgramRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GetProgramResponse) => void
  ): grpcWeb.ClientReadableStream<GetProgramResponse>;

  updateProgram(
    request: UpdateProgramRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UpdateProgramResponse) => void
  ): grpcWeb.ClientReadableStream<UpdateProgramResponse>;

  startDebugger(
    request: StartDebugRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<StartDebugResponse>;

  sendDebugMessage(
    request: DebugRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<DebugResponse>;

}

export class DevPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  listPrograms(
    request: ListProgramsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ListProgramsResponse>;

  getProgram(
    request: GetProgramRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GetProgramResponse>;

  updateProgram(
    request: UpdateProgramRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UpdateProgramResponse>;

  startDebugger(
    request: StartDebugRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<StartDebugResponse>;

  sendDebugMessage(
    request: DebugRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<DebugResponse>;

}

