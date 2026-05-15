# GraphDocument v1 (Read-only Visual Contract)

`GraphDocument v1` is the source-level payload used by standalone visual tooling.

## Scope

- Source-level AST projection (not IR)
- Read-only visualization and navigation
- Stable IDs for future overlays/debug adapters

## Top-level shape

```json
{
  "version": "v1",
  "workspace": {"id": "...", "rootPath": "...", "packageIds": [], "anchor": {}},
  "packages": []
}
```

## Entity coverage

- `WorkspaceGraph` / `PackageGraph` / `FileGraph`
- `ComponentGraph` (nodes/edges/inports/outports)
- `InterfaceGraph`
- `TypeDecl`, `ConstDecl`, `ImportRef`
- `SourceAnchor` on all major elements

## ID stability

IDs are deterministic hashes of source identity tuples (`workspace/module/package/file/entity/...`).

- Repeated projection over unchanged source must produce identical IDs.
- IDs are opaque and must be treated as stable references by clients.

## CLI entrypoint

Use `neva visual [workspace]`.

- Serves read-only API: `/api/graph/workspace`, `/api/graph/file?id=...`, `/api/graph/component?id=...`
- Serves a minimal standalone explorer UI at `/`
