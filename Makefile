.PHONY: runtimesdk
runtimesdk:
	protoc api/runtimesdk.proto \
		--go_out=pkg/runtimesdksdk

.PHONY: devserversdk
devserversdk:
	protoc api/devserver.proto \
		--js_out=import_style=commonjs,binary:web/src/sdk \
		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:web/src/sdk \
		--go_out=pkg/devserversdk \
		--go-grpc_out=pkg/devserversdk

.PHONY: plugins
plugins:
	rm -rf plugins/*
	go build -o plugins/print.so -buildmode=plugin -gcflags="all=-N -l" internal/runtime/operators/io/print.go
	go build -o plugins/lock.so -buildmode=plugin -gcflags="all=-N -l" internal/runtime/operators/flow/lock.go