# Proposal for new stdlib operations

This is a proposal for several APIs that could be added to stdlib.

## `builtin` candidates

|Entity|Signature|Notes|
|---|---|---|
|`Not`|`Not(x bool) (x bool)`|Return logical not operation applied to input `x`|
|`Pair<A, B>`|`struct { first A; second B }`|A generic pair of items. Can stand in for other pair types we've introduced such as `ZipResult`|
|`OrDefault<T>`|`OrDefault<T>(t T)`|Try to return `t` or `default` if `t` is unavailable. `default` can be injected with `IInjectable<T>`|
|`Unzip<T, R>`|`Unzip<T, R>(seq stream<ZipResult<T, R>>) (first stream<T>, second stream<R>)`|The inverse of `Zip<T, R>` which is already in `builtin`|

## DI Interfaces

Inject these no-inport interfaces rather than having it in the signature.

If lambda components are ever supported, we can write this more nicely.

```
// Before
rand.Bool
// ...
0.3 -> bool:p
bool -> :out

// After
rand.Bool{() -> 0.3}
// ...
bool -> :out
```

There are multiple benefits to this pattern:

* Allows more components to satisfy simpler interfaces, while still retaining configurability.
* Moving configurability to the compilation step instead of runtime.
* Remove ambiguity of ports that are only "skimmed" that is read from once at component initialization and then never again. i.e. any port that would be skimmed should be an injectable.

|Entity|Signature|Notes|
|---|---|---|
|`IBool`|`() -> (bool)`|An interface for components that return a `bool`|
|`IInt`|`() -> (int)`|An interface for components that return a `int`|
|`IFloat`|`() -> (float)`|An interface for components that return a `float`|
|`IString`|`() -> (string)`|An interface for components that return a `string`|
|`IInjectable<T>`|`() -> (T)`|A generic interface for components that return a `T`|

## Map

|Entity|Signature|Notes|
|---|---|---|
|`IBijectionMapper<T, R>`|`(T) -> (R)`|Interface for mapping `T` to exactly one `R`|
|`IMapper<T, R>`|`(T) -> (stream<R>)`|Interface for mapping `T` to one or more `R`; Use `IBijectionMapper<T, R>` to make a one-to-one mapping|
|`Map<T, R>`|`Map(seq stream<T>) (stream<R>)`|Map `seq` using the `map`|
|`MapBiject<T, R>`|`MapBiject(seq stream<T>) (stream<R>)`|Map `seq` using the `map`. `seq` OUTPORT has the same number of elements|
|`ToFloat`|`ToFloat(x int) (x float)`|Convert `x` to float. Implements `IBijectionMapper<int, float>`|
|`ToInt`|`ToInt(x float) (x int)`|Convert `x` to float. Implements `IBijectionMapper<float, int>`|

## `rand`

A new stdlib package for pseudo-randomness.

|Entity|Signature|Notes|
|---|---|---|
|`Bool`|`Bool() (x bool)`|Returns random bools. To bias the outputs use `p` (chance of `true`). Can inject `p` with an `IFloat`|
|`Exp`|`Exp() (x float)`|Returns random exponentially distributed floats with lambda 1. Use `Div` to scale by a new rate. Can inject lambda with an `IFloat`|
|`Int`|`Int() (x int)`|Returns random unscaled uniform ints|
|`Intn`|`Intn(n int) (x int, err error)`|Returns uniform ints between `[0, n)`|
|`Norm`|`Norm() (x float)`|Returns normal floats with mean 0 and stddev 1. Can inject mean, stddev as `IFloat`s|
|`Sig`|`Sig() (sig any)`|Produce Poisson signals at the average rate of `1/lambda` per second. Author note: I don't have a use case but think this is cool. Can inject the rate parameter as an `IFloat`|

## Reduce

