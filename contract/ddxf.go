package contract

import (
	"time"

	"github.com/zhiqiangxu/ddxf"
)

// SellerItemInfo for ddxf
type SellerItemInfo struct {
	Item        ddxf.DTokenItem
	ResourceDDO string
	Sold        uint32
	Owners      map[ddxf.OntID]uint32
}

// DDXFContract for ddxf
type DDXFContract struct {
	sellerItems map[ddxf.OntID]SellerItemInfo
}

// DTokenSellerPublish is called by DTokenSeller
func (c *DDXFContract) DTokenSellerPublish(resourceID ddxf.OntID, resourceDDO string, item ddxf.DTokenItem) {
	if _, ok := c.sellerItems[resourceID]; ok {
		panic("resourceID already exists")
	}

	if resourceDDO == "" {
		panic("resourceDDO empty")
	}

	c.sellerItems[resourceID] = SellerItemInfo{Item: item, ResourceDDO: resourceDDO, Owners: make(map[ddxf.OntID]uint32)}

}

// BuyDToken is called by DTokenBuyer
func (c *DDXFContract) BuyDToken(resourceID ddxf.OntID, n uint32, buyerAccount ddxf.OntID) {
	if !c.checkWitness(buyerAccount) {
		panic("buyerAccount no witness")
	}

	itemInfo, ok := c.sellerItems[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	if time.Now().Unix() > itemInfo.Item.ExpiredDate {
		panic("resourceID already expired")
	}

	if itemInfo.Sold >= itemInfo.Item.Stocks {
		panic("resourceID already sold out")
	}

	if itemInfo.Sold+n >= itemInfo.Item.Stocks {
		panic("resourceID not enough")
	}

	if !c.transferFeeFromAccount(buyerAccount, itemInfo.Item.Fee, n) {
		panic("balance not enough")
	}

	itemInfo.Sold += n
	itemInfo.Owners[buyerAccount] += n
	c.sellerItems[resourceID] = itemInfo
}

// UseDToken is called by buyer
func (c *DDXFContract) UseDToken(resourceID, account ddxf.OntID, n uint32) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemInfo, ok := c.sellerItems[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	ownCount, ok := itemInfo.Owners[account]
	if !ok {
		panic("resourceID not owned by account")
	}

	if n > ownCount {
		panic("resourceID owned not enough")
	}

	ownCount -= n
	if ownCount == 0 {
		delete(itemInfo.Owners, account)
	} else {
		itemInfo.Owners[account] = ownCount
	}

	c.sellerItems[resourceID] = itemInfo

}

func (c *DDXFContract) checkWitness(account ddxf.OntID) bool {
	return true
}

func (c *DDXFContract) transferFeeFromAccount(account ddxf.OntID, fee ddxf.Fee, n uint32) bool {
	return false
}

// BuyerItemInfo for ddxf
type BuyerItemInfo struct {
	Item       ddxf.CrowdSouring
	PublishDDO string
	Collected  uint32
	Providers  map[ddxf.OntID]uint32
}

// DTokenBuyerPublish is called by DTokenBuyer
func (c *DDXFContract) DTokenBuyerPublish(publishID ddxf.OntID, publishDDO string, item ddxf.CrowdSouring) {

}
