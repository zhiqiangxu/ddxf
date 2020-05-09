package ddxf

import (
	"crypto/sha256"
	"encoding/json"
	"hash/crc32"

	"github.com/zhiqiangxu/util"
)

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

// Object2Bytes converts any object to alphabetical byte sequence
func Object2Bytes(obj interface{}) (bytes []byte, err error) {
	bytes, err = json.Marshal(obj)
	if err != nil {
		return
	}

	var sorted interface{}
	err = json.Unmarshal(bytes, &sorted)
	if err != nil {
		return
	}

	bytes, err = json.Marshal(sorted)
	return
}

// HashObject converts any object to a hash string
// map keys are sorted asc
// struct fields are sorted by field order
// obj: {"f1":"v1",...}
func HashObject(obj interface{}) (h string, err error) {
	bytes, err := Object2Bytes(obj)
	if err != nil {
		return
	}

	h = Sha256Bytes(bytes)
	return
}

// Sha256Bytes converts any byte sequence to a sha256 hash string
func Sha256Bytes(bytes []byte) string {
	h := sha256.New()
	h.Write(bytes)
	hash := h.Sum(nil)
	return util.String(hash)
}

// Crc32Bytes converts any byte sequence to a crc32 hash string
func Crc32Bytes(data []byte) []byte {
	s := crc32.ChecksumIEEE(data)
	return []byte{byte(s >> 24), byte(s >> 16), byte(s >> 8), byte(s)}
}
