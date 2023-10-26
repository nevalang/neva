# Requirements

- Go: https://go.dev/doc/install
- Make: https://www.gnu.org/software/make/#download
- NodeJS and NPM: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm/
- Antlr: `pip install antlr4-tools`
- Protoc: https://grpc.io/docs/protoc-installation/
- Tygo: `go install github.com/gzuidhof/tygo@latest`

## VSCode extensions

These are not really required but recommended in order you're using VSCode

- https://marketplace.visualstudio.com/items?itemName=nevalang.vscode-nevalang
- https://marketplace.visualstudio.com/items?itemName=mike-lischke.vscode-antlr4
- https://marketplace.visualstudio.com/items?itemName=pedro-w.tmlanguage
- https://marketplace.visualstudio.com/items?itemName=tooltitudeteam.tooltitude
- https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid
- https://marketplace.visualstudio.com/items?itemName=yoavbls.pretty-ts-errors
- https://marketplace.visualstudio.com/items?itemName=Plex.vscode-protolint

# Development

See [ARCHITECTURE.md](./ARCHITECTURE.md) and [Makefile](./Makefile). Also many directories contain `README.md` files.

## VSCode Extension

Check out [tygo.yaml](./tygo.yaml). It depends on types defined in the `src` and `ts` packages and thus it's dangerous to rename those types. If you gonna do so make sure you don't brake TS types generation.

# Naming conventions

## Tests

Use `_` instead of space in for test-case names because go turns spaces into underscores and makes it hard to find specific case.

# Learning Resources

## Common

- [The Bible from the FBP father mr. J.Paul Morrison](https://jpaulmorrison.com/fbp/1stedchaps.html)
- [Dataflow and Reactive Programming Systems: A Practical Guide](https://www.amazon.com/Dataflow-Reactive-Programming-Systems-Practical/dp/1497422442)
- [Bob Martin's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Grokking Algorithms: An Illustrated Guide for Programmers and Other Curious People](https://www.amazon.com/Grokking-Algorithms-illustrated-programmers-curious/dp/1617292230)
- [Code: The Hidden Language of Computer Hardware and Software](https://www.amazon.com/Code-Language-Computer-Hardware-Software/dp/0735611319)

## Golang

### Must Read

- [How To Write Go Code](https://go.dev/doc/code)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Proverbs](https://go-proverbs.github.io/)
- [50 Shades Of Go](http://golang50shad.es/)
- [Darker Corners Of Go](https://rytisbiel.com/2021/03/06/darker-corners-of-go/)

### Highly Recommended

- [Concurrency is not parallelism](https://go.dev/blog/waza-talk)
- [Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Errors Are Values](https://go.dev/blog/errors-are-values)
- [Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)

### Nice To Know

- [The Go Programming Language (The Book)](https://www.amazon.com/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440)
- [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings)
- Concurrency
    - [Go Concurrency Patterns: Timing out, moving on](https://go.dev/blog/concurrency-timeouts)
    - [Go Concurrency Patterns: Context](https://go.dev/blog/context)
    - [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)

## JavaScript

- https://developer.mozilla.org/en-US/
- https://www.typescriptlang.org/
- https://react.dev/
- https://github.com/getify/You-Dont-Know-JS

## VSCode Extensions

- [Custom Editors](https://code.visualstudio.com/api/extension-guides/custom-editors)
- [Webviews](https://code.visualstudio.com/api/extension-guides/webview)
- [Language Servers (LSP)](https://code.visualstudio.com/api/language-extensions/language-server-extension-guide)
- [Syntax Highlighter](https://code.visualstudio.com/api/language-extensions/syntax-highlight-guide)
- [LSP Overview](https://microsoft.github.io/language-server-protocol/)
- [Go library for LSP implementation](https://github.com/tliron/glsp)

## Video

- ["Stop Writing Dead Programs" by Jack Rusher](https://youtu.be/8Ab3ArE8W3s?feature=shared)
- ["The Mess We're In" by Joe Armstrong](https://youtu.be/lKXe3HUG2l4?feature=shared)
- ["Propositions as Types" by Philip Wadler](https://youtu.be/IOiZatlZtGU?feature=shared)
- ["Outperforming Imperative with Pure Functional Languages" by Richard Feldman](https://youtu.be/vzfy4EKwG_Y?feature=shared)
- ["What Is a Strange Loop and What is it Like To Be One?" by Douglas Hofstadter (2013)](https://youtu.be/UT5CxsyKwxg?feature=shared)

# Community

Here you can find help

- [Flow-Based Programming Discord](https://discord.gg/JHWRuZQJ)
- [r/ProgrammingLanguages](https://www.reddit.com/r/ProgrammingLanguages/)
