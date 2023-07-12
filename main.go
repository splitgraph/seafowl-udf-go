package main

import (
	"fmt"
	"unsafe"
)

func doAdd(left, right int64) int64 {
	return left + right
}

//export addi64
func addi64(inputPtr unsafe.Pointer) unsafe.Pointer {
	wrappedUdf := func(args []interface{}) (interface{}, error) {
		// Decode the arguments
		left := args[0].(int64)
		right := args[1].(int64)

		// Call the function
		result := doAdd(left, right)
		fmt.Println("result", result)

		// Encode the result
		return result, nil
	}

	// Call WrapUDF with the wrapped function
	return WrapUDF(inputPtr, wrappedUdf)
}

// placeholder main function b/c it's required for WASM compilation
func main() {}
