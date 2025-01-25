# Program Structure

Nevalang programs have a simple, strict structure. This document outlines the key elements.

High level schema of the nevalang program:

```
build {
    modules {
        packages {
            files {
                imports
                entities
            }
        }
    }
}
```

## Build

A build is the highest-level abstraction for the compiler, including the entry module and all dependencies. It combines local and remote modules into a single structure for analysis and code generation. The compiler downloads dependencies recursively, checking compatibility based on language versions. Each build contains at least two modules: stdlib and entry. Each module in build has unique reference that used to resolve imports. While users work with modules, the compiler operates on the complete build.

## Module

A set of (packages) (directories with `*.neva` files) and a manifest (`neva.yml` or `neva.yaml`file at the root). Minimal nevalang module with manifest and one package main:

```
project
├── main
│   └── main.neva
└── neva.yaml
```

Module is usually a git-repo but not necessary. Module that isn't published in git cannot be used as a dependency by other modules, but can be used as an entry module (compilation entry point).

### Manifest File

The manifest defines the module's minimum supported language version and dependencies. Here's an example manifest with a dependency on the Nevalang compiler version `0.30.2` and a third-party module:

```yaml
neva: 0.30.2
deps:
  github.com/nevalang/x:
    path: github.com/nevalang/x
    version: 0.0.16
```

The `deps` field is a map where each dependency has an alias. When adding dependencies via CLI (e.g., `neva get github.com/nevalang/x`), the package manager automatically inserts a key-value pair. Third-party dependencies must have a valid git-clone path and a fixed semver version. The package manager uses git to download the repo and looks for the corresponding git-tag. The alias typically defaults to the module's path, but custom aliases allow multiple versions of the same module:

> WIP: CLI tool planned for CI/CD to verify module's backward compatibility

```yaml
neva: 0.30.2
deps:
  github.com/nevalang/x@0-0-12:
    path: github.com/nevalang/x
    version: 0.0.16
  github.com/nevalang/x@0-0-11:
    path: github.com/nevalang/x
    version: 0.0.11
```

Package manager creates aliases automatically, but manual additions are possible. Running `neva build` or `neva run` is sufficient, as the compiler checks for dependencies that need downloading.

> WIP: Future compiler to suggest downloading and adding undefined dependencies to manifest, instead of just throwing errors.

### Module Reference

Module references uniquely identify modules in a build, used by the compiler to resolve imports. It consists of a required path and version. We've seen module references in the manifest file:

```yaml
path: github.com/nevalang/x
version: 0.0.16
```

### Entry Module

Module that serves as the compilation entry point, typically the user's program. It's determined by `neva build` or `neva run`. Compilation begins by reading all files from the entry module, then downloading its dependencies recursively. The result is an executable, so it must contain a main package. Dependency modules may lack main packages if they are libraries.

Entry module is the only module that doesn't used as a dependency by any other module in the build. However, compiler must assign some reference to it to be able to resolve local imports. This is how reference to entry module looks like:

```yaml
path: @
version: ""
```

Also, entry and stdlib modules are only modules without versions in their references.

### Stdlib Module

Each module implicitly depends on the `stdlib` module, which the compiler automatically includes in every build. The compiler contains embedded stdlib files, eliminating the need for downloads. As a result, the stdlib module is the only one without dependencies, while all other modules implicitly depend on it. Stdlib source code is located [here](../std/).

Here's how inserted dependency looks like:

```yaml
std:
  path: std
  version: ""
```

### Manifest Lookup

The module is determined by the location of the manifest file (neva.yaml or neva.yml). For example, with the project structure shown above:

```
project
├── main
│   └── main.neva
└── neva.yaml
```

If we run this command from the root:

```shell
> neva build main
```

The compiler will recognize `neva.yaml` and identify this directory as the entry module. If run from the `main` directory containing `main.neva`:

```shell
> cd main && neva build .
```

The compiler will search for the manifest file in the current directory. If not found, it will look in the parent directory. In this case, it will find the manifest one level up.

After finding the manifest, all directories containing `*.neva` files are part of the module visible to the compiler. Modules can have any depth. Here's an example of a slightly nested module structure:

```
project
├── main
│   └── main.neva
├── foo
│   ├── bar
│   │   └── image.png
│   │   └── baz.neva
│   │   └── bax.neva
│   └── foo.neva
└── neva.yaml
```

All `*.neva` files and the `neva.yaml` manifest are part of the module, excluding non-Nevalang files like `image.png`.

The same logic applies when building dependencies. A manifest file is expected at the root of your repository to allow the compiler to build the module when another module declares it as a dependency.

## Package

A set of `*.neva` files in a single directory. Example:

```
project
├── main
│   └── main.neva
├── foo
│   ├── bar
│   │   └── image.png
│   │   └── baz.neva
│   │   └── bax.neva
│   └── foo.neva
└── neva.yaml
```

Here we see 3 packages: `main`, `foo` and `foo/bar`. Note that `foo/bar` is a separate package from `foo`, despite being located inside it. Packages don't have automatic access to each other's entities; they must import public entities explicitly. Nesting packages can be convenient for semantically related code, but it doesn't affect visibility scope.

### Main Package

The main package must contain a `Main` component with a specific signature and no exports.

```neva
def Main(start any) (stop any) {
    // ...
}
```

The main package is identified by its use as the compilation entry point, not by its name. Different entry points can be specified using the CLI:

```shell
> neva build foo
> neva build foo/bar
```

## File

A `.neva` file contains imports and entities. Files organize packages for readability without their own visibility scope. Entities in one file can be referenced from another within the same package:

```neva
// pkg/foo.neva
pub type myFloat float

// pkg/bar.neva
def Bar(mf myFloat) (sig) {
    // ...
}
```

Here, `myFloat` from `pkg/bar.neva` is used without import as it's defined in `pkg/foo.neva` within the same `pkg` package.

## Imports

To reference entities from other packages, imports are used. Imports are grouped inside curly braces `{}`:

```neva
import {}
```

Single import can be one-liner `import { strings }`. Multiple imports:

```neva
import {
  lists
  strings
  streams
}
```

Each import consists of `[<module_name>:]<package_path>`. The module name and colon are omitted for stdlib imports. Example:

```neva
imports {
  strings // module omitted, package is "strings"
  @:lib // module is "@" (entry), package is "lib"
  github.com/nevalang/x:foo/bar/baz // module is "github.com/nevalang/x" (third-party), package is "foo/bar/baz"
}
```

Each import adds a namespace to the file for accessing imported entities. The namespace corresponds to the last part of the package path. For example, with imports `strings`, `@:lib`, and `github.com/nevalang/x:foo/bar/baz`, we have namespaces `strings`, `lib`, and `baz`. We can reference public entities (prefixed with `pub`) using `<namespace>.<entity_name>`, like `strings.Replace`, `lib.Stuff`, and `baz.Bax`.

You can create a namespace manually using an import alias. Let's say you need to import `github.com/nevalang/x:foo/bar/baz` and your local package `baz`:

```neva
import {
  @:baz
  github.com/nevalang/x:foo/bar/baz
}
```

This would create a namespace collision as both packages creates `baz` namespace. Referencing `baz.Bax` would be ambiguous. The compiler prevents this by throwing an error. To fix, assign an alias to one of the imports:

```neva
import {
  my_baz @:baz
  github.com/nevalang/x:foo/bar/baz
}
```

Now it's clear that `baz.Bax` refers to the `github.com/nevalang/x:foo/bar/baz` package, while `my_baz.Bax` refers to your local `@:baz` package.

Imports are categorized into three types:

- std imports
- local imports
- third-party imports

### Stdlib Imports

Stdlib imports are imports of packages from the `std` module. For stdlib imports, omit the module name and `:` separator. Instead of `import { std:strings }`, use `import { strings }`. The compiler will automatically prefix it with `std`.

#### Builtin Package

Every file implicitly imports the `std:builtin` package. When referring to builtin entities, omit the namespace (e.g., use `Add` instead of `builtin.Add`). The compiler automatically prefixes with `builtin.`. Note that local entities with the same name will shadow builtin ones, but it's best to avoid such situations. Never manually import the builtin package as the compiler does this automatically. To learn what's inside the builtin explore [stdlib source code](../std/builtin/).

### Local Imports

Packages within the entry module can be imported using `@`. Consider this project structure:

```

project
├── main
│ └── main.neva
├── foo
│ ├── bar
│ │ └── image.png
│ │ └── baz.neva
│ │ └── bax.neva
│ └── foo.neva
└── neva.yaml

```

Now let's import entities from `foo` to `foo/bar` and both into `main`:

```neva
// foo/foo.neva
pub const p int = 3.14

// foo/bar/bar.neva
import { @/foo }

pub def AddP(data float) (res float) {
  (:data + $foo.p) -> :res
}

// main/main.neva
import {
  @/foo
  @/foo/bar
}

def Main(start any) (stop any) {
  bar.AddP, Println
  ---
  :start -> { $foo.p -> addP:right -> println -> :stop }
}
```

### Third-Party Imports

Third-party imports are imports of packages located in third-party modules - modules that are defined as a dependency in module manifest file. It's not possible to import transitive dependency - dependency of your dependency. To import third-party package you need to refer to its module the same way you have key in your manifest's `deps` map. Example:

```yaml
deps:
  github.com/nevalang/x:
    path: github.com/nevalang/x
    version: 0.0.16
```

Then when you `import { github.com/nevalang/x }` compiler will know exactly path and version of the module you are referring to.

## Entity

A Nevalang package consists of interrelated entity definitions. There are four types of entities, each with its own documentation:

1. **Types** - message shape definition
2. **Constants** - reusable messages with static values
3. **Interfaces** - definition of IO of abstract components
4. **Components** - computation units

Entities can be public (`pub`) or private, determining their visibility outside the package.

## Entity Reference

Entity references include an optional package name and the entity name. Package names can be omitted for entities in the same package or in `std/builtin`. Local entities with the same name as builtin entities shadow them.
