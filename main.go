package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Expected two arguments")
		os.Exit(1)
	}
	var err error
	N := make([]int, 2)
	for i := 0; i < 2; i++ {
		if N[i], err = strconv.Atoi(args[i]); err != nil {
			panic(err)
		}
	}

	fmt.Println(AddInts(N[0], N[1]))
}