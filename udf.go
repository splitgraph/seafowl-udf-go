package main

import (
	"fmt"
	"unsafe"
	// "tinygo-msgpack"
)

// func encode(buf []byte) error {
// 	return UnsafeString(buf)
// }

// func decode(buf *[]byte) {
// 	return msgpack.from()
// }

// Required by Seafowl
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
