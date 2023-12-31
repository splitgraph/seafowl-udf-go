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
	sizeBuf := (*[SizeNumBytes]byte)(ptr)[:]

	// Assign length of the message to inputSize
	inputSize := binary.LittleEndian.Uint32(sizeBuf)

	// Skip SizeNumBytes to the actual value
	ptr = unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + SizeNumBytes)

	// Convert the pointer to a Go slice of the correct length
	inputBuf := (*[1 << 30]byte)(ptr)[:inputSize:inputSize]

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

func writeOutput(val interface{}) unsafe.Pointer {
	// Determine length
	sizer := msgpack.NewSizer()
	sizer.WriteAny(val)
	outputSize := sizer.Len()

	// Make a buffer with space for the size at the beginning
	resultBuffer := make([]byte, SizeNumBytes+uintptr(outputSize))

	// Write the size to the beginning of the buffer
	binary.LittleEndian.PutUint32(resultBuffer, uint32(outputSize))

	// Serialize the value
	err := Encode(val, resultBuffer[SizeNumBytes:])
	if err != nil {
		log.Fatal(fmt.Errorf("error encoding output: %w", err))
	}

	return unsafe.Pointer(&resultBuffer[0])
}

func WrapUDF(inputPtr unsafe.Pointer, f func([]interface{}) (interface{}, error)) unsafe.Pointer {
	// Read the input
	input := ReadInput(inputPtr)

	// Call the function
	output, err := f(input)
	if err != nil {
		log.Fatal(fmt.Errorf("error applying function: %w", err))
	}

	// Return the output
	return writeOutput(output)
}
