# AGENTS.md

Follow these instructions.

## 1. 🤖 Operating Protocol

1. Use `context7` MCP server (when available) to fetch libraries API documentation.
2. For GitHub context, try `gh` (GitHub CLI) first; if unavailable or insufficient, fall back to `curl` (e.g., issues/PR descriptions/comments).
3. Run `golangci-lint` and `go test`. Fix warnings.
4. If uncertainty > 10%, ask user.
5. Update this file if changes to process, architecture, or rules.
6. For any `*.neva` change (parser/stdlib/examples/e2e), explicitly follow `docs/style_guide.md` (node aliases `lower_snake_case`, tabs, import ordering, omit explicit port names when unambiguous). Use `go.mod` as source of truth for Go imports.
7. Plan -> Review -> Execute -> Review.
8. Refactor: Actively identify and resolve unnecessary complexity or duplication. Prioritize code clarity and long-term maintainability over chasing theoretical perfection.
9. Use targeted tests and cap long-running commands to ~5 minutes unless explicitly requested otherwise.
10. PR review workflow: when user asks to address PR comments, apply the requested code/doc changes first, then post a reply to each addressed comment via `gh`; do not resolve review threads unless user explicitly asks.
11. For generated tests, add short intent comments so the expected behavior is obvious without reverse-engineering fixtures.

## 2. 📈 Self-Improvement Protocol

**After each session** (bug fix, feature, brainstorm), update this file (`AGENTS.md`) with:

- **Language semantics** learned (e.g., how connections work, port behavior).
- **Common patterns** discovered (e.g., typical error causes, debugging approaches).
- **Architecture insights** gained (e.g., how compiler phases interact).
- **Gotchas** encountered (e.g., edge cases, non-obvious behaviors).

**Balance**: Keep it concise. Every line must earn its place. Remove outdated info. Split into workflows/rules if sections grow large.

**Goal**: Build perfect context so future sessions start smarter.

### Session Notes (2026-01-23)

- **Language semantics**: Chained connections nest; receiver type constraints must use the sender from the same chain link (not the outer sender).
- **Language semantics**: Network senders now parse any `constLit`; only primitives/union literals get a type expr without additional analyzer inference.
- **Language semantics**: Union tag/type compatibility relies on structural union checks via normal sender/receiver subtype validation.
- **Common patterns**: Overload/generic resolution relies on `deriveNodeConstraintsFromNetwork`; nested chains need explicit sender-to-receiver pairing to avoid union/any leakage.
- **Architecture insights**: Switch case output type inference is used in analyzer (before desugaring) for both overload resolution and network type checks.
- **Architecture insights**: Union node tag/data validation is driven by `buildUnionTagInfos` + `validateUnionDataReceivers` to keep Union<T> wiring consistent.
- **Gotchas**: Portless outports require explicit port names when multiple outports exist; generic std nodes like `fmt.Println` need explicit type args if type inference is unavailable.

### Session Notes (2026-01-26)

- **Language semantics**: Union:tag currently assumes a single sender; fan-in is explicitly rejected to keep wiring semantics unambiguous.
- **Common patterns**: Avoid extra union-type equivalence checks at tag-binding time; rely on standard subtype validation for compatibility.
- **Gotchas**: In sender position, `[...]` is parsed as fan-in (multiple senders), so list literals should not be used as senders without explicit disambiguation.
- **Gotchas**: Literal senders are limited to primitives/union literals; list/dict/struct literals are rejected for now (use const refs).

### Session Notes (2026-01-27)

- **Language semantics**: Literal senders must be in a chain triggered by a signal (e.g. `:start -> true -> ...`).
- **Language semantics**: Fan-out (`-> [...]`) can contain chained connections, but the chain head is the sender before the fan-out (e.g. `:start -> [foo, bar -> baz]`).
- **Language semantics**: `std/core.New<T>` now requires a `sig` inport; explicit uses must be triggered (e.g. `:start -> new:sig`).
- **Common patterns**: Stdlib Switch setups should seed case values via `:start` to avoid autonomous constant senders.
- **Architecture insights**: `analyzeSender` rejects const senders when `prevChainLink` is empty (analyzer enforces the rule).
- **Gotchas**: Stdlib modules can fail compilation after analyzer rule changes; update `std/*` constants accordingly.

### Session Notes (2026-01-28)

