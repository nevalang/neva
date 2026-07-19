# Vision

Neva is a general-purpose, statically typed, compiled dataflow language.

Its core value is not a single feature, but a combination of properties that
work together:

1. Hybrid programming model (text + visual tooling).
2. Concurrency-first execution model.
3. Strong static semantics and reliability.
4. AI-native development ergonomics.

## Core Principles

### 1) Hybrid Programming (Text + Visual)

Neva source remains a first-class text language, but the language model is also
optimized for visual programming workflows.

This is not "visual-only" positioning. It is a unified language that supports
both manual coding and visual graph workflows without semantic mismatch.

### 2) Concurrency-First By Design

Neva is designed around explicit node/edge dataflow, where concurrent execution
is default behavior rather than an advanced add-on.

Goal: make scalable multi-core utilization natural at language level, while
preserving predictable reasoning about program behavior.

### 3) Reliability Through Static Semantics

Neva prioritizes compile-time guarantees:

- strict static typing;
- semantic analysis before runtime;
- explicit data movement contracts;
- minimized surface for unsafe/implicit behavior.

The language should be easy to analyze and hard to misuse accidentally.

### 4) AI-Native, Without Sacrificing Human Authoring

Neva should be ergonomic for both:

1. human-written code (readable, maintainable, explicit);
2. AI-generated code (predictable, structurally consistent, easy to validate).

AI-native direction must not degrade manual development quality.

## Product/Application Tracks

Neva targets multiple practical domains where dataflow and concurrency are
valuable:

1. Web development (frontend/backend/fullstack workflows).
2. Streaming and event-driven processing.
3. Data transformation pipelines, including ETL-style workloads.
4. Network/distributed and integration-heavy systems.
5. ML-adjacent workloads where static contracts and pipelines matter.

These tracks are complementary. No single track defines Neva alone.
