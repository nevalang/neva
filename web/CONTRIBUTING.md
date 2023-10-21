# Contributing

Extension consist of 3 parts:

- `syntaxes` - syntax highlighting (textmate grammars)
  - `neva.tmLanguage.yml` - source code grammar
  - `neva.tmLanguage.json` - build artefact, generated from yaml
  - `tests` - \*.neva files to test syntax highlighting
  - `yaml2json` and `watch-yaml2json` commands in `package.json`
  - `language-configuration.json` - out of the box vscode's highlighting for simple things like brackets, etc
- `src` - vscode extension itself
  - `extension.ts` - high-level vscode api usage, registers custom editor and commands
  - `editor.ts` - glue between vscode api and webview, custom text editor interface implementation
  - `api/proto/src.ts` - generated protobuf sdk for transmitting source code between vscode extension and golang language server
- `webview` webview-ui (react app)
  - `dist` - build artefact, used by extension

## Development

You cannot develop extension and webview at the same time

### Extension

Webview must be built (`webview/dist` directory must contain generated static files)

```bash
npm watch-yaml2json # only if you gonna make changes to grammar
npm watch-ts # only if you gonna make changes to files in ./src folder
```

Then **start debugger** (`Launch process` task in `./vscode/launch.json`)

### Webview

Webview itself does not know anyting about extension, it's just a react app

```bash
npm start # vite devserver
```

## Production

```bash
cd webview && npm run build # build webview
cd ../ && npm run build # build extension
vsce package # create package for marketplace
```

## FAQ

- [Why not WASM?](https://github.com/nevalang/neva/discussions/374#discussioncomment-7345045)