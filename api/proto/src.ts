/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "srcproto";

export enum EntityKind {
  ENTITY_KIND_UNSPECIFIED = 0,
  ENTITY_KIND_CONST = 1,
  ENTITY_KIND_TYPE_DEF = 2,
  ENTITY_KIND_INTERFACE = 3,
  ENTITY_KIND_COMPONENT = 4,
  UNRECOGNIZED = -1,
}

export function entityKindFromJSON(object: any): EntityKind {
  switch (object) {
    case 0:
    case "ENTITY_KIND_UNSPECIFIED":
      return EntityKind.ENTITY_KIND_UNSPECIFIED;
    case 1:
    case "ENTITY_KIND_CONST":
      return EntityKind.ENTITY_KIND_CONST;
    case 2:
    case "ENTITY_KIND_TYPE_DEF":
      return EntityKind.ENTITY_KIND_TYPE_DEF;
    case 3:
    case "ENTITY_KIND_INTERFACE":
      return EntityKind.ENTITY_KIND_INTERFACE;
    case 4:
    case "ENTITY_KIND_COMPONENT":
      return EntityKind.ENTITY_KIND_COMPONENT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return EntityKind.UNRECOGNIZED;
  }
}

export function entityKindToJSON(object: EntityKind): string {
  switch (object) {
    case EntityKind.ENTITY_KIND_UNSPECIFIED:
      return "ENTITY_KIND_UNSPECIFIED";
    case EntityKind.ENTITY_KIND_CONST:
      return "ENTITY_KIND_CONST";
    case EntityKind.ENTITY_KIND_TYPE_DEF:
      return "ENTITY_KIND_TYPE_DEF";
    case EntityKind.ENTITY_KIND_INTERFACE:
      return "ENTITY_KIND_INTERFACE";
    case EntityKind.ENTITY_KIND_COMPONENT:
      return "ENTITY_KIND_COMPONENT";
    case EntityKind.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface File {
  imports: { [key: string]: string };
  entities: { [key: string]: Entity };
}

export interface File_ImportsEntry {
  key: string;
  value: string;
}

export interface File_EntitiesEntry {
  key: string;
  value: Entity | undefined;
}

export interface Entity {
  exported: boolean;
  kind: EntityKind;
  const: Const | undefined;
  typeDef: TypeDef | undefined;
  interface: Interface | undefined;
  component: Component | undefined;
}

export interface Const {
  ref: EntityRef | undefined;
  value: Msg | undefined;
}

export interface Msg {
  typeExpr: string;
  bool: boolean;
  int: number;
  float: number;
  str: string;
  vecs: Const[];
  map: { [key: string]: Const };
}

export interface Msg_MapEntry {
  key: string;
  value: Const | undefined;
}

export interface TypeDef {
  params: TypeParam[];
  bodyExpr: string;
}

export interface Component {
  interface: Interface | undefined;
  nodes: { [key: string]: Node };
  connections: Connection[];
}

export interface Component_NodesEntry {
  key: string;
  value: Node | undefined;
}

export interface Interface {
  typeParams: TypeParam[];
  io: IO | undefined;
}

export interface TypeParam {
  name: string;
  constr: string;
}

export interface IO {
  ins: Port[];
  outs: Port[];
}

export interface Port {
  typeExpr: string;
  isArray: boolean;
}

export interface Node {
  entityRef: EntityRef | undefined;
  typeArgs: string[];
  componentDis: { [key: string]: Node };
}

export interface Node_ComponentDisEntry {
  key: string;
  value: Node | undefined;
}

export interface EntityRef {
  pkg: string;
  name: string;
}

export interface Connection {
  senderSide: SenderConnectionSide | undefined;
  receiverSides: ReceiverConnectionSide[];
}

export interface ReceiverConnectionSide {
  portAddr: PortAddr | undefined;
  selectors: string[];
}

export interface SenderConnectionSide {
  constRef: EntityRef | undefined;
  portAddr: PortAddr | undefined;
  selectors: string[];
}

export interface PortAddr {
  node: string;
  port: string;
  idx: number;
}

function createBaseFile(): File {
  return { imports: {}, entities: {} };
}

export const File = {
  encode(message: File, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.imports).forEach(([key, value]) => {
      File_ImportsEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    Object.entries(message.entities).forEach(([key, value]) => {
      File_EntitiesEntry.encode({ key: key as any, value }, writer.uint32(18).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): File {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFile();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          const entry1 = File_ImportsEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.imports[entry1.key] = entry1.value;
          }
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          const entry2 = File_EntitiesEntry.decode(reader, reader.uint32());
          if (entry2.value !== undefined) {
            message.entities[entry2.key] = entry2.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): File {
    return {
      imports: isObject(object.imports)
        ? Object.entries(object.imports).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
      entities: isObject(object.entities)
        ? Object.entries(object.entities).reduce<{ [key: string]: Entity }>((acc, [key, value]) => {
          acc[key] = Entity.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: File): unknown {
    const obj: any = {};
    if (message.imports) {
      const entries = Object.entries(message.imports);
      if (entries.length > 0) {
        obj.imports = {};
        entries.forEach(([k, v]) => {
          obj.imports[k] = v;
        });
      }
    }
    if (message.entities) {
      const entries = Object.entries(message.entities);
      if (entries.length > 0) {
        obj.entities = {};
        entries.forEach(([k, v]) => {
          obj.entities[k] = Entity.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<File>, I>>(base?: I): File {
    return File.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<File>, I>>(object: I): File {
    const message = createBaseFile();
    message.imports = Object.entries(object.imports ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = globalThis.String(value);
      }
      return acc;
    }, {});
    message.entities = Object.entries(object.entities ?? {}).reduce<{ [key: string]: Entity }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = Entity.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseFile_ImportsEntry(): File_ImportsEntry {
  return { key: "", value: "" };
}

export const File_ImportsEntry = {
  encode(message: File_ImportsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): File_ImportsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFile_ImportsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): File_ImportsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? globalThis.String(object.value) : "",
    };
  },

  toJSON(message: File_ImportsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== "") {
      obj.value = message.value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<File_ImportsEntry>, I>>(base?: I): File_ImportsEntry {
    return File_ImportsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<File_ImportsEntry>, I>>(object: I): File_ImportsEntry {
    const message = createBaseFile_ImportsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseFile_EntitiesEntry(): File_EntitiesEntry {
  return { key: "", value: undefined };
}

export const File_EntitiesEntry = {
  encode(message: File_EntitiesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Entity.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): File_EntitiesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFile_EntitiesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = Entity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): File_EntitiesEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? Entity.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: File_EntitiesEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = Entity.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<File_EntitiesEntry>, I>>(base?: I): File_EntitiesEntry {
    return File_EntitiesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<File_EntitiesEntry>, I>>(object: I): File_EntitiesEntry {
    const message = createBaseFile_EntitiesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? Entity.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseEntity(): Entity {
  return { exported: false, kind: 0, const: undefined, typeDef: undefined, interface: undefined, component: undefined };
}

export const Entity = {
  encode(message: Entity, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exported === true) {
      writer.uint32(8).bool(message.exported);
    }
    if (message.kind !== 0) {
      writer.uint32(16).int32(message.kind);
    }
    if (message.const !== undefined) {
      Const.encode(message.const, writer.uint32(26).fork()).ldelim();
    }
    if (message.typeDef !== undefined) {
      TypeDef.encode(message.typeDef, writer.uint32(34).fork()).ldelim();
    }
    if (message.interface !== undefined) {
      Interface.encode(message.interface, writer.uint32(42).fork()).ldelim();
    }
    if (message.component !== undefined) {
      Component.encode(message.component, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Entity {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntity();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.exported = reader.bool();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.kind = reader.int32() as any;
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.const = Const.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.typeDef = TypeDef.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.interface = Interface.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.component = Component.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Entity {
    return {
      exported: isSet(object.exported) ? globalThis.Boolean(object.exported) : false,
      kind: isSet(object.kind) ? entityKindFromJSON(object.kind) : 0,
      const: isSet(object.const) ? Const.fromJSON(object.const) : undefined,
      typeDef: isSet(object.typeDef) ? TypeDef.fromJSON(object.typeDef) : undefined,
      interface: isSet(object.interface) ? Interface.fromJSON(object.interface) : undefined,
      component: isSet(object.component) ? Component.fromJSON(object.component) : undefined,
    };
  },

  toJSON(message: Entity): unknown {
    const obj: any = {};
    if (message.exported === true) {
      obj.exported = message.exported;
    }
    if (message.kind !== 0) {
      obj.kind = entityKindToJSON(message.kind);
    }
    if (message.const !== undefined) {
      obj.const = Const.toJSON(message.const);
    }
    if (message.typeDef !== undefined) {
      obj.typeDef = TypeDef.toJSON(message.typeDef);
    }
    if (message.interface !== undefined) {
      obj.interface = Interface.toJSON(message.interface);
    }
    if (message.component !== undefined) {
      obj.component = Component.toJSON(message.component);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Entity>, I>>(base?: I): Entity {
    return Entity.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Entity>, I>>(object: I): Entity {
    const message = createBaseEntity();
    message.exported = object.exported ?? false;
    message.kind = object.kind ?? 0;
    message.const = (object.const !== undefined && object.const !== null) ? Const.fromPartial(object.const) : undefined;
    message.typeDef = (object.typeDef !== undefined && object.typeDef !== null)
      ? TypeDef.fromPartial(object.typeDef)
      : undefined;
    message.interface = (object.interface !== undefined && object.interface !== null)
      ? Interface.fromPartial(object.interface)
      : undefined;
    message.component = (object.component !== undefined && object.component !== null)
      ? Component.fromPartial(object.component)
      : undefined;
    return message;
  },
};

function createBaseConst(): Const {
  return { ref: undefined, value: undefined };
}

export const Const = {
  encode(message: Const, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ref !== undefined) {
      EntityRef.encode(message.ref, writer.uint32(10).fork()).ldelim();
    }
    if (message.value !== undefined) {
      Msg.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Const {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConst();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.ref = EntityRef.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = Msg.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Const {
    return {
      ref: isSet(object.ref) ? EntityRef.fromJSON(object.ref) : undefined,
      value: isSet(object.value) ? Msg.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Const): unknown {
    const obj: any = {};
    if (message.ref !== undefined) {
      obj.ref = EntityRef.toJSON(message.ref);
    }
    if (message.value !== undefined) {
      obj.value = Msg.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Const>, I>>(base?: I): Const {
    return Const.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Const>, I>>(object: I): Const {
    const message = createBaseConst();
    message.ref = (object.ref !== undefined && object.ref !== null) ? EntityRef.fromPartial(object.ref) : undefined;
    message.value = (object.value !== undefined && object.value !== null) ? Msg.fromPartial(object.value) : undefined;
    return message;
  },
};

function createBaseMsg(): Msg {
  return { typeExpr: "", bool: false, int: 0, float: 0, str: "", vecs: [], map: {} };
}

export const Msg = {
  encode(message: Msg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.typeExpr !== "") {
      writer.uint32(10).string(message.typeExpr);
    }
    if (message.bool === true) {
      writer.uint32(16).bool(message.bool);
    }
    if (message.int !== 0) {
      writer.uint32(24).int64(message.int);
    }
    if (message.float !== 0) {
      writer.uint32(33).double(message.float);
    }
    if (message.str !== "") {
      writer.uint32(42).string(message.str);
    }
    for (const v of message.vecs) {
      Const.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    Object.entries(message.map).forEach(([key, value]) => {
      Msg_MapEntry.encode({ key: key as any, value }, writer.uint32(58).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Msg {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.typeExpr = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.bool = reader.bool();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.int = longToNumber(reader.int64() as Long);
          continue;
        case 4:
          if (tag !== 33) {
            break;
          }

          message.float = reader.double();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.str = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.vecs.push(Const.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          const entry7 = Msg_MapEntry.decode(reader, reader.uint32());
          if (entry7.value !== undefined) {
            message.map[entry7.key] = entry7.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Msg {
    return {
      typeExpr: isSet(object.typeExpr) ? globalThis.String(object.typeExpr) : "",
      bool: isSet(object.bool) ? globalThis.Boolean(object.bool) : false,
      int: isSet(object.int) ? globalThis.Number(object.int) : 0,
      float: isSet(object.float) ? globalThis.Number(object.float) : 0,
      str: isSet(object.str) ? globalThis.String(object.str) : "",
      vecs: globalThis.Array.isArray(object?.vecs) ? object.vecs.map((e: any) => Const.fromJSON(e)) : [],
      map: isObject(object.map)
        ? Object.entries(object.map).reduce<{ [key: string]: Const }>((acc, [key, value]) => {
          acc[key] = Const.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Msg): unknown {
    const obj: any = {};
    if (message.typeExpr !== "") {
      obj.typeExpr = message.typeExpr;
    }
    if (message.bool === true) {
      obj.bool = message.bool;
    }
    if (message.int !== 0) {
      obj.int = Math.round(message.int);
    }
    if (message.float !== 0) {
      obj.float = message.float;
    }
    if (message.str !== "") {
      obj.str = message.str;
    }
    if (message.vecs?.length) {
      obj.vecs = message.vecs.map((e) => Const.toJSON(e));
    }
    if (message.map) {
      const entries = Object.entries(message.map);
      if (entries.length > 0) {
        obj.map = {};
        entries.forEach(([k, v]) => {
          obj.map[k] = Const.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Msg>, I>>(base?: I): Msg {
    return Msg.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Msg>, I>>(object: I): Msg {
    const message = createBaseMsg();
    message.typeExpr = object.typeExpr ?? "";
    message.bool = object.bool ?? false;
    message.int = object.int ?? 0;
    message.float = object.float ?? 0;
    message.str = object.str ?? "";
    message.vecs = object.vecs?.map((e) => Const.fromPartial(e)) || [];
    message.map = Object.entries(object.map ?? {}).reduce<{ [key: string]: Const }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = Const.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseMsg_MapEntry(): Msg_MapEntry {
  return { key: "", value: undefined };
}

export const Msg_MapEntry = {
  encode(message: Msg_MapEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Const.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Msg_MapEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsg_MapEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = Const.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Msg_MapEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? Const.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Msg_MapEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = Const.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Msg_MapEntry>, I>>(base?: I): Msg_MapEntry {
    return Msg_MapEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Msg_MapEntry>, I>>(object: I): Msg_MapEntry {
    const message = createBaseMsg_MapEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null) ? Const.fromPartial(object.value) : undefined;
    return message;
  },
};

function createBaseTypeDef(): TypeDef {
  return { params: [], bodyExpr: "" };
}

export const TypeDef = {
  encode(message: TypeDef, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.params) {
      TypeParam.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.bodyExpr !== "") {
      writer.uint32(18).string(message.bodyExpr);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TypeDef {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTypeDef();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.params.push(TypeParam.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.bodyExpr = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TypeDef {
    return {
      params: globalThis.Array.isArray(object?.params) ? object.params.map((e: any) => TypeParam.fromJSON(e)) : [],
      bodyExpr: isSet(object.bodyExpr) ? globalThis.String(object.bodyExpr) : "",
    };
  },

  toJSON(message: TypeDef): unknown {
    const obj: any = {};
    if (message.params?.length) {
      obj.params = message.params.map((e) => TypeParam.toJSON(e));
    }
    if (message.bodyExpr !== "") {
      obj.bodyExpr = message.bodyExpr;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TypeDef>, I>>(base?: I): TypeDef {
    return TypeDef.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TypeDef>, I>>(object: I): TypeDef {
    const message = createBaseTypeDef();
    message.params = object.params?.map((e) => TypeParam.fromPartial(e)) || [];
    message.bodyExpr = object.bodyExpr ?? "";
    return message;
  },
};

function createBaseComponent(): Component {
  return { interface: undefined, nodes: {}, connections: [] };
}

export const Component = {
  encode(message: Component, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.interface !== undefined) {
      Interface.encode(message.interface, writer.uint32(10).fork()).ldelim();
    }
    Object.entries(message.nodes).forEach(([key, value]) => {
      Component_NodesEntry.encode({ key: key as any, value }, writer.uint32(18).fork()).ldelim();
    });
    for (const v of message.connections) {
      Connection.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Component {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComponent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.interface = Interface.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          const entry2 = Component_NodesEntry.decode(reader, reader.uint32());
          if (entry2.value !== undefined) {
            message.nodes[entry2.key] = entry2.value;
          }
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.connections.push(Connection.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Component {
    return {
      interface: isSet(object.interface) ? Interface.fromJSON(object.interface) : undefined,
      nodes: isObject(object.nodes)
        ? Object.entries(object.nodes).reduce<{ [key: string]: Node }>((acc, [key, value]) => {
          acc[key] = Node.fromJSON(value);
          return acc;
        }, {})
        : {},
      connections: globalThis.Array.isArray(object?.connections)
        ? object.connections.map((e: any) => Connection.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Component): unknown {
    const obj: any = {};
    if (message.interface !== undefined) {
      obj.interface = Interface.toJSON(message.interface);
    }
    if (message.nodes) {
      const entries = Object.entries(message.nodes);
      if (entries.length > 0) {
        obj.nodes = {};
        entries.forEach(([k, v]) => {
          obj.nodes[k] = Node.toJSON(v);
        });
      }
    }
    if (message.connections?.length) {
      obj.connections = message.connections.map((e) => Connection.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Component>, I>>(base?: I): Component {
    return Component.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Component>, I>>(object: I): Component {
    const message = createBaseComponent();
    message.interface = (object.interface !== undefined && object.interface !== null)
      ? Interface.fromPartial(object.interface)
      : undefined;
    message.nodes = Object.entries(object.nodes ?? {}).reduce<{ [key: string]: Node }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = Node.fromPartial(value);
      }
      return acc;
    }, {});
    message.connections = object.connections?.map((e) => Connection.fromPartial(e)) || [];
    return message;
  },
};

function createBaseComponent_NodesEntry(): Component_NodesEntry {
  return { key: "", value: undefined };
}

export const Component_NodesEntry = {
  encode(message: Component_NodesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Node.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Component_NodesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseComponent_NodesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = Node.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Component_NodesEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? Node.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Component_NodesEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = Node.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Component_NodesEntry>, I>>(base?: I): Component_NodesEntry {
    return Component_NodesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Component_NodesEntry>, I>>(object: I): Component_NodesEntry {
    const message = createBaseComponent_NodesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null) ? Node.fromPartial(object.value) : undefined;
    return message;
  },
};

function createBaseInterface(): Interface {
  return { typeParams: [], io: undefined };
}

export const Interface = {
  encode(message: Interface, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.typeParams) {
      TypeParam.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.io !== undefined) {
      IO.encode(message.io, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Interface {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInterface();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.typeParams.push(TypeParam.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.io = IO.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Interface {
    return {
      typeParams: globalThis.Array.isArray(object?.typeParams)
        ? object.typeParams.map((e: any) => TypeParam.fromJSON(e))
        : [],
      io: isSet(object.io) ? IO.fromJSON(object.io) : undefined,
    };
  },

  toJSON(message: Interface): unknown {
    const obj: any = {};
    if (message.typeParams?.length) {
      obj.typeParams = message.typeParams.map((e) => TypeParam.toJSON(e));
    }
    if (message.io !== undefined) {
      obj.io = IO.toJSON(message.io);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Interface>, I>>(base?: I): Interface {
    return Interface.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Interface>, I>>(object: I): Interface {
    const message = createBaseInterface();
    message.typeParams = object.typeParams?.map((e) => TypeParam.fromPartial(e)) || [];
    message.io = (object.io !== undefined && object.io !== null) ? IO.fromPartial(object.io) : undefined;
    return message;
  },
};

function createBaseTypeParam(): TypeParam {
  return { name: "", constr: "" };
}

export const TypeParam = {
  encode(message: TypeParam, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.constr !== "") {
      writer.uint32(18).string(message.constr);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TypeParam {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTypeParam();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.constr = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TypeParam {
    return {
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      constr: isSet(object.constr) ? globalThis.String(object.constr) : "",
    };
  },

  toJSON(message: TypeParam): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.constr !== "") {
      obj.constr = message.constr;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TypeParam>, I>>(base?: I): TypeParam {
    return TypeParam.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TypeParam>, I>>(object: I): TypeParam {
    const message = createBaseTypeParam();
    message.name = object.name ?? "";
    message.constr = object.constr ?? "";
    return message;
  },
};

function createBaseIO(): IO {
  return { ins: [], outs: [] };
}

export const IO = {
  encode(message: IO, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.ins) {
      Port.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.outs) {
      Port.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IO {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIO();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.ins.push(Port.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.outs.push(Port.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IO {
    return {
      ins: globalThis.Array.isArray(object?.ins) ? object.ins.map((e: any) => Port.fromJSON(e)) : [],
      outs: globalThis.Array.isArray(object?.outs) ? object.outs.map((e: any) => Port.fromJSON(e)) : [],
    };
  },

  toJSON(message: IO): unknown {
    const obj: any = {};
    if (message.ins?.length) {
      obj.ins = message.ins.map((e) => Port.toJSON(e));
    }
    if (message.outs?.length) {
      obj.outs = message.outs.map((e) => Port.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<IO>, I>>(base?: I): IO {
    return IO.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<IO>, I>>(object: I): IO {
    const message = createBaseIO();
    message.ins = object.ins?.map((e) => Port.fromPartial(e)) || [];
    message.outs = object.outs?.map((e) => Port.fromPartial(e)) || [];
    return message;
  },
};

function createBasePort(): Port {
  return { typeExpr: "", isArray: false };
}

export const Port = {
  encode(message: Port, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.typeExpr !== "") {
      writer.uint32(10).string(message.typeExpr);
    }
    if (message.isArray === true) {
      writer.uint32(16).bool(message.isArray);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Port {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePort();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.typeExpr = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.isArray = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Port {
    return {
      typeExpr: isSet(object.typeExpr) ? globalThis.String(object.typeExpr) : "",
      isArray: isSet(object.isArray) ? globalThis.Boolean(object.isArray) : false,
    };
  },

  toJSON(message: Port): unknown {
    const obj: any = {};
    if (message.typeExpr !== "") {
      obj.typeExpr = message.typeExpr;
    }
    if (message.isArray === true) {
      obj.isArray = message.isArray;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Port>, I>>(base?: I): Port {
    return Port.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Port>, I>>(object: I): Port {
    const message = createBasePort();
    message.typeExpr = object.typeExpr ?? "";
    message.isArray = object.isArray ?? false;
    return message;
  },
};

function createBaseNode(): Node {
  return { entityRef: undefined, typeArgs: [], componentDis: {} };
}

export const Node = {
  encode(message: Node, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entityRef !== undefined) {
      EntityRef.encode(message.entityRef, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.typeArgs) {
      writer.uint32(18).string(v!);
    }
    Object.entries(message.componentDis).forEach(([key, value]) => {
      Node_ComponentDisEntry.encode({ key: key as any, value }, writer.uint32(26).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Node {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.entityRef = EntityRef.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.typeArgs.push(reader.string());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          const entry3 = Node_ComponentDisEntry.decode(reader, reader.uint32());
          if (entry3.value !== undefined) {
            message.componentDis[entry3.key] = entry3.value;
          }
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Node {
    return {
      entityRef: isSet(object.entityRef) ? EntityRef.fromJSON(object.entityRef) : undefined,
      typeArgs: globalThis.Array.isArray(object?.typeArgs) ? object.typeArgs.map((e: any) => globalThis.String(e)) : [],
      componentDis: isObject(object.componentDis)
        ? Object.entries(object.componentDis).reduce<{ [key: string]: Node }>((acc, [key, value]) => {
          acc[key] = Node.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Node): unknown {
    const obj: any = {};
    if (message.entityRef !== undefined) {
      obj.entityRef = EntityRef.toJSON(message.entityRef);
    }
    if (message.typeArgs?.length) {
      obj.typeArgs = message.typeArgs;
    }
    if (message.componentDis) {
      const entries = Object.entries(message.componentDis);
      if (entries.length > 0) {
        obj.componentDis = {};
        entries.forEach(([k, v]) => {
          obj.componentDis[k] = Node.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Node>, I>>(base?: I): Node {
    return Node.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Node>, I>>(object: I): Node {
    const message = createBaseNode();
    message.entityRef = (object.entityRef !== undefined && object.entityRef !== null)
      ? EntityRef.fromPartial(object.entityRef)
      : undefined;
    message.typeArgs = object.typeArgs?.map((e) => e) || [];
    message.componentDis = Object.entries(object.componentDis ?? {}).reduce<{ [key: string]: Node }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = Node.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseNode_ComponentDisEntry(): Node_ComponentDisEntry {
  return { key: "", value: undefined };
}

export const Node_ComponentDisEntry = {
  encode(message: Node_ComponentDisEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Node.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Node_ComponentDisEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNode_ComponentDisEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.value = Node.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Node_ComponentDisEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? Node.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Node_ComponentDisEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = Node.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Node_ComponentDisEntry>, I>>(base?: I): Node_ComponentDisEntry {
    return Node_ComponentDisEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Node_ComponentDisEntry>, I>>(object: I): Node_ComponentDisEntry {
    const message = createBaseNode_ComponentDisEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null) ? Node.fromPartial(object.value) : undefined;
    return message;
  },
};

function createBaseEntityRef(): EntityRef {
  return { pkg: "", name: "" };
}

export const EntityRef = {
  encode(message: EntityRef, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pkg !== "") {
      writer.uint32(10).string(message.pkg);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EntityRef {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEntityRef();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.pkg = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EntityRef {
    return {
      pkg: isSet(object.pkg) ? globalThis.String(object.pkg) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
    };
  },

  toJSON(message: EntityRef): unknown {
    const obj: any = {};
    if (message.pkg !== "") {
      obj.pkg = message.pkg;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EntityRef>, I>>(base?: I): EntityRef {
    return EntityRef.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EntityRef>, I>>(object: I): EntityRef {
    const message = createBaseEntityRef();
    message.pkg = object.pkg ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseConnection(): Connection {
  return { senderSide: undefined, receiverSides: [] };
}

export const Connection = {
  encode(message: Connection, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.senderSide !== undefined) {
      SenderConnectionSide.encode(message.senderSide, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.receiverSides) {
      ReceiverConnectionSide.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Connection {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConnection();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.senderSide = SenderConnectionSide.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.receiverSides.push(ReceiverConnectionSide.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Connection {
    return {
      senderSide: isSet(object.senderSide) ? SenderConnectionSide.fromJSON(object.senderSide) : undefined,
      receiverSides: globalThis.Array.isArray(object?.receiverSides)
        ? object.receiverSides.map((e: any) => ReceiverConnectionSide.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Connection): unknown {
    const obj: any = {};
    if (message.senderSide !== undefined) {
      obj.senderSide = SenderConnectionSide.toJSON(message.senderSide);
    }
    if (message.receiverSides?.length) {
      obj.receiverSides = message.receiverSides.map((e) => ReceiverConnectionSide.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Connection>, I>>(base?: I): Connection {
    return Connection.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Connection>, I>>(object: I): Connection {
    const message = createBaseConnection();
    message.senderSide = (object.senderSide !== undefined && object.senderSide !== null)
      ? SenderConnectionSide.fromPartial(object.senderSide)
      : undefined;
    message.receiverSides = object.receiverSides?.map((e) => ReceiverConnectionSide.fromPartial(e)) || [];
    return message;
  },
};

function createBaseReceiverConnectionSide(): ReceiverConnectionSide {
  return { portAddr: undefined, selectors: [] };
}

export const ReceiverConnectionSide = {
  encode(message: ReceiverConnectionSide, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.portAddr !== undefined) {
      PortAddr.encode(message.portAddr, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.selectors) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReceiverConnectionSide {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReceiverConnectionSide();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.portAddr = PortAddr.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.selectors.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReceiverConnectionSide {
    return {
      portAddr: isSet(object.portAddr) ? PortAddr.fromJSON(object.portAddr) : undefined,
      selectors: globalThis.Array.isArray(object?.selectors)
        ? object.selectors.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: ReceiverConnectionSide): unknown {
    const obj: any = {};
    if (message.portAddr !== undefined) {
      obj.portAddr = PortAddr.toJSON(message.portAddr);
    }
    if (message.selectors?.length) {
      obj.selectors = message.selectors;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ReceiverConnectionSide>, I>>(base?: I): ReceiverConnectionSide {
    return ReceiverConnectionSide.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ReceiverConnectionSide>, I>>(object: I): ReceiverConnectionSide {
    const message = createBaseReceiverConnectionSide();
    message.portAddr = (object.portAddr !== undefined && object.portAddr !== null)
      ? PortAddr.fromPartial(object.portAddr)
      : undefined;
    message.selectors = object.selectors?.map((e) => e) || [];
    return message;
  },
};

function createBaseSenderConnectionSide(): SenderConnectionSide {
  return { constRef: undefined, portAddr: undefined, selectors: [] };
}

export const SenderConnectionSide = {
  encode(message: SenderConnectionSide, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.constRef !== undefined) {
      EntityRef.encode(message.constRef, writer.uint32(10).fork()).ldelim();
    }
    if (message.portAddr !== undefined) {
      PortAddr.encode(message.portAddr, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.selectors) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SenderConnectionSide {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSenderConnectionSide();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.constRef = EntityRef.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.portAddr = PortAddr.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.selectors.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SenderConnectionSide {
    return {
      constRef: isSet(object.constRef) ? EntityRef.fromJSON(object.constRef) : undefined,
      portAddr: isSet(object.portAddr) ? PortAddr.fromJSON(object.portAddr) : undefined,
      selectors: globalThis.Array.isArray(object?.selectors)
        ? object.selectors.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: SenderConnectionSide): unknown {
    const obj: any = {};
    if (message.constRef !== undefined) {
      obj.constRef = EntityRef.toJSON(message.constRef);
    }
    if (message.portAddr !== undefined) {
      obj.portAddr = PortAddr.toJSON(message.portAddr);
    }
    if (message.selectors?.length) {
      obj.selectors = message.selectors;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SenderConnectionSide>, I>>(base?: I): SenderConnectionSide {
    return SenderConnectionSide.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SenderConnectionSide>, I>>(object: I): SenderConnectionSide {
    const message = createBaseSenderConnectionSide();
    message.constRef = (object.constRef !== undefined && object.constRef !== null)
      ? EntityRef.fromPartial(object.constRef)
      : undefined;
    message.portAddr = (object.portAddr !== undefined && object.portAddr !== null)
      ? PortAddr.fromPartial(object.portAddr)
      : undefined;
    message.selectors = object.selectors?.map((e) => e) || [];
    return message;
  },
};

function createBasePortAddr(): PortAddr {
  return { node: "", port: "", idx: 0 };
}

export const PortAddr = {
  encode(message: PortAddr, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.node !== "") {
      writer.uint32(10).string(message.node);
    }
    if (message.port !== "") {
      writer.uint32(18).string(message.port);
    }
    if (message.idx !== 0) {
      writer.uint32(24).int32(message.idx);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PortAddr {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePortAddr();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.node = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.port = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.idx = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PortAddr {
    return {
      node: isSet(object.node) ? globalThis.String(object.node) : "",
      port: isSet(object.port) ? globalThis.String(object.port) : "",
      idx: isSet(object.idx) ? globalThis.Number(object.idx) : 0,
    };
  },

  toJSON(message: PortAddr): unknown {
    const obj: any = {};
    if (message.node !== "") {
      obj.node = message.node;
    }
    if (message.port !== "") {
      obj.port = message.port;
    }
    if (message.idx !== 0) {
      obj.idx = Math.round(message.idx);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PortAddr>, I>>(base?: I): PortAddr {
    return PortAddr.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PortAddr>, I>>(object: I): PortAddr {
    const message = createBasePortAddr();
    message.node = object.node ?? "";
    message.port = object.port ?? "";
    message.idx = object.idx ?? 0;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(globalThis.Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
