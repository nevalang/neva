# About

"99 Bottles of Beer" is a classical programming task that utilizes loops, conditions and io. You can see the details [here](https://www.99-bottles-of-beer.net).

## Implementation Details

1. It seems obvious to use `range` with `for`, but without topology-level loop we will have concurrency at the level of `Next2Lines`, even though it's a pipeline. See [github issue](https://github.com/nevalang/neva/issues/754) for details. That's why we used `While` instead, it was implemented exactly to solve this problem with truly sequential looping.
2. Because of use of `While` we need to implement our `FirstLine` and `SecondLine` in a way that they pass their input further _after_ they finish their job. However, because of impossibility to reuse same sender, we would have to use explicit `Lock` or `Pass`. That's why we use `Relay` - it's a HOC that does that for us.
3. `FirstLine` and `SecondLine` implemented in a way that it would be possible to have race, but thanks to `While`, it's not. Without sequential iteration guarantee at the level of `Main` (if we dould have concurrency there) - it's possible to see prints in different orders (e.g. we send `99` to print, then `98` to print, but `98` is printed before `99`, etc.)
