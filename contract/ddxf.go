package contract

import (
	"fmt"
	"time"

	"strings"

	"github.com/zhiqiangxu/ddxf"
)

// DTokenItem for ddxf
type DTokenItem struct {
	Fee         ddxf.Fee
	ExpiredDate int64
	Stocks      uint32
	Templates   map[string]struct{}
}

// RT for resource type
type RT int

const (
	// RTStaticFile for static file
	RTStaticFile RT = iota
)

// ResourceDDO is ddo for resource
type ResourceDDO struct {
	Manager  ddxf.OntID // data owner id
	Endpoint string     // data service provider uri
	RT       RT
}

// SellerItemInfo for ddxf
// immutable
type SellerItemInfo struct {
	Item        DTokenItem
	ResourceDDO ResourceDDO
}

// SellerItemStatus for ddxf
// mutable
type SellerItemStatus struct {
	Sold   uint32
	Owners map[ddxf.OntID]map[string]uint32
}

// DDXFContract for ddxf
type DDXFContract struct {
	sellerItemInfo   map[string]SellerItemInfo
	sellerItemStatus map[string]SellerItemStatus
}

var emptyResourceDDO = ResourceDDO{}

// DTokenSellerPublish is called by DTokenSeller
func (c *DDXFContract) DTokenSellerPublish(resourceID string, resourceDDO ResourceDDO, item DTokenItem) {
	if _, ok := c.sellerItemInfo[resourceID]; ok {
		panic("resourceID already exists")
	}

	if resourceDDO == emptyResourceDDO {
		panic("resourceDDO empty")
	}

	switch resourceDDO.RT {
	case RTStaticFile:
		for tokenHash := range item.Templates {
			// desc hash : data hash
			if len(strings.Split(tokenHash, ":")) != 2 {
				panic(fmt.Sprintf("invalid tokenHash %s", tokenHash))
			}
		}
	}

	c.sellerItemInfo[resourceID] = SellerItemInfo{Item: item, ResourceDDO: resourceDDO}
	c.sellerItemStatus[resourceID] = SellerItemStatus{Owners: make(map[ddxf.OntID]map[string]uint32)}

}

// BuyDTokenFromReseller is called by DTokenBuyer to buy dtoken from another buyer(reseller)
func (c *DDXFContract) BuyDTokenFromReseller(resourceID string, n uint32, buyerAccount, resellerAccount ddxf.OntID) {
	if !c.checkWitness(buyerAccount) {
		panic("buyerAccount no witness")
	}
	if !c.checkWitness(resellerAccount) {
		panic("resellerAccount no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	itemStatus := c.sellerItemStatus[resourceID]
	resellerTokens, ok := itemStatus.Owners[resellerAccount]
	if !ok {
		panic("resourceID not owned by resellerAccount")
	}

	for tokenHash := range itemInfo.Item.Templates {
		if resellerTokens[tokenHash] < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
	}

	if !c.transferFeeFromAccount(buyerAccount, resellerAccount, itemInfo.Item.Fee, n) {
		panic("balance not enough")
	}

	ownedTokens := itemStatus.Owners[buyerAccount]
	if ownedTokens == nil {
		ownedTokens = make(map[string]uint32)
		itemStatus.Owners[buyerAccount] = ownedTokens
	}

	for tokenHash := range itemInfo.Item.Templates {
		resellerTokenCount := resellerTokens[tokenHash]
		if resellerTokenCount == n {
			delete(resellerTokens, tokenHash)
		} else {
			resellerTokens[tokenHash] -= n
		}

		ownedTokens[tokenHash] += n
	}
}

// BuyDToken is called by DTokenBuyer
func (c *DDXFContract) BuyDToken(resourceID string, n uint32, buyerAccount ddxf.OntID) {
	if !c.checkWitness(buyerAccount) {
		panic("buyerAccount no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	if time.Now().Unix() > itemInfo.Item.ExpiredDate {
		panic("resourceID already expired")
	}

	itemStatus := c.sellerItemStatus[resourceID]
	if itemStatus.Sold >= itemInfo.Item.Stocks {
		panic("resourceID already sold out")
	}

	if itemStatus.Sold+n >= itemInfo.Item.Stocks {
		panic("resourceID not enough")
	}

	if !c.transferFeeFromAccount(buyerAccount, itemInfo.ResourceDDO.Manager, itemInfo.Item.Fee, n) {
		panic("balance not enough")
	}

	itemStatus.Sold += n
	ownedTokens := itemStatus.Owners[buyerAccount]
	if ownedTokens == nil {
		ownedTokens = make(map[string]uint32)
		itemStatus.Owners[buyerAccount] = ownedTokens
	}
	for tokenHash := range itemInfo.Item.Templates {
		ownedTokens[tokenHash] += n
	}

}

// UseDToken is called by buyer
func (c *DDXFContract) UseDToken(resourceID string, account ddxf.OntID, tokenHash string, n uint32) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemStatus, ok := c.sellerItemStatus[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	ownedTokens, ok := itemStatus.Owners[account]
	if !ok {
		panic("resourceID not owned by account")
	}

	tokenCount := ownedTokens[tokenHash]
	if tokenCount < n {
		panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
	}

	if tokenCount == n {
		delete(ownedTokens, tokenHash)
		if len(ownedTokens) == 0 {
			delete(itemStatus.Owners, account)
		}
	} else {
		ownedTokens[tokenHash] = tokenCount - n
	}
}

// UseDTokenSuit is called by buyer
func (c *DDXFContract) UseDTokenSuit(resourceID string, account ddxf.OntID, n uint32) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	itemStatus, ok := c.sellerItemStatus[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	ownedTokens, ok := itemStatus.Owners[account]
	if !ok {
		panic("resourceID not owned by account")
	}

	for tokenHash := range itemInfo.Item.Templates {
		if ownedTokens[tokenHash] < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
	}

	var toDelete []string
	for tokenHash, ownCount := range ownedTokens {
		if ownCount == n {
			toDelete = append(toDelete, tokenHash)
		} else {
			ownedTokens[tokenHash] = ownCount - n
		}
	}
	for _, tokenHash := range toDelete {
		delete(ownedTokens, tokenHash)
	}

	if len(ownedTokens) == 0 {
		delete(itemStatus.Owners, account)
	}
}

func (c *DDXFContract) checkWitness(account ddxf.OntID) bool {
	return true
}

func (c *DDXFContract) transferFeeFromAccount(buyerAccount, sellerAccount ddxf.OntID, fee ddxf.Fee, n uint32) bool {
	return false
}

// // BuyerItemInfo for ddxf
// type BuyerItemInfo struct {
// 	Item       ddxf.CrowdSouring
// 	PublishDDO string
// 	Collected  uint32
// 	Providers  map[ddxf.OntID]uint32
// }

// // DTokenBuyerPublish is called by DTokenBuyer
// func (c *DDXFContract) DTokenBuyerPublish(publishID ddxf.OntID, publishDDO string, item ddxf.CrowdSouring) {

// }
