# Tagged Unions Implementation Analysis

## Overview

**PR #830**: Major language feature change introducing tagged unions to replace enums. **53 commits**, **406 files**, **+7,080/-5,579 lines**.

**Issues Resolved**: #751 (Tagged Unions), #747 (Pattern Matching), #726 (Match Syntax), #725 (Switch Statement), #749 (Type Assertions)

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
- **CRITICAL**: Conditional operator type checking (overloaded vs non-overloaded)

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
- **CRITICAL**: Proper union-to-union subtype checking

## Current Status & Issues

### ✅ COMPLETED

1. **Type System Crashes**: Fixed null pointer dereference in `Expr.String()` method
2. **Operator Overloading**: Implemented conditional logic for overloaded vs non-overloaded operators
3. **Parser Grammar**: Complete ANTLR grammar rewrite for tagged unions
4. **Runtime Functions**: Union wrapper implementations
5. **Union Sender Type Resolution**: Fixed critical parser bug where union tag-only syntax was incorrectly treated as constants

### 🚨 HIGH PRIORITY - CURRENT FOCUS

1. **✅ Union Sender Type Resolution Bug**: **FIXED** - Critical issue with union tag-only syntax parsing

   - **Problem**: Union tag `Day::Friday` was incorrectly being treated as a constant instead of a union sender
   - **Root Cause**: Parser was creating `Const` objects instead of properly setting the `Union` field in `ConnectionSender`
   - **Fix Applied**: Modified `parseSingleSender` in `internal/compiler/parser/listener_helpers.go` to:
     - Create `UnionSender` objects instead of `Const` objects for union senders
     - Properly handle wrapped data for union senders with values (e.g., `Day::Friday(42)`)
     - Store union sender data in the `Union` field of `ConnectionSender`
   - **Impact**: Union tag-only syntax now properly parsed and type-resolved
   - **Status**: ✅ **RESOLVED** - Parser tests passing, union senders correctly identified

2. **Expression Resolution Validation**: Core expression validation preventing basic compilation

   - Error: `expression must be valid in order to be resolved: expr must be ether literal or instantiation, not both and not neither`
   - Location: `internal/compiler/sourcecode/typesystem/validator.go:40`
   - Root Cause: Unknown. Needs to be figured out.

3. **Struct Field Compatibility**: Struct subtype checking failures
   - Error: `Subtype struct is missing field of supertype: body`
   - Impact: HTTP response handling and struct operations
   - Root Cause: Unknown. Needs to be figured out.

### ⏳ MEDIUM PRIORITY

3. **Dependency Module Resolution**: Intermittent empty `modRef` causing "dependency module not found:" errors
   - Location: `internal/compiler/sourcecode/scope.go:155`
   - Root Cause: Unknown. Needs to be figured out.

### ⏳ LOW PRIORITY

4. **Function Signature Mismatches**: Parameter count mismatches throughout codebase
5. **Import and Module Issues**: Missing runtime imports, node references
6. **Standard Library Components**: Missing network definitions (`---` sections)
7. **E2E Test Recovery**: 100% failure rate in e2e tests

## Migration Requirements

### Breaking Changes

1. **Enum → Union**: All enum usage must migrate to union syntax
2. **Generic Operators → Overloading**: Replace generic operator calls with specific overloaded functions
3. **ParseNum → Specific Functions**: Use `Atoi`, `ParseInt`, `ParseFloat` instead of `ParseNum<T>`

### Migration Examples

```neva
// Enum → Union
type Status enum { Success, Error } → type Status union { Success, Error }

// Generic → Overloaded
strconv.ParseNum<int> → strconv.Atoi
strconv.ParseNum<float> → strconv.ParseFloat

// Pattern Matching
def HandleResult(result Result) (output string) {
    match result {
        Success(msg) -> processSuccess:data
        Error(err) -> processError:data
    }
}
```

## Files Changed Summary

- **Core Language**: 15 files (parser, type system, analyzer)
- **Standard Library**: 8 files (operators, runtime functions)
- **Examples/Tests**: 350+ files (all examples, e2e tests, neva.yml)
- **Documentation**: 6 files (comparison, terminology docs)
- **Infrastructure**: 27 files (version, CI/CD, build tools)

## Union Sender Architecture & Fix Details

### The Two-Approach Problem

The codebase had two different structures for handling unions, leading to confusion:

1. **UnionLiteral** (for constants in `MsgLiteral`):

```go
type UnionLiteral struct {
    EntityRef core.EntityRef `json:"entityRef,omitempty"`
    Tag       string         `json:"tag,omitempty"`
    Data      *ConstValue    `json:"data,omitempty"`  // wraps another const value
    Meta      core.Meta      `json:"meta,omitempty"`
}
```

2. **UnionSender** (for network connections in `ConnectionSender`):

```go
type UnionSender struct {
    EntityRef core.EntityRef    `json:"entityRef,omitempty"`
    Tag       string            `json:"tag,omitempty"`
    Data      *ConnectionSender `json:"data,omitempty"`  // wraps another sender
    Meta      core.Meta         `json:"meta,omitempty"`
}
```

### The Critical Question: Const vs Union Sender

**Original Issue**: When parsing `Day::Friday`, should it be:

- A **const sender** with union message literal value? (stored in `sender.Const.Value.Message.Union`)
- A **union sender** with tag-only data? (stored in `sender.Union`)

### The Fix: Union Sender Approach

**Decision**: Use the **Union Sender** approach because:

1. **Network Semantics**: `Day::Friday` in a connection like `Day::Friday -> receiver` is fundamentally a network sender, not a constant
2. **Wrapping Capability**: Union senders can wrap other senders: `Day::Friday(port_sender)` or `Day::Friday(42)`
3. **Type Resolution**: The analyzer's `getResolvedSenderType` can properly handle union senders through the `Union` field
4. **Consistency**: Aligns with the four union sender cases defined in the grammar

### Implementation Changes

**File**: `internal/compiler/parser/listener_helpers.go`

**Before** (incorrect):

```go
if unionSender != nil {
    // ... parse union data ...
    constant = &src.Const{
        Value: src.ConstValue{
            Message: &src.MsgLiteral{
                Union: &src.UnionLiteral{...},
            },
        },
    }
}
```

**After** (correct):

```go
if unionSender != nil {
    // ... parse union data ...
    unionSenderData = &src.UnionSender{
        EntityRef: parsedUnionRef,
        Tag:       unionSender.IDENTIFIER().GetText(),
        Data:      wrappedSender,  // if wrapped data exists
        Meta:      core.Meta{...},
    }
}
```

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

### Test Updates

Updated parser tests in `internal/compiler/parser/parser_test.go` to use the new `Union` field instead of the old `Const.Value.Message.Union` path.

## Next Steps for AI Agents

### ⚠️ CRITICAL GUIDELINES

- **Focus on single issues**: Never fix multiple issues simultaneously, wait for the input after issue is fixed
- **Think before fixing**: Analyze root cause, don't patch symptoms. Avoid adding mindless if-else checks to avoid panics, nil pointer dereferences, etc, unless root-cause is not obvious.
- **Preserve operator syntax**: Never replace operators with components.
