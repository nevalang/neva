# Formatter

`neva fmt` renders valid Neva source in the canonical style. It is part of the
Neva toolchain and has no configuration: the formatter version always matches
the parser version installed with Neva.

```text
neva fmt [flags] [file or directory ...]
```

With no path, `neva fmt` reads one source file from standard input and writes
the formatted source to standard output. With file paths, it writes formatted
source to standard output in argument order. A directory is walked recursively
in lexicographic path order. The walk skips `.git`, `.neva`, `node_modules`,
and `vendor` directories; an explicit file argument must have a `.neva`
extension.

## Modes

Only one output mode may be selected.

| Flag | Behavior |
| --- | --- |
| `-w` | Atomically replace each changed source file, preserving its permissions. |
| `-d` | Write a unified diff for each changed file. |
| `-l` | List each changed file. |
| `-check` | List each changed file and exit with status 1 when any need formatting. |
| `-e` | Continue after errors and report every independent file error. |

Without `-w`, `-d`, `-l`, or `-check`, formatting never changes a file.

Examples:

```sh
# Format an editor buffer.
neva fmt < main.neva

# Rewrite all source files below the current directory.
neva fmt -w .

# Inspect the proposed changes.
neva fmt -d ./src

# Enforce canonical formatting in CI.
neva fmt -check .
```

## Scope

The formatter is file-local and syntax-only. It does not resolve modules,
type-check, build a program, or rewrite its meaning. It rejects invalid source
and never writes that input.

It enforces the local layout rules in the [Style Guide](./style_guide.md):
indentation, whitespace, imports, canonical safe line wrapping, comma-separated
forms, and struct/union declaration blocks. Rules that need semantic knowledge,
such as naming or whether a port name may be omitted, belong to `neva lint` or
an explicit refactoring.

## Exit status

- `0`: every input was processed successfully; with `-check`, all were already
  formatted.
- `1`: a source or I/O error occurred, or `-check` found unformatted source.
- `2`: command-line invocation is invalid.
