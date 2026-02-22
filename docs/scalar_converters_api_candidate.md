# Scalar Converters API Candidate (Strict Go-Parity Split)

This document proposes a concrete API shape for scalar converters in Neva,
following strict Go parity as the default direction:

- `builtin`: language-style casts (Go-like conversion semantics)
- `strconv`: text parsing/formatting APIs (Go-like package role)

It is a candidate API only (signatures + behavior contract), not an
implementation plan.

## Goals

1. Keep conversion UX simple and predictable.
2. Keep core language syntax unchanged (components-only model).
3. Align with Go behavior, not just naming.
4. Avoid magic truthiness/implicit coercions.

## Non-goals

1. No new language syntax for casts.
2. No automatic `bool <-> number` conversions.
3. No generic "convert everything to everything" API.

## Builtin Candidate

Go-like language casts.

```neva
// std/builtin/scalars.neva (candidate)

// Int converts float to int by truncating fractional part toward zero.
// See semantics section for NaN/Inf/out-of-range behavior.
#extern(int_from_float)
pub def Int(data float) (res int)

// Float converts int to float.
#extern(float_from_int)
pub def Float(data int) (res float)

// String converts integer code point to a UTF-8 string (Go-like string(int)).
// This is not decimal formatting.
#extern(string_from_int_codepoint)
pub def String(data int) (res string)
```

### Builtin Semantics

`Float(int) -> float` (Go-like numeric conversion):

1. Convert `int64` to `float64` with IEEE754 rounding as needed.
2. Always succeeds.

`Int(float) -> int` (Go-like numeric conversion):

1. Fractional part is discarded (toward zero).
2. For `NaN`, `+Inf`, `-Inf`, and out-of-range finite values, behavior follows
   Go conversion semantics.

`String(int) -> string` (Go-like code-point conversion):

1. Converts integer to Unicode code point UTF-8 representation.
2. This is not `"42"` decimal formatting; for decimal use `strconv.Itoa`.

## Strconv Candidate

Text parsing/formatting for scalars (Go-like `strconv` role).

```neva
// std/strconv/strconv.neva (candidate additions)

// ParseBool parses textual booleans.
#extern(parse_bool)
pub def ParseBool(data string) (res bool, err error)

// FormatBool formats bool as "true" or "false".
#extern(format_bool)
pub def FormatBool(data bool) (res string)

// Itoa formats int in base 10.
#extern(itoa)
pub def Itoa(data int) (res string)

// FormatInt formats int with explicit base.
// Go-style API: no error outport.
#extern(format_int)
pub def FormatInt(data int, base int) (res string)

// FormatFloat formats float in canonical decimal form.
// Candidate keeps API minimal; no fmt-verb surface here.
#extern(format_float)
pub def FormatFloat(data float) (res string)
```

Existing (already present):

```neva
#extern(atoi)
pub def Atoi(data string) (res int, err error)

#extern(parse_int)
pub def ParseInt(data string, base int, bits int) (res int, err error)

#extern(parse_float)
pub def ParseFloat(data string, bits int) (res float, err error)
```

### Strconv Semantics

`ParseBool`:

1. Accept `true/false` plus case variants and `1/0` (Go-like).
2. Any other input returns `err`.

`Itoa`:

1. Base-10 integer text.

`FormatInt`:

1. Valid bases are 2..36.
2. Out-of-range base follows Go-style behavior (no `err` outport).

`FormatFloat`:

1. Returns canonical string form intended for machine-oriented conversion flows.
2. Contract should target stable round-trip behavior with `ParseFloat`.

## Explicitly Not Included

These are intentionally excluded from this candidate:

1. `Bool(int)`, `Bool(float)`
2. `Int(bool)`, `Float(bool)`
3. `String(float)` and `String(bool)` builtin casts

Rationale: follow Go conversion surface and keep boolean/number coercions
explicitly out of builtin casts.

## Example Usage

```neva
import {
    strconv
}

def Main(start any) (stop any) {
    to_float Float
    parse_bool strconv.ParseBool
    int_to_str strconv.Itoa
    cp_to_str String
    ---
    :start -> [
        42 -> to_float,
        'true' -> parse_bool,
        123 -> int_to_str,
        42 -> cp_to_str
    ]
}
```

## Migration Notes

If this candidate is accepted:

1. Keep `Atoi/ParseInt/ParseFloat` as-is (no breaking changes).
2. Add new `builtin` scalar cast components.
3. Add `strconv` parse/format helpers incrementally.
4. Document precise semantics in `docs/qna.md` and `docs/book`.
