package dataid

import "github.com/zhiqiangxu/ddxf"

// DataID ...
type DataID interface {
	UpdateDesc(desc string) (block uint32, err error)
	GetDesc(block uint32) (desc string, err error)
	GetDescByHash(descHash string) (desc string, err error)
	GetDescHash(block uint32) (descHash string, err error)
}

// GenerateTokenHash ...
func GenerateTokenHash(tokenDesc string, dataDesc []string) string {

	hash := ddxf.Sha256Bytes([]byte(tokenDesc))
	for _, desc := range dataDesc {
		hash += ddxf.Sha256Bytes([]byte(desc))
	}

	return ddxf.Sha256Bytes([]byte(hash))
}
