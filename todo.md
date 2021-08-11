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
