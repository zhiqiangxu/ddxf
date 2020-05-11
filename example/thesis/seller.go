package thesis

import (
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ddxf/contract"
)

// Seller ..
type Seller interface {
	OntID() ddxf.OntID

	RememberResourceAndTokens(resourceID string, resourceDDO contract.ResourceDDO, item contract.DTokenItem) error
	ForgetResourceAndTokens(resourceID string, resourceDDO contract.ResourceDDO, item contract.DTokenItem) error

	PublishDtoken(resourceID string, resourceDDO contract.ResourceDDO, item contract.DTokenItem) error
}
