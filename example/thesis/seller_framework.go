package thesis

import "github.com/zhiqiangxu/ddxf/contract"

// SellerFramework ...
type SellerFramework interface {
	Publish(resourceID string, resourceDDO contract.ResourceDDO, item contract.DTokenItem) error
}

// NewSellerFramework ...
func NewSellerFramework(seller Seller) SellerFramework {
	return &sellerFramework{seller: seller}
}

type sellerFramework struct {
	seller Seller
}

func (sf *sellerFramework) Publish(resourceID string, resourceDDO contract.ResourceDDO, item contract.DTokenItem) (err error) {

	resourceDDO.Manager = sf.seller.OntID()

	err = sf.seller.RememberResourceAndTokens(resourceID, resourceDDO, item)
	if err != nil {
		return
	}

	err = sf.seller.PublishDtoken(resourceID, resourceDDO, item)
	if err != nil {
		err = sf.seller.ForgetResourceAndTokens(resourceID, resourceDDO, item)
		return
	}

	return
}
