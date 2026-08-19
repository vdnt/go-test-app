package main

import "fmt"

func main() {
	fmt.Println(Greet("rollout"))
}

func Greet(name string) string {
	if name == "" {
		name = "world"
	}
	return "hello, " + name + "!"
}
