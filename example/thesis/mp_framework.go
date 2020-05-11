package thesis

import "github.com/zhiqiangxu/ddxf/contract"

// MPFramework ...
type MPFramework interface {
	ValidateAndPublish(resourceID string, resourceDDO contract.ResourceDDO, item contract.DTokenItem) (mpItemID string, err error)
	MakeTxForItem(mpItemID string) (tx string, err error)
}
