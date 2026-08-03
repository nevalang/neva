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

### Composite Literals

Keep an empty list or struct literal compact: `[]` and `{}`. Keep a non-empty
literal on one line when its complete surrounding line fits under 80
characters. Otherwise write one item or field per line, with a trailing comma
after the final item and the closing delimiter on its own line.

Apply the same compact-or-vertical rule to comma-separated type parameters,
type arguments, interface ports, fan-in, and fan-out. Every vertical sequence
has a trailing comma after its final item. Union variants are structural
declarations: write each on its own line without commas.

For struct and union type declarations, keep an empty or single-item body
compact. Write two or more fields or variants as a block: place the opening
and closing braces on their own lines and indent every item once.

### Imports

Group imports by type: stdlib, third-party, local. An import without a module prefix is stdlib, an `@:` prefix is local, and every other explicit module prefix is third-party. Separate groups with newlines if any group has more than 2 imports. Sort alphabetically within groups.

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
- Omit the name on each interface side that has exactly one port. This keeps
  the interface structural and lets components use their own port names.
- Use specific inport names if you have more than one - e.g. `(filename, data)` for `io.WriteAll`.
- Use type-parameters to preserve type info between input and output if needed.

## Networks

- Keep components small and focused. Aim for about 3 nodes and 5 connections;
  split at 5 nodes or 10 connections unless the flat graph has a clear reason
  to stay together.
- Omit a node's port name when its resolved interface has exactly one port in
  that direction. This lets implementations rename that port without updating
  the network.
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

Good comments explain the entity's purpose and the behavior a user needs to
compose it correctly.

Use a leading `//` block immediately above an entity. Exported entities should
have at least a short comment explaining their purpose.

- Free text describes intent, constraints, and other entity-level behavior.
- Use `@inport <name> <text>` and `@outport <name> <text>` for port behavior.
  Describe when that port receives or sends messages, what those messages mean,
  and the port's role in the component. Include ordering, completion, errors,
  or side effects when they matter to composition. Do not use port tags merely
  to restate a port's type or to describe an abstract value detached from the
  port's behavior.
- For an interface's single anonymous input or output port, use `_` in place
  of its name: `@inport _ <text>` or `@outport _ <text>`.
- When an interface or component uses port tags, document every one of its
  inports and outports.
- Use `@example <text>` for an external usage example when it makes the
  component easier to apply. Multiple `@example` lines are allowed.
- Separate meaningful sections with an empty commented line (`//`).

Example:

```neva
// Normalizes each input message after a start signal arrives.
// Repeated start messages begin independent normalization requests.
//
// @inport start Starts one normalization request.
// @inport data Supplies the message normalized by that request.
//
// @outport res Sends the normalized message when processing succeeds.
// @outport err Sends the error for a failed request instead of `res`.
//
// @example :start -> process:start
// @example ' hello ' -> process:data
// @example process:res -> :stop
def Process(start any, data string) (res string, err error)
```

## Engineering Rules

### Prefer The Simplest Solution

Always prefer the simplest solution.

1. First, find the simplest solution.
2. Prove that it is insufficient.
3. Add complexity only after you proved it is necessary.
