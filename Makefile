# === Development ===
# build neva cli for host OS and put to the PATH
.PHONY: install
# install neva cli to the PATH
user := $(shell whoami)
root := $("root")
install-neva:
    ifeq ($(user), root)
		${shell GOBIN=/usr/local/bin go install -ldflags="-s -w" `pwd`/cmd/neva}
    else
		${shell go install -ldflags="-s -w" `pwd`/cmd/neva}
	endef
    endif
# install neva lsp to the PATH
install-lsp:
	@${shell ln -s `pwd`/cmd/lsp/ `pwd`/cmd/neva-lsp}
# use soft link to change the executable file name
	@${shell cd cmd/neva-lsp}
    ifeq ($(user), root)
		${shell GOBIN=/usr/local/bin go install -ldflags="-s -w" `pwd`/cmd/neva-lsp}
    else
		${shell go install -ldflags="-s -w" `pwd`/cmd/neva-lsp}
	endef
    endif
# cleanup
	@${shell rm -r cmd/neva-lsp}
install:
	$(MAKE) install-neva
	$(MAKE) install-lsp
.PHONY: uninstall
uninstall:
	$(MAKE) uninstall-neva
	$(MAKE) uninstall-lsp
UNINSTALL_PATH ?= $(or $(shell go env GOPATH)/bin,$(GOBIN))
uninstall-neva:
    ifeq ($(wildcard $(UNINSTALL_PATH)/neva), $(UNINSTALL_PATH)/neva)
		${shell rm -f "$(UNINSTALL_PATH)/neva"}
		${shell echo "Uninstalled neva from $(UNINSTALL_PATH)"}
    else
		${shell echo "neva cli was not installed or not found in $(UNINSTALL_PATH)."; \}
	endef
    endif
uninstall-lsp:
    ifeq ($(wildcard $(UNINSTALL_PATH)/neva-lsp), , $(UNINSTALL_PATH)/neva-lsp)
		${shell rm -f "$(UNINSTALL_PATH)/neva-lsp"}
		${shell echo "Uninstalled neva-lsp from $(UNINSTALL_PATH)"}
    else
		${shell echo "neva language server was not installed or not found in $(UNINSTALL_PATH)."; \}
	endef
    endif
# generate go parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated
# make clean
.PHONY: clean
clean:
	- rm neva neva-*

# clean install
.install:
	$(MAKE) clean
	$(MAKE) uninstall
	$(MAKE) install
# generate ts types from go src pkg
.PHONY: tygo
tygo:
	@tygo generate

# === Release Artifacts ===

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
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o neva-darwin-amd64 ./cmd/neva

# build neva cli for arm64 mac
.PHONY: build-darwin-arm64
build-darwin-arm64:
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o neva-darwin-arm64 ./cmd/neva

# build neva cli for amd64 linux
.PHONY: build-linux-amd64
build-linux-amd64:
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o neva-linux-amd64 ./cmd/neva

# build neva cli for arm64 linux
.PHONY: build-linux-arm64
build-linux-arm64:
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o neva-linux-arm64 ./cmd/neva

# build neva cli for loong64 linux
.PHONY: build-linux-loong64
build-linux-loong64:
	@GOOS=linux GOARCH=loong64 go build -ldflags="-s -w" -o neva-linux-loong64 ./cmd/neva

# build neva cli for amd64 windows
.PHONY: build-windows-amd64
build-windows-amd64:
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o neva-windows-amd64.exe ./cmd/neva

# build neva cli for arm64 windows
.PHONY: build-windows-arm64
build-windows-arm64:
	@GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o neva-windows-arm64.exe ./cmd/neva

# === Tool Artifacts ===

# build neva lsp for all target platforms
.PHONY: build-lsp
build-lsp:
	$(MAKE) build-lsp-darwin-amd64
	$(MAKE) build-lsp-darwin-arm64
	$(MAKE) build-lsp-linux-amd64
	$(MAKE) build-lsp-linux-arm64
	$(MAKE) build-lsp-linux-loong64
	$(MAKE) build-lsp-windows-amd64
	$(MAKE) build-lsp-windows-arm64

# build neva lsp for amd64 mac
.PHONY: build-lsp-darwin-amd64
build-lsp-darwin-amd64:
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o neva-lsp-darwin-amd64 ./cmd/lsp

# build neva lsp for arm64 mac
.PHONY: build-lsp-darwin-arm64
build-lsp-darwin-arm64:
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o neva-lsp-darwin-arm64 ./cmd/lsp

# build neva lsp for amd64 linux
.PHONY: build-lsp-linux-amd64
build-lsp-linux-amd64:
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o neva-lsp-linux-amd64 ./cmd/lsp

# build neva lsp for arm64 linux
.PHONY: build-lsp-linux-arm64
build-lsp-linux-arm64:
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o neva-lsp-linux-arm64 ./cmd/lsp

# build neva lsp for loong64 linux
.PHONY: build-lsp-linux-loong64
build-lsp-linux-loong64:
	@GOOS=linux GOARCH=loong64 go build -ldflags="-s -w" -o neva-lsp-linux-loong64 ./cmd/lsp
# build neva lsp for amd64 windows

.PHONY: build-lsp-windows-amd64
build-lsp-windows-amd64:
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o neva-lsp-windows-amd64.exe ./cmd/lsp

# build neva lsp for arm64 windows
.PHONY: build-lsp-windows-arm64
build-lsp-windows-arm64:
	@GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o neva-lsp-windows-arm64.exe ./cmd/lsp
