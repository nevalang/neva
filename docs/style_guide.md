# Style Guide

This guide sets standards for Nevalang code to ensure consistency and readability.

## Formatting

### Line Length

Keep lines under 80 characters.

- Good for split-screen
- Fits larger fonts without scrolling
- Leaves space for IDE features (code lens, git blame, inline-hints, etc.)
- Allows reading full lines with eye movement

### Indentation

Use tabs over spaces.

- Tabs let users set their preferred width
- Tabs reduce file size

### Imports

Use a single import block per file.

Group imports by type: stdlib, third-party, local. Separate groups with newlines if any group has more than 2 imports. Sort alphabetically within groups.

## Naming Conventions

Names should inherit context from parent scope. Good naming eliminates need for comments. Names generally rather short than long.

- **Packages/Files**: `lower_snake_case`
- **Types**: `CamelCase`
- **Interfaces**: `CamelCase` with `I` prefix
- **Constants**: `lower_snake_case`
- **Components**: `CamelCase` noun
- **Nodes**: `lower_snake_case`
- **Ports**: `lowercase`
- **Union Tags**: `CamelCase`

### Node Instantiation

- Prefer giving a node the same name as the component used to instantiate it (e.g. `println fmt.Println`).
- When wrapping a component in a higher-order component, mention both to retain clarity (e.g. `for_each_println ForEach{fmt.Println}`).

## Interfaces

- Use outports to separate data flows, not for destructuring.
- Use `res` for primary output and `err` for failures.
- Never use `data` as an outport name.
- Use `data` as an inport name only when the input is truly generic.
- Prefer domain names for inports when they add clarity (e.g. `url`, `filename`, `left`, `right`).
- Outport `err` must be of type `error`.
- Ports `data` and `res` of type `any` are interpreted as signals.
- Use name `sig` if you have _extra_ trigger-inport.
- Use names `then` and `else` if you implement boolean branching.
- Use specific inport names if you have more than one - e.g. `(filename, data)` for `io.WriteAll`.
- Use type-parameters to preserve type info between input and output if needed.

## Networks

- Keep components small and focused. Aim for about 3 nodes and 5 connections;
  split at 5 nodes or 10 connections unless the flat graph has a clear reason
  to stay together.
- Omit port names when possible. It enables renaming of ports without updating
  the networks.
- Use `?` to propagate errors unless custom error handling is needed.
- Prefer chaining connections inline when possible
  (e.g. `c -> switch:case[0] -> println`) to keep the dataflow compact and
  easier to scan.
- Treat dense networks with more than 5-6 connections as a smell. Prefer
  extracting a named helper component when it improves scanability.
- Prefer standard flow names: `sig` for trigger inputs, `res` for success
  outputs, and `err` for errors.

Example:

```neva
read:res -> fromBytes -> :res
```

## Comments

Good comments explain why.

Use leading `//` block immediately above entity.

Exported entities should have at least a short leading comment explaining their purpose.

- Free text is allowed and should describe intent/constraints.
- Use `@inport <name> <text>` for inport semantics.
- Use `@outport <name> <text>` for outport semantics.
- Use `@example <text>` for external usage examples (how to use component from outside).
- Multiple `@example` lines are allowed.
- Separate logical sections with an empty commented line (`//`).

Example:

```neva
// Processes payload and returns normalized result.
// Keeps stable behavior for repeated start signals.
//
// @inport start Trigger signal.
// @inport data Input payload.
//
// @outport res Normalized payload.
// @outport err Processing error.
//
// @example :start -> process:start
// @example 'hello' -> process:data
// @example process:res -> :stop
def Process(start any, data string) (res string, err error)
```

## Engineering Rules

### Prefer The Simplest Solution

Always prefer the simplest solution.

1. First, find the simplest solution.
2. Prove that it is insufficient.
3. Add complexity only after you proved it is necessary.