|Entity|Signature|Notes|
|---|---|---|
|`Add<T>`|`Add<T>(a T, b T) T`|Returns `a+b`. Iplements `IFolder<T, T>`. We need similar `IFolder` implementations for the other binary math ops|
|`IPredicate<T>`||Boolean predicate over `T`. Alias for `IReducer<T, bool>`|
|`IReducer<T, R>`|`(stream<T>) -> (R)`| Reduce stream of `T` to `R`|
|`IFolder<T, R>`|`(T, T) -> (R)`|Fold two values into a single one|
|`FoldLeft<T, R>`|`FoldLeft<T, R>(seq stream<T>, init R) R`|Combine values starting from the first element and `init`|
|`FoldRight<T, R>`|`FoldRight<T, R>(seq stream<T>, init R) R`|Combine values starting from the last element and `init`|
|`Group<T, R>`|`Group<T>(seq stream<T>) (map<R, stream<T>>)`|Group sequence items using a `IBijectionMapper<T, R>`|
|`Partition<T>`|`Partition<T>(seq stream<T>) (inner stream<T>, outer stream<T>)`|Uses an `IPredicate<T>` to partition the stream into inner and outer components|
|`Reduce<T, R>`|`Reduce<T, R>(seq stream<T>) R`|Reduce `seq` using `reduce`|
|`ReducePartial<T, R>`|`Map`|alias for `Map`|
|`First<T>`|`First<T>(a T, b T) T`|Returns `a`. Implements `IFolder<T, T>`|
|`Second<T>`|`Second<T>(a T, b T) T`|Returns `b`. Implements `IFolder<T, T>`|
|`Max<T>`|`Max<T>(seq stream<T>) T`|Implements `IReducer<T, T>`. `T` should be comparable, like `int` or `float`|
|`Min<T>`|`Min<T>(seq stream<T>) T`|Implements `IReducer<T, T>`. `T` should be comparable, like `int` or `float`|

## `streams`

|Entity|Signature|Notes|
|---|---|---|
|`Chain<T>`|`Chain<T>(seq [stream<T>]) stream<T>`|Exhausts input streams in order and outputs a combined stream|
|`Collate<T>`|`Collate<T>(first stream<T>, second stream<T>) (stream<T>)`|Runs collate on the input sequences. See `ICollator<T>`|
|`Drop<T>`|`Drop<T>(n int, seq stream<T>) (stream<T>)`|Discard the first `n` elements from `seq`|
|`ICollator<T>`|`(stream<T>, stream<T>) -> (stream<T>)`|Interface for components which combine multiple streams. This is the "collate" operation from the FBP book.|
|`Head<T>`|`Head<T>(seq stream<T>) T`|Return the first element from `seq` same as `Take(1, seq)`|
|`Flatten<T>`|`Flatten<T>(seq stream<stream<T>>) (stream<T>)`|Flatten a stream of streams into a single stream|
|`FromElem<T>`|`FromElem<T>(elem T) stream<T>`|Return a stream of length one from the first `elem`|
|`Forever<T>`|`Forever<T>(elem T) stream<T>`|Return a infinite stream from elements of `elem`| 
|`RepeatElem<T>`|`RepeatElem<T>(n int, elem T)`|Returns a stream repeating the first `elem` `n` times|
|`Tail<T>`|`Tail<T>(seq stream<T>) (stream<T>)`|All but the first element of `seq`. Same as `Drop(1, seq)`|
|`Take<T>`|`Take<T>(n int, seq stream<T>) (stream<T>)`|Return the first `n` elements from `seq`|

## `time`

|Entity|Signature|Notes|
|---|---|---|
|`Duration`|`Duration(duration string) (nsec int, err error)`|Helper for making durations in humanized format (e.g. `"1hr6s"`)|
|`Deadline<T>`|`Deadline<T>(nsec int, t T) (T, err error)`|If waiting more than the given time to receive an item, return `ErrDeadline`|
|`ErrDeadline`|`error{"Deadline exceeded"}`|A error returned from `Deadline<T>`|
|`ErrThrottled`|`error{"Throttled"}`|A error returned from `RateLimit<T>`|
|`ErrBackoff`|`error{"Backoff"}`|A error returned from `Backoff<T>`|
|`Every<T>`|`Every<T>(nsec int, t T) (t T)`|Send `t` at most once per duration|
|`RateLimit<T>|`RateLimit<T>(nsec int, t T) (t T, err error)`|Try to return `t` but no more than once per duration or return `ErrThrottled`|
|`Backoff<T>`|`Backoff<T>(nsec int, exp float, t T) (t T, err error)`|Try to return `t` or do exponential backoff if `t` is unavailable and return `ErrBackoff`. Backoff resets when returning `t`|