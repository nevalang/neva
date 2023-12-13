# Naming Conventions

CamelCase is used everywhere except packages and constants.

Names are rather short and always start with lowercase except entities.

For camelCase (both lower and upper) traditional style (not Go style) is used. For example it's `xmlHttpRequest` and not `xmlHTTPRequest`.

## Packages

Package names are in `lower_snake_case`. short 1-3 words, perfectly one word. Examples: `http` or `business_logic`. They inherit context from their path so it's `net/http` and not `net_http`.

## Entities

`CommonCamelCase` is used for types, interfaces and components and `UPPER_SNAKE_CASE` is used for constants.

Entity names should be relatively long (up to 5 words) and descriptive. It's important
because other names (e.g. ports) must be short and names of the entities they represent will serve as a documentation.

Abbreviation is ok if there is a generally accepted one. Or the name turns out to be extremely long (more than 3 words).

For example, `AsynchronousFileReader` is bad because there is generally accepted abbreviation for "Asynchronous", it's "async".

On the other hand `AsyncFRdr` is bad too because words "File" and "Reader" are short enough already and there is no need to shorten them. Besides there's no accepted abbreviation "F" and "Rdr" for them.

Perfect name would be `AsyncFileReader`.

Another bad example would be `GeneralPurposeReadonlyLinuxSocketStream`. It consist of 6 words which is too much, one of them must be omitted. `GeneralReadonlyLinuxSocketStream` is acceptable but given how long this name is, `GeneralReadLinuxSockStream` is better.

### Types

No special rules for types. Examples: `User`, `UserId`, `OrderDetail`, `HttpResponse`, `DayOfWeek`, `ResponseCode`, `FileType`.

Enum elements are named exactly the same way: `{ Monday, Tuesday, ... }`.

Struct fields are named this way too except they start with the lower case:

```neva
User struct {
    age int
    name string
}
```

### Components

Component names typically ends with "er" which describe a performed action. Examples would be `Printer` and `Logger`.

### Interfaces

Interfaces are named exactly like components except their names are prefixed with `I`. Examples: `IReader`, `IStreamer`.

### Constants

Examples: `DEFAULT_TIMEOUT`, `API_ENDPOINT`, `STOCK_MARKET_CLOSE_TIME`.

## Ports

Use short (1-5 characters) names for ports.

Don't shorten names unnecessarily. For example `file` is better than `f` and `value` is bettern than `v`.

5 characters is not a lot though so you have to shorten most of the time if there's more then 1 word. So `userID` is too long and should be `uid`.

## Nodes

Nodes are generally named exactly like components but in `lowerCamelCase`. E.g. for `FileReader` we would have `fileReader`.

Except if we have several instances of the same component. Then we must stress out why there's several of them and what is the difference. For instance for a component that needs two Adder instances it could be:

```neva
nodes {
    firstAdder Adder<int>
    secondAdder Adder<int>
}
```
