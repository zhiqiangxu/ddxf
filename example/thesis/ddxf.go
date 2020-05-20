package thesis

import (
	"sync"

	"github.com/zhiqiangxu/ddxf/contract"
)

var (
	dc     *contract.DDXFContract
	dcLock sync.Mutex
)

// DDXF ...
func DDXF() *contract.DDXFContract {
	if dc != nil {
		return dc
	}

	dcLock.Lock()
	defer dcLock.Unlock()

	if dc != nil {
		return dc
	}

	dc = contract.NewDDXFContract(contract.NewDTokenContract())

	return dc
}
