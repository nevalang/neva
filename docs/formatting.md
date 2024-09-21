# General

Don't try to be as short as possible (and don't try to be as verbose as possible). Seek balance instead. Readability is king. Code that is understandable at a glance is good code. There are multiple ways to write the same thing, sometimes shorter, sometimes longer. You should choose the most appropriate form for each specific usecase.

# Nodes

write nodes in one line if

- < 80 chars
- no aliases

good

```neva
flow Main(start) (stop) {
    Foo, Bar, Baz, Bax
    ---
    :start -> foo -> bar -> baz -> bax -> :stop
}
```

bad

```neva
flow Main(start) (stop) {
    f Foo, Bar, Baz, Bax
    ---
    :start -> foo -> bar -> baz -> bax -> :stop
}
```