- **Common patterns**: Enabling `govet`'s `fieldalignment` in golangci-lint surfaces many existing struct-order warnings; plan for a staged rollout (only-new-issues) or bulk field reordering.
- **Common patterns**: `betteralign` is best run as a separate make target with `-fix` since it cannot integrate into golangci-lint.
- **Gotchas**: `Select` waits for all `then` array inputs on each `if`; if any `then` slot is only sent conditionally it can deadlock.
- **Common patterns**: golangci-lint v2.5.0 no longer exposes `gosimple`, `stylecheck`, or `tenv` as standalone linters; rely on `staticcheck` instead.
- **Gotchas**: `nolintlint` flags stale `//nolint:golint`; replace with a real doc comment or an active linter name.
- **Common patterns**: `gosec` G115 still flags checked conversions; use a guard plus `// #nosec G115` on the cast line.
- **Common patterns**: `gosec` G306 expects 0600 perms; for readable artifacts, keep 0644 and add `// #nosec G306` with rationale.
- **Architecture insights**: Runtime cancellation now flows via a typed context key helper to avoid raw string keys.
- **Common patterns**: `gocritic`'s commentedOutCode flags comments that resemble code; remove them instead of paraphrasing.
- **Common patterns**: `golangci-lint` thelper prefers `t.Helper()` inside helper constructors used in tests.
- **Common patterns**: `govet` fieldalignment warnings can be suppressed with `//nolint:govet` on struct or inline struct literal lines when reordering is too invasive.
- **Common patterns**: `gocyclo` can be temporarily suppressed with `//nolint:gocyclo` and a short rationale, but consider later refactor.

### Session Notes (2026-01-29)

- **Common patterns**: `copyloopvar` flags redundant loop-var shadowing (`name := name`) in Go 1.22; remove the copy unless needed for older versions.
- **Gotchas**: Runtime now panics on invalid array/case indices (negative or uint8 overflow) to surface compiler bugs early.
- **Common patterns**: Resolve union elements in deterministic key order to avoid nondeterministic error paths in tests.
- **Common patterns**: `golangci-lint` config uses a strict allowlist with targeted test/go:generate exclusions; tune `tparallel`/`gocyclo` per package if flakiness or complexity hotspots appear.
- **Common patterns**: Consider opt-in linters like `depguard`, `gofumpt`, `goimports`, `gochecknoglobals`, `gochecknoinits`, `wrapcheck`, and `forbidigo` based on team appetite; add them gradually to avoid noisy rollouts.
- **Common patterns**: `govet` fieldalignment can be satisfied by moving pointer-heavy fields before strings to reduce pointer bytes and GC scan range.
- **Common patterns**: `depguard` can enforce stdlib-only boundaries with separate rules for `internal/runtime/*.go` and `internal/runtime/funcs/**` when funcs may import `internal/runtime`.
- **Gotchas**: `gochecknoglobals` will flag template strings, expected outputs in tests, and shared helpers; plan exclusions or refactors before enabling widely.
- **Common patterns**: To satisfy `gochecknoglobals`, switch static template vars to `const`, and move shared values into helper funcs; use per-program counters instead of package globals.
- **Common patterns**: If performance trumps `gochecknoglobals`, allow a scoped global counter with `//nolint:gochecknoglobals` rather than threading it through every outport.

### Session Notes (2026-01-31)

- **Language semantics**: Import blocks now allow inline comments; struct literals accept trailing commas.
- **Language semantics**: Const/literal senders are only valid when chained from a triggering signal.
- **Language semantics**: `Main` must have exactly one `start` inport and one `stop` outport, both `any` and non-array.
- **Common patterns**: Use explicit stdlib components (`operators`, `routers.Switch`, `streams.Range`, `Union`) instead of composite sender syntax.
- **Architecture insights**: Debug runtime validation code is generated by the Go backend only when the CLI flag is enabled.
- **Gotchas**: Dotenv loaders are signal-style and do not return dictionaries; use `os.Environ` to read environment values.

### Session Notes (2026-02-01)

