.PHONY: sdk
sdk:
	make go_sdk && make ts_sdk

.PHONY: go_sdk
go_sdk:
	docker run --rm \
		-v ${PWD}:/app openapitools/openapi-generator-cli generate \
		-i /app/api/api.yml \
		-g go-server  \
		-o /app/generated_go_sdk \
		--additional-properties=isGoSubmodule=false,packageName=sdk,featureCORS=true
		
# rm -rf pkg/sdk && \
# mv generated_go_sdk/go pkg/sdk && \
# rm -rf generated_go_sdk

.PHONY: ts_sdk
ts_sdk:
	docker run --user $(id -u):$(id -g) --rm \
		-v ${PWD}:/app openapitools/openapi-generator-cli generate \
		-i /app/api/api.yml \
		-o /app/generated_ts_sdk \
		-g typescript-axios \
		--additional-properties=supportsES6=true

# rm -rf web/src/sdk 
# mv generated_ts_sdk  web/src/sdk 
# sudo rm -rf generated_ts_sdk

.PHONY: go_plugins
go_plugins:
	go version
	go build -o plugins/and.so -buildmode=plugin internal/runtime/operators/and/main.go
	go build -o plugins/filter.so -buildmode=plugin internal/runtime/operators/filter/main.go
	go build -o plugins/more.so -buildmode=plugin internal/runtime/operators/more/main.go
	go build -o plugins/mul.so -buildmode=plugin internal/runtime/operators/mul/main.go
	go build -o plugins/or.so -buildmode=plugin internal/runtime/operators/or/main.go
	go build -o plugins/remainder.so -buildmode=plugin internal/runtime/operators/remainder/main.go