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
	Templates   map[string] /*tokenHash*/ struct{}
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
	Manager      ddxf.OntID // data owner id
	Endpoint     string     // data service provider uri
	DescHash     string     // required if len(Templates) > 1
}

// SellerItemInfo for ddxf
// immutable
type SellerItemInfo struct {
	Item        DTokenItem
	ResourceDDO ResourceDDO
}

// CountAndAgent for ddxf
type CountAndAgent struct {
	Count  uint32
	Agents map[ddxf.OntID]uint32
}

// DecCount for decrease Count
func (caa *CountAndAgent) DecCount(n uint32) (usedup bool) {
	caa.Count -= n
	usedup = caa.Count == 0
	return
}

// DecCountByAgent for decrease Count by agent
func (caa *CountAndAgent) DecCountByAgent(n uint32, agent ddxf.OntID) (usedup bool) {
	caa.Count -= n
	caa.Agents[agent] -= n
	if caa.Agents[agent] == 0 {
		delete(caa.Agents, agent)
	}

	usedup = caa.Count == 0
	return
}

// SellerItemStatus for ddxf
// mutable
type SellerItemStatus struct {
	Sold   uint32
	Owners map[ddxf.OntID]map[string]*CountAndAgent
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
		panic("ResourceDDO.Hash empty for batched template")
	}

	c.sellerItemInfo[resourceID] = SellerItemInfo{Item: item, ResourceDDO: resourceDDO}
	c.sellerItemStatus[resourceID] = SellerItemStatus{Owners: make(map[ddxf.OntID]map[string]*CountAndAgent)}

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
		if resellerTokens[tokenHash].Count < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
	}

	if !c.transferFeeFromAccount(buyerAccount, resellerAccount, itemInfo.Item.Fee, n) {
		panic("balance not enough")
	}

	ownedTokens := itemStatus.Owners[buyerAccount]
	if ownedTokens == nil {
		ownedTokens = make(map[string]*CountAndAgent)
		itemStatus.Owners[buyerAccount] = ownedTokens
	}

	for tokenHash := range itemInfo.Item.Templates {
		resellerToken := resellerTokens[tokenHash]
		if resellerToken.DecCount(n) {
			delete(resellerTokens, tokenHash)
		}

		ownedTokens[tokenHash].Count += n
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
		ownedTokens = make(map[string]*CountAndAgent)
		itemStatus.Owners[buyerAccount] = ownedTokens
	}
	for tokenHash := range itemInfo.Item.Templates {
		ownedTokens[tokenHash].Count += n
	}

}

// UseDToken is called by seller, but signed by buyer
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

	ownedToken := ownedTokens[tokenHash]
	if ownedToken == nil || ownedToken.Count < n {
		panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
	}

	if ownedToken.DecCount(n) {
		delete(ownedTokens, tokenHash)
		if len(ownedTokens) == 0 {
			delete(itemStatus.Owners, account)
		}
	}
}

// UseDTokenSuitByAgent is called by agent
func (c *DDXFContract) UseDTokenSuitByAgent(resourceID string, account, agent ddxf.OntID, n uint32) {
	if !c.checkWitness(account) {
		panic("agent no witness")
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
		ownedToken := ownedTokens[tokenHash]
		if ownedToken == nil || ownedToken.Count < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
		if ownedToken.Agents[agent] < n {
			panic(fmt.Sprintf("resourceID not allowed enough(%d) for token:%s", ownedToken.Agents[agent], tokenHash))
		}
	}

	var toDelete []string
	for tokenHash, ownedToken := range ownedTokens {
		if ownedToken.DecCountByAgent(n, agent) {
			toDelete = append(toDelete, tokenHash)
		}
	}
	for _, tokenHash := range toDelete {
		delete(ownedTokens, tokenHash)
	}

	if len(ownedTokens) == 0 {
		delete(itemStatus.Owners, account)
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
		ownedToken := ownedTokens[tokenHash]
		if ownedToken == nil || ownedToken.Count < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
	}

	var toDelete []string
	for tokenHash, ownedToken := range ownedTokens {
		if ownedToken.DecCount(n) {
			toDelete = append(toDelete, tokenHash)
		}
	}
	for _, tokenHash := range toDelete {
		delete(ownedTokens, tokenHash)
	}

	if len(ownedTokens) == 0 {
		delete(itemStatus.Owners, account)
	}
}

// SetDTokenSuitAgents is called by buyer
func (c *DDXFContract) SetDTokenSuitAgents(resourceID string, account ddxf.OntID, agents []ddxf.OntID, n uint32) {
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

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	for tokenHash := range itemInfo.Item.Templates {
		ownedToken := ownedTokens[tokenHash]
		if ownedToken == nil || ownedToken.Count < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
	}

	for tokenHash := range itemInfo.Item.Templates {
		agentAndCount := make(map[ddxf.OntID]uint32)
		for _, agent := range agents {
			agentAndCount[agent] = n
		}
		ownedTokens[tokenHash].Agents = agentAndCount
	}

	return
}

// AddDTokenSuitAgents is called by buyer
func (c *DDXFContract) AddDTokenSuitAgents(resourceID string, account ddxf.OntID, agents []ddxf.OntID, n uint32) {
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

	itemInfo, ok := c.sellerItemInfo[resourceID]
	if !ok {
		panic("resourceID not exists")
	}

	for tokenHash := range itemInfo.Item.Templates {
		ownedToken := ownedTokens[tokenHash]
		if ownedToken == nil || ownedToken.Count < n {
			panic(fmt.Sprintf("resourceID owned not enough for token:%s", tokenHash))
		}
	}

	for tokenHash := range itemInfo.Item.Templates {
		ownedToken := ownedTokens[tokenHash]
		for _, agent := range agents {
			ownedToken.Agents[agent] += n
		}
	}
	return
}

// RemoveDTokenSuitAgents is called by buyer
func (c *DDXFContract) RemoveDTokenSuitAgents(resourceID string, account ddxf.OntID, agents []ddxf.OntID) {
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

	for _, caa := range ownedTokens {
		for _, agent := range agents {
			delete(caa.Agents, agent)
		}
	}
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
