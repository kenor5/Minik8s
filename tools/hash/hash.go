package hash

import (
	"hash/crc32"
)

func HASH(byte2 []byte) uint32 {
	h := crc32.NewIEEE()
	//data := []byte("Hello, World!")
	data := byte2
	h.Write(data)

	hashValue := h.Sum32()
	return hashValue
	//fmt.Printf("Hash value: %d\n", hashValue)
}
