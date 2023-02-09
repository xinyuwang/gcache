package gcache

import (
	"hash/fnv"
	"unsafe"
)

func fnv1(key string) uint32 {
	h := fnv.New32()
	h.Write(unsafeStringToBytes(key))
	return h.Sum32()
}

func unsafeStringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
