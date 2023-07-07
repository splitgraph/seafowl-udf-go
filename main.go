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
		left, ok := args[0].(int64)
		if !ok {
			return nil, fmt.Errorf("expected first argument to be int64, got %T", args[0])
		}
		right, ok := args[1].(int64)
		if !ok {
			return nil, fmt.Errorf("expected second argument to be int64, got %T", args[1])
		}

		// Call the function
		result := doAdd(left, right)

		// Encode the result
		return result, nil
	}

	// Call WrapUdf with the wrapped function
	return WrapUDF(inputPtr, wrappedUdf)
}

// placeholder main function b/c it's required for WASM compilation
func main() {}

/*
make_scalar_function_wasm_messagepack
	Called at function registration time, NOT when the UDF is invoked

line 556 is where the UDF gets called

scale is # of decimal places

add WASI if only for the logging

we transpose from Datafusion columsn into a struct/object
that is a Rust repr of whatever can be seralized to msgpack

instance.call gets the whole result (table)

in order for the host and wasm VM to communicate, they need to write into the memory

0. figure out how much memory do we need to allocate
1. allocate a chunk of memory in WASM land
		function defined in the WASM module that we call
		"reserve 64 bytes and give me the pointer for the start of that memory"
		""
2. write msgpack serialized data into that memory

pass a pointer, _and_ the size
since we can only return a single value, we return a pointer to a segment of WASM
memory and the first value has the size
read_udf_output()

Write 4 bytes (a 32bit integer)

Seafowl is considered the "host"

wasmtime by itself is just a library
here's a blob with bytecode
i want you to load it, find the funciton named "X"
take these args and run it

WASI: basically a library with two parts
Part 1 sits in the host process (seafowl)
		you can call it from within the WASM program
		if you have callback fns defined from within the WASM
Part 2 is functions that get compiled into the WASM
	wasi symbols are in your assemblyscript module
	they don't work unless the library that registers this stuff knows about them
`wasmtime` is a CLI app that maps current env to the wasm module


*/
