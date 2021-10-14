.PHONY: sdk
sdk:
	make _go_sdk && make _ts_sdk

.PHONY: go_sdk
go_sdk:
	docker run --rm \
		-v ${PWD}:/app openapitools/openapi-generator-cli generate \
		-i /app/api/api.yml \
		-g go-server  \
		-o /app/generated_go_sdk \
		--additional-properties=isGoSubmodule=false,packageName=sdk,featureCORS=true && \
		
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