- **Gotchas**: Overload resolution can fail when constraints are collected from neighboring overloaded nodes; treat ambiguous neighbor port types as non-constraints to avoid eliminating all candidates.
- **Language semantics**: Array-bypass now uses `[*]` on both sides of a normal connection; index `255` is reserved for the wildcard.
- **Architecture insights**: Array-bypass is handled as a normal connection with a sentinel slot index (`ArrayBypassIdx`) instead of a dedicated AST connection type.
- **Gotchas**: Using `[255]` directly is now a parser error; always use `[*]` for array-bypass.
- **Common patterns**: Stream/list conversions use `streams.FromList<T>` and `lists.FromStream<T>`; `StreamToList`/`ListToStream` were removed.
- **Architecture insights**: Stream helpers live under `std/streams` (map/filter/reduce in dedicated files; `ForEach` in `std/streams/for_each.neva`).
- **Gotchas**: `std/time` only provides delays/durations; there is no wall-clock date/time source, so pass timestamps/durations into components that need "today."
- **Common patterns**: `lists.FromArray<T>` builds a list by chaining `streams.FromArray<T>` into `lists.FromStream<T>`.
- **Common patterns**: README community CTAs use shields.io badges for Open Collective links to keep style consistent.
- **Common patterns**: `e2e/` tests are separate Neva modules with Go e2e tests; `examples/` is a single Neva module where all examples must compile together.
- **Gotchas**: `e2e/` includes verbose mirror variants that can duplicate `examples/` coverage; keep them only when they assert distinct semantics.
- **Architecture insights**: Avoid composite/control-flow syntax to preserve 1:1 visual/text graph mapping; express control flow via explicit stdlib nodes.
- **Common patterns**: Use HOCs and stdlib components for looping/branching since Neva is static and code is not data.
- **Gotchas**: Composite senders collided with chain-only sender/type-system rules; explicit node wiring avoids those conflicts.

### Session Notes (2026-02-02)

- **Common patterns**: `go fix ./...` in Go 1.25 updates build tags to `//go:build` when needed; check for removed legacy `+build` lines in tests.
- **Common patterns**: Prefer adding a `toolchain` line in `go.mod` to pin the Go patch version for consistent local builds.
- **Common patterns**: Align CI `go-version` with the `toolchain` patch level to avoid mismatched builds.
- **Common patterns**: Full `go test ./...` can exceed 5 minutes; rerun with a higher timeout or target slower `e2e/` packages.
- **Common patterns**: In CI, combine `go test -race` with `-count=1` and `-shuffle=on` to surface concurrency flakiness earlier.
- **Common patterns**: Use targeted benchmark jobs (`go test -run=^$ -bench=. -benchmem`) to catch CPU/memory regressions without slowing the main test lane.
- **Gotchas**: `go test ./...` can exceed 5 minutes when e2e packages dominate; split jobs or shard packages to keep CI timeouts stable.
- **Common patterns**: Runtime must remain stdlib-only; use build tags to keep leak-detection helpers/test deps out of runtime packages.
- **Architecture insights**: E2E tests that invoke `go run` exercise compiler + runtime; add `go test`-driven e2e using built binaries for runtime-only coverage.
- **Language semantics**: Connection AST now uses `Connection` directly; `NormalConnection` wrapper was removed (chained/deferred connections reuse the same struct).
- **Architecture insights**: Parser `connDef` now directly parses sender/receiver; `chainedNormConn` points to `connDef` (no intermediate rule).
- **Common patterns**: Array-bypass detection lives in analyzer; desugarer/irgen treat `[*]` connections as validated 1:1 wiring.
- **Gotchas**: Array-bypass consumes the whole port; mixing `[*]` with other slot usage is rejected; index 255 remains reserved (max index 254).
- **Gotchas**: Desugaring anonymous DI args must preserve `ErrGuard`/`OverloadIndex`; dropping `ErrGuard` makes portless senders pick `err`, creating implicit fan-out and flaky emit errors.

### Session Notes (2026-02-04)

- **Architecture insights**: Deferred connections were pure syntax sugar implemented in parser/analyzer/desugarer and lowered to `builtin.Lock`.
- **Common patterns**: Replace `a -> { b -> c }` with explicit lock wiring: `a -> lock:sig`, `b -> lock:data`, `lock:data -> c`.
- **Gotchas**: Deferred syntax usage lived in parser smoke tests and `e2e/simple_fan_out`; both must be migrated when removing the feature.

### Session Notes (2026-02-10)

