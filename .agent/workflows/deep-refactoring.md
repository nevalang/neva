---
description: A systematic approach to refactoring complex codebases.
---

# Deep Refactoring Workflow

Use this workflow when a component has become too complex to understand or modify safely, or when architectural flaws (e.g. duplicated logic) are identified.

## 1. Phase 1: Knowledge Mapping

The goal is to move from "mystery code" to "well-defined intent".

- **Plain English Specification**: Write a document describing what the function/component is _supposed_ to do in strictly plain English. Do not look at the code yet; use the current problem statement and requirements.
- **Dependency Visualization**: Use Mermaid flowcharts to map out how functions call each other. Identify:
  - Indirect recursion loops.
  - Points of high coupling.
  - Redundant entry points.
- **Data Flow Analysis**: Trace how data (especially complex objects) transforms through the system.

## 2. Phase 2: Decoupling

The goal is to separate _what_ is being done from _how_ and _where_.

- **Extract Pure Logic**: Identify "dirty" functions that do multiple things. Extract the pure logic into standalone, side-effect-free functions, when makes sense.
- **Unify Entry Points**: Ensure there is a single source of truth for specific operations, unless clearly needed.

## 3. Phase 3: Tuning & Validation

The goal is to polish, optimize, and ensure correctness.

- **Complexity Reduction**: Target non-optimal code paths.
- **Readability Pass**: Use the "plain English" descriptions from Phase 1 to rename variables and update comments so they match the intent.
- **Verification**:
  - Run existing tests.
  - Add/Remove/Change tests when/if needed and/or contributes to code quality.