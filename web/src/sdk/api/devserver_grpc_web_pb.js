/**
 * @fileoverview gRPC-Web generated client stub for devserver
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.devserver = require('./devserver_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.devserver.DevServerClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.devserver.DevServerPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.devserver.ListProgramsRequest,
 *   !proto.devserver.ListProgramsResponse>}
 */
const methodDescriptor_DevServer_ListPrograms = new grpc.web.MethodDescriptor(
  '/devserver.DevServer/ListPrograms',
  grpc.web.MethodType.UNARY,
  proto.devserver.ListProgramsRequest,
  proto.devserver.ListProgramsResponse,
  /**
   * @param {!proto.devserver.ListProgramsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.ListProgramsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.devserver.ListProgramsRequest,
 *   !proto.devserver.ListProgramsResponse>}
 */
const methodInfo_DevServer_ListPrograms = new grpc.web.AbstractClientBase.MethodInfo(
  proto.devserver.ListProgramsResponse,
  /**
   * @param {!proto.devserver.ListProgramsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.ListProgramsResponse.deserializeBinary
);


/**
 * @param {!proto.devserver.ListProgramsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.devserver.ListProgramsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.ListProgramsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerClient.prototype.listPrograms =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/devserver.DevServer/ListPrograms',
      request,
      metadata || {},
      methodDescriptor_DevServer_ListPrograms,
      callback);
};


/**
 * @param {!proto.devserver.ListProgramsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.devserver.ListProgramsResponse>}
 *     A native promise that resolves to the response
 */
proto.devserver.DevServerPromiseClient.prototype.listPrograms =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/devserver.DevServer/ListPrograms',
      request,
      metadata || {},
      methodDescriptor_DevServer_ListPrograms);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.devserver.GetProgramRequest,
 *   !proto.devserver.GetProgramResponse>}
 */
const methodDescriptor_DevServer_GetProgram = new grpc.web.MethodDescriptor(
  '/devserver.DevServer/GetProgram',
  grpc.web.MethodType.UNARY,
  proto.devserver.GetProgramRequest,
  proto.devserver.GetProgramResponse,
  /**
   * @param {!proto.devserver.GetProgramRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.GetProgramResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.devserver.GetProgramRequest,
 *   !proto.devserver.GetProgramResponse>}
 */
const methodInfo_DevServer_GetProgram = new grpc.web.AbstractClientBase.MethodInfo(
  proto.devserver.GetProgramResponse,
  /**
   * @param {!proto.devserver.GetProgramRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.GetProgramResponse.deserializeBinary
);


/**
 * @param {!proto.devserver.GetProgramRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.devserver.GetProgramResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.GetProgramResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerClient.prototype.getProgram =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/devserver.DevServer/GetProgram',
      request,
      metadata || {},
      methodDescriptor_DevServer_GetProgram,
      callback);
};


/**
 * @param {!proto.devserver.GetProgramRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.devserver.GetProgramResponse>}
 *     A native promise that resolves to the response
 */
proto.devserver.DevServerPromiseClient.prototype.getProgram =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/devserver.DevServer/GetProgram',
      request,
      metadata || {},
      methodDescriptor_DevServer_GetProgram);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.devserver.UpdateProgramRequest,
 *   !proto.devserver.UpdateProgramResponse>}
 */
const methodDescriptor_DevServer_UpdateProgram = new grpc.web.MethodDescriptor(
  '/devserver.DevServer/UpdateProgram',
  grpc.web.MethodType.UNARY,
  proto.devserver.UpdateProgramRequest,
  proto.devserver.UpdateProgramResponse,
  /**
   * @param {!proto.devserver.UpdateProgramRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.UpdateProgramResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.devserver.UpdateProgramRequest,
 *   !proto.devserver.UpdateProgramResponse>}
 */
const methodInfo_DevServer_UpdateProgram = new grpc.web.AbstractClientBase.MethodInfo(
  proto.devserver.UpdateProgramResponse,
  /**
   * @param {!proto.devserver.UpdateProgramRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.UpdateProgramResponse.deserializeBinary
);


/**
 * @param {!proto.devserver.UpdateProgramRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.devserver.UpdateProgramResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.UpdateProgramResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerClient.prototype.updateProgram =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/devserver.DevServer/UpdateProgram',
      request,
      metadata || {},
      methodDescriptor_DevServer_UpdateProgram,
      callback);
};


/**
 * @param {!proto.devserver.UpdateProgramRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.devserver.UpdateProgramResponse>}
 *     A native promise that resolves to the response
 */
proto.devserver.DevServerPromiseClient.prototype.updateProgram =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/devserver.DevServer/UpdateProgram',
      request,
      metadata || {},
      methodDescriptor_DevServer_UpdateProgram);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.devserver.StartDebugRequest,
 *   !proto.devserver.StartDebugResponse>}
 */
const methodDescriptor_DevServer_StartDebugger = new grpc.web.MethodDescriptor(
  '/devserver.DevServer/StartDebugger',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.devserver.StartDebugRequest,
  proto.devserver.StartDebugResponse,
  /**
   * @param {!proto.devserver.StartDebugRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.StartDebugResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.devserver.StartDebugRequest,
 *   !proto.devserver.StartDebugResponse>}
 */
const methodInfo_DevServer_StartDebugger = new grpc.web.AbstractClientBase.MethodInfo(
  proto.devserver.StartDebugResponse,
  /**
   * @param {!proto.devserver.StartDebugRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.StartDebugResponse.deserializeBinary
);


/**
 * @param {!proto.devserver.StartDebugRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.StartDebugResponse>}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerClient.prototype.startDebugger =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/devserver.DevServer/StartDebugger',
      request,
      metadata || {},
      methodDescriptor_DevServer_StartDebugger);
};


/**
 * @param {!proto.devserver.StartDebugRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.StartDebugResponse>}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerPromiseClient.prototype.startDebugger =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/devserver.DevServer/StartDebugger',
      request,
      metadata || {},
      methodDescriptor_DevServer_StartDebugger);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.devserver.DebugRequest,
 *   !proto.devserver.DebugResponse>}
 */
const methodDescriptor_DevServer_SendDebugMessage = new grpc.web.MethodDescriptor(
  '/devserver.DevServer/SendDebugMessage',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.devserver.DebugRequest,
  proto.devserver.DebugResponse,
  /**
   * @param {!proto.devserver.DebugRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.DebugResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.devserver.DebugRequest,
 *   !proto.devserver.DebugResponse>}
 */
const methodInfo_DevServer_SendDebugMessage = new grpc.web.AbstractClientBase.MethodInfo(
  proto.devserver.DebugResponse,
  /**
   * @param {!proto.devserver.DebugRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.devserver.DebugResponse.deserializeBinary
);


/**
 * @param {!proto.devserver.DebugRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.DebugResponse>}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerClient.prototype.sendDebugMessage =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/devserver.DevServer/SendDebugMessage',
      request,
      metadata || {},
      methodDescriptor_DevServer_SendDebugMessage);
};


/**
 * @param {!proto.devserver.DebugRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.devserver.DebugResponse>}
 *     The XHR Node Readable Stream
 */
proto.devserver.DevServerPromiseClient.prototype.sendDebugMessage =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/devserver.DevServer/SendDebugMessage',
      request,
      metadata || {},
      methodDescriptor_DevServer_SendDebugMessage);
};


module.exports = proto.devserver;

