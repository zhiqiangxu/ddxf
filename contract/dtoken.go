package contract

import (
	"fmt"

	"github.com/zhiqiangxu/ddxf"
)

// DTokenContract for dtoken
type DTokenContract interface {
	GenerateDToken(account ddxf.OntID, resourceID string, Templates TokenTemplates, n uint32)
	UseToken(account ddxf.OntID, resourceID string, tokenHash string, n uint32)
	UseTokenByAgent(account, agent ddxf.OntID, resourceID string, tokenHash string, n uint32)
	// for reseller
	TransferDToken(fromAccount, toAccount ddxf.OntID, resourceID string, Templates TokenTemplates, n uint32)

	SetAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID, n uint32)
	SetTokenAgents(account ddxf.OntID, resourceID, tokenHash string, agents []ddxf.OntID, n uint32)
	AddAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID, n uint32)
	AddTokenAgents(account ddxf.OntID, resourceID, tokenHash string, agents []ddxf.OntID, n uint32)
	RemoveAgents(account ddxf.OntID, resourceID string, agents []ddxf.OntID)
	RemoveTokenAgents(account ddxf.OntID, resourceID, tokenHash string, agents []ddxf.OntID)
}

type dTokenContract struct {
	Owners map[string]map[ddxf.OntID]map[string]*CountAndAgent
}

// NewDTokenContract is default implmentation for DTokenContract
func NewDTokenContract() DTokenContract {
	return &dTokenContract{Owners: make(map[string]map[ddxf.OntID]map[string]*CountAndAgent)}
}

func (d *dTokenContract) GenerateDToken(account ddxf.OntID, resourceID string, templates TokenTemplates, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		owners = make(map[ddxf.OntID]map[string]*CountAndAgent)
		d.Owners[resourceID] = owners
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		accountTokens = make(map[string]*CountAndAgent)
		owners[account] = accountTokens
	}

	for tokenHash := range templates {
		caa := accountTokens[tokenHash]
		if caa == nil {
			caa = &CountAndAgent{Agents: make(map[ddxf.OntID]uint32)}
			accountTokens[tokenHash] = caa
		}
		caa.Count += n
	}
	return
}

func (d *dTokenContract) UseToken(account ddxf.OntID, resourceID string, tokenHash string, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenHash]
	if caa == nil {
		panic("account has no tokenHash")
	}

	if !caa.CanDecCount(n) {
		panic("tokenHash not enough")
	}

	if caa.DecCount(n) {
		delete(accountTokens, tokenHash)
		if len(accountTokens) == 0 {
			delete(owners, account)
		}
	}
	return
}

func (d *dTokenContract) UseTokenByAgent(account, agent ddxf.OntID, resourceID string, tokenHash string, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenHash]
	if caa == nil {
		panic("account has no tokenHash")
	}

	if !caa.CanDecCount(n) {
		panic("tokenHash not enough")
	}
	if !caa.CanDecCountByAgent(n, agent) {
		panic("agent count not enough")
	}

	if caa.DecCountByAgent(n, agent) {
		delete(accountTokens, tokenHash)
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
		toAccountTokens = make(map[string]*CountAndAgent)
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

func (d *dTokenContract) SetTokenAgents(account ddxf.OntID, resourceID, tokenHash string, agents []ddxf.OntID, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenHash]
	if caa == nil {
		panic("account has no tokenHash")
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

func (d *dTokenContract) AddTokenAgents(account ddxf.OntID, resourceID, tokenHash string, agents []ddxf.OntID, n uint32) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenHash]
	if caa == nil {
		panic("account has no tokenHash")
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

func (d *dTokenContract) RemoveTokenAgents(account ddxf.OntID, resourceID, tokenHash string, agents []ddxf.OntID) {
	owners := d.Owners[resourceID]
	if owners == nil {
		panic("resourceID not exists")
	}
	accountTokens := owners[account]
	if accountTokens == nil {
		panic("account has no resourceID")
	}

	caa := accountTokens[tokenHash]
	if caa == nil {
		panic("account has no tokenHash")
	}

	caa.RemoveAgents(agents)
	return
}
