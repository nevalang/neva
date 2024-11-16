# Style Guide

This guide sets standards for Nevalang code to ensure consistency and readability.

## Formatting

### Line Length

Keep lines under 80 characters.

- Comfortable for split-screen viewing
- Accommodates larger font sizes without horizontal scrolling
- Leaves room for IDE features (code lens, git blame, inline-hints, etc.)
- Enables reading full lines with eye movement alone

### Indentation

Use tabs over spaces.

- Tabs allow users to customize indentation width according to their preferences
- Tabs are more efficient in terms of file size

### Imports

Group imports by type: stdlib, local, third-party. Separate groups with newlines if at least 1 group has more than 2 imports. Sort alphabetically within groups.

## Naming Conventions

Names should inherit context from parent scope. Good naming eliminates need for comments. Names generally rather short than long.

- **Packages/Files**: lower_snake_case up to 3 words
- **Types**: CamelCase up to 3 words
- **Interfaces**: CamelCase with `I` prefix up to 3 words
- **Constants**: lower_snake_case up to 3 words
- **Components**: CamelCase noun up to 3 words
- **Nodes**: lowerCamelCase up to 3 words
- **Ports**: lowercase, 1 word up to 5 letters

## Interfaces

- Use outports to separate data flows, not for destructuring
- Use `data` for input with payload, `sig` for input without payload (trigger), `res` for output with payload, `sig` for output without payload (success), and `err` for failures.
- `err` outport must be `error` type. `sig` inport must be `any` type for flexibility. For outports, `sig` should be `struct{}` in concrete components and `any` in interfaces. Never use `any` for `res` outport.
- Don't send input data downstream - parent already knows it
- Use type-parameters when need to preserve type info between input and output

## Networks

- Omit port names when possible. It enables renaming of ports without updating the networks.
