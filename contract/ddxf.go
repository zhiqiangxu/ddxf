package contract

import (
	"fmt"
	"time"

	"strings"

	"github.com/zhiqiangxu/ddxf"
)

// TokenTemplates for ddxf
type TokenTemplates map[string] /*tokenHash*/ struct{}

// DTokenItem for ddxf
type DTokenItem struct {
	Fee         ddxf.Fee
	ExpiredDate int64
	Stocks      uint32
	Templates   TokenTemplates
}

// RT for resource type
type RT int

const (
	// RTStaticFile for static file
	RTStaticFile RT = iota
)

// ResourceDDO is ddo for resource
type ResourceDDO struct {
	ResourceType RT
	Manager      ddxf.OntID     // data owner id
	Endpoint     string         // data service provider uri
	DescHash     string         // required if len(Templates) > 1
	dtc          DTokenContract // can be empty
}

// SellerItemInfo for ddxf
// immutable
type SellerItemInfo struct {
	Item        DTokenItem
	ResourceDDO ResourceDDO
}

// DDXFContract for ddxf
type DDXFContract struct {
	sellerItemInfo map[string]SellerItemInfo
	sellerItemSold map[string]uint32
	dftDtc         DTokenContract
}

var emptyResourceDDO = ResourceDDO{}

// NewDDXFContract is ctor for DDXFContract
func NewDDXFContract(dftDtc DTokenContract) *DDXFContract {
	return &DDXFContract{sellerItemSold: make(map[string]uint32), dftDtc: dftDtc}
}

// DTokenSellerPublish is called by DTokenSeller
func (c *DDXFContract) DTokenSellerPublish(resourceID string, resourceDDO ResourceDDO, item DTokenItem) {
	if _, ok := c.sellerItemInfo[resourceID]; ok {
		panic("resourceID already exists")
	}

	if resourceDDO == emptyResourceDDO {
		panic("resourceDDO empty")
	}

	if len(item.Templates) == 0 {
		panic("template empty")
	}

	switch resourceDDO.ResourceType {
	case RTStaticFile:
		for tokenHash := range item.Templates {
			// desc hash : data hash
			if len(strings.Split(tokenHash, ":")) != 2 {
				panic(fmt.Sprintf("invalid tokenHash %s", tokenHash))
			}
		}
	}

	if len(item.Templates) > 1 && resourceDDO.DescHash == "" {
		panic("ResourceDDO.DescHash empty for batched template")
	}

	c.sellerItemInfo[resourceID] = SellerItemInfo{Item: item, ResourceDDO: resourceDDO}

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

	if !c.transferFeeFromAccount(buyerAccount, resellerAccount, itemInfo.Item.Fee, n) {
		panic("buyerAccount balance not enough")
	}

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.TransferDToken(resellerAccount, buyerAccount, resourceID, itemInfo.Item.Templates, n)

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

	sold := c.sellerItemSold[resourceID]
	if sold >= itemInfo.Item.Stocks {
		panic("resourceID already sold out")
	}

	if sold+n >= itemInfo.Item.Stocks {
		panic("resourceID not enough")
	}

	if !c.transferFeeFromAccount(buyerAccount, itemInfo.ResourceDDO.Manager, itemInfo.Item.Fee, n) {
		panic("balance not enough")
	}

	c.sellerItemSold[resourceID] += n

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.GenerateDToken(buyerAccount, resourceID, itemInfo.Item.Templates, n)

}

// UseToken is called by buyer
func (c *DDXFContract) UseToken(resourceID string, account ddxf.OntID, tokenHash string, n uint32) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.UseToken(account, resourceID, tokenHash, n)
}

// UseTokenByAgent is called by agent
func (c *DDXFContract) UseTokenByAgent(resourceID string, account, agent ddxf.OntID, tokenHash string, n uint32) {
	if !c.checkWitness(agent) {
		panic("agent no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.UseTokenByAgent(account, agent, resourceID, tokenHash, n)
}

// SetDTokenAgents is called by buyer
func (c *DDXFContract) SetDTokenAgents(resourceID string, account ddxf.OntID, agents []ddxf.OntID, n uint32) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.SetAgents(account, resourceID, agents, n)

	return
}

// AddDTokenAgents is called by buyer
func (c *DDXFContract) AddDTokenAgents(resourceID string, account ddxf.OntID, agents []ddxf.OntID, n uint32) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.AddAgents(account, resourceID, agents, n)

	return
}

// RemoveDTokenAgents is called by buyer
func (c *DDXFContract) RemoveDTokenAgents(resourceID string, account ddxf.OntID, agents []ddxf.OntID) {
	if !c.checkWitness(account) {
		panic("account no witness")
	}

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	dtc := itemInfo.ResourceDDO.dtc
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.RemoveAgents(account, resourceID, agents)
	return
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
