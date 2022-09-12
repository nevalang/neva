# Memory conclusions

- pointer: 8 bytes
- interface: 16 bytes
- slice: 24 bytes
- string: 16 bytes
- struct: sum of fields

# Chan info

- Safe to receive from closed chan
- Unsafe to send to closed chan (panic)

# Program structure

```yaml
project:
  - space.yml: (compiler version, other spaces as deps)
  - user_domain (package):
      - pkg.yml (imports, pkg-wide types and values)
      - get_age.src.yml
      - get_last_name.src.yml
  - usecase:
      - pkg.yml
  - port:
      grpc:
        - handlers.src.yml

# === space.yml ===

compilerVersion: <semver_str>
deps:
  shared:
    path: <path> # is std tracked here?
    version: <semver_str>

# === pkg.yml ===

# ... (imports?, exported)

# === <name>.src.yml ===

modules:
  <name>:
    genericParams: [<generic_name>, ...] # move to io?

    io:
      in:
        <name>:
          isArray: bool
          msgType:
            typeName: <type_name>
            genericArgs:
              - typeName: <type_name>
                genericArgs: [<generic_name>, int]
      out: ...

    nodes:
      structBuilders:
        <name>:
          msgType: ... # must be struct

      consts:
        <name>:
          msgType: ...
          value:
            structValue: ...

        deps:
            interfaces:
                <name>:
                    scope:

      workers:
        <name>:
          refType: component | interface
           componentRef:
                scope: pkg-wide | space-wide
                componentName: ""
                spaceWide:
                    type: local | thirdParty
                    thirdPartySpace: ""
                    pkgName: ""

          # interfaceRef: ???

    net:
      in:
        <name>:
          - nodeType: worker
            name: <worker_name>
            meta:
              connectionType: normal | structReading
              structPath: [foo, bar] # value.foo.bar
```

=========

- кто инстанцирует компонент с дженериками?
