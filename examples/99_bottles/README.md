# About

"99 Bottles of Beer" is a classical programming task that utilizes loops, conditions and io. You can see the details [here](https://www.99-bottles-of-beer.net).

## Implementation

1. It seems obvious to use `range` with `for`, but without topology-level loop we will have concurrency at the level of `Next2Lines`, even though it's a pipeline. See [github issue](https://github.com/nevalang/neva/issues/754) for details. That's why we used `While` instead, it was implemented exactly to solve this problem with truly sequential looping.
2. Because of use of `While` we need to implement our `FirstLine` and `SecondLine` in a way that they receive data, perform side-effect and _then_ pass that data further downstream. However, because of impossibility to same sender twice, we would have to use explicit locks. That's why we use `Relay` - it's a HOC that does that for us under the hood.
