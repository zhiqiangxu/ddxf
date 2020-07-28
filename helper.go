package ddxf

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"hash/crc32"

	"github.com/ontio/ontology/common"
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
func HashObject(obj interface{}) (h [sha256.Size]byte, err error) {
	bytes, err := Object2Bytes(obj)
	if err != nil {
		return
	}

	// str := hex.EncodeToString(bytes)
	// fmt.Println("bytes", bytes, "\nstr", str)

	h = Sha256Bytes(bytes)
	return
}

// HashObject2U256 returns u256
func HashObject2U256(obj interface{}) (h common.Uint256, err error) {
	hash, err := HashObject(obj)
	if err != nil {
		return
	}

	// fmt.Println("after HashObject", hex.EncodeToString(hash[:]))
	h, err = common.Uint256ParseFromBytes(hash[:])
	return
}

// HashObject2Hex returns u256 hex string
func HashObject2Hex(obj interface{}) (result string, err error) {
	h, err := HashObject2U256(obj)
	if err != nil {
		return
	}

	result = hex.EncodeToString(h[:])

	// fmt.Println("after ToHexString", result)
	return
}

// Sha256Bytes converts any byte sequence to a sha256 hash string
func Sha256Bytes(bytes []byte) [sha256.Size]byte {
	// fmt.Println("Sha256Bytes", bytes)
	return sha256.Sum256(bytes)
}

// Crc32Bytes converts any byte sequence to a crc32 hash string
func Crc32Bytes(data []byte) [4]byte {
	s := crc32.ChecksumIEEE(data)
	return [4]byte{byte(s >> 24), byte(s >> 16), byte(s >> 8), byte(s)}
}
