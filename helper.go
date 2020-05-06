package ddxf

import (
	"crypto/sha256"

	"github.com/zhiqiangxu/util"
)

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func hash(bytes []byte) string {
	h := sha256.New()
	h.Write(bytes)
	hash := h.Sum(nil)
	return util.String(hash)
}
