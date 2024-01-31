---
Title: Compiler directives
---


# Compiler Directives

Compiler directives are special instructions for compiler. They are not intended to be used on a daily basis by regular user but good nevalang programmer must understand how they work because they are base for how many language features operate.

- `#runtime_func` - directive for _component_ that tells the compiler that given component does not have implementation in source code and instead a runtime function call must be created.
- `#runtime_func_msg` - Directive for _node_ that tells the compiler to insert given message to corresponding runtime function call. Can only be used with nodes that are instantiated with components with `runtime_func` directive.
- `#struct_inports` - tells the compiler that inports for this component are not defined in nevalang source code but instead must be derived from it's type-argument which is the structure. Inports will correspond to structure fields.
