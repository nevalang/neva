# Const Literal Senders: Analyzer Implications

## Summary

We updated the grammar so any `constLit` is a valid network sender (instead of only `primitiveConstLit` + `unionLit`). This makes list/struct/dict literals syntactically valid senders. Union literals remain valid and can wrap any const literal.

## Why This Matters

The analyzer currently expects literal senders to be primitive (or union literal with a known union type). It relies on `ConnectionSender.Const.TypeExpr` being set so it can resolve sender types in `getConstSenderType` and enforce subtype checks. For list/struct/dict literals we do not currently infer a type, so analysis will fail with "Literal sender type is empty" unless we add inference or explicit restrictions.

## Immediate Impact

- **Primitives**: continue to work (type expr set to `bool/int/float/string`).
- **Union literals**: continue to work (type expr set to union entity ref). Union literals can already wrap non-primitive values.
- **List/struct/dict literals**: now parse as senders, but analyzer lacks type inference. This will likely raise errors during analysis.

## Open Design Questions

1) **Type inference for complex literals**
   - Lists: need element type(s) or an `any`-like list type.
   - Struct/dict: syntax is shared, so we need a rule to distinguish or define a unified type.
   - Nested literals: need recursive inference.

2) **Semantics for dict vs struct**
   - The grammar uses `structLit` for both struct and dict syntax.
   - Analyzer would need a rule to pick `struct` vs `dict`, or treat them as one until type is known.

3) **Strictness**
   - Option A: implement inference for all const literal senders.
   - Option B: explicitly reject non-primitive literal senders in analyzer for now (even though grammar allows them).

## Current Parser Behavior

The parser now accepts `constLit` on sender side and builds a `Const` with:

- `TypeExpr` set for primitives and union literals.
- `TypeExpr` empty for list/struct/dict literals (needs inference or explicit rejection later).

## Next Steps (Decision Needed)

- Decide whether to implement full inference for complex literal senders, or to forbid them in analyzer until inference is ready.
- If forbidding, add a clear diagnostic error ("non-primitive literal sender not supported") so the grammar stays permissive but semantics remain strict.
