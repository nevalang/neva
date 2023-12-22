# Best Practices

> This section is under heavy development

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

### Between Entities

```neva
types {
    Foo int

    Bar str
}
```

## Design

### Avoid Multiple Parents

When designing component's networks prefer _Pipe_ or _Tree_ structure over graphs where single child have multiple parents. This makes network hard to understand visually. Sometimes it's impossible to avoid that though.

### Keep Number of Ports Small

Try to keep the **number of inports no bigger than 3** and **outports no bigger than 5**.

Sometimes we _pass data on_ - we use our outports to not only send the results but also to pass our inputs down the stream so downstream nodes can also use them without being connected to upstream nodes which makes network hard to read both visually and textually.

<!-- Outports are optional, inports are not. This means if you have 5 outports, user of your component might not use all of them. On the other hand if you have 3 or more inports - user of your component will be forced to use them all in order to compile the program. -->

## Naming Conventions

CamelCase is used everywhere except packages and constants.

Names are rather short and always start with lowercase except entities.

For camelCase (both lower and upper) traditional style (not Go style) is used. For example it's `xmlHttpRequest` and not `xmlHTTPRequest`.

### Packages

Package names are in `lower_snake_case`. short 1-3 words, perfectly one word. Examples: `http` or `business_logic`. They inherit context from their path so it's `net/http` and not `net_http`.

### Entities

`CommonCamelCase` is used for types, interfaces and components and `UPPER_SNAKE_CASE` is used for constants.

Entity names should be relatively long (up to 5 words) and descriptive. It's important
because other names (e.g. ports) must be short and names of the entities they represent will serve as a documentation.

Abbreviation is ok if there is a generally accepted one. Or the name turns out to be extremely long (more than 3 words).

For example, `AsynchronousFileReader` is bad because there is generally accepted abbreviation for "Asynchronous", it's "async".

On the other hand `AsyncFRdr` is bad too because words "File" and "Reader" are short enough already and there is no need to shorten them. Besides there's no accepted abbreviation "F" and "Rdr" for them.

Perfect name would be `AsyncFileReader`.

Another bad example would be `GeneralPurposeReadonlyLinuxSocketStream`. It consist of 6 words which is too much, one of them must be omitted. `GeneralReadonlyLinuxSocketStream` is acceptable but given how long this name is, `GeneralReadLinuxSockStream` is better.

#### Types

No special rules for types. Examples: `User`, `UserId`, `OrderDetail`, `HttpResponse`, `DayOfWeek`, `ResponseCode`, `FileType`.

Enum elements are named exactly the same way: `{ Monday, Tuesday, ... }`.

Struct fields are named this way too except they start with the lower case:

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

Examples: `DEFAULT_TIMEOUT`, `API_ENDPOINT`, `STOCK_MARKET_CLOSE_TIME`.

### Ports

Use short (1-5 characters) names for ports.

Don't shorten names unnecessarily. For example `file` is better than `f` and `value` is bettern than `v`.

5 characters is not a lot though so you have to shorten most of the time if there's more then 1 word. So `userID` is too long and should be `uid`.

### Nodes

Nodes are generally named exactly like components but in `lowerCamelCase`. E.g. for `FileReader` we would have `fileReader`.

Except if we have several instances of the same component. Then we must stress out why there's several of them and what is the difference. For instance for a component that needs two Adder instances it could be:

```neva
nodes {
    firstAdder Adder<int>
    secondAdder Adder<int>
}
```
