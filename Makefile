# === Development ===

# build neva cli for host OS and put to the PATH with `go install`
.PHONY: install
install:
	@go install ./cmd/neva

# generate go parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# fix struct field ordering for optimal padding/pointer data
.PHONY: align
betteralign-fix:
	betteralign -fix ./...

# apply gofix rewrites across the repo
.PHONY: gofix
gofix:
	go fix ./...

# check that gofix produces no diff (CI-friendly)
.PHONY: gofix-check
gofix-check:
	go fix ./...
	git diff --exit-code

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0 run ./...

.PHONY: test-unit
test-unit:
	go list ./... \
		| grep -Ev '^github.com/nevalang/neva/(e2e|examples)(/|$$)' \
		| xargs -r go test -v

.PHONY: vulncheck
vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# check potential nil-derefs
.PHONY: nilaway
nilaway:
	nilaway \
		-exclude-test-files \
		-include-pkgs="github.com/nevalang/neva/internal" \
		./...

.PHONY: bench-runtime
bench-runtime:
	go test ./benchmarks -run=^$$ -bench BenchmarkRuntimeE2E -benchtime=1x -count=1

.PHONY: bench-runtime-routing-complex
bench-runtime-routing-complex:
	go test ./benchmarks -run=^$$ -bench '^BenchmarkRuntimeE2E/(simple_routers|simple_selectors|complex_)' -benchtime=1x -count=1

# === Release Build ===

# Release flags:
# -trimpath removes host paths from debug metadata (smaller, reproducible artifacts)
# -buildvcs=false skips VCS stamping to keep output deterministic and compact
# build neva cli for all target platforms
.PHONY: build
build:
	$(MAKE) build-darwin-amd64
	$(MAKE) build-darwin-arm64
	$(MAKE) build-linux-amd64
	$(MAKE) build-linux-arm64
	$(MAKE) build-linux-loong64
	$(MAKE) build-windows-amd64
	$(MAKE) build-windows-arm64

# build neva cli for amd64 mac
.PHONY: build-darwin-amd64
build-darwin-amd64:
	@GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-darwin-amd64 ./cmd/neva

# build neva cli for arm64 mac
.PHONY: build-darwin-arm64
build-darwin-arm64:
	@GOOS=darwin GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-darwin-arm64 ./cmd/neva

# build neva cli for amd64 linux
.PHONY: build-linux-amd64
build-linux-amd64:
	@GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-linux-amd64 ./cmd/neva

# build neva cli for arm64 linux
.PHONY: build-linux-arm64
build-linux-arm64:
	@GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-linux-arm64 ./cmd/neva

# build neva cli for loong64 linux
.PHONY: build-linux-loong64
build-linux-loong64:
	@GOOS=linux GOARCH=loong64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-linux-loong64 ./cmd/neva

# build neva cli for amd64 windows
.PHONY: build-windows-amd64
build-windows-amd64:
	@GOOS=windows GOARCH=amd64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-windows-amd64.exe ./cmd/neva

# build neva cli for arm64 windows
.PHONY: build-windows-arm64
build-windows-arm64:
	@GOOS=windows GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags="-s -w" -o neva-windows-arm64.exe ./cmd/neva
