# Contributing

> WARNING: Don't forget **it's JS**. _Sometimes_ you need to `rm -rf node_modules` and install deps from scratch and everything magically gonna work. Also if you don't see changes in webview while developing, try to kill `watch` task in terminal and re-rerun debuggin.

## Structure

This extension consist of 4 parts:

1. Language server. This is go application under `cmd/lsp/`
2. VSCode extension (bridge) in `web` (you are here)
3. Webview-UI (React-application) in `web/webview`
4. Syntax highlighting grammar under `web/syntaxes` directory

## Development

Simply `VSCode Extension` debug task. See [launch.json](../.vscode/launch.json) and [tasks.json](../.vscode/tasks.json) to figure out what's going on.

### Webview

See

## Production

```bash
npm run build # build textmate grammar, webview and vscode extension
vsce package # pack everything into a VSIX package
```

## FAQ

- [Why not use WASM for LSP integration?](https://github.com/nevalang/neva/discussions/374#discussioncomment-7345045)
