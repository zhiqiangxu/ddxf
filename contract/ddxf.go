package contract

import (
	"fmt"
	"time"

	"crypto/sha256"
	"hash/crc32"

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
	Manager           ddxf.OntID // data owner id
	ResourceType      RT
	Endpoint          string            // data service provider uri
	TokenEndpoint     map[string]string // endpoint for tokens
	TokenResourceType map[string]RT     // RT for tokens
	DescHash          string            // required if len(Templates) > 1
	DTC               DTokenContract    // can be empty
	MP                Marketplace       // can be empty
	Split             SplitPolicy       // can be empty
}

// ResourceTypeForToken ...
func (ddo *ResourceDDO) ResourceTypeForToken(tokenHash string) RT {
	rt, ok := ddo.TokenResourceType[tokenHash]
	if ok {
		return rt
	}

	return ddo.ResourceType
}

// EndpointForToken ...
func (ddo *ResourceDDO) EndpointForToken(tokenHash string) string {
	ep, ok := ddo.TokenEndpoint[tokenHash]
	if ok {
		return ep
	}

	return ddo.Endpoint
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

// NewDDXFContract is ctor for DDXFContract
func NewDDXFContract(dftDtc DTokenContract) *DDXFContract {
	return &DDXFContract{sellerItemSold: make(map[string]uint32), dftDtc: dftDtc}
}

// DTokenSellerPublish is called by DTokenSeller
func (c *DDXFContract) DTokenSellerPublish(resourceID string, resourceDDO ResourceDDO, item DTokenItem) {
	if resourceDDO.Manager == "" {
		panic("manager empty")
	}
	if !c.checkWitness(resourceDDO.Manager) {
		panic("manager no witness")
	}

	if _, ok := c.sellerItemInfo[resourceID]; ok {
		panic("resourceID already exists")
	}

	if resourceDDO.Endpoint == "" {
		if len(resourceDDO.TokenEndpoint) == 0 {
			panic("endpoint empty")
		}

		for tokenHash := range item.Templates {
			if resourceDDO.TokenEndpoint[tokenHash] == "" {
				panic(fmt.Sprintf("endpoint empty not tokenHash:%s", tokenHash))
			}
		}
	}

	if len(item.Templates) == 0 {
		panic("template empty")
	}

	for tokenHash := range item.Templates {
		rt := resourceDDO.ResourceTypeForToken(tokenHash)

		switch rt {
		case RTStaticFile:
			// desc hash + data hash
			if len(tokenHash) != sha256.Size+crc32.Size {
				panic(fmt.Sprintf("invalid tokenHash %s", tokenHash))
			}
		}
	}

	if len(item.Templates) > 1 && len(resourceDDO.DescHash) != sha256.Size {
		panic("ResourceDDO.DescHash invalid for batched template")
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

	if !c.transferFeeFromAccount(buyerAccount, resellerAccount, nil, nil, itemInfo.Item.Fee, n) {
		panic("buyerAccount balance not enough")
	}

	dtc := itemInfo.ResourceDDO.DTC
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

	if !c.transferFeeFromAccount(buyerAccount, itemInfo.ResourceDDO.Manager, itemInfo.ResourceDDO.MP, itemInfo.ResourceDDO.Split, itemInfo.Item.Fee, n) {
		panic("balance not enough")
	}

	c.sellerItemSold[resourceID] += n

	dtc := itemInfo.ResourceDDO.DTC
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

	dtc := itemInfo.ResourceDDO.DTC
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

	dtc := itemInfo.ResourceDDO.DTC
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

	dtc := itemInfo.ResourceDDO.DTC
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

	dtc := itemInfo.ResourceDDO.DTC
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

	dtc := itemInfo.ResourceDDO.DTC
	if dtc == nil {
		dtc = c.dftDtc
	}
	dtc.RemoveAgents(account, resourceID, agents)
	return
}

func (c *DDXFContract) checkWitness(account ddxf.OntID) bool {
	return true
}

func (c *DDXFContract) transferFeeFromAccount(buyerAccount, sellerAccount ddxf.OntID, mp Marketplace, split SplitPolicy, fee ddxf.Fee, n uint32) bool {
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
