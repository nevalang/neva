# Module versioning

Suppose there was fixes in the module:

- no changes in ports: patch
- add new outport: minor (user _can_ use it)
- add new inport / remove old import / rename old inport - major change (user _have to_ make changes)

# Package versioning

- no changes in modules: patch
- new exported module added: minor
- old module removed or its api changed (see previous section) - major

# Tooling

There should be a tool to check backwards compatibility and generate versions.
Idea of using git for that may not be good because there is no way to prevent user from uploading broken stuff.
Only special repo-service could give that ability.
