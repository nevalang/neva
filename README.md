# builtin types

```
/ / non-generic types
bool - true and false
int - signed 64 bit integer
float - signed 64 bit float
str - immutable string

// generic types
list<T> - dynamic array of elements T
dict<T> - map with string keys and T values
maybe<T> = some<T, null>
struct<T...> - restricted multi-type dict
```

## struct example

all keys in such dict must be present and contain values of compatible types

```
user<T, Y> {
    name str
    age int
    friends []user // list<user>

    something T
    somethingElse Y
}
```

## maybe type

```
node<T> {
    parent node? // maybe<node>
    value T
}
```

# Interfaces

deps {
direct {
http net.http
}
indirect {
http2 net.httpInterface
asIs<T> {
in {
data T
}
out {
data T
}
}
}
}

nodes {
effects {
funcs {
print
read
}
}
interfaces {

    }
    components {

    }

}

рекурсивные типы допустимы благодаря отсутствию констрейнтов

type expr resolution

- найти type ref в таблице типов
- сравнить количество generics из expr с type
- если не совпадают - вернуть ошибку
- проделать рекурсивно для каждого type expr в args

```
// list and dict are base types

a<t> = list<dict<t>>
b<t> = dict<list<t>>

a<b<int>>

===

a<b<int>> = list<dict<b<int>>>

list<dict<b<int>>> = list<dict<dict<list<int>>>>

===

Resolved Type Expr - all typerefs points to base types

list<
    dict<
        dict<
            list<int>
        >
    >
>

===

obj<t> = { l list<t> }

d = dict<obj<int>>

---

dict<obj<int>> = dict<{l list<t>}>
```

base type это bool, int, float, str, list, dict и любой тип, который ведёт на struct
сам тип struct является абстрактным и допускается к использованию напрямую, без алиаса
даже структурам без дженериков нужны алиасы, чтобы к ним можно было как-то обратиться

```
root: root

types: {}

messages: {}

interfaces:
    pkgWideInterface:
        params: [t]
        io:
            in:
                sig: t
            out: {}

components:
    root:
        interface:
            params: [t]
            io:
                in:
                    sig: t
                out: {}

        di:
            printer: std.io.printer
            local: pkgWideInterface

        config: {}

        nodes:
            print:
                di: printer
            local:
                di: local
            userBuilder:
                component: std.builder
                config:
                    name:
                        type: str
                        strValue: ""
            greetGiver:
                component: std.giver
                args: [int]
                config:
                    v:
                        type: int
                        intValue: 42

        net:
            in.sig:
                - di.printer

    exports:
        globalInterface:
            entityType: "interface"
            localName: pkgWideInterface
```

```
compiler {
    pkgs { [string]pkg }
    rootPkg string
}

pkg {
    entities {[string]entity}
    exports {string:string}
    rootComponent *string
}

pkgEntity {
    typ | component | message | datatype | interface
    component component
    message message
    datatype datatypeExpr
    interface  interface
}

component {
    ports interface
    interfacesDeps { [string]: interfaceRef }
    nodes {[string]: nodeDef}
    net ...
}

message {
    datatypeExpr {typeref|datatypeExpr}
    value
}

datatype datatypeExpr

entityRef {
    name: pkg, name string
}

nodeDef {
    type interface|direct
    ref string
    iips {
        [inport:string]: msgref
    }
}

runtime {
    ports []{portRef, portDef}
    network []{from portRef, to portRef[]}
    effects {
        giver (generator) {
                msg ...
                ports {}
        }
        void (absorber/destructor) {
            ports []portRef
        }
        operator {
            ref string
            io {}
        }
    }
    startPort portRef
}

portDef {
    bufSize uint8
}

portRef {
    path, node, name string
}
```

# Array ports bypass rules

_Контекст узла_ при компиляции это специфическая информация о том, как узел-родитель использует данный компонент в своей подсети:

```
nodeCtx {
    arrPortsSize {
        in {
            x: 3
        }
        out {
            y: 2
        }
    }
}
```

Контекст (арпортов) должен быть всегда ясен

байпас для сабноды разрешен тогда и только тогда
когда он идёт в инпорт или аутпорт компонента
потому что иначе возможна ситуация, при которых мы остаёмся без контекста
и, следовательно, не можем вычислить кол-во слотов для портов
это ситуация - сабнода.арраутпорт -> сабнода.арринпорт

таким образом, в руте байпас сабноды невозможен в принципе
т.к. нет арринпоров и арпортов

следовательно, где бы, на каком уровне вложенности, не началась бы цепочка байпасов
она начинается с индексных ссылок на конкретные слоты
следовательно, всегда можно вычислить контекст

# Valid network

Валидная сеть это сеть в которой

- Порт с IIP не имеет входных соединений
- Рутовый компонент не имеет дженериков (кроме sig<T>)
- Любой компонент кроме корневого должен иметь как инпорты так и аутпорты (как минимум по 1 штуке)


