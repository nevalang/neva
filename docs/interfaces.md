# Interfaces

Interface is a _component signature_ that describes _abstract component_ - its input and output _ports_ and optional type parameters. Interfaces are used for _dependency injection_.

## Ports

Port definition consist of a _type expression_ describing the data-type port expects and a flag that describes whether the port is an _array_ or _single_ port. Type expression can refer to interface's type parameters. If no type paremeter given then `any` is implicitly used.

## Single Ports

Single port is port with one _slot_. Reference to such ports should not include slot index.

## Array Ports

Array port is port with multiple (up to 255) _slots_. Such ports must be referenced either via slot indexes or in _array-bypass connection_ expressions.
