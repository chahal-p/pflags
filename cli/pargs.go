package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	args := os.Args[1:]
	fmt.Printf("All args: %v\n", args)
	intOrExtArgDel := slices.Index(args, "--")
	if intOrExtArgDel == -1 {
		fmt.Println("No internal and external delimitor ('--') present.")
		return
	}
	println(intOrExtArgDel)
	intArgs := args[:intOrExtArgDel]
	extArgs := args[intOrExtArgDel+1:]
	fmt.Println(intArgs)
	fmt.Println(extArgs)
}
