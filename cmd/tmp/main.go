package main

type Box[T any] struct {
	x T
}

func main() {
	b := Box[Box[int32]]{}
	print(b)
}
