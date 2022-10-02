.PHONY: runtimesdk
runtimesdk:
	protoc api/runtime.proto \
		--go_out=pkg/runtimesdk

.PHONY: devserversdk
devserversdk:
	protoc api/devserver.proto \
		--js_out=import_style=commonjs,binary:web/src/sdk \
		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:web/src/sdk \
		--go_out=pkg/devserversdk \
		--go-grpc_out=pkg/devserversdk

.PHONY: debugplugins
debugplugins:
	rm -rf plugins/*
	go build -o plugins/io.so -buildmode=plugin -gcflags="all=-N -l" internal/operators/io/io.go