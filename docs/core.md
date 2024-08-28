# Core Concepts

## Build

User usually not think in terms of build, but building is the first compiler stage. Compiler downloads (if not cached) all the needed dependencies and bundle them together into a object, that compiler can then analyze. It is known at compilation step which module in the build is the _main module_.

## Module

Module is a set of packages on a file-system that have single root directory with `neva.yml` file that called manifest or manifesto. Module's manifest defines: minimal language version it depends on and all the 3rd-party dependencies. Each dependency has path and version. If there are several deps with the same name we can add alias to them.

### Main Module

Main module is the one that contains main package with main component. In other words it's the entry point for compiler to generate a program.

## Package

Package consist of one or more files. Files don't affect the visibility scope of the entities. Each entity defined in one file is visible from another file in the same package. In terms of the file-system package is set of files in a single directory. Each package has unique path where module (project) is the root. It doesn't matter how you include packages into eachother, all it affects is names of the packages. You can, however, (and should) to express your architecture with your directory structure. It is especially important because your directory structure dictates your namespaces.

### Main Package

Main package is the one that can serve as an entry-point. It must contain component `Main` with special signature and no exported entities. As any other package it can consist of any amount of files. It's ok for any module in a build to contain main package, but only one main pkg is actually used in the build. Other name for main package could be "executable" package.

## File

File with `.neva` extension that contains entities and is a part of the package. File consist of imports and entities. Imports are special instructions that allow to depend on other entities from other packages. All imports are explicit except imports form `builtin` package which is imported implicitly to each file at desugaring stage.

## Imports

Each file can contain implicit import of `builtin` and also can contain explicit imports added by user.

There are 3 types of imports: std, 3rd-party and local. Std looks like `strings` or `net/http` - just pkg name. 3d-party looks like `<mod_name>/<pkg_name>` where `mod_name` is the name defined in the module's manifest file `neva.yml`. Finally local imports are of the form `@/<pkg_name>` where `@` is the module's root.

## Entity

Nevalang package is set of entities of 4 types: _types, constants, interfaces and components_. Each entity has unique name across package namespace. Entities could be public or private which mean if they can be used outside of the package (imported). Nevalang is 100% declarative language so all you can do is define entities with some relations.
