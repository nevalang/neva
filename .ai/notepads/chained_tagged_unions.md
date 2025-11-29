# Chained Tagged Unions

# Chained Union Senders

There two unsupported cases that we need to support

1. `-> Input::Int ->` - chained, tag-only
2. `-> Input::Int(42) ->` - chained,  with value/wrapped sender

## Case 1 `-> Input::Int ->`

1. create virtual const of type union with specified tag
2. create `new`  node from `NewV2` and bind const to it
3. connect original sender to `new` node and connect `new` node to original receiver

```neva
const __const__1 Input = Input::Int
...
#bind(__const__1)
__new__1 = NewV2<Input>
...
-> __new__1:sig
__new__1:res -> ...
```

## Case 2 `-> Input::Int(foo) ->`

1. create virtual const of type union with specified tag
2. create `union` node from `UnionWrapV2` component and bind const to it
3. connect sender that union wraps to inport `data` of the `union` node
4. connect original sender of the chain to `sig` inport of the `union` node
5. connect `res` outport of the `union` node to original receiver

```neva
const __const__1 string = 'Int'
...
#bind(__const__1)
__union_1__ = UnionWrapV2<Input>
...
foo -> __union_1__:data
... -> __union_1__:sig
__union_1__:res -> ...
```

---

*Notepad ID: 39c1e8d5-ba05-43cd-b55e-fae0273234d3*

*Created: 12/7/2024, 10:27:07 PM*

