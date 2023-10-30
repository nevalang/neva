- our goal is to send node's IOs to the vscode extension
- to do so we want to extend LSP in a way where it "indexes" workdir and, at the time of a request, know all packages, files, entities in a workdir in a parsed and analyzed way
- to do so we need to change the way analyzer and (maybe) compiler works:
    - compiler must be able to avoid ir-generation completely
        - we can achieve this by processing case with irgen==nil
        - or making another method (better)
    - compiler also must be able to process (compile) programs without specified entry point (specified main package)
        - there are 2 compilation modes: compilation of a directory (workspace, module) and compilation of a executable (with specified entry point / main package)
        - there could be many executable packages in a module, even in "library/module" compilation mode analyzer must be able to point to errors that are made in main packages so we can process it at the lsp/vscode level

Other stuff

- maybe we should make compiler more simple to abstract away repo and make it receive raw packages and return bytes or IR or even IR structures
- we can optimize reading from disk stdlib by embedding it into the binary
    - need to think whether that could be the problem for development
    - this needs we gonna recompile compiler every time we make changes to std
- builder algorithm (how to turn directory into program)
    - without 3rd party:
        - 
    - with 3rd party deps:
        - read neva.mod
        - fetch every dep
        - for every dep repeat
        - do until there's no deps (think about other cases where we need to stop)