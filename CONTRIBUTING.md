# Requirements

## System

- Go: https://go.dev/doc/install
- Make: https://www.gnu.org/software/make/#download
- NodeJS and NPM: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm/
- Antlr: `pip install antlr4-tools`
- Protoc: https://grpc.io/docs/protoc-installation/
- Tygo: `go install github.com/gzuidhof/tygo@latest`

## VSCode

These are not really required but recommended in order you're using VSCode

- [nevalang](https://marketplace.visualstudio.com/items?itemName=nevalang.vscode-nevalang)
- [antlr4](https://marketplace.visualstudio.com/items?itemName=mike-lischke.vscode-antlr4)
- [tmlanguage](https://marketplace.visualstudio.com/items?itemName=pedro-w.tmlanguage)
- [tooltitude](https://marketplace.visualstudio.com/items?itemName=tooltitudeteam.tooltitude)
- [markdown-mermaid](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid)
- [ts-errors](https://marketplace.visualstudio.com/items?itemName=yoavbls.pretty-ts-errors)
- [vscode-protolint](https://marketplace.visualstudio.com/items?itemName=Plex.vscode-protolint)

# Development

First read the [language specification](docs/spec.md), [design document](./ARCHITECTURE.md) and see what can [Makefile](./Makefile) do.

Remember that many go packages contain [doc comments](https://tip.golang.org/doc/comment). Do not fear to read the source code, leave the comments for unintuitive parts.

Try to follow [clean code](https://github.com/Pungyeon/clean-go-article) and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) best practices in general. Follow [SOLID](https://en.wikipedia.org/wiki/SOLID).

Discuss your idea via github discussions or issues before implementing it. Write tests, avoid using `nolint`. Leave the comments (but only when you have to), update documentation.

## Github Issues

Issues must only be created for known bugs and understandable architecture issues. Not ideas, suggestions or feature requests. Discussions must be used for that instead.

## VSCode Extension

Check out [tygo.yaml](./tygo.yaml). It depends on types defined in the `src` and `typesystem` packages and thus it's dangerous to rename those types. If you gonna do so make sure you don't brake TS types generation. Check [web/CONTRIBUTING.md](./web/CONTRIBUTING.md).

## ANTLR Grammar

Don't forget to open `neva.g4` file before debugging with VSCode!

# Naming conventions

## Tests

Use `_` instead of space in for test-case names because go turns spaces into underscores and makes it hard to find specific case.

# Learning Resources

## FBP/DataFlow

- [Elements of Dataflow and Reactive Programming Systems](https://youtu.be/iFlT93wakVo?feature=shared)
- [The origins of Flow Based Programming with J Paul Morrison](https://youtu.be/up2yhNTsaDs?feature=shared)
- [Dataflow and Reactive Programming Systems: A Practical Guide](https://www.amazon.com/Dataflow-Reactive-Programming-Systems-Practical/dp/1497422442)
- [Flow-Based Programming: A New Approach to Application Development](https://jpaulmorrison.com/fbp/1stedchaps.html)
- [Samuel Smith - "Flow Based Programming"](https://youtu.be/j3cP8uwf5YM?feature=shared)

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

- [Strings, bytes, runes and characters in Go](https://go.dev/blog/strings)
- Concurrency
  - [Go Concurrency Patterns: Timing out, moving on](https://go.dev/blog/concurrency-timeouts)
  - [Go Concurrency Patterns: Context](https://go.dev/blog/context)
  - [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)

## JavaScript

- [MDN](https://developer.mozilla.org/en-US/)
- [TypeScript docs](https://www.typescriptlang.org/)
- [React docs](https://react.dev/)
- [You don't know JS books](https://github.com/getify/You-Dont-Know-JS)

## VSCode Extensions API Docs

- [Custom Editors](https://code.visualstudio.com/api/extension-guides/custom-editors)
- [Webviews](https://code.visualstudio.com/api/extension-guides/webview)
- [Language Servers (LSP)](https://code.visualstudio.com/api/language-extensions/language-server-extension-guide)
- [Syntax Highlighter](https://code.visualstudio.com/api/language-extensions/syntax-highlight-guide)
- [LSP Overview](https://microsoft.github.io/language-server-protocol/)
- [Go library for LSP implementation](https://github.com/tliron/glsp)
- [LSP official docs](https://microsoft.github.io/language-server-protocol/)

## Subjective Recommendations

### Videos

- ["Stop Writing Dead Programs" by Jack Rusher](https://youtu.be/8Ab3ArE8W3s?feature=shared)
- ["The Mess We're In" by Joe Armstrong](https://youtu.be/lKXe3HUG2l4?feature=shared)
- ["Propositions as Types" by Philip Wadler](https://youtu.be/IOiZatlZtGU?feature=shared)
- ["Outperforming Imperative with Pure Functional Languages" by Richard Feldman](https://youtu.be/vzfy4EKwG_Y?feature=shared)
- ["What Is a Strange Loop and What is it Like To Be One?" by Douglas Hofstadter (2013)](https://youtu.be/UT5CxsyKwxg?feature=shared)
- ["The Economics of Programming Languages" by Evan Czaplicki (Strange Loop 2023)](https://youtu.be/XZ3w_jec1v8?feature=shared)
- [Why Isn't Functional Programming the Norm? – Richard Feldman](https://youtu.be/QyJZzq0v7Z4?feature=shared)
- https://www.youtube.com/watch?v=SxdOUGdseq4 "Simple Made Easy" - Rich Hickey (2011)

### Books And Articles

- [Grokking Algorithms: An Illustrated Guide for Programmers and Other Curious People](https://www.amazon.com/Grokking-Algorithms-illustrated-programmers-curious/dp/1617292230)
- [Bob Martin's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [The Go Programming Language](https://www.amazon.com/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440)
- [Designing Data-Intensive Applications: The Big Ideas Behind Reliable, Scalable, and Maintainable Systems](https://www.amazon.com/Designing-Data-Intensive-Applications-Reliable-Maintainable/dp/1449373321)
- [Code: The Hidden Language of Computer Hardware and Software](https://www.amazon.com/Code-Language-Computer-Hardware-Software/dp/0735611319)

### Other

- [Go By Example](https://gobyexample.com/)

# Community

Here you can find help

- [Flow-Based Programming Discord](https://discord.gg/JHWRuZQJ)
- [r/ProgrammingLanguages](https://www.reddit.com/r/ProgrammingLanguages/)
