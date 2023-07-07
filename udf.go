package main

import (
	"fmt"
	"log"
	"unsafe"

	msgpack "github.com/wapc/tinygo-msgpack"
)

func Encode(v interface{}, buf []byte) error {
	encoder := msgpack.NewEncoder(buf)
	encoder.WriteAny(v)
	return encoder.Err()
}

func Decode(buf []byte) (interface{}, error) {
	decoder := msgpack.NewDecoder(buf)
	return decoder.ReadByteArray()
}

// Required, see https://seafowl.io/docs/guides/custom-udf-wasm#the-wasm-udf-protocol
//
//export alloc
func alloc(size uintptr) unsafe.Pointer {
	buffer := make([]byte, 0, size)
	pointer := unsafe.Pointer(&buffer[0])
	return pointer
}

// Go is GC'd and doesn't support dealloc, but this is required along with alloc
//
//export dealloc
func dealloc(pointer uintptr, capacity int32) {}

func decodeI64(v interface{}) (int64, error) {
	if i64, ok := v.(int64); ok {
		return i64, nil
	}
	return 0, fmt.Errorf("Expected to find int64 value, but received %v instead", v)
}

const SizeNumBytes = int(unsafe.Sizeof(int32(0))) // intended to match Rust's std::mem::size_of::<i32>();

//export ReadInput
func ReadInput(ptr unsafe.Pointer) []interface{} {
	// Convert pointer to byte slice
	sizeBuf := (*[SizeNumBytes]byte)(ptr)[:]
	fmt.Println("sizeBuf", sizeBuf)

	// Convert byte slice to integer
	var inputSize int
	for i, b := range sizeBuf {
		inputSize |= int(b) << (8 * i)
	}
	fmt.Println("inputSize", inputSize)

	// Convert the pointer to a Go slice of the correct length
	inputBuf := (*[1 << 30]byte)(ptr)[:inputSize:inputSize]
	fmt.Println("inputBuf", inputBuf)

	// Decode the input buffer
	inputValue, err := Decode(inputBuf)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading input buffer: %w", err))
	}

	// Convert input value to array
	inputArray, ok := inputValue.([]interface{})
	if !ok {
		log.Fatal(fmt.Errorf("error reading input buffer as array, found instead: %v", inputValue))
	}

	return inputArray
}

//export WriteOutput
func WriteOutput(val interface{}) unsafe.Pointer {
	// Make a buffer with space for the size at the beginning
	// TODO pretty sure int64's need to be 8 bytes
	serializedOutput := make([]byte, SizeNumBytes)

	// Serialize the value
	err := Encode(val, serializedOutput[SizeNumBytes:])
	if err != nil {
		log.Fatal(fmt.Errorf("error encoding output: %w", err))
	}

	// Write the size to the beginning of the buffer
	outputSize := len(serializedOutput) - SizeNumBytes
	for i := 0; i < SizeNumBytes; i++ {
		serializedOutput[i] = byte(outputSize >> (8 * i))
	}

	return unsafe.Pointer(&serializedOutput[0])
}

// export WrapUDF
func WrapUDF(inputPtr unsafe.Pointer, f func([]interface{}) (interface{}, error)) unsafe.Pointer {
	fmt.Println("WrapUDF()")
	fmt.Printf("Pointer value: %v, type: %T, address: %p\n", inputPtr, inputPtr, &inputPtr)

	// Read the input
	input := ReadInput(inputPtr)

	// Call the function
	output, err := f(input)
	if err != nil {
		log.Fatal(fmt.Errorf("error applying function: %w", err))
	}

	// Return the output
	return WriteOutput(output)
}
