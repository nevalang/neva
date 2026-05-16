# Visual View v1 (Internal)

Contract namespace for clients: `neva/view/*`.

## RPC surface (MVP)

- `neva/view/getProgram`
- `neva/view/getFileView`
- `neva/view/resolveEntityRef`

Legacy alias:

- `resolve_file` (deprecated; remove next minor)

## Identity policy

- IDs are readable path-like IDs (`module/package/file/entity/...`).
- Node/entity IDs are semantic-path based.
- Connection IDs use endpoint-signature + chain-path + duplicate-ordinal.
- Deterministic ordering is required in all payloads.

## Current status (implementation)

- `pkg/view` projections are AST-first and deterministic for:
  - modules/packages/files (sorted keys),
  - entities per file (sorted keys),
  - component overloads (sorted overload indices),
  - connections (signature + anchor ordering + duplicate ordinal).
- Node `entityRef` is pre-resolved during projection into canonical references
  (`ResolvedRef`) so `resolveEntityRef` transport can be a direct lookup path.

## Ownership

- `neva`: AST-first projection primitives in `pkg/view`.
- `neva-lsp`: transport/backend for `neva/view/*` and standalone readonly explorer.
- `vscode-neva`: client consuming the same methods.
