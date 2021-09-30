# Modules visualisation

Im talking about map of modules


# Imported module could be a program descriptor!
# Imported module could be a program descriptor!
# Imported module could be a program descriptor!

# Pure GENIUS shit

# DATA EDITOR

Data editor is a mind-map-like GUI
that allowes create graph
where one can leads to another

a way to visualize message interface creation


# Module resolver

start

- load deps to map
- load imports to map
- find root name with "find mod" procedure

find mod

- use "check path" to check if path local or remote
- if local use "find local mod"
- otherwise use "find remote mod"

check path

- ...

find local mod

- use local path to find file on a disc
- if it's there return its bytes
- otherwise throw err

find remote mod

- use remote resolver to find remote file
- return bytes or err

# BLACK ADAPTERS MAGIC!!!

Модуль, который динамически создаёт другие модули.
Кейс - динамическое создание адаптера между компонентом, который принимает список
и компонентом, который имеет аррай-портс интерфейс.

При старте такой модуль создаёт аррай портс с кол-ом словов соотв. длине списка.
При получении значения он пишет в этот порт.

# Mock autogen

Every module depends on components via interfaces.
So it should be possible to generate mock modules,
that would allow to program behaviour.

## Motivation

Test is simply a program that uses (e.g. `std/testing`) test utils

## Mock API (go:generate-ish?)

???

```
prog1.yml
prog2.yml
common/
  mod1.yml
prog1/
  mod2.yml
prog2/
  mod3.yml

neva run prog1.yml
```

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

ЭВРИКА!!! (Уже не уверен....)

```
Меня только что осенило, что парадокс "как быть с аррай-инпортами корневых модулей?" наконец-то разрешён.

Проблемы нет - компилятор может понять, сколько каналов создать под такой инпорт (вычислить его size для NodeMeta) на основании сети.

Всё дело в том, как работает референция таких портов в схеме - чтобы сослаться на аррай-порт, нужен индекс.
Таким образом, в сети рутового модуля будут ссылки на индексы. Согласно правилам валидации, между указателей на слоты таких портов
не должно быть "дыр".

Стало быть, максимальный индекс плюс один это есть размер.
```
