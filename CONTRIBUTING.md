# Requirements

Make sure you've read the https://nevalang.org/docs from cover to cover

## System

- Go: https://go.dev/doc/install
- Make: https://www.gnu.org/software/make/#download
- NodeJS and NPM: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm/
- Antlr: `pip install antlr4-tools`
- Tygo: `go install github.com/gzuidhof/tygo@latest`

## VSCode

These are not really required but recommended in order you're using VSCode

- [nevalang](https://marketplace.visualstudio.com/items?itemName=nevalang.vscode-nevalang)
- [antlr4](https://marketplace.visualstudio.com/items?itemName=mike-lischke.vscode-antlr4)
- [tmlanguage](https://marketplace.visualstudio.com/items?itemName=pedro-w.tmlanguage)
- [markdown-mermaid](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid)

# Development

See [architecture high level overview](./ARCHITECTURE.md) and what [Makefile](./Makefile) can do.

## VSCode Extension

Extension depends on types defined in the `src` and `typesystem` packages so it's dangerous to rename those. If you going to do so, make sure you did't brake TS types generation.

Check out [tygo.yaml](./tygo.yaml). and CONTRIBUTING.md in vscode-neva repo.

## ANTLR Grammar

Don't forget to open `neva.g4` file before debugging with VSCode

# Naming conventions

## Tests

Use `_` instead of space in for test-case names because go turns spaces into underscores and makes it hard to find specific case.

# Learning Resources

## FBP/DataFlow

- [Flow-Based Programming: A New Approach to Application Development](https://jpaulmorrison.com/fbp/1stedchaps.html)
- [Dataflow and Reactive Programming Systems: A Practical Guide](https://www.amazon.com/Dataflow-Reactive-Programming-Systems-Practical/dp/1497422442)

## Golang

Advanced golang knowledge is required. Especially understanding of concurrency.

- [Concurrency is not parallelism](https://go.dev/blog/waza-talk)
- [Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Go Concurrency Patterns: Timing out, moving on](https://go.dev/blog/concurrency-timeouts)
- [Go Concurrency Patterns: Context](https://go.dev/blog/context)
- [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)

## Subjective Recommendations

### Videos

- ["Stop Writing Dead Programs" by Jack Rusher](https://youtu.be/8Ab3ArE8W3s?feature=shared)
- ["The Mess We're In" by Joe Armstrong](https://youtu.be/lKXe3HUG2l4?feature=shared)
- ["Propositions as Types" by Philip Wadler](https://youtu.be/IOiZatlZtGU?feature=shared)
- ["Outperforming Imperative with Pure Functional Languages" by Richard Feldman](https://youtu.be/vzfy4EKwG_Y?feature=shared)
- ["What Is a Strange Loop and What is it Like To Be One?" by Douglas Hofstadter (2013)](https://youtu.be/UT5CxsyKwxg?feature=shared)
- ["The Economics of Programming Languages" by Evan Czaplicki (Strange Loop 2023)](https://youtu.be/XZ3w_jec1v8?feature=shared)
- [Why Isn't Functional Programming the Norm? – Richard Feldman](https://youtu.be/QyJZzq0v7Z4?feature=shared)

### Books And Articles

- [Bob Martin's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Code: The Hidden Language of Computer Hardware and Software](https://www.amazon.com/Code-Language-Computer-Hardware-Software/dp/0735611319)

# Community

Check out https://nevalang.org/community to find out where you can get help