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





