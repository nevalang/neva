# Style Guide

This guide sets standards for Nevalang code organization, formatting, and naming conventions to ensure consistency and readability.

## Code Organization

- **File Length**: Aim for files under 300 lines. Longer files should be split.

## Formatting

- **Line Length**: Keep lines under 80 characters.
- **File Size**: Aim for files under 100 lines.
- **Comments**: Avoid comments. Use clear naming instead. If necessary, keep them short and simple.
- **Indentation**: Use tabs over spaces.

## API Design

- **Generics**: Use when data type consistency across ports is important.
- **Components**: Use outports for different data paths, structs for related data.
- **Network**: Prefer simple topologies (pipes, trees) over complex networks.
- **Ports**: Limit to 3 inports and 5 outports max.
- **Interfaces**: Keep small with minimal ports.

## Naming Conventions

Names should inherit context from parent scope. Good naming eliminates need for comments.

- **Packages/Files**: `lower_snake_case`, up to 3 words.
- **Entities**: `CamelCase` for types, interfaces, components. `lowerCamelCase` for constants.
- **Interfaces**: `CamelCase` with `I` prefix.
- **Components**: Noun, like for functions in most languages.
- **Ports**: Short `lowerCamelCase`, up to 5 letters.
- **Nodes**: `lowerCamelCase`, distinguish instances by meaning.
- **Enums**: Singular form.