- **Common patterns**: Visual editor should preserve 1:1 mapping between text and graph; avoid hidden nodes/edges that break this invariant.
- **Common patterns**: IR visualization backends (DOT/Mermaid/ThreeJS) already group ports by component path; reuse this hierarchy for visual editor clustering.
- **UI/UX patterns**: Node editors commonly offer a palette with categorized nodes and wire-splice insertion; Neva should mirror this for discoverability and fast wiring.
- **UI/UX patterns**: Canvas annotations/notes are first-class (n8n-style), useful for comments and documentation in the visual view.

### Session Notes (2026-02-13)

- **Architecture insights**: Source-level model for visualization is `internal/compiler/ast/flowast.go`; older `sourcecode`/`pkg/lsp` path assumptions are stale in this repo state.
- **Architecture insights**: LSP `resolve_file` request/response types exist in `nevalang/neva-lsp/cmd/lsp/server`, but `GetFileView` is not wired yet; treat it as a stabilization task before adding new visual endpoints.
- **Architecture insights**: Current in-repo runtime interception exposes `Sent`/`Received` only; visual debugger overlay hooks should map to these events first.
- **Common patterns**: Keep one canonical planning document for visual-editor architecture and mark older drafts as superseded to avoid conflicting implementation guidance.
- **Common patterns**: Keep exactly one active plan file per task area in `.agent/plans/` to avoid drift between parallel markdown plans.

### Session Notes (2026-02-14)

- **Common patterns**: Keep planning docs in plain Markdown without agent-specific frontmatter so they stay readable and transferable across tools.
- **Common patterns**: In architecture plans, define abstract payload/model names (e.g., `VisualDocument`) in a terminology section to avoid conflating them with LSP method names.
### Session Notes (2026-02-13)

- **Language semantics**: Receiver/sender pairing for chained connections must recurse even when a chain head has a concrete `PortAddr`, otherwise downstream receivers lose constraints.
- **Common patterns**: For selector senders in overload-constraint collection (for example `:state -> .rate -> mul:left`), pass previous chain-link senders into sender-type derivation.
- **Architecture insights**: `deriveNodeConstraintsFromNetwork` needs both direct senders and previous chain-link context to infer selector output types before desugaring.
- **Gotchas**: Without selector-aware constraints, overloaded std operators can appear ambiguous and force typed wrapper hacks in examples.

- **Common patterns**: New LSP features in `nevalang/neva-lsp/cmd/lsp/server` should include short function doc comments plus targeted inline comments for recursive AST traversal/encoding code.
- **Common patterns**: For LSP handlers returning `any`, prefer typed empty results (e.g., empty completion/symbol/location lists) or `false` for `PrepareRename` fallback instead of `nil, nil` to satisfy `nilnil`.
- **Common patterns**: For LSP file lookup failures (`findFile`), propagate the error rather than returning success with nil payload to avoid `nilerr`.
- **Common patterns**: `gosec` G115 in LSP token/range encoding should use explicit int bounds checks plus `// #nosec G115` on checked casts.
- **Gotchas**: `nilnil` can be intentionally valid for LSP nullable results (e.g., hover/rename); use narrowly scoped `//nolint:nilnil` with a protocol rationale when required.
- **Architecture insights**: Source-level model for visualization is `internal/compiler/ast/flowast.go`; older `sourcecode`/`pkg/lsp` path assumptions are stale in this repo state.
- **Architecture insights**: LSP `resolve_file` request/response types exist in `nevalang/neva-lsp/cmd/lsp/server`, but `GetFileView` is not wired yet; treat it as a stabilization task before adding new visual endpoints.
- **Architecture insights**: Current in-repo runtime interception exposes `Sent`/`Received` only; visual debugger overlay hooks should map to these events first.
- **Common patterns**: Keep one canonical planning document for visual-editor architecture and mark older drafts as superseded to avoid conflicting implementation guidance.
- **Common patterns**: Keep exactly one active plan file per task area in `.agent/plans/` to avoid drift between parallel markdown plans.

### Session Notes (2026-02-14)

