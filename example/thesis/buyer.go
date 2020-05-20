package thesis

import (
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ddxf/contract"
)

// Buyer ...
type Buyer struct {
}

const (
	buyerID = ddxf.OntID("buyer xxx")
)

type (
	// BuyDtokenInput ...
	BuyDtokenInput struct {
		ResourceID string
		N          uint32
	}

	// BuyDtokenOutput ...
	BuyDtokenOutput struct {
		TokenTemplates contract.TokenTemplates
	}
)

// NewBuyer ...
func NewBuyer() *Buyer {
	return &Buyer{}
}

// BuyDtoken ...
func (b *Buyer) BuyDtoken(input BuyDtokenInput) (output BuyDtokenOutput) {
	output.TokenTemplates = DDXF().BuyDToken(input.ResourceID, input.N, buyerID)
	return
}
