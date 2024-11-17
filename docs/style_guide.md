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

- **Packages/Files**: `lower_snake_case` up to 3 words
- **Types**: `CamelCase` up to 3 words
- **Interfaces**: `CamelCase` with `I` prefix up to 3 words
- **Constants**: `lower_snake_case` up to 3 words
- **Components**: `CamelCase` noun up to 3 words
- **Nodes**: `lowerCamelCase` up to 3 words
- **Ports**: `lowercase`, up to 5 characters

## Interfaces

- Use outports to separate data flows, not for destructuring.
- Use `data` for input with payload, `sig` for input without payload (trigger), `res` for output with payload, `sig` for output without payload (success), and `err` for failures.
- `err` outport must be of type `error`. `sig` inport must be of type `any` for flexibility. Never use `any` for `res` outport.
- Don't send input data downstream; the parent already has it.
- Use type-parameters to preserve type info between input and output.

## Networks

- Omit port names when possible. It enables renaming of ports without updating the networks.
