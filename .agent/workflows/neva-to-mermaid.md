---
description: Convert Neva programs to valid Mermaid diagrams
---

# Neva to Mermaid Skill

Use this workflow to convert any Neva source code into a Mermaid flowchart.

## Instructions

1. Read the Neva code to be converted.
2. Generate a Mermaid `flowchart TB` diagram using the following rules:

   - **Layout**: Always include this header at the top:
     ```mermaid
     ---
     config:
       layout: elk
     ---
     flowchart TB
     ```

   - **Ports**:
     - Use **stadium** shapes for external ports `:in` / `:out`:
       - `([":data"])`, `([":res"])`, `([":err"])`, etc.
     - In edges, label with the **port name** (`data`, `res`, `err`, `sig`, `case[0]`, etc).

   - **Components**:
     - Use unique Mermaid node IDs.
     - Render components as regular boxes: `id["name: Type"]`.

   - **Component nodes (IMPORTANT)**: 
     - Component nodes MUST be rendered using **only their declared instance name**.
     - Do NOT annotate component nodes with type names, generic parameters, or usage information.

   - **Literals**:
     - Render numeric and string literals as rectangles using Mermaid shape syntax:
       - `zero@{ label: "0", shape: rect }`
       - `dot@{ label: "'.'", shape: rect }`

   - **Connections**:
     - Use `-- <port> -->` labels for the wire’s port name.
     - For `a:out -> b:in`, label the edge with the **output port name** (`out`), unless Neva explicitly uses the input port name in the wiring (e.g., `b:sig`).
     - Only label an edge with a port name when the Neva source explicitly names a port on that connection.
       - Examples that MUST be labeled:
       - `x:res -> y` → label edge as `res`
       - `x -> y:sig` → label edge as `sig`
       - `x:case[0] -> y` → label edge as `case[0]`
    - If a port name is **omitted** in Neva, do **not** infer or guess it. Render an unlabeled edge:
       - `x -> y` → `x --> y`

   - **Fan-out / Fan-in (IMPORTANT)**:
     - Mermaid edges do not create real “ports as objects”, so list-wiring must be represented explicitly with **virtual junction nodes**.
     - A junction node is a small circle: `jX(( ))` where `jX` is a unique ID.

     **Fan-out**: `src -> [a, b, :out]`
     - Insert a junction node `jX(( ))`.
     - Wire: `src -- <port> --> jX`
     - Then wire: `jX --> a`, `jX --> b`, `jX --> out_port`

     **Fan-in**: `[a, b] -> dst:port`
     - Insert a junction node `jX(( ))`.
     - Wire: `a --> jX` and `b --> jX`
     - Then wire: `jX -- <port> --> dst`

   - **Grouping (`&`)**:
     - You MAY use `&` to reduce clutter only when it does **not** hide fan-in/fan-out topology.
     - If Neva uses list-wiring (`[...]`), prefer junction nodes over `&`.

3. Don't lookup otherfiles in repo and do not bother calling tools like websearch, this text describes absolutely everything you need.
4. Ensure the output is ready to be pasted into the Mermaid playground.

## Reference Example

### Neva

```neva
pub def Tap<T>(data T) (res T, err error) {
    pass1 Pass<T>
    pass2 Pass<T>
    lock Lock<T>
    handler ITapHandler<T>
    ---
    :data -> [lock:data, handler]
    handler:res -> pass1
    handler:err -> [pass2, :err]
    [pass1, pass2] -> lock:sig
    lock -> :res
}