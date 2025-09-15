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
6. **Std Module Dependency Bug**: Fixed "dependency module not found: std" error by correcting relative imports in std module

   - **Problem**: `std/errors/errors.neva` had `import { runtime }` which tried to reference the `runtime` package within the same `std` module, but the std module didn't have a dependency on itself
   - **Root Cause**: The analyzer couldn't resolve `runtime` because it was looking for it in the std module's dependencies, which was empty
   - **Solution**: Changed import to `import { @:runtime }` to use relative import syntax within the same module
   - **Impact**: Resolves intermittent "dependency module not found: std" panic during compilation
   - **Status**: ✅ **RESOLVED** - Uses existing relative import functionality, conceptually correct approach

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

2. **🚨 NEW PRIORITY: Go Backend Template Port Channel Resolution**: Critical template execution failure preventing program execution

   - **Error**: `template: tpl.go:43:14: executing "tpl.go" at <getPortChanNameByAddr "in" "start">: error calling getPortChanNameByAddr: port chan not found: in:start`
   - **Command**: `neva run examples/hello_world`
   - **Location**: Go backend template execution in `internal/compiler/backend/golang/`
   - **Root Cause**: The `getPortChanNameByAddr` function cannot find the port channel for `in:start` during template execution
   - **Impact**: Prevents any program from running, blocking basic functionality
   - **Status**: 🚨 **NEW HIGH PRIORITY** - Needs immediate investigation

3. **🚨 NEW PRIORITY: HTTP Response Struct Field Compatibility**: Critical struct subtype checking failure preventing compilation

   - **Error**: `advanced_error_handling/main.neva:24:13: Incompatible types: http_get -> .body: Subtype struct is missing field of supertype: body`
   - **Location**: `advanced_error_handling/main.neva:24:13`
   - **Root Cause**: HTTP response struct subtype checking is failing - the subtype struct is missing the `body` field that the supertype expects
   - **Impact**: Prevents compilation of programs using HTTP operations, blocking network functionality
   - **Status**: 🚨 **NEW HIGH PRIORITY** - Needs immediate investigation

4. **🚨 NEW PRIORITY: Type System Subtyping Logic Inconsistency**: Critical type system logic failure with tagged unions

   - **Error**: `add_numbers_from_stdin/main.neva:33:27: Incompatible types: atoi -> out:res: Subtype and supertype must both be either literals or instances, except if supertype is union: expression { child maybe<{ text string, child maybe<error> }>, text string }, constraint int`
   - **Location**: `add_numbers_from_stdin/main.neva:33:27`
   - **Root Cause**: **FUNDAMENTAL TYPE SYSTEM DESIGN ISSUE** - The error message reveals a logical contradiction in the type system's subtyping rules for tagged unions
   - **Impact**: Prevents compilation of programs using string-to-int conversion with union types, blocking basic I/O operations
   - **Status**: 🚨 **NEW HIGH PRIORITY** - Needs immediate investigation

   #### **DETAILED ANALYSIS: The Type System Logic Contradiction**

   **The Error Message Breakdown**:

   ```
   "Subtype and supertype must both be either literals or instances, except if supertype is union"
   ```

   **The Contradiction**:

   1. **Error says**: "except if supertype is union" (suggesting it should allow the conversion)
   2. **But then rejects it anyway** (the conversion is still blocked)
   3. **The constraint is `int`**, not a union, so the "except if supertype is union" rule shouldn't even apply

   **The Real Issue**: The type system contains **legacy untagged union logic** that is inconsistent with the new tagged union reality.

   #### **Untagged vs Tagged Union Subtyping Paradigm Shift**

   **OLD SYSTEM (Untagged Unions)**:

   ```neva
   // This was valid in untagged unions:
   int IS-A (int | string)  // ✅ True - int is a subtype of int|string
   ```

   **Why it worked**: Untagged unions were "any of these types" - so `int` could flow into `int|string` because it was one of the valid options.

   **NEW SYSTEM (Tagged Unions)**:

   ```neva
   // This is NOT valid in tagged unions:
   int IS-A union { Foo int, Bar string }  // ❌ False - int is not a tagged union
   ```

   **Why it doesn't work**: Tagged unions are **distinct types** with explicit variant tags. `int` is not a tagged union at all - it's a primitive type.

   #### **The Type System Logic Problem**

   The error message suggests the type system is still using **old untagged union logic**:

   1. **Old Logic**: "If supertype is union, then subtype can be any of the union's constituent types"
   2. **New Logic**: "If supertype is union, then subtype must ALSO be a union (or a specific tagged variant)"

   #### **What Should the Type System Do?**

   **Option 1: Strict Tagged Union Subtyping**

   ```neva
   // Only these should be valid:
   union { Foo int, Bar string } IS-A union { Foo int, Bar string, Baz bool }  // ✅
   union { Foo int } IS-A union { Foo int, Bar string }  // ✅ (subset of variants)
   int IS-A union { Foo int, Bar string }  // ❌ (primitive is not a union)
   ```

   **Option 2: Hybrid Approach (Current Attempt?)**

   ```neva
   // Maybe the system is trying to allow:
   int IS-A union { Foo int, Bar string }  // ❌ But this breaks tagged union semantics
   ```

   #### **Research Areas for Type System Investigation**

   **Files to investigate**:

   - `internal/compiler/sourcecode/typesystem/validator.go` - Core type validation logic
   - `internal/compiler/sourcecode/typesystem/subtype.go` - Subtyping rules implementation
   - `internal/compiler/analyzer/` - Type checking during analysis phase

   **Key Questions**:

   1. **What are the current subtyping rules for tagged unions?**
   2. **How should primitives interact with tagged unions?**
   3. **What error messages should be shown when these rules are violated?**
   4. **Is the type system trying to maintain backward compatibility with untagged union logic?**

   **The Current Error Message is Misleading**: It talks about "literals vs instances" when the real issue is "primitives vs tagged unions".

