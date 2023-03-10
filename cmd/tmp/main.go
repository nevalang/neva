// generator/main.go
package main

import "fmt"

// func main() {
// 	// Read in the runtime code as a byte slice.
// 	runtimeBytes, err := internal.RuntimeFiles.ReadFile("runtime/runtime.go")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Generate code...
// 	// ...

// 	// Create a buffer to hold the generated code.
// 	var buf bytes.Buffer

// 	// Write the runtime code to the buffer.
// 	buf.WriteString("// Runtime code:\n")
// 	buf.Write(runtimeBytes)
// 	buf.WriteString("\n")

// 	// Write the generated code to the buffer.
// 	buf.WriteString("// Generated code:\n")
// 	buf.WriteString("// ... ")

// 	// Output the buffer to the console.
// 	fmt.Println(buf.String())
// }

func main() {
	a := make([]int, 0, 3)
	b := make([]int, 3, 3)
	c := make([]int, 3)
	fmt.Println(a, b, c)
}
