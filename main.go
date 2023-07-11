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
		left, err := DecodeI64(&args[0])
		if err != nil {
			return nil, fmt.Errorf("expected first argument to be int64, got %T", args[0])
		}
		right, err := DecodeI64(&args[1])
		if err != nil {
			return nil, fmt.Errorf("expected second argument to be int64, got %T", args[1])
		}

		// Call the function
		result := doAdd(left, right)

		// Encode the result
		return result, nil
	}

	// Call WrapUDF with the wrapped function
	return WrapUDF(inputPtr, wrappedUdf)
}

// placeholder main function b/c it's required for WASM compilation
func main() {}
