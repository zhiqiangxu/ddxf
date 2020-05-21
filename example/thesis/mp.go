package thesis

import (
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ddxf/contract"
)

// MPThesis ...
type MPThesis struct {
	Title    string
	Body     string
	Extra    string
	DataHash string
}

// Context ...
func (t *MPThesis) Context() map[string]interface{} {
	return map[string]interface{}{
		"Title":    "string",
		"Body":     "string",
		"Extra":    "string",
		"DataHash": "string",
	}
}

// MP ...
type MP struct {
	MPC contract.Marketplace
}

// NewMP ...
func NewMP() *MP {
	return &MP{MPC: contract.NewMarketplace()}
}

type (
	// PublishMPInput ...
	PublishMPInput struct {
		Fee         ddxf.Fee
		ExpiredDate int64
		Stocks      uint32
		MPDesc      string
		MPTheses    []MPThesis
	}
	// PublishMPOutput ...
	PublishMPOutput struct {
		OK     bool
		ItemID string
	}
)

// PublishMP ...
func (mp *MP) PublishMP(input PublishMPInput) (output PublishMPOutput) {
	return
}

type (
	// SetResourceIDInput ...
	SetResourceIDInput struct {
		ItemID     string
		ResourceID string
	}
	// SetResourceIDOutput ...
	SetResourceIDOutput struct {
		OK bool
	}
)

// SetResourceID ...
func (mp *MP) SetResourceID(input SetResourceIDInput) (output SetResourceIDOutput) {
	return
}

// JSONLDTypes ...
func (mp *MP) JSONLDTypes() []string {
	return []string{"thesis"}
}

type (
	// JSONLDInput ...
	JSONLDInput struct {
		Type string
	}
	// JSONLDOutput ...
	JSONLDOutput struct {
		Context map[string]interface{}
	}
)

// JSONLD ...
func (mp *MP) JSONLD(input JSONLDInput) (output JSONLDOutput) {
	if input.Type != "thesis" {
		return
	}

	output.Context = map[string]interface{}{
		"Fee":         "http://schema.org/fee",
		"ExpiredDate": "http://schema.org/expireData",
		"Stocks":      "http://schema.org/stocks",
		"MPDesc":      "http://schema.org/desc",
		"MPTheses":    "http://schema.org/MPTheses",
	}
	return
}