- **Common patterns**: In overload-constraint collection, treat empty or ambiguous neighbor candidate type sets as "no constraint" and let other edges disambiguate overloads.
- **Common patterns**: Reuse shared helpers for type dedupe (`appendUniqueType`) and unambiguous-single-type checks (`singleUnambiguousType`) to avoid drift across analyzer paths.
- **Architecture insights**: `getPossibleSenderTypes` is best-effort for overload filtering; errors in this path should avoid panics and defer user-facing diagnostics to regular validation.
- **Gotchas**: Selector senders without `prevChainLink` context are invalid and must yield no inferred constraint during overload filtering.
- **Common patterns**: Prefer explicit `(type, ok)` helper return (`singleUnambiguousType`) over slice/nil signaling when representing optional constraints.
- **Common patterns**: Avoid local function aliases in analyzer hot paths when a direct helper call keeps intent clearer (`singleUnambiguousType(...)` directly at call sites).
- **Common patterns**: For LSP package-entity traversal, prefer `src.Package.Entities()` (range-func iterator) over ad-hoc flattening maps.
- **Architecture insights**: Keep shared index snapshot lock access (`getBuild`/`setBuild`) on `Server` to avoid mutex handling drift across feature files.
- **Common patterns**: CodeLens `Data` payload should validate explicit kind enums (`references`/`implementations`) instead of relying on default switch fallbacks.
- **Gotchas**: `go test ./...` can hang or run very long in this repo due broad e2e/examples coverage; prioritize targeted package tests after confirming lint for quick iteration.
- **Language semantics**: For MVP LSP, interface implementation lookup can be approximated structurally from IO port names/array flags/type strings, which is enough for actionable navigation without analyzer-level complexity.
- **Common patterns**: Keep `implementations` codelenses interface-only; component lenses stay on `references` and can include references to interfaces they structurally implement.
- **Architecture insights**: VS Code TextMate grammar (`vscode-neva/syntaxes`) coexists with LSP semantic tokens; no mandatory grammar removal is needed for MVP rollout.
- **Common patterns**: LSP tests are easiest to scale with small in-memory build fixtures plus focused handler-level assertions (`TextDocumentCodeLens`, `CodeLensResolve`) instead of end-to-end editor wiring.
- **Architecture insights**: After PR #1020 merge, core LSP language features are available in `nevalang/neva-lsp/cmd/lsp/server`; visual-editor planning can assume this baseline while treating `resolve_file`/visual endpoints as separate follow-up wiring.

### Session Notes (2026-02-14, AST/Core extraction prep)

- **Architecture insights**: `ast` and `core` are now externalized under `pkg/` (`pkg/ast`, `pkg/core`) so other Go modules (like a split LSP repo) can import them without `internal/` restrictions.
- **Common patterns**: For large package moves, do physical file moves first, then repo-wide import rewrites, then trim formatting-only churn to keep PR review focused.
- **Gotchas**: Broad `gofmt` runs can introduce noisy doc-comment/newline changes in unrelated files; revert those hunks unless they are intentional.


### Session Notes (2026-02-15)

- **Architecture insights**: LSP workspace scanning moved into `pkg/indexer` as an explicit intermediate extraction step for splitting LSP out of the main repository.
- **Common patterns**: Expose tooling-facing errors via a package-local adapter (`indexer.Error`) that wraps `compiler.Error` but stabilizes surfaced fields (`Message`, `Meta`) and behavior (`Error()`, `Unwrap()`).
- **Gotchas**: LSP diagnostics code must guard `Meta == nil` when deriving URIs/ranges from indexer errors to avoid nil-pointer crashes.
- **Architecture insights**: `pkg/indexer.NewDefault` centralizes default parser/builder/analyzer wiring, so LSP binaries can initialize indexing without importing `internal/*`.
- **Architecture insights**: `pkg/typesystem` now re-exports core type-system types for tooling, allowing LSP code to stop importing `internal/compiler/typesystem` directly.

### Session Notes (2026-02-18)

- **Architecture insights**: `pkg/indexer.FullScan` now always analyzes workspace snapshots in library mode (`mainPkgName == ""`) so LSP is independent from runnable-entry assumptions.
- **Common patterns**: Keep workspace indexing semantics library-first in `pkg/indexer`; executable-specific `Main` validation belongs to CLI/compiler executable flow.
- **Gotchas**: Passing executable `MainPkg` into workspace scans can surface `main package not found` and then cascade into LSP `file not found in build`.

### Session Notes (2026-02-14, visual plan decisions)

