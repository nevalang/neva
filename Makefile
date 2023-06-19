# .PHONY: devserversdk
# devserversdk:
# 	protoc api/devserver.proto \
# 		--js_out=import_style=commonjs,binary:web/src/sdk \
# 		--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:web/src/sdk \
# 		--go_out=pkg/devserversdk \
# 		--go-grpc_out=pkg/devserversdk

.PHONY: antlr
antlr:
	cd internal/compiler/frontend && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated