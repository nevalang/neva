.PHONY: devserversdk
devserversdk:
	protoc api/devserver.proto \
		--js_out=import_style=commonjs,binary:web/src/sdk \
		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:web/src/sdk \
		--go_out=pkg/devserversdk \
		--go-grpc_out=pkg/devserversdk

.PHONY: vulncheck
vulncheck:
	govulncheck ./...

.PHONY: race
race:
	go test -race ./...