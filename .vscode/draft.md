This document explores the syntax-level features needed to handle dataflow use cases, beyond the already supported capabilities like fan-in/out, binary/ternary operations, ranges, and struct selectors.

We focus on two main categories of features:

1. Routers - Components that choose a direction for message flow
2. Selectors - Components that choose a specific value

Routers:

- If - Basic conditional routing
- Switch - Multi-branch routing
- Race - Order-based routing

Selectors:

- Match - Value-based selection
- Select - Order-based selection

These components can be further categorized by how they make decisions:

1. Value-based - Uses the content of messages to make choices (If, Switch, Match)
2. Order-based - Uses message arrival order/source to make choices (Race, Select)

An important design principle is that all these features must be race-condition free and intuitive to use. For example, instead of extended switch form, deferred connections could theoretically be used, but they introduce intermediate lock nodes with latency that can cause race conditions. By providing first-class routing constructs instead, we ensure safe and predictable behavior without needing complex locking mechanisms.

# Routers

## If

### Basic

If is the simplest router that receives a boolean message and routes it into two different directions based on its value.

```
sender -> if {
    then -> then_receiver
    else -> else_receiver
}
```

Example

```
true -> if {
    then -> println
    else -> panic
}
```

### One Branch, Many Receivers

If supports multiple receivers per branch, which are handled like a fan-out. Keep in mind that the If component will not receive new messages until the previous ones are handled. Therefore, the sender of the If will be blocked by the slowest receiver in a branch.

```
sender -> if {
    then -> [then_receiver1, then_receiver2]
    ...
}
```

Example

```
job:failed -> if {
    then -> log:fatal
    else -> [log:info, db:write]
}
```

### With Final Receiver

The If component has an extended form where its body is connected to a receiver. In this case, If sends a message to the selected branch receiver first, and then, after that receiver has received the message, it sends to the final receiver.

It is important to note that the basic If, after receiving a boolean message, is blocked only by the selected branch receiver(s). However, If with a final receiver is also blocked by that final receiver. Similarly, final receiver is blocked by the selected branch receiver.

This setup allows answering the question "when has the selected branch received the message?" but not "when has the selected node finished processing?" For example, sending a message to `println` will start printing, but in `if {...} -> final_receiver`, the final receiver has no guarantee that printing is finished.

```
sender -> if {
    then -> then_receiver
    else -> else_receiver
} -> final_receiver
```

One might think that this form of if could be emulated by putting `final_receiver` into branches

```
// this is incorrect
sender -> if {
    then -> [then_receiver, final]
    else -> [else_receiver, final]
}
```

This is not true. It has different semantics because the selected and final receivers will receive messages concurrently, not sequentially. Additionally, it is forbidden to refer to the same receiver twice.

Example

```
(:user.age > 18) -> if {
    then -> greet
    else -> show_banner
} -> println
```

However, an If component with a final receiver can also have multiple receivers per branch, and even multiple final receivers. They are handled like a fan-out and affect the If component in the same way in terms of blocking.

```
sender -> if {
    then -> [then_receiver1, then_receiver2]
    else -> else_receiver
} -> [final_receiver1, final_receiver1]
```

Example

```
:user.paid -> if {
    then -> [greet, println]
    else -> log:fatal
} -> [db, cache]
```

> All `If` features, including those that are described below, can be combined in various ways. For brevity, we won't demonstrate all combinations (only some of them), but all forms of `If` follow this semantics: The `If` component always waits for its inports first, then waits for the selected receivers, and finally waits for the final receivers. It does not receive the next messages until the previous iteration is complete.

### With Two Inports

Previous forms of If had only one inport, which is for the condition. It's possible to use if with 2 inports: one for the condition and one for the data. The purpose is to route one message based on the value of another message.

> The condition sender has been moved from the left part of the expression to the right. This might seem confusing at first, but it is done to avoid even greater confusion, given how we are accustomed to seeing if-statements in other programming languages.

```
data_sender -> if cond_sender {
    then -> then_receiver
    else -> else_receiver
}
```

Example

```
42 -> if true {
    then -> println // prints 42, not true
    else -> panic // would panic with 42, not false
}
```

The If component can have both a final receiver and 2 inports at the same time:

```
data_sender -> if cond_sender {
    then -> then_receiver
    else -> else_receiver
} -> final_receiver
```

Example

```
:user -> if :theme {      // if doesn't care if what order user and theme came, it waits for both
    dark -> respect:inc   // either foo or bar receives 42, not true or false
    light -> respect:decr // after foo or bar recieved, IfResult<int> is sent to println
} -> println              // after println receiver, if is able to receive next condition and data
```

The message sent to the final receiver in this case will have the type `IfResult`. This type allows access to both the condition and the data:

```
pub type IfResult<T> struct {
    cond bool
    data T
}
```

Example

```
{ cond: true, data: 42 }
```
