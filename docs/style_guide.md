# Style Guide

This guide sets standards for Nevalang code organization, formatting, and naming conventions to ensure consistency and readability.

## Formatting

### Line Length

Keep lines under 80 characters.

- Comfortable for split-screen viewing
- Accommodates larger font sizes without horizontal scrolling
  - Better for visual accessibility
- Leaves room for IDE features (code lens, git blame, inline-hints, etc.)
- Enables reading full lines with eye movement alone

### Indentation

Use tabs over spaces.

- Tabs allow users to customize indentation width according to their preferences
  - More accessible for users with visual impairments who may need larger indentation
- Tabs are more efficient in terms of file size

## Naming Conventions

Names should inherit context from parent scope. Good naming eliminates need for comments. Names generally rather short than long.

- **Packages/Files**: lower_snake_case up to 3 words
- **Types**: CamelCase up to 3 words
- **Interfaces**: CamelCase with `I` prefix up to 3 words
- **Constants**: lowerCase up to 3 words
- **Components**: CamelCase noun up to 3 words
- **Nodes**: lowerCamelCase up to 3 words
- **Ports**: lowercase, 1 word up to 5 letters

## Interfaces

- Use type-parameters when need to preserve type information between input and output.
- Limit to 3 inports and outports max.

## Networks

- Prefer simple topologies over complex networks.
- Omit port names when possible
