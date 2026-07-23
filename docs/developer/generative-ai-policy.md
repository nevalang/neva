# Generative AI Policy

This policy applies only to the `nevalang/neva` repository. Other Nevalang
repositories may adopt different policies.

## Contributions

The use of generative AI is neither required nor discouraged. It is a tool,
like an editor, language server, formatter, autocomplete, or visual editor.

What matters is authorship responsibility. By opening a pull request, a
contributor represents that they understand every proposed change well enough
to explain its behavior, trade-offs, limitations, and tests, and to maintain it
after review.

We do not accept unexamined generated code. In this repository, “vibe coding”
means delegating implementation to an AI system without personally
understanding and taking responsibility for the resulting change. That is
incompatible with review and therefore with merging.

The origin of a change does not lower the bar for correctness, clarity,
maintainability, testing, or review.

## Neva for AI-Assisted Development

Neva is a general-purpose language, not an AI-first language. Its design is
intended to remain good for people writing code directly.

At the same time, we want Neva to be an excellent language for coding agents
and other AI-assisted workflows. The qualities that make code easier for people
to understand and validate also make it easier for language models to produce
and repair:

- a small, opinionated core;
- explicit dataflow and strong static semantics;
- predictable compilation and useful diagnostics;
- substantial validation performed by the compiler rather than optional
  external tooling.

These properties narrow the space of invalid programs and give an agent a
clear feedback loop: generate, compile, diagnose, and improve. AI assistance
must never come at the expense of human readability or sound language design.

## Applied GenAI

This is distinct from using AI to write software. We are also interested in
Neva as a language for building applications that use generative AI: models,
providers, tokens, context management, long-lived memory, agent protocols, and
related integrations.

Neva’s core will remain general-purpose. It will not acquire language features
that make it a specialized GenAI language.

However, thoughtfully designed standard-library packages for applied GenAI may
be appropriate alongside other practical facilities such as networking, JSON,
and I/O. Any such API must earn its place through a clear, stable, and
idiomatic design. No particular package design is promised today, and this work
should be approached deliberately rather than rushed.
