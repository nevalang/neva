.PHONY: sdk
sdk:
	protoc api/devserver.proto \
		--js_out=import_style=commonjs,binary:web/src/sdk \
		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:web/src/sdk \
		--go_out=pkg/devserversdk \
		--go-grpc_out=pkg/devserversdk

.PHONY: go_plugins
go_plugins:
	go build -o plugins/and.so -buildmode=plugin internal/runtime/operators/and/main.go
	go build -o plugins/filter.so -buildmode=plugin internal/runtime/operators/filter/main.go
	go build -o plugins/select.so -buildmode=plugin internal/runtime/operators/select/main.go
	go build -o plugins/more.so -buildmode=plugin internal/runtime/operators/more/main.go
	go build -o plugins/mul.so -buildmode=plugin internal/runtime/operators/mul/main.go
	go build -o plugins/or.so -buildmode=plugin internal/runtime/operators/or/main.go
	go build -o plugins/remainder.so -buildmode=plugin internal/runtime/operators/remainder/main.go