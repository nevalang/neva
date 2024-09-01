# Program Structure

```yaml
build:
  modules:
    manifest:
      language_version: string

    packages:
      files:
        imports: {}
        entities: {}
```

## Build

Build is the most high level abstraction and it's the first compiler stage. Building (the program) means collecting all the source code together in one place, so it's possible to work with it as with single thing.

Compiler downloads (those that are not cached) all the needed _dependencies_ and bundle them together into a object, that it can analyze. It is known at compilation time which _module_ in the build is the _main module_.

## Module

Module is a set of _packages_ on a file-system that have single root directory with `neva.yml` file that called _manifest_ or _manifesto_. Manifest defines: minimal language version it supports and all the _3rd-party_ modules it depends on. Each dependency has _path_ and _version_. When compiler _builds_ the program it recursively downloads each module with all its dependencies.

### Main Module

Main module is the one that contains _main package_ with _main component_. In other words it's the entry point for compiler to generate a program.

## Package

Package consist of one or more files. Files don't affect the visibility scope of the entities. Each entity defined in one file is visible from another file in the same package. In terms of the file-system package is set of files in a single directory. Each package has unique path where module (project) is the root. It doesn't matter how you include packages into eachother, all it affects is names of the packages. You can, however, (and should) to express your architecture with your directory structure. It is especially important because your directory structure dictates your namespaces.

### Main Package

Main package is the one that can serve as an entry-point. It must contain _component_ `Main` with special signature and no exported entities. As any other package it can consist of any amount of files. It's ok for any module in a Build to contain main package, but only one main package is actually used as an entry point.

## File

File with `.neva` extension that contains entities and is a part of the package. File consist of imports and entities. Imports are special instructions that allow to depend on other entities from other packages. All imports are explicit except imports form `builtin` package which is imported implicitly to each file at desugaring stage.

## Imports

### Path and Alias

Import statement consist of _path_ and optional _alias_. Alias means how we refer that package in source code. E.g. `import { foo/bar }` imports `foo/bar` package and assigns `bar` alias to it. That is, by default last part of the path equals alias. If you want custom alias, which is only allowed if you have imports with duplicated aliases like `foo/bar, baz/bar`.

There are 3 types of imports: _std_, _3rd-party_ and _local_. Std looks like `strings` or `net/http` - just pkg name. 3d-party looks like `<mod_name>/<pkg_name>` where `mod_name` is the name defined in the module's manifest file `neva.yml`. Finally local imports are of the form `@/<pkg_name>` where `@` is the module's root.

### Entity Reference Resolution

When you reference some entity by its name and (optional) package modifier (`Foo` is reference without package modifier and `foo.Bar` is with package modifier), Nevalang's compiler _resolves_ that _entity reference_.

If you refer without modifier like `Foo` then it first looks for Foo entity inside the package - this or other files in that directory.

**Each file contains implicit import** of `builtin`.
Note that so if file is not found in local package, compiler tries to resolve it like `builtin.Foo`. This is how compiler resolves builtin references without package modifier like `int`.

## Entity

Entity is very important abstraction and the most important entity is _Component_. Nevalang is 100% declarative so all you actually do is defining entities and their relations. Entities is what forms a Nevalang package. Nevalang program is essentially a bunch of entities. **Each entity has unique name across package namespace**.

There are 4 _kinds_ of entities:

1. Types - message shape definition
2. Constants - reusable messages with static values
3. Interfaces - abstract components
4. Components - computation units

### Visibility Scope

Entities could be public or private what means - can they be used (imported) outside of the package or not. Keyword `pub` is used to denote that. `Foo` is private entity and `pub Foo` is public entity.
