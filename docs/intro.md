# Introduction

_Package_ is a set of _entities_ separated by files all located in the same directory. Name of the directory is the name of the package.

There are special rules dictating how package should be separated into files, enforced by _linter_, but it only serve maintanability purpose. It doesn't matter for _compiler_ whether package is one big file or many small ones. We well discuss these rules later.

There are 4 types of _entity_ package can contain: _type_, _interface_, _constant_ and _component_, each having correspondinging keyword: `type`, `io`, `const` and `comp`.

Each entity can be either _private_ or _public_. Every entity is private by default and could be made public by using `pub` keyword. Private entity means it cannot be used outside of the package while public means it could be _imported_ by using `use` keyword.

Set of packages located in the same directory form a _project_. There's usually one project per git repository, but that's not necessary. Every project must have `proj.yml` in root directory. This file describe wanted compiler version and dependencies.

Every project must contain at least one _executable_ package. Executable package is a package that

# Maintainability guide

## Naming

- Components named in UpperCamelCase with "er" ending like `Printer`. Name must express that component is _doing_ something
- Types are also named in UpperCamelCase, not necessarily with "er" ending. E.g. `MyInt`
- Interfaces are named just like components but with `I` prefix. E.g. `IReader`
- Constants are named in `UPPER_SNAKE_CASE`. E.g. `PI` or `DB_TIMEOUT`
- Nodes are named just like components but in lowerCamelCase. E.g. `printer = Printer()`. Similar to how class instances usually named in OOP languages
- Packages are named in lower_snake_case. E.g. `json` or `build_info`
- Port are named shortly, following _port naming convention_, if possible:
    - `v` for _value_
    - `sig` for _signal_
    - `err` for _error_
    - First letter of the type otherwise: E.g. port with `User` type must me named `u`
    - Port that accepts array or vector must be named like regular port but with doubled letter or with "s" on the end. E.g. `vv`, `uu` or `errs`
    - If following this convention makes code less readable then an exception could be made. Common sense is still there. There could be something wrong with the design though.


Before we discussed that there are _special rules_ about how entities are located in the package. Again, these rules enforced by linter and not compiler, which means _program will work just fine without these_. But one of the main goals of this language is to make programs as maintainable as possible.

- If there's more than 1 entity of the same type in the file, then they must be grouped.
- Order of entities in the file must be: `use, type, io, const, comp`
- Each file must contain one component and entities that are used by that component. Those entities could also be `pub` though. One file should not refer to entities defined in other files in the same package.
- One exception to previous rule is `*.shared.neva` files. There could be one `shared.neva` file or many of them e.g. `interfaces.shared.neva` and `const.shared.neva`
- Entities defined in this _shared_ files could be referenced from any other file in the same package including other shared files. But there must be (it won't compile otherwise) no dependency cycle
-
