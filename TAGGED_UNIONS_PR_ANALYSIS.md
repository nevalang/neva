# Tagged Unions Implementation Analysis

## Core Language Changes

### 1. Enums → Tagged Unions

**Problem**: Untagged unions made runtime type identification impossible, preventing proper pattern matching.

**Solution**: Tagged unions with explicit variant tagging.

#### Syntax Migration

```neva
// OLD (enum)
type Day enum { Monday, Tuesday, Wednesday }

// NEW (tagged union)
type Day union { Monday, Tuesday, Wednesday }
```

#### Union Sender Syntax

Four supported cases for pattern matching and value construction:

1. `Input::Int ->` (non-chained, tag-only)
2. `-> Input::Int ->` (chained, tag-only)
3. `Input::Int(foo) ->` (non-chained, with value)
4. `-> Input::Int(42) ->` (chained, with value)

#### Grammar Changes (`internal/compiler/parser/neva.g4`)

```antlr
// OLD
typeLitExpr: enumTypeExpr | structTypeExpr;
enumTypeExpr: 'enum' NEWLINE* '{' NEWLINE* IDENTIFIER (',' NEWLINE* IDENTIFIER)* NEWLINE* '}';

// NEW
typeLitExpr: structTypeExpr | unionTypeExpr;
unionTypeExpr: 'union' NEWLINE* '{' NEWLINE* unionFields? '}';
unionFields: unionField ((',' NEWLINE* | NEWLINE+) unionField)*;
unionField: IDENTIFIER typeExpr? NEWLINE*;
```

### 2. Pattern Matching & Control Flow

**New Features**: `match` and `switch` statements with exhaustive case handling.

#### Pattern Matching Forms

1. **Route Selection**: `src -> match { pattern -> receiver, _ -> default }`
2. **Value Selection**: `src -> match { pattern: value, _: else } -> dst`
3. **Safe Connections**: `src -> match { pattern: value -> receiver }`

#### Runtime APIs

- `MatchV1<T>(src T, [pattern] any) ([dst] T, else T)`
- `MatchV2<T, Y>(src T, [pattern] T, [value] Y, else Y) (dst T)`

### 3. Operator Overloading Refactor

**Problem**: Generic operators with untagged unions created type constraint issues.

**Solution**: Real function overloading replacing generic operators.

#### Before (Generic)

```neva
#extern(int int_add, float float_add, string string_add)
pub def Add<T int | float | string>(left T, right T) (res T)
```

#### After (Overloaded)

```neva
#extern(int_add)
pub def Add(left int, right int) (res int)
#extern(float_add)
pub def Add(left float, right float) (res float)
#extern(string_add)
pub def Add(left string, right string) (res string)
```

## Implementation Architecture

### Parser Level (`internal/compiler/parser/`)

- Union type expression parsing
- Union sender syntax parsing
- Union literal constant parsing
- Complete ANTLR grammar rewrite

### Analyzer Level (`internal/compiler/analyzer/`)

- Union sender validation
- Subtype validation for union tags
- Exhaustive case handling verification
- Conditional operator type checking (overloaded vs non-overloaded)

### Desugarer Level (`internal/compiler/desugarer/`)

- Union sender desugaring (4 cases)
- Integration with overloaded components

### Runtime Level (`internal/runtime/`)

- `union_wrapper_v1.go` and `union_wrapper_v2.go`
- Type system integration for tagged unions

## Critical Implementation Details

### Component Overloading Support

- **Entity Structure**: `Component Component` → `Component []Component`
- **Directive Changes**: `Directives map[Directive][]string` → `Directives map[Directive]string`
- **Node Overload Index**: Added `OverloadIndex *int` field to `Node` struct
- **Overload Resolution**: `getNodeOverloadIndex` function for type-based resolution

### Type System Integration

- Tagged union type definitions and subtype checking
- Union type validation and resolution
- Pattern matching type safety
- Updated union-to-union subtype checking

## Files Changed Summary

- **Core Language**: 15 files (parser, type system, analyzer)
- **Standard Library**: 8 files (operators, runtime functions)
- **Examples/Tests**: 350+ files (all examples, e2e tests, neva.yml)
- **Documentation**: 6 files (comparison, terminology docs)
- **Infrastructure**: 27 files (version, CI/CD, build tools)

### Implementation Changes

**ConnectionSender Construction**:

```go
parsedSender := src.ConnectionSender{
    PortAddr:       senderSidePortAddr,
    Const:          constant,           // for actual constants like 42, "hello"
    Union:          unionSenderData,    // for union senders like Day::Friday
    Range:          rangeExpr,
    StructSelector: senderSelectors,
    Ternary:        ternaryExpr,
    Binary:         binaryExpr,
    Meta:           core.Meta{...},
}
```

### ⚠️ CRITICAL GUIDELINES

- **Focus on single issues**: Never fix multiple issues simultaneously, wait for the input after issue is fixed
- **Think before fixing**: Analyze root cause, don't patch symptoms. Avoid adding mindless if-else checks to avoid panics, nil pointer dereferences, etc, unless root-cause is not obvious.
- **Preserve operator syntax**: Never replace operators with components.
- **Investigate std module issues**: If encountering "dependency module not found: std" errors, check for incorrect imports in std module files. Use relative imports (`@:package`) for intra-module references.
