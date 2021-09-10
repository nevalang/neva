# Debugger

## Undo/Redo

Instead of keeping log of all sends/receives, keep only previous values.

## Editing of messages values

## Live changing of networks

- Обмен байткодом между рантаймами
- DEBUGGER (Обёртка над компилятором и рантаймом, в рантайме, вероятно, мидлварь)
- FBP SHELL
- mocker?
- type system (типы должны быть максимально просты и совместимы с `gRPC`, `graphQL` и `json schema`)
- Close all the ports when there are no senders to receive a message from.

pkg.yml

```yaml
import:
  pow:
    repo: https://github.com/emil14/pow
    version: 4.20.69

use:
  +: operators
  ^: pow
```

repl:

```
set in.x int
set out.y int

set deps.*.in.nums []int
set deps.*.out.mul int

set workers.multi *

connect in.x multi.in[0]
connect in.x multi.in[1]
connect multi.out.mul out.x

get in
get out
get in.x
get out.y

get deps
get deps.plus
get deps.plus.in
get deps.plus.in.nums

get workers
get workers.multi
get.workers.multi.in
get.workers.multi.out
get deps.workers.multi.in.nums
get deps.workers.multi.out.sum

get env[deps.plus]

```
