---
Title: Style guide
---

Neva language has official style guide and it's described in this document.

## Code Organization

### Keep Files Small

Try to keep files less than **300** lines.

Exceptions _sometimes_ allowed but files that are longer that 300 lines usually could be splitted. Files that are larger than 500 lines are definitely a sign of a problem.

Perfectly it's about 100-200 lines per file.

## Formatting

### Use Tabs Instead of Spaces

1.  Anyone can configure their editor to render tabs exactly like they want whether it's 2, 4 or 8 spaces. Yet the source code itself is the same for everyone. With spaces you have to decide how many spaces you use. Someone likes 2, someone likes 4 and so on.
2.  Tabs use less data for encoding the source code. E.g. for tabs that looks like a 4 spaces there's x4 characters economy.

### Keep Lines Short

Avoid lines that are longer than **80** characters.

1. Even though we have big screens these days, the space that we can focus on with our eyes is limited to about 80 characters. With this line length we can read code without moving the head and with less amount of eyes movement.
2. Professional programmers tend to use split-screens a lot. 120 characters is ok if you only have one editor and probably one sidebar but what about 2 editors? What about 2 editors and sidebar? 2 editors and two sidebars? You got the point.
3. Many people use big fonts because they care about their eyes. Long lines means horizontal scroll for them and scroll itself takes time.
4. Professional programmers sometimes uses different kinds of "code lenses" in their IDEs so e.g. the line that is technically 120 characters becomes actually 150 or 200. Think of argument names lenses for instance.

### Add Newlines

#### Between Sections

```neva
import {
    ...
}

types {
    ...
}

interfaces {
    ...
}

const {
    ...
}

components {
    ...
}
```

#### Between Entities

```neva
types {
    Foo int

    Bar str
}
```

## Design

### Use generics when necessary

You need generics (type parameters) when you need to preserve data type on the output side.

E.g. `Destructor` component doesn't have outports so it doesn't matter what you passed in. That's why `Destructor` accepts `any` instead of `T`.

On the other hand `Blocker` needs to know the type of `data` on the input so the type of the `data` on the output is preserved. That's why it's `Destructor<T>`.

### Separate downstream flow with outports

When you have _structured_ data data use `struct`, when you want to _separate flow_ - create outports.

Example: `NumParser` sends `res` or `err` but never both at the same time. It's case for outports.

On the other hand when you have e.g. two pieces of data `foo` and `bar`, but their firing condition is always the same - use `struct { foo T1, bar T2 }`. This way user of your component could simply use `struct selector` if needed.

### Avoid Multiple Parents When Possible

When designing component's networks prefer _Pipe_ or _Tree_ structure over graphs where single child have multiple parents. This makes network hard to understand visually. Sometimes it's impossible to avoid that though.

### Keep Number of Ports Small

Try to keep the **number of inports no bigger than 3** and **outports no bigger than 5**.

Sometimes we _pass data on_ - we use our outports to not only send the results but also to pass our inputs down the stream so downstream nodes can also use them without being connected to upstream nodes which makes network hard to read both visually and textually.

## Naming Conventions

CamelCase is used everywhere except for package names.

Names are rather short, but their length depends on their scope.

For camelCase (both lower and upper) traditional style (not Go style) is used. For example it's `xmlHttpRequest` and not `xmlHTTPRequest`.

### Packages

Package names are short (1-3 words) `lower_snake_case` strings. Examples: `http` or `business_logic`. They inherit context from their path so it's `net/http` and not `net_http`.

### Entities

`CommonCamelCase` is used for _types_, _interfaces_ and _components_ and `lowerCamelCase` for constants.

Entity names can be relatively long (up to 3 words) and descriptive. It's important because port names must be short and names of the entities they represent will serve as a documentation.

Abbreviation ok if there is a generally accepted one. Or the name turns out to be extremely long (more than 3 words).

For example, `AsynchronousFileReader` is bad because there is generally accepted abbreviation for "Asynchronous", it's "async".

On the other hand `AsyncFRdr` is bad too because words "File" and "Reader" are short enough already and there is no need to shorten them. Besides there's no accepted abbreviation "F" and "Rdr" for them.

Perfect name would be `AsyncFileReader`.

#### Types

Types are generally CamelCase nouns. Examples: `User`, `UserId`, `OrderDetail`, `HttpResponse`, `DayOfWeek`, `ResponseCode`, `FileType`.

Enum elements are named exactly the same way: `{ Monday, Tuesday, ... }`.

Struct fields are named this way except they start with the lower case:

```neva
User struct {
    age int
    name string
}
```

#### Components

Component names typically ends with "er" which describe a performed action. Examples would be `Printer` and `Logger`.

#### Interfaces

Interfaces are named exactly like components except their names are prefixed with `I`. Examples: `IReader`, `IStreamer`.

#### Constants

Constants are lowerCamelCase `defaultTimeout`, `apiEndpoint`, `stockMarketTime`.

### Ports

Use short (up to 5 chars) names for ports.

Single-letter names are ok when it's obvious what they mean based on context. E.g. `f` is okay when it's type is `File`. You'll find many examples of `s` for `string`, `b` for `bool` or `l` for `list`, etc in stdlib.

It's a good practice though to give meaningful names when possible e.g. `res` instead of `s` when it's not just string but rather result of some operation.

Also try to follow patterns from stdlib like `res, err`, `ok, miss`, `some, none`, etc.

And remember that 5 characters is not a lot so you have to shorten most of the time if there's more than 1 word. So `userID` is too long and should be `uid`.

It's important to have short port names because of the visual programming. When you work with graphs of nodes you'll see a lot of connections. It will become unreadable very quickly if there's a lot of ports per node and/or array ports involved.

### Nodes

Nodes are generally named exactly like their components but in `lowerCamelCase`. E.g. for `FileReader` we would have `fileReader`.

If node is abstract (it's instantiated with interface instead of component) then `I` prefix is omitted. So it's `reader IReader` and not `iReader`.

Except if we have several instances of the same component. Then we must stress out why there's several of them and what is the difference. For instance for a component that needs two Adder instances it could be:

```neva
nodes {
    adder1 Adder<int>
    adder2 Adder<int>
}
```
