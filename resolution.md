```
type mapAlias<t,y> = map<t,y>
type intMap<t> = y<t, int>
type myStr = str

const intMapMsg intMap<myStr> = {'': 0}

io FooIo = ...
component Foo<t> {...}
```

The process of resolving is this
- `intMap<myStr> -> map<str, int>` 
- `Foo<intMap<myStr>> -> Foo<map<str, int>>`

- Take entity
    - 