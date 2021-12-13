import * as jspb from "google-protobuf"

export class DebugRequest extends jspb.Message {
  getMsg(): Msg | undefined;
  setMsg(value?: Msg): void;
  hasMsg(): boolean;
  clearMsg(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DebugRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DebugRequest): DebugRequest.AsObject;
  static serializeBinaryToWriter(message: DebugRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DebugRequest;
  static deserializeBinaryFromReader(message: DebugRequest, reader: jspb.BinaryReader): DebugRequest;
}

export namespace DebugRequest {
  export type AsObject = {
    msg?: Msg.AsObject,
  }
}

export class DebugResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DebugResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DebugResponse): DebugResponse.AsObject;
  static serializeBinaryToWriter(message: DebugResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DebugResponse;
  static deserializeBinaryFromReader(message: DebugResponse, reader: jspb.BinaryReader): DebugResponse;
}

export namespace DebugResponse {
  export type AsObject = {
  }
}

export class StartDebugRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartDebugRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartDebugRequest): StartDebugRequest.AsObject;
  static serializeBinaryToWriter(message: StartDebugRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartDebugRequest;
  static deserializeBinaryFromReader(message: StartDebugRequest, reader: jspb.BinaryReader): StartDebugRequest;
}

export namespace StartDebugRequest {
  export type AsObject = {
    path: string,
  }
}

export class StartDebugResponse extends jspb.Message {
  getEvent(): DebugEvent | undefined;
  setEvent(value?: DebugEvent): void;
  hasEvent(): boolean;
  clearEvent(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartDebugResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartDebugResponse): StartDebugResponse.AsObject;
  static serializeBinaryToWriter(message: StartDebugResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartDebugResponse;
  static deserializeBinaryFromReader(message: StartDebugResponse, reader: jspb.BinaryReader): StartDebugResponse;
}

export namespace StartDebugResponse {
  export type AsObject = {
    event?: DebugEvent.AsObject,
  }
}

export class DebugEvent extends jspb.Message {
  getFrom(): string;
  setFrom(value: string): void;

  getTo(): string;
  setTo(value: string): void;

  getMsg(): Msg | undefined;
  setMsg(value?: Msg): void;
  hasMsg(): boolean;
  clearMsg(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DebugEvent.AsObject;
  static toObject(includeInstance: boolean, msg: DebugEvent): DebugEvent.AsObject;
  static serializeBinaryToWriter(message: DebugEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DebugEvent;
  static deserializeBinaryFromReader(message: DebugEvent, reader: jspb.BinaryReader): DebugEvent;
}

export namespace DebugEvent {
  export type AsObject = {
    from: string,
    to: string,
    msg?: Msg.AsObject,
  }
}

export class Msg extends jspb.Message {
  getType(): ValueType;
  setType(value: ValueType): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Msg.AsObject;
  static toObject(includeInstance: boolean, msg: Msg): Msg.AsObject;
  static serializeBinaryToWriter(message: Msg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Msg;
  static deserializeBinaryFromReader(message: Msg, reader: jspb.BinaryReader): Msg;
}

export namespace Msg {
  export type AsObject = {
    type: ValueType,
  }
}

export class ListProgramsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProgramsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProgramsRequest): ListProgramsRequest.AsObject;
  static serializeBinaryToWriter(message: ListProgramsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProgramsRequest;
  static deserializeBinaryFromReader(message: ListProgramsRequest, reader: jspb.BinaryReader): ListProgramsRequest;
}

export namespace ListProgramsRequest {
  export type AsObject = {
  }
}

export class ListProgramsResponse extends jspb.Message {
  getPathsList(): Array<string>;
  setPathsList(value: Array<string>): void;
  clearPathsList(): void;
  addPaths(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProgramsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProgramsResponse): ListProgramsResponse.AsObject;
  static serializeBinaryToWriter(message: ListProgramsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProgramsResponse;
  static deserializeBinaryFromReader(message: ListProgramsResponse, reader: jspb.BinaryReader): ListProgramsResponse;
}

export namespace ListProgramsResponse {
  export type AsObject = {
    pathsList: Array<string>,
  }
}

export class GetProgramRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProgramRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProgramRequest): GetProgramRequest.AsObject;
  static serializeBinaryToWriter(message: GetProgramRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProgramRequest;
  static deserializeBinaryFromReader(message: GetProgramRequest, reader: jspb.BinaryReader): GetProgramRequest;
}

export namespace GetProgramRequest {
  export type AsObject = {
    path: string,
  }
}

export class GetProgramResponse extends jspb.Message {
  getProgram(): Program | undefined;
  setProgram(value?: Program): void;
  hasProgram(): boolean;
  clearProgram(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProgramResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProgramResponse): GetProgramResponse.AsObject;
  static serializeBinaryToWriter(message: GetProgramResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProgramResponse;
  static deserializeBinaryFromReader(message: GetProgramResponse, reader: jspb.BinaryReader): GetProgramResponse;
}

export namespace GetProgramResponse {
  export type AsObject = {
    program?: Program.AsObject,
  }
}

export class UpdateProgramRequest extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  getProgram(): Program | undefined;
  setProgram(value?: Program): void;
  hasProgram(): boolean;
  clearProgram(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProgramRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProgramRequest): UpdateProgramRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProgramRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProgramRequest;
  static deserializeBinaryFromReader(message: UpdateProgramRequest, reader: jspb.BinaryReader): UpdateProgramRequest;
}

export namespace UpdateProgramRequest {
  export type AsObject = {
    path: string,
    program?: Program.AsObject,
  }
}

export class UpdateProgramResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProgramResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProgramResponse): UpdateProgramResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProgramResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProgramResponse;
  static deserializeBinaryFromReader(message: UpdateProgramResponse, reader: jspb.BinaryReader): UpdateProgramResponse;
}

export namespace UpdateProgramResponse {
  export type AsObject = {
  }
}

export class Program extends jspb.Message {
  getScopeMap(): jspb.Map<string, Component>;
  clearScopeMap(): void;

  getRoot(): string;
  setRoot(value: string): void;

  getImports(): ProgramImports | undefined;
  setImports(value?: ProgramImports): void;
  hasImports(): boolean;
  clearImports(): void;

  getMeta(): ProgramMeta | undefined;
  setMeta(value?: ProgramMeta): void;
  hasMeta(): boolean;
  clearMeta(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Program.AsObject;
  static toObject(includeInstance: boolean, msg: Program): Program.AsObject;
  static serializeBinaryToWriter(message: Program, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Program;
  static deserializeBinaryFromReader(message: Program, reader: jspb.BinaryReader): Program;
}

export namespace Program {
  export type AsObject = {
    scopeMap: Array<[string, Component.AsObject]>,
    root: string,
    imports?: ProgramImports.AsObject,
    meta?: ProgramMeta.AsObject,
  }
}

export class ProgramMeta extends jspb.Message {
  getCompilerVersion(): string;
  setCompilerVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgramMeta.AsObject;
  static toObject(includeInstance: boolean, msg: ProgramMeta): ProgramMeta.AsObject;
  static serializeBinaryToWriter(message: ProgramMeta, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProgramMeta;
  static deserializeBinaryFromReader(message: ProgramMeta, reader: jspb.BinaryReader): ProgramMeta;
}

export namespace ProgramMeta {
  export type AsObject = {
    compilerVersion: string,
  }
}

export class ProgramImports extends jspb.Message {
  getStdMap(): jspb.Map<string, string>;
  clearStdMap(): void;

  getGlobalMap(): jspb.Map<string, string>;
  clearGlobalMap(): void;

  getLocalMap(): jspb.Map<string, string>;
  clearLocalMap(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProgramImports.AsObject;
  static toObject(includeInstance: boolean, msg: ProgramImports): ProgramImports.AsObject;
  static serializeBinaryToWriter(message: ProgramImports, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProgramImports;
  static deserializeBinaryFromReader(message: ProgramImports, reader: jspb.BinaryReader): ProgramImports;
}

export namespace ProgramImports {
  export type AsObject = {
    stdMap: Array<[string, string]>,
    globalMap: Array<[string, string]>,
    localMap: Array<[string, string]>,
  }
}

export class Component extends jspb.Message {
  getType(): ComponentType;
  setType(value: ComponentType): void;

  getOperator(): Operator | undefined;
  setOperator(value?: Operator): void;
  hasOperator(): boolean;
  clearOperator(): void;

  getModule(): Module | undefined;
  setModule(value?: Module): void;
  hasModule(): boolean;
  clearModule(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Component.AsObject;
  static toObject(includeInstance: boolean, msg: Component): Component.AsObject;
  static serializeBinaryToWriter(message: Component, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Component;
  static deserializeBinaryFromReader(message: Component, reader: jspb.BinaryReader): Component;
}

export namespace Component {
  export type AsObject = {
    type: ComponentType,
    operator?: Operator.AsObject,
    module?: Module.AsObject,
  }
}

export class Operator extends jspb.Message {
  getIo(): ComponentIO | undefined;
  setIo(value?: ComponentIO): void;
  hasIo(): boolean;
  clearIo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Operator.AsObject;
  static toObject(includeInstance: boolean, msg: Operator): Operator.AsObject;
  static serializeBinaryToWriter(message: Operator, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Operator;
  static deserializeBinaryFromReader(message: Operator, reader: jspb.BinaryReader): Operator;
}

export namespace Operator {
  export type AsObject = {
    io?: ComponentIO.AsObject,
  }
}

export class Module extends jspb.Message {
  getIo(): ComponentIO | undefined;
  setIo(value?: ComponentIO): void;
  hasIo(): boolean;
  clearIo(): void;

  getDepsMap(): jspb.Map<string, ComponentIO>;
  clearDepsMap(): void;

  getWorkersMap(): jspb.Map<string, string>;
  clearWorkersMap(): void;

  getNetList(): Array<ModuleConnection>;
  setNetList(value: Array<ModuleConnection>): void;
  clearNetList(): void;
  addNet(value?: ModuleConnection, index?: number): ModuleConnection;

  getConstMap(): jspb.Map<string, ModuleConst>;
  clearConstMap(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Module.AsObject;
  static toObject(includeInstance: boolean, msg: Module): Module.AsObject;
  static serializeBinaryToWriter(message: Module, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Module;
  static deserializeBinaryFromReader(message: Module, reader: jspb.BinaryReader): Module;
}

export namespace Module {
  export type AsObject = {
    io?: ComponentIO.AsObject,
    depsMap: Array<[string, ComponentIO.AsObject]>,
    workersMap: Array<[string, string]>,
    netList: Array<ModuleConnection.AsObject>,
    constMap: Array<[string, ModuleConst.AsObject]>,
  }
}

export class ModuleConst extends jspb.Message {
  getType(): ValueType;
  setType(value: ValueType): void;

  getIntValue(): number;
  setIntValue(value: number): void;

  getStrValue(): string;
  setStrValue(value: string): void;

  getBoolValue(): boolean;
  setBoolValue(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ModuleConst.AsObject;
  static toObject(includeInstance: boolean, msg: ModuleConst): ModuleConst.AsObject;
  static serializeBinaryToWriter(message: ModuleConst, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ModuleConst;
  static deserializeBinaryFromReader(message: ModuleConst, reader: jspb.BinaryReader): ModuleConst;
}

export namespace ModuleConst {
  export type AsObject = {
    type: ValueType,
    intValue: number,
    strValue: string,
    boolValue: boolean,
  }
}

export class ModuleConnection extends jspb.Message {
  getFrom(): ConnectionAddr | undefined;
  setFrom(value?: ConnectionAddr): void;
  hasFrom(): boolean;
  clearFrom(): void;

  getTo(): ConnectionAddr | undefined;
  setTo(value?: ConnectionAddr): void;
  hasTo(): boolean;
  clearTo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ModuleConnection.AsObject;
  static toObject(includeInstance: boolean, msg: ModuleConnection): ModuleConnection.AsObject;
  static serializeBinaryToWriter(message: ModuleConnection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ModuleConnection;
  static deserializeBinaryFromReader(message: ModuleConnection, reader: jspb.BinaryReader): ModuleConnection;
}

export namespace ModuleConnection {
  export type AsObject = {
    from?: ConnectionAddr.AsObject,
    to?: ConnectionAddr.AsObject,
  }
}

export class ConnectionAddr extends jspb.Message {
  getNode(): string;
  setNode(value: string): void;

  getPort(): string;
  setPort(value: string): void;

  getIdx(): string;
  setIdx(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConnectionAddr.AsObject;
  static toObject(includeInstance: boolean, msg: ConnectionAddr): ConnectionAddr.AsObject;
  static serializeBinaryToWriter(message: ConnectionAddr, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConnectionAddr;
  static deserializeBinaryFromReader(message: ConnectionAddr, reader: jspb.BinaryReader): ConnectionAddr;
}

export namespace ConnectionAddr {
  export type AsObject = {
    node: string,
    port: string,
    idx: string,
  }
}

export class ComponentIO extends jspb.Message {
  getInMap(): jspb.Map<string, Port>;
  clearInMap(): void;

  getOutMap(): jspb.Map<string, Port>;
  clearOutMap(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ComponentIO.AsObject;
  static toObject(includeInstance: boolean, msg: ComponentIO): ComponentIO.AsObject;
  static serializeBinaryToWriter(message: ComponentIO, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ComponentIO;
  static deserializeBinaryFromReader(message: ComponentIO, reader: jspb.BinaryReader): ComponentIO;
}

export namespace ComponentIO {
  export type AsObject = {
    inMap: Array<[string, Port.AsObject]>,
    outMap: Array<[string, Port.AsObject]>,
  }
}

export class Port extends jspb.Message {
  getIsArray(): boolean;
  setIsArray(value: boolean): void;

  getType(): ValueType;
  setType(value: ValueType): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Port.AsObject;
  static toObject(includeInstance: boolean, msg: Port): Port.AsObject;
  static serializeBinaryToWriter(message: Port, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Port;
  static deserializeBinaryFromReader(message: Port, reader: jspb.BinaryReader): Port;
}

export namespace Port {
  export type AsObject = {
    isArray: boolean,
    type: ValueType,
  }
}

export enum ComponentType { 
  COMPONENT_TYPE_OPERATOR = 0,
  COMPONENT_TYPE_MODULE = 1,
}
export enum ValueType { 
  VALUE_TYPE_INT = 0,
  VALUE_TYPE_STR = 1,
  VALUE_TYPE_BOOL = 2,
}
