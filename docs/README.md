# Documentation

Welcome to Nevalang's documentation!

It describes general purpose dataflow compiled language with static types. Here you'll find motivation behind the language, it's philosophy, differences between dataflow and controlflow paradigms and of course language abstractions and relations between them.

## About this Document

1. This document describes finite state of the language. Some features (e.g. visual editor) might not be implemented at the time of writing, but they yet important for concept so they will be mentioned
2. It was written as an attempt to create at least something. I needed to dump all the information I collected about the language with years. So it's far from perfect
3. This document doesn't teach you how to write Nevalang programs, because Nevalang is very immature language and it is changins all the time. This is especially true for stdlib components. However, a lot about Nevalang is already clear and will never change. We are talking about more fundamental stuff such as philosophy, abstractions, execution model, etc.

## Table of contents

- [About](./about.md)
- [Motivation](./motivation.md)
- [Paradigm](./paradigm.md)
  - [Flow-Based-Programming](./fbp.md)
- [Program Structure](./program_structure.md)
  - [Type](./type_entity.md)
  - [Constant](./const_entity.md)
  - [Interface](./interface_entity.md)
  - [Component](./component_entity.md)
    - [IO](./component_io.md)
    - [Network](./component_net.md)
