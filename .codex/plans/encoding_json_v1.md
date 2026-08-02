# `std/encoding/json` v1 implementation plan

Status: implementation in progress; PR 1 merged, PR 2 active

Last synchronized: 2026-08-02

Canonical tracking issue: https://github.com/nevalang/neva/issues/933

This file preserves the implementation plan across chat compactions. GitHub
issues remain the public decision log. Do not silently turn a decision gate
below into an implementation choice.

## Goal

Ship the first useful, statically typed JSON API:

```neva
pub def Marshal<T>(data T) (res bytes, err error)
pub def Unmarshal<T>(data bytes) (res T, err error)
```

The implementation uses the latest Go toolchain supported by the repository
and `encoding/json/v2`. The compiler remains unaware of JSON.

## Confirmed decisions

### Public API and scope

- The package path is `std/encoding/json`.
- v1 exposes only `Marshal<T>` and static `Unmarshal<T>`.
- Dynamic decoding is deferred to #1164.
- Streaming framing/cancellation is deferred to #1159 and must reuse the same
  per-value mapping.
- Full runtime introspection is deferred to #749.
- JSON-specific codec logic stays in `internal/runtime/funcs`; it does not
  belong in the general `internal/runtime/messages` domain package.

### Type metadata bridge

- Both `Marshal<T>` and `Unmarshal<T>` need the compiler-resolved structural
  form of `T`. Decoding needs it to construct a value. Encoding also needs it
  once the JSON mapping distinguishes `maybe<T>` from an ordinary tagged union:
  a runtime `UnionMsg` contains only its current tag/payload and cannot reveal
  the complete static union shape. Future serialization metadata reinforces the
  same requirement.
- The compiler injects that form into the existing static `FuncCall.Config` as
  one ordinary immutable Neva message. There is no new Go `Msg` subtype and no
  metadata attached to ordinary runtime values.
- The portable format is owned by a low-level `std/reflect` package:
  - `reflect.Type` is the complete finite descriptor passed around;
  - `reflect.TypeNode` is one node in the descriptor;
  - conceptually `Type` is `list<TypeNode>`, root at index 0;
  - every composite edge stores an integer index into that same list: list/dict
    item, struct field type, and tagged-union payload type;
  - recursion is an ordinary edge back to an existing node. There is no
    separate `TypeNode.Ref` variant.
- The initial wire fields are:
  - `Type = list<TypeNode>`;
  - `TypeNode.Struct = list<StructField>` where a field is
    `{ name string, node int }`;
  - `TypeNode.Union = list<UnionCase>` where a case is
    `{ tag string, data maybe<int> }`.
- The descriptor contains resolved JSON-relevant shape only: scalar kinds,
  bytes, list/dict item indexes, struct fields, and tagged-union cases.
  It omits source aliases, constraints, and generic parameters.
- IR contains only an ordinary message matching the `reflect.Type` wire shape;
  it does not contain a runtime reference to the Neva `reflect` package.

### Directive contract

- Working syntax: `#bind_type(<TypeExpr>)` on a component declaration.
- `<TypeExpr>` is any valid type expression resolvable in component scope, not
  only a declared type parameter.
- It may appear at most once, only on a component with `#extern`.
- v1 binds exactly one resolved type. Multiple bound types are deferred.
- It cannot be combined with a node-level static `#bind` config for the same
  runtime call because `FuncCall.Config` currently carries one message.
- The parser/AST own the directive vocabulary and syntactic argument shape.
  Unknown directive names are parse diagnostics, not arbitrary strings in AST.
- The analyzer owns semantic validation: placement, cardinality, combinations,
  scope, and successful type resolution.
- The current grammar accepts only `IDENTIFIER` inside a directive argument.
  Therefore #1163 must extend the grammar/parser for a structured `TypeExpr`;
  listener-side string reparsing is not acceptable.

### JSON mapping already accepted

- `bytes` encode as a base64 JSON string and decode from that form.
- `NaN`, `+Inf`, and `-Inf` are errors in v1. JSON has no number syntax for
  them; a future explicit format option may encode them as strings.
- Struct field matching is exact and case-sensitive.
- Duplicate JSON object names are errors.
- Unknown struct fields are ignored in v1.
- The future serialization-metadata design must be able to opt into rejecting
  unknown fields; the concrete syntax (comments, attributes, or another form)
  is intentionally not chosen yet.
- Missing Neva struct fields are errors. Neva has neither partial struct
  construction nor Go-style zero values to fill omitted fields.
- JSON field names initially equal Neva field names. Renaming, omission,
  defaults, field-level optionality/strictness, and format-specific metadata
  are follow-up work.
- Ordinary tagged unions keep the existing Neva JSON shape: a payload case is
  `{ "tag": "Case", "data": ... }`; a tag-only case is
  `{ "tag": "Case" }`. The structurally recognized `maybe<T>` shape is the
  sole v1 exception under discussion below.
- `null` is not and will not become a Neva language value or type. It must not
  be added to syntax, AST, the type system, or the runtime scalar set.
- `null` support is required for the first useful package, so #1168 is a v1
  design prerequisite rather than a post-v1 enhancement. The current
  recommendation maps the resolved structural `maybe<T>` shape
  (`Some T | None`) to JSON value/null, rejects null for non-maybe static
  targets, and accepts the conventional narrow lossiness of nested `maybe`.

### Runtime construction and lifecycle

- Native decoding may use mutable Go accumulators while parsing, but publishes
  only a fully built immutable Neva message.
- Static descriptor/config defects are compiler bugs. Runtime-function `Create`
  may fail fast while constructing the handler. The event handler must return a
  normal `err` output for malformed user JSON or target/value mismatch and must
  not panic during normal execution.
- Add a read-only allocation-free `messages.StructRange` (exact name can follow
  repository API convention) for external field iteration. Use it in every
  call site where it removes a temporary map/projection without obscuring a
  more direct internal loop.

## Decisions that remain explicit gates

### Dict key order during marshal

The original deterministic choice had real motivation:

- reproducible bytes for logs, snapshots, caches, diffs, and signatures;
- compatibility with Go `encoding/json` v1 and the current Neva message JSON
  path;
- less flaky byte-level testing.

It also costs a key slice, `O(n log n)` sorting, and delays streaming output.
Go JSON v2 deliberately traverses maps non-deterministically by default and
offers deterministic encoding as an option.

Current recommendation: follow JSON v2 and do not sort dict keys in the default
`Marshal`; preserve declaration order for structs without sorting. The public
contract promises JSON semantic equivalence, not stable bytes. Before coding,
confirm this recommendation explicitly. If deterministic/canonical output is
needed, add a separate option/component rather than silently taxing every call.

### Compiled type plan

This mechanism is accepted. “Lookup table” and “cache” were misleading names:
there is no process-global type registry. The portable `reflect.Type` is
already a flat indexed graph. At runtime-function `Create`, convert that one
ordinary Neva message into one private flat native Go representation for that
runtime node. Child types remain integer indexes into the same node slice;
direct maps may additionally map struct field or union case names to indexes.
The handler closure owns the compiled plan, and normal Go reachability/GC
releases it with the program.

The purpose is modest and implementation-local: validate/extract the static
message once and avoid repeated `Msg` interface traversal, union extraction,
type assertions, and boxed map/list access for every encoded or decoded JSON
value. It is not a second portable type model and not a cross-program cache.

Do not add a generic `ReflectTypeView` that merely wraps and retains the raw
message; it adds indirection without solving the hot-path cost. Keep each
compiled plan private to `internal/runtime/funcs` because it is a JSON codec
execution plan, not the portable `reflect.Type` API.

Still benchmark it against direct descriptor traversal. Measure `Create`
time/allocation, single-decode latency, repeated-decode throughput, and retained
bytes for small and nested schemas. These measurements validate the expected
optimization; they no longer block the basic representation choice.

### JSON null interoperability

Rejecting `null` is implementable but probably does not meet the “first useful
package” goal. Before JSON codec implementation, #1168 must choose a mapping
without a language-level null. Evidence from Rust Serde and OCaml deriving
supports treating the option/maybe constructor as a special JSON mapping:

- `None` maps to JSON `null`;
- `Some(v)` maps to the JSON representation of `v`;
- an ordinary tagged union that is not structurally `Some T | None` keeps its
  tagged representation;
- a future dynamic JSON value may represent Null as one case of an ordinary
  tagged union;
- null for a non-maybe static target is an error;
- missing struct fields remain errors unless future serialization metadata
  explicitly changes that contract.

Neva is structurally typed, so the candidate can recognize the fully resolved
two-case shape instead of relying on a discarded source alias. Whether every
structurally equivalent union should receive this mapping must be confirmed.

Nested `maybe` must have a stated reversibility rule. Options are: accept
Serde-like lossiness; reject only values such as `Some(None)` that JSON cannot
distinguish from `None`; or use a context-dependent tagged inner encoding. The
current recommendation is to follow the established Rust Serde/OCaml option
mapping and accept the narrowly documented lossiness. This preserves a simple,
compositional rule at every nesting level and avoids making a value's encoding
depend on its surrounding type depth. The round-trip law must explicitly
exclude values collapsed by the JSON mapping, such as nested `Some(None)`.

`any` does not itself solve JSON null. There is no runtime `AnyMsg`, and a
static target of `any` gives the decoder no basis for inventing
`maybe<any>.None`. Until a dynamic JSON value carrier is designed, null under
an `any` target remains an error. A future dynamic carrier may represent Null
as an ordinary `encoding/json.Value` tagged-union case hidden behind `any`; that
is separate from the static `maybe<T>` mapping and belongs with #1164.

Special encoding follows the static descriptor. If a `maybe` value is hidden
behind `any`, `Marshal<any>` cannot prove that it is `maybe` and must use the
ordinary tagged-union representation instead of the null/value exception.

### Static `any` decoding

The exact non-null JSON-to-`any` mapping still needs confirmation before the
codec PR. The likely rule is self-describing decoding: bool/string to their
scalar messages, syntactically integral in-range numbers to int, other finite
numbers to float, arrays to untyped `list<any>`, and objects to `dict<any>`.
Required gates are integer overflow/precision and whether v1 supports `any`
targets at all before the dynamic `encoding/json.Value` carrier. JSON null
under `any` remains an error under the current static design.

## Why Rust Serde is not automatically reversible

Serde treats `Option` specially for JSON. For example, these distinct Rust
values both produce `null`:

```text
Option<Option<Int>>::None
Option<Option<Int>>::Some(Option<Int>::None)
```

The first serializes the outer `None`; the second serializes `Some(...)` by
serializing the inner `None`. The decoder cannot infer which value produced
`null`. This is not generally catastrophic: it matters only when a program
needs those two states to remain distinct, or when the codec promises exact
round trips. Ordinary Rust enums use an enum representation; `Option` is a
special serialization case. The current Neva recommendation deliberately
accepts the same documented nested ambiguity in exchange for one simple,
compositional rule.

## Struct field policy rationale

- Unknown: ignore. This supports forward-compatible producers and matches
  Neva's structural "consumer needs a subset" intuition. A later strict option
  may reject unknown fields.
- Missing: error. Go leaves omitted fields at their pre-existing destination
  values, which are usually zero values for a fresh struct. Neva returns a new
  immutable value and has no zero/default-field semantics, so copying Go here
  would invent a new language behavior.
- Duplicate: error. JSON does not define duplicate-member semantics; rejecting
  ambiguity follows Go v2 and avoids security/interoperability splits.
- Exact case: require. This follows Go v2 and ordinary Neva field identity.

## Pull request sequence

### PR 1: structured compiler directives (#1163)

Status: merged in PR #1173.

Scope:

- move concrete directive kinds and structured directive AST nodes to `pkg/ast`;
- extend grammar/parser to represent directive-specific arguments, including a
  full type expression for future `#bind_type`;
- reject unknown directive names during parsing;
- preserve `#extern`, `#bind`, and `#autoports` behavior;
- keep semantic placement/combinations in analyzer.

Tests: grammar/parser table tests for no argument, identifier argument, type
expression, malformed input, unknown directive, source locations; existing
compiler/e2e suite for old directives.

### PR 2: `std/reflect` descriptor contract

Status: active on `codex/reflect-type-descriptor`.

Scope:

- add Neva `reflect.Type` and `reflect.TypeNode` declarations;
- add the minimum Go-side message construction/reading helpers needed for the
  portable shape, without JSON logic or public `TypeOf`;
- cover recursion through cyclic node indexes, deterministic struct/union
  traversal, out-of-bounds indexes, and empty descriptors.

Tests: Neva package tests, message/helper unit tests, finite recursive examples,
descriptor equality/round-trip fixtures. Update developer docs and #749 only to
state the narrow descriptor substrate; do not claim general reflection.

### PR 3: generic `#bind_type` compiler bridge (#1160)

Scope:

- add directive/analyzer invariants listed above;
- resolve arbitrary valid `TypeExpr` after analysis;
- lower resolved type into one `reflect.Type`-shaped config message;
- carry it through IR and generated Go without JSON knowledge;
- align repository and generated modules to the latest supported Go toolchain;
- enable required Go experiments uniformly for all generated native/wasm child
  builds without dropping other configured experiments.

Tests: analyzer positive/negative table tests; lowering fixtures for scalars,
empty typed containers, nested composites, tagged unions, recursion, aliases,
and literal `any`; backend compile/e2e contract under `e2e/`.

### PR 4: JSON codec and public package (#933)

Scope:

- add `std/encoding/json` declarations and extern registrations; both generic
  components request their resolved `T` through `#bind_type(T)`;
- implement marshal/unmarshal logic inside `internal/runtime/funcs`, sharing
  private helpers between creators;
- add allocation-free struct iteration and reuse it where beneficial;
- compile each component's `reflect.Type` config once in `Create` into a private
  immutable encode/decode plan;
- implement accepted bytes, floats, struct-field, duplicate, and the #1168 null
  policy;
- apply the explicit dict-order decision made at the gate.

Tests:

- table-driven runtime-function unit tests;
- package/component tests;
- several root `e2e/` modules, not one overloaded example;
- native and wasm coverage where runtime build paths differ;
- malformed JSON, duplicate fields, unknown fields, missing fields, case
  mismatch, type mismatch, empty typed containers, nested lists/dicts/structs,
  tagged unions with/without payload, recursion, alias-equivalent targets,
  bytes, finite/non-finite floats, and null at every nesting position;
- property/round-trip coverage for JSON-representable values.

### PR 5: performance gate and documentation

This may be folded into PR 4 only if review remains manageable; it must land
before declaring v1 complete.

- benchmark direct descriptor traversal versus the compiled codec plan;
- benchmark dict marshal with and without sorting on small/medium/large dicts;
- benchmark flat and nested encode/decode, typed containers, structs, unions,
  and repeated decode through one runtime node;
- compare against current printing/MarshalJSON baseline where semantically
  meaningful;
- document public mapping, errors, ordering guarantee, `null` limitation, and
  runtime/compiler ownership;
- update #933 with results and close only after the chosen contracts and tests
  are merged.

## Test laws

For every JSON-representable `v : T`:

```text
Unmarshal<T>(Marshal(v)) == v
```

“JSON-representable” excludes values collapsed by an intentionally lossy
mapping. Under the current nullable recommendation, `None` and nested
`Some(None)` both encode as `null`, so exact equality is not promised for that
pair. All other successful mappings remain subject to the law.

For accepted JSON input, marshal after unmarshal may change whitespace and
object-member order but must preserve JSON semantics. This law permits order
normalization; it does not itself require nondeterministic output for repeated
marshaling of the same value.

## Follow-up issue map

- #933: static `std/encoding/json` v1 and canonical decision record.
- #1163: structured directive AST/parser work.
- #1160: directive model and `#bind_type` bridge.
- #749: future runtime type assertion/introspection.
- #1159: streaming JSON decoder.
- #1164: dynamic non-stream decoding.
- #1167: serialization metadata/tags and per-format policies.
- #1168: JSON `null` mapping without language-level null.

## Naming note outside JSON scope

`internal/builder.Builder` loads the module graph, injects stdlib dependencies,
downloads dependencies, and returns a `compiler.RawBuild`. `ModuleLoader` or
`BuildResolver` would be more descriptive than `Builder`, but renaming it is
unrelated to JSON and should not share these PRs.

## Next implementation action

Finish PR 2 from current `origin/main`, then implement the generic
`#bind_type` bridge in PR 3. Do not begin JSON codec code before both the
descriptor contract and compiler bridge land. The remaining decisions that
must be closed before PR 4 are: default dict ordering, final nested-maybe
lossiness confirmation, and static `any` decoding policy.
