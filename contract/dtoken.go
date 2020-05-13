package contract

import (
	"fmt"

	"github.com/zhiqiangxu/ddxf"
)

// DTokenContract for dtoken
type DTokenContract interface {
	GenerateDToken(account ddxf.OntID, resourceID string, Templates TokenTemplates, n uint32)
	UseToken(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, n uint32)
	UseTokenByAgent(account, agent ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, n uint32)
	// for reseller
	TransferDToken(fromAccount, toAccount ddxf.OntID, resourceID string, Templates TokenTemplates, n uint32)

	SetAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID, n uint32)
	SetTokenAgents(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, agents []ddxf.OntID, n uint32)
	AddAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID, n uint32)
	AddTokenAgents(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, agents []ddxf.OntID, n uint32)
	RemoveAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID)
	RemoveTokenAgents(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, agents []ddxf.OntID)
}

// CountAndAgent for ddxf
type CountAndAgent struct {
	Count  uint32
	Agents map[ddxf.OntID]uint32
}

// IncCount for increase Count
func (caa *CountAndAgent) IncCount(n uint32) {
	caa.Count += n
}

// CanDecCount checks whether can DecCount
func (caa *CountAndAgent) CanDecCount(n uint32) bool {
	return caa.Count >= n
}

// ClearAgents clears all agents
func (caa *CountAndAgent) ClearAgents() {
	for agent := range caa.Agents {
		delete(caa.Agents, agent)
	}
}

// RemoveAgents for CountAndAgent
func (caa *CountAndAgent) RemoveAgents(agents []ddxf.OntID) {
	for _, agent := range agents {
		delete(caa.Agents, agent)
	}
}

// AddAgents for CountAndAgent
func (caa *CountAndAgent) AddAgents(agents []ddxf.OntID, n uint32) {
	for _, agent := range agents {
		caa.Agents[agent] += n
	}
}

// DecCount for decrease Count
func (caa *CountAndAgent) DecCount(n uint32) (usedup bool) {
	caa.Count -= n
	usedup = caa.Count == 0
	return
}

// CanDecCountByAgent checks whether agent can decrease Count
func (caa *CountAndAgent) CanDecCountByAgent(n uint32, agent ddxf.OntID) bool {
	return caa.Agents[agent] >= n
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

type dTokenContract struct {
	Owners map[string]map[ddxf.OntID]map[TokenTemplate]*CountAndAgent
}

// NewDTokenContract is default implmentation for DTokenContract
func NewDTokenContract() DTokenContract {
	return &dTokenContract{Owners: make(map[string]map[ddxf.OntID]map[TokenTemplate]*CountAndAgent)}
}

func (d *dTokenContract) GenerateDToken(account ddxf.OntID, resourceID string, templates TokenTemplates, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		owners = make(map[ddxf.OntID]map[TokenTemplate]*CountAndAgent)
		d.Owners[resourceID] = owners
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		accountTokens = make(map[TokenTemplate]*CountAndAgent)
		owners[account] = accountTokens
	}

	for tokenTemplate := range templates {
		caa := accountTokens[tokenTemplate]
		if caa == nil {
			caa = &CountAndAgent{Agents: make(map[ddxf.OntID]uint32)}
			accountTokens[tokenTemplate] = caa
		}
		caa.Count += n
	}
	return
}

func (d *dTokenContract) UseToken(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenTemplate]
	if caa == nil {
		panic("account has no tokenTemplate")
	}

	if !caa.CanDecCount(n) {
		panic("tokenTemplate not enough")
	}

	if caa.DecCount(n) {
		delete(accountTokens, tokenTemplate)
		if len(accountTokens) == 0 {
			delete(owners, account)
		}
	}
	return
}

func (d *dTokenContract) UseTokenByAgent(account, agent ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenTemplate]
	if caa == nil {
		panic("account has no tokenTemplate")
	}

	if !caa.CanDecCount(n) {
		panic("tokenTemplate not enough")
	}
	if !caa.CanDecCountByAgent(n, agent) {
		panic("agent count not enough")
	}

	if caa.DecCountByAgent(n, agent) {
		delete(accountTokens, tokenTemplate)
		if len(accountTokens) == 0 {
			delete(owners, account)
		}
	}
}

func (d *dTokenContract) TransferDToken(fromAccount, toAccount ddxf.OntID, resourceID string, templates TokenTemplates, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	fromAccountTokens := owners[fromAccount]
	if fromAccountTokens == nil {
		panic("account has no resourceID")
	}

	// check first
	for tokenHash := range templates {
		caa := fromAccountTokens[tokenHash]
		if caa == nil {
			panic(fmt.Sprintf("fromAccount has no tokenHash:%s", tokenHash))
		}
		if !caa.CanDecCount(n) {
			panic(fmt.Sprintf("fromAccount has no enough tokenHash:%s", tokenHash))
		}
	}

	toAccountTokens := owners[toAccount]
	if toAccountTokens == nil {
		toAccountTokens = make(map[TokenTemplate]*CountAndAgent)
		owners[toAccount] = toAccountTokens
	}

	// then transfer
	for tokenHash := range templates {
		fromAccountTokens[tokenHash].DecCount(n)
		toCaa := toAccountTokens[tokenHash]
		if toCaa == nil {
			toCaa = &CountAndAgent{Agents: make(map[ddxf.OntID]uint32)}
			toCaa.IncCount(n)
			toAccountTokens[tokenHash] = toCaa
		}
	}

	return
}

func (d *dTokenContract) SetAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	for _, caa := range accountTokens {
		caa.ClearAgents()
		caa.AddAgents(agents, n)
	}

	return
}

func (d *dTokenContract) SetTokenAgents(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, agents []ddxf.OntID, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenTemplate]
	if caa == nil {
		panic("account has no tokenTemplate")
	}

	caa.ClearAgents()
	caa.AddAgents(agents, n)
	return
}

func (d *dTokenContract) AddAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	for _, caa := range accountTokens {
		caa.AddAgents(agents, n)
	}

	return
}

func (d *dTokenContract) AddTokenAgents(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, agents []ddxf.OntID, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenTemplate]
	if caa == nil {
		panic("account has no tokenTemplate")
	}

	caa.AddAgents(agents, n)
	return
}

func (d *dTokenContract) RemoveAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	for _, caa := range accountTokens {
		caa.RemoveAgents(agents)
	}

	return
}

func (d *dTokenContract) RemoveTokenAgents(account ddxf.OntID, resourceID string, tokenTemplate TokenTemplate, agents []ddxf.OntID) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenTemplate]
	if caa == nil {
		panic("account has no tokenTemplate")
	}

	caa.RemoveAgents(agents)
	return
}
