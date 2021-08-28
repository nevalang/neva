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
set dep.plus +
set workers.a plus

connect in.x a.in[0] a.in[1]

show

show io
show io.in
show io.out

show deps
show deps.plus
show deps.plus.in
show deps.plus.out
show deps.plus.in.nums
show deps.plus.out.sum

show workers
show workers.a
show.workers.a.in
show.workers.a.out
show deps.workers.a.in.nums
show deps.workers.a.out.sum
```
