package main

import (
	"encoding/binary"
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
	arraySize, err := decoder.ReadArraySize()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, arraySize)

	// Loop over the elements of the array and decode each one.
	for i := 0; i < int(arraySize); i++ {
		// Decode the next value. The exact function to call here will depend on
		// the types of values you're expecting in the array.
		// For example, if you're expecting integers, you might call decoder.DecodeInt().
		value, err := decoder.ReadAny()
		if err != nil {
			return nil, err
		}

		// Add the decoded value to the result slice.
		result[i] = value
	}

	// Return the slice of decoded values.
	return result, nil
}

// Required, see https://seafowl.io/docs/guides/custom-udf-wasm#the-wasm-udf-protocol
//
//export alloc
func alloc(size uintptr) unsafe.Pointer {
	buffer := make([]byte, size, size)
	if len(buffer) == 0 {
		log.Fatal(fmt.Errorf("buffer is empty"))
	}
	pointer := unsafe.Pointer(&buffer[0])
	return pointer
}

// Go is GC'd and doesn't support dealloc, but this is required along with alloc
//
//export dealloc
func dealloc(pointer uintptr, capacity int32) {}

func DecodeI64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case []byte:
		if len(v) < 8 {
			return 0, fmt.Errorf("Expected byte slice of length 8 for int64, but received length %d instead", len(v))
		}
		u64 := binary.BigEndian.Uint64(v)
		return int64(u64), nil
	default:
		return 0, fmt.Errorf("Expected to find []byte value for int64, but received %T instead", v)
	}
}

const SizeNumBytes = unsafe.Sizeof(int32(0)) // intended to match Rust's std::mem::size_of::<i32>();

func ReadInput(ptr unsafe.Pointer) []interface{} {
	// Convert pointer to byte slice
	// sizeBuf := (*[SizeNumBytes]byte)(ptr)[:]
	sizeBuf := []byte{3, 0, 0, 0}
	fmt.Println("sizeBuf", sizeBuf)

	inputSize := binary.LittleEndian.Uint32(sizeBuf)
	inputSize = 3
	fmt.Println("inputSize", inputSize)
	// Convert the pointer to a Go slice of the correct length
	// inputBuf := (*[1 << 30]byte)(ptr)[:inputSize:inputSize]
	inputBuf := []byte{146, 1, 103}
	fmt.Println("inputBuf", inputBuf)

	// Decode the input buffer
	inputValue, err := Decode(inputBuf)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading input buffer: %w", err))
	}
	fmt.Println("inputValue", inputValue)

	// Convert input value to array
	inputArray, ok := inputValue.([]interface{})
	if !ok {
		log.Fatal(fmt.Errorf("error reading input buffer as array, found instead: %v", inputValue))
	}
	fmt.Println("inputArray", inputArray)

	return inputArray
}

func writeOutput(val interface{}) unsafe.Pointer {
	// Make a buffer with space for the size at the beginning
	serializedOutput := make([]byte, SizeNumBytes+unsafe.Sizeof(int64(0)))

	// Write the size to the beginning of the buffer
	outputSize := len(serializedOutput) - int(SizeNumBytes)
	binary.LittleEndian.PutUint32(serializedOutput, uint32(outputSize))

	// Serialize the value
	err := Encode(val, serializedOutput[SizeNumBytes:])
	if err != nil {
		log.Fatal(fmt.Errorf("error encoding output: %w", err))
	}

	fmt.Println("outputSize", outputSize)
	fmt.Println("serializedOutput", serializedOutput, "uint32(outputSize)", uint32(outputSize))

	return unsafe.Pointer(&serializedOutput[0])
}

func WrapUDF(inputPtr unsafe.Pointer, f func([]interface{}) (interface{}, error)) unsafe.Pointer {
	fmt.Printf("WrapUDF() pointer value: %v, type: %T\n", inputPtr, inputPtr)

	// Read the input
	input := ReadInput(inputPtr)
	fmt.Println("input", input)

	// Call the function
	output, err := f(input)
	if err != nil {
		log.Fatal(fmt.Errorf("error applying function: %w", err))
	}
	fmt.Println("output", output)

	// Return the output
	return writeOutput(output)
}
