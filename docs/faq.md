# Frequently Asked Questions (FAQ)

## When to use generics (type-parameters/arguments)?

When you need to preserve data-type on the output.

E.g. `Void` component doesn't have outports so it doesn't matter what you passed in. That's why `Void` accepts `any` instead of `T`.

On the other hand `Lock` needs to know the type of `v` on the input so the type of the `v` on the output is preserved. That's why it's `Void<T>`.

## When to use `struct` and when to create different outports?

When you have _structured_ data data use `struct`, when you ant to _separate dataflows_ - create outports.

Example: `ParseInt` sends value and error but never at the same time. It's always one of them firing depending on the input.