package ddxf

import (
	"crypto/sha256"
	"encoding/json"

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

	h = HashBytes(bytes)
	return
}

// HashBytes converts any byte sequence to a hash string
func HashBytes(bytes []byte) string {
	h := sha256.New()
	h.Write(bytes)
	hash := h.Sum(nil)
	return util.String(hash)
}
