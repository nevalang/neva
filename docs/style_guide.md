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

### Node Instantiation

- Prefer giving a node the same name as the component used to instantiate it (e.g. `println fmt.Println`).
- When wrapping a component in a higher-order component, mention both to retain clarity (e.g. `for_println For{fmt.Println}`).

## Interfaces

- Use outports to separate data flows, not for destructuring.
- Use `data` for input, `res` for output and `err` for failures.
- Outport `err` must be of type `error`.
- Ports `data` and `res` of type `any` are interpreted as signals.
- Use name `sig` if you have _extra_ trigger-inport.
- Use names `then` and `else` if you implement boolean branching.
- Use specific inport names if have more than one - e.g. `(filename, data)` for `io.WriteAll`.
- Use type-parameters to preserve type info between input and output if needed.

## Networks

- Omit port names when possible. It enables renaming of ports without updating the networks.
- Use `?` for to propogate errors except custom error handling is needed.
- Prefer chaining connections inline when possible (e.g. `c -> switch:case[0] -> fmt.Println`) to keep the dataflow compact and easier to scan.
