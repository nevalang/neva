# Memory conclusions

- pointer: 8 bytes
- interface: 16 bytes
- slice: 24 bytes
- string: 16 bytes
- struct: sum of fields

# Chan info

- Safe to receive from closed chan (default value)
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

---

компилятор пытается сгенерировать программу таким образом
чтобы у каналов были буфферы, релевантные их использованию (продумать)
чтобы рантайму не пришлось подключать механику очередей

однако рантайм способен обслуживать каналы таким образом
чтобы они никогда не блокировались
используя ручные очереди

=== компилятор

выходные порты никогда не имеют буфера
т.к. рантайм вычитывает их в неблокирующем режиме
иными словами, рантайм всегда сможет вычитать сообщение из выходного порта так быстро
как это только возможно 

для входных портов размер буфера равен количеству пишущих в них выходных портов
это предположение, основанное на том, что чем у получателя больше отправителей
тем быстрее может наполняться его буфер

конечно, реальность сложнее
один быстрый отправитель может заполнить буфер быстрее чем несколько медленных
однако, лучше иметь буфер, чем не иметь его
а его размер должен быть как-то обоснован

рантайм, однако, устроен таким образом
что при переполнении буфера подключаются "ручные очереди"
которыми он управляет сам
благодаря чему порты, которые способны принять сообщение, будут тут же принимать его
не будучи никак заблокированы операциями, с ними не связанными

=== рантайм

от отправителя пришло сообщение
это сообщение надо передать всем 
надо сделать это так, чтобы передача сообщения одному получателю не была заблокирована другой предачей
для этого рантайм использует очереди

=== 