# About

"99 Bottles of Beer" is a [classical programming task](https://www.99-bottles-of-beer.net) that requires looping, conditionals and IO.

This example solves this task with streams. We generate stream of numbers from `99` to `1` and process each, one by one

<!-- ## Implementation Details

This simple example contains something very important to understand about concurrency.

These two components `PrintFirst` and `PrintSecond` are not concurrency-safe themselves, because branches of `switch` are concurrent. Race-condition can happen if next `data` arrive before we successfully sent `res` (or `err`):

```neva
def PrintFirst(data int) (res any, err error) {
	p1, p2, p3 fmt.Println
	format Format
	---
	:data -> switch {
		0: 'No more bottles of beer on the wall, no more bottles of beer.' -> p1
		1: '1 bottle of beer on the wall, 1 bottle of beer.' -> p2
		_ -> format -> p3
	}
	[p1, p2, p3] -> :res
}

def PrintSecond(data int) (res any, err error) {
	p1, p2, p3, p4 fmt.Println<any>
	format Format
	---
	:data -> switch {
		-1: 'Go to the store and buy some more, 99 bottles of beer on the wall.' -> p1
		0: 'Take one down and pass it around, no more bottles of beer on the wall.\n' -> p2 
		1: 'Take one down and pass it around, 1 bottle of beer on the wall.\n' -> p3
		_: -> format -> p4
	}
	[p1, p2, p3, p4] -> :res
}
```

> To avoid race-conditions, we _must_ process each element one by one
>
> TIP: if you have heavy stream-processors, it's better to split operation in several steps. This way you'll be able to utilize more resources, because more nodes will be able to run in parallel. But this is of cource depends on - if parent component can receive next values, until previous ones are processed. For example, if your parent is wrapped into `streams.For`, like here, then it won't help. In other words, separating `PrintFirst` into 2 sub-nodes (e.g. `get_string` and `print_string`) won't do anything.
>
> Pro-TIP: readability is usually more important than performance, so always understand performance, yet focus on readability, until performance issues are faced. TLDR: do not pre-optimize.

Thankfully both of them are part of the `PrintLines` parent-component, that is used by its parent `Main` like this: `streams.For<int>{PrintLines}`.

```neva
def PrintLines(data int) (res any, err error) {
	print_first sync.Tap<int>{PrintFirst}?
	dec math.Decrement<int>
	print_second PrintSecond
	---
	:data -> print_first -> dec -> print_second -> :res
}
```

`PrintLines` is also not concurrency-safe because if new `data` comes before we send previous `res/err` we might get race conditions e.g. `print_first, print_first, print_second`.

> This is okay when we deal with pure components, that doesn't do side-effects. But here we do printing, so order is critical.

```neva
def Main(start any) (stop any) {
	print_lines streams.For<int>{PrintLines}
	range streams.Range
	wait streams.Wait<any>
	---
	:start -> [
		99 -> range:from,
		-1 -> range:to
	]
	range -> print_lines
	print_lines:res -> wait -> :stop
	print_lines:err -> panic
}
```

We're in luch because `PrintLines` is wrapped into `streams.For`. Component `For` guarantees that its dependency will never receive new message, while previous one wasn't fully processed (while its result or error wasn't successfully received). Of course except for the first stream element, because nothing comes before it.

## Even Better

```neva
import {
    fmt
	math
    runtime
    streams
	sync
}

def Main(start any) (stop any) {
	print_lines streams.Map<int>{MapLines}
   println fmt.Println
	range streams.Range
	print_lines streams.Map<int>{MapSecond}
	wait streams.Wait<any>
	---
	:start -> [
		99 -> range:from,
		-1 -> range:to
	]
	range -> map_lines
	map_lines:err -> panic
   map_lines:res -> println -> wait -> :stop
}

def MapLines(data int) (res string, err error) {
   map_first sync.Enrich<int, string>{MapFirst}?
   dec operators.Decrement
   map_second MapSecond?
   ---
   :data -> map_first
   map_first:data -> dec -> map_second
   '${map_first:res}\n${map_second}' -> :res
}

...
``` -->