5. **🚨 NEW PRIORITY: Error Struct Type Resolution Bug**: Critical type resolution failure in `add_numbers_from_stdin` example

   - **Error**: `add_numbers_from_stdin/main.neva:33:27: Incompatible types: atoi -> out:res: ... expression { child maybe<{ text string, child maybe<error> }>, text string }, constraint int`
   - **Location**: `add_numbers_from_stdin/main.neva:33:27`
   - **Root Cause**: **COMPILER LOGIC ERROR** - The compiler thinks we're sending an error struct to the `res` outport which expects `int`, but this is incorrect
   - **Impact**: Prevents compilation of programs using string-to-int conversion, blocking basic I/O operations
   - **Status**: 🚨 **NEW HIGH PRIORITY** - Needs immediate investigation

   #### **DETAILED ANALYSIS: The Error Struct Problem**

   **The Issue**: The compiler is incorrectly resolving the type of the expression being sent to `atoi -> out:res`. It thinks the expression is:

   ```
   { child maybe<{ text string, child maybe<error> }>, text string }
   ```

   **What This Means**: The compiler believes we're trying to send a complex struct (possibly an error struct) to the `res` outport, but `res` expects an `int`.

   #### **Possible Root Causes**

   **1. Type Resolution Error in Analyzer**:

   - The analyzer might be incorrectly resolving the type of the input to `atoi`
   - Could be in `internal/compiler/analyzer/` type resolution logic

   **2. Connection Sender Type Resolution**:

   - The connection sender type might be incorrectly resolved
   - Could be in `internal/compiler/sourcecode/typesystem/` sender type resolution

   **3. Union Type Handling**:

   - The complex union type `{ child maybe<{ text string, child maybe<error> }>, text string }` suggests union type resolution is failing
   - Could be in union type handling logic

   **4. Error Handling Logic**:

   - The `error` type in the union suggests error handling logic is involved
   - Could be in error type resolution or error handling components

   #### **Research Areas for Error Struct Investigation**

   **Files to investigate**:

   - `internal/compiler/analyzer/` - Type resolution during analysis
   - `internal/compiler/sourcecode/typesystem/` - Type system logic
   - `examples/add_numbers_from_stdin/main.neva` - The actual example to understand the intended flow
   - `std/` - Standard library components for `atoi` and error handling

   **Key Questions**:

   1. **What is the actual type of the input to `atoi`?**
   2. **Why is the compiler resolving it as a complex error struct?**
   3. **Where in the compiler pipeline is this incorrect type resolution happening?**
   4. **Is this related to the union type subtyping issue above?**

   **The Real Question**: The compiler is making a fundamental mistake in understanding what type is being sent to `atoi -> out:res`. We need to trace through the compiler pipeline to find where this incorrect type resolution occurs.

6. **Expression Resolution Validation**: Core expression validation preventing basic compilation

   - Error: `expression must be valid in order to be resolved: expr must be ether literal or instantiation, not both and not neither`
   - Location: `internal/compiler/sourcecode/typesystem/validator.go:40`
   - Root Cause: Unknown. Needs to be figured out.

7. **Struct Field Compatibility**: General struct subtype checking failures
   - Error: `Subtype struct is missing field of supertype: body`
   - Impact: HTTP response handling and struct operations
   - Root Cause: Unknown. Needs to be figured out.

### ⏳ MEDIUM PRIORITY

8. **Dependency Module Resolution**: Intermittent empty `modRef` causing "dependency module not found:" errors
   - Location: `internal/compiler/sourcecode/scope.go:155`
   - Root Cause: Unknown. Needs to be figured out.

### ⏳ LOW PRIORITY

9. **Function Signature Mismatches**: Parameter count mismatches throughout codebase
10. **Import and Module Issues**: Missing runtime imports, node references
11. **Standard Library Components**: Missing network definitions (`---` sections)
12. **E2E Test Recovery**: 100% failure rate in e2e tests

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
- **Investigate std module issues**: If encountering "dependency module not found: std" errors, check for incorrect imports in std module files. Use relative imports (`@:package`) for intra-module references.
