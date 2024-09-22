# Program Structure

## Build

Build is the highest-level abstraction and the first compiler stage. It collects all source code, including dependencies, into a single object for analysis.

## Module

A set of packages with a root directory containing `neva.y(a)ml` (manifest). The manifest defines the minimum supported language version and dependencies. The main module contains the main package with the main component.

## Package

A set of files in a single directory. Entities in one file are visible from other files in the same package. The main package must contain the `Main` component with a special signature and no exported entities.

## File

A `.neva` file containing entities and imports. Imports allow dependencies on entities from other packages. The `builtin` package is implicitly imported.

## Imports

Import statements consist of a path and optional alias. There are three types: std, 3rd-party, and local. Entity references are resolved first within the package, then in `builtin`.

## Entity

The core abstraction in Nevalang. There are four kinds:

1. Types - message shape definition
2. Constants - reusable messages with static values
3. Interfaces - abstract components
4. Components - computation units

Entities can be public (`pub`) or private, determining their visibility outside the package.

## Entity Reference

Entity references include an optional package name and the entity name. Package names can be omitted for entities in the same package or in `std/builtin`. Local entities with the same name as builtin entities shadow them.