- **Architecture insights**: Visual projection should be raw-AST-first, with analyzer/desugarer data as optional overlays rather than the base graph source.
- **Common patterns**: Prefer explicit `neva/get*` visual methods over reviving incomplete `resolve_file`/`GetFileView` bridge paths when no compatibility obligations exist.
- **UI/UX patterns**: MVP IDE integration can start as a command-based readonly preview (Markdown-preview style) with component focus/fullscreen, while custom-editor/side-panel/inline variants stay as explicit follow-up decisions.
- **Architecture insights**: Standalone visual app should depend on Neva LSP transport, not editor-specific APIs, so IDE and standalone consume the same graph contract.

### Session Notes (2026-02-18)

- **Architecture insights**: `cmd/lsp` is removed from this repository; canonical LSP source now lives in `nevalang/neva-lsp`.
- **Common patterns**: Before deleting mirrored code, normalize module import paths and diff against the extracted repo to verify functional parity.
- **Gotchas**: Removing `cmd/lsp` requires cleaning `Makefile` `build-lsp*` targets and running `go mod tidy` to drop stale `glsp` module dependencies.

### Session Notes (2026-02-21)

- **Gotchas**: Post-processing JSON via global `strings.ReplaceAll` (`:` / `,`) corrupts string payloads (for example `"a:b,c"` becomes `"a: b, c"`).
- **Common patterns**: Keep readable one-line JSON spacing by scanning marshaled bytes and adding spaces only outside JSON strings (track quote/escape state).
- **Common patterns**: For union string formatting, marshal `data` to JSON first (for correct quoting/escaping) and then apply spacing policy to preserve readability.
- **Language semantics**: Node aliases are required for top-level component node declarations (validated in analyzer); anonymous shorthand is still accepted for DI argument node blocks.
- **Gotchas**: DI argument aliases are semantic (must match dependency node names like `handler`/`predicate`/`reducer`); arbitrary renames can break IR generation.
- **Architecture insights**: Treat `std/builtin` as two logical layers: compiler-lowering primitives (`New/Lock/Field/FanIn/FanOut/Del`) vs user-prelude ergonomics; keeping them in one namespace increases API drift and confusion.
- **Common patterns**: Port naming works best with a two-tier rule: default to `data/res/err/sig`, but keep domain-specific single inport names for boundary APIs (`url`, `filename`) when they carry protocol meaning.
- **Common patterns**: Stream terminal helpers like `streams.Len` are safe additions; slicing/take-style helpers should be specified together with producer-cancellation semantics (issue #666) to avoid misleading performance expectations.
- **Gotchas**: Keep `std/* #extern(...)` names in lockstep with `internal/runtime/funcs/registry.go`; missing mappings silently create dead APIs (current gaps include `map_len`, `string_len`, `string_slice`, `list_slice`, `int_neg`, `float_neg`, and `*_is_*_or_equal` string/float variants).
- **Gotchas**: Renaming stdlib ports can break compiler bootstrap via `internal/compiler/utils/generated`; update `internal/compiler/utils/utils.neva` and regenerate generated exports in the same change.
- **Common patterns**: Regenerate compiler utils with the repository CLI (`go run ../../cmd/neva ...` from `internal/compiler`) to avoid stale global `neva` binaries producing incompatible generated code.
- **Common patterns**: When stdlib API or naming conventions change (for example `For` -> `ForEach`, `data` outport -> `res`), update `docs/style_guide.md` and `docs/qna.md` in the same change to prevent documentation drift.

### Session Notes (2026-02-21, converters brainstorm)

- **Language semantics**: Neva remains strictly explicit for type conversion; operators/binary expressions do not perform implicit coercion.
- **Common patterns**: Collection conversion idiom in stdlib is target-package + `FromSource` naming (`streams.FromList`, `lists.FromStream`, `streams.FromArray`, `lists.FromArray`).
- **Architecture insights**: Runtime conversion funcs preserve sequencing contracts (`stream<T>` emits `idx`/`last`; array-port converters iterate slots in deterministic order).
- **Gotchas**: Conversions with non-total semantics (parsing, narrowing) should expose `err error` and rely on `?` propagation to stay idiomatic.

### Session Notes (2026-02-22)

- **Common patterns**: Prefer documenting conversion policy in `docs/qna.md` as "components in stdlib, not language syntax" to preserve the explicit graph model.
- **Architecture insights**: Treat `builtin` as language/prelude boundary for primitives and compiler-coupled contracts; put policy-heavy converters in regular std packages.
- **Architecture insights**: `Union` is explicitly analyzer-coupled; `Struct` is coupled via `#autoports` flow/desugaring conventions, so both belong to compiler-contract surface.
- **Language semantics**: builtin `Type` (`type Type any`) is a semantic marker for intentionally heterogeneous ports (for example `Union`/`Switch`), not a distinct runtime type.
- **Common patterns**: Keep `strconv` for canonical value conversion/parsing contracts and `fmt` for presentation/I/O formatting, even when both are deterministic.
- **Common patterns**: `docs/qna.md` should capture stable rationale only, not open-task/status notes.
- **Common patterns**: Go-like scalar split: total numeric casts without errors in builtin; text parsing/formatting in `std/strconv`; avoid builtin bool<->number magic.
- **Common patterns**: Keep concrete API candidates separate from `docs/qna.md`; Q&A should contain stable rationale, not implementation drafts.
- **Language semantics**: For strict Go parity, allow builtin `String(int)` as Unicode code-point cast; keep decimal/bool/float text formatting in `std/strconv`.
- **Common patterns**: When addressing PR review comments, apply requested changes first, then reply to each comment in GitHub; do not resolve threads unless explicitly requested.
- **Common patterns**: Keep collection converters in `target-package + FromSource` shape (`streams.FromDict`, `dicts.FromStream`) for consistency with existing list/stream APIs.
- **Language semantics**: `dicts.FromStream` uses last-write-wins for duplicate keys by assigning incoming entries into the resulting dict.
- **Gotchas**: `streams.FromDict` iterates Go maps, so output stream order is non-deterministic and should not be asserted directly in tests.

### Session Notes (2026-02-23)

- **Language semantics**: `bytes` is a dedicated builtin binary payload type; `string` is for text semantics.
- **Common patterns**: Keep text/binary boundaries explicit with `bytes.FromString` / `strings.FromBytes` rather than implicit coercions.
- **Architecture insights**: Adding a new payload primitive requires synchronized changes in IR `MsgType`, runtime `Msg` variants, backend type mapping, and stdlib/runtime extern boundaries.
- **Common patterns**: Binary stdlib boundaries (`io`, `http`, `image`) should use `bytes` to avoid repeated `string`/`[]byte` conversions.
- **Gotchas**: Full `go test ./...` can run very long; prioritize targeted validation for touched compiler/runtime packages plus changed `examples/*` and `e2e/*`.

### Session Notes (2026-02-24)

- **Common patterns**: For Go 1.26 migration, keep `go.mod` (`go` + `toolchain`) and CI `actions/setup-go` patch versions aligned to avoid toolchain skew.
- **Common patterns**: Integrate `go fix` as a dedicated CI check (`go fix ./...` + `git diff --exit-code`) so failures clearly signal "run gofix and commit".
- **Gotchas**: `go fix` can apply broad modernizations in one pass; run lint/tests immediately after and fix secondary lints introduced by rewrites.
- **Language semantics**: `streams.FromString` currently emits one `string` per Unicode code point (rune semantics), not bytes or grapheme clusters.
- **Common patterns**: Keep converter behavior contracts in both stdlib signatures and runtime func comments (order/duplication/materialization) so policy is visible at API and implementation layers.
- **Common patterns**: When posting issue updates with code-style identifiers, use `gh ... --body-file` to avoid shell substitution on backticks.
- **Gotchas**: `lists.FromStream` / `dicts.FromStream` remain blocking materializers and can grow memory with long streams; document this explicitly in comments/docs.
- **Language semantics**: `bytes` has no literal/const syntax; use explicit converters (`bytes.FromString`, `strings.FromBytes`) in networks.
- **Common patterns**: For `*.neva` edits, enforce `lower_snake_case` node aliases and verify style against `docs/style_guide.md`.
- **Architecture insights**: `BytesMsg` is represented as `[]byte`; avoid defensive copies in runtime hot paths and treat immutability as a usage convention.
- **Gotchas**: New e2e module manifests should use explicit language version (`0.34.0`), not shorthand aliases.
### Session Notes (2026-02-24, numeric/bytes direction)

- **Language semantics**: For numeric expansion discussions, prefer Go-style names (`int8..int64`, `uint8..uint64`, `float32/float64`) over Rust-style (`i8/u8/f32`).
- **Language semantics**: Keep ergonomic `int`/`float` in user-facing APIs even if fixed-width families are introduced.
- **Language semantics**: `byte` should remain an alias of `uint8`; avoid architecture-dependent `uint` semantics unless explicitly fixed and documented.
- **Architecture insights**: Numeric-width gains are often secondary to message/container representation costs; plan #904 together with #28 (`bytes`) to realize practical low-level performance wins.
### Session Notes (2026-02-25)

- **Common patterns**: `pkg/e2e.Run` should always execute commands with an explicit per-run timeout (default 30s) rather than relying only on global `go test -timeout`.
- **Architecture insights**: For e2e command chains (`*.test -> neva -> generated output`), cancellation must target the whole process group on Unix (`Setpgid` + group kill) to avoid orphaned child processes.
- **Gotchas**: Test-runner interruption/timeouts can leave `neva_run_*/output` descendants alive, which users may report as “zombies” and observe as sustained CPU heat on macOS.
## 3. ⚡ Core Concepts

- **Dataflow**: Programs are graphs. Nodes process data; edges transport it.
- **Implicit Parallelism**: Every node runs in parallel.
- **Type System**: Static, structural, with generics and tagged-unions.
- **Visibility**: Entities are private by default. Export with `pub`.
- **Entities**:
  - **Components**: Logic containers, node blueprints.
  - **Interfaces**: Port definitions, abstract components.
  - **Types**: Type definitions.
  - **Constants**: Fixed values.
- **Program Hierarchy**:
  - **Module**: Root unit (has `neva.yml`).
  - **Package**: Directory with `*.neva` files.
  - **Component**: The building block.

## 4. 🧠 Architecture

### Compiler (`internal/compiler/`)

1. **Parser**: ANTLR4 -> Raw AST.
2. **Analyzer**: Semantics & Type checking.
3. **Desugaring**: Syntactic sugar -> Canonical AST.
4. **IR Gen**: Canonical AST -> IR.
5. **Backend**: IR -> Target Code.

### Runtime (`internal/runtime/`)

The runtime is a library embedded into every compiled program.

- **FuncRunner**: Executes "native components" (runtime functions).
- **Connector**: Manages message passing.
- **Extensibility**:
  - **Native Functions**: `internal/runtime/funcs`.
  - **Func Registry**: `internal/runtime/funcs/registry.go`.
  - **Func Interface**: `runtime.Func` & `FuncCreator`.

### Standard Library (`std/`)

The standard library provides components for all programs. Some are implemented in Neva, some use runtime functions written in Go via `#extern`.

## 5. 🛠️ Debugging Tips

**Debug Compiler Output**:

- **IR**: `neva run --target ir <pkg>`
- **Trace**: `neva run --emit-trace <pkg>`
- **Runtime validation**: `neva run --debug-runtime-validation <pkg>` or `neva build --debug-runtime-validation <pkg>` (compiler-only check that prints unconnected senders/receivers to validate runtime wiring; intended for language developers inspecting compiler output)

**Debug the CLI/Compiler**:

- **Logs**: Use `fmt.Printf`, remove before finishing.
- **Tests**: `go test ./...`

## 6. 🗺️ Key Locations

- `cmd/neva/`: CLI Entry point.
- `internal/cli/`: CLI implementation.
- `internal/compiler/parser/neva.g4`: Grammar Definition.
- `e2e/`: End-to-End Tests.
- `examples/`: Example programs.
- `pkg/`: Shared utilities.

## 7. 🎨 Coding Standards

- **Go Idioms**:
  - Comments: Every function should have a short doc comment. If it relates to Neva semantics, include a tiny Neva example when helpful.
  - Use `any` instead of `interface{}`.
  - TD tests: `tests := []struct{ name string ... }`
  - Test case names: `lower_snake_case`
  - KISS: simpler code > complex abstractions
  - Runtime funcs naming: prefer `<role>Msg` for variables of type `runtime.Msg` and explicit names like `streamItemMsg`.
  - Utils: `pkg/` for shared utils (EXCEPT `runtime`)
    - If duplicated in 3+ places, move it to `pkg/` (except `runtime`).

**Workflow**:

1. `make build` (Verify compilation).
2. `golangci-lint run ./...` then `go test ./...` (Verify lint + tests).
3. `make antlr` (Regenerate parser if `.g4` changed).
