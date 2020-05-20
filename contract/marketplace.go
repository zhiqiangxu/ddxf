package contract

import "github.com/zhiqiangxu/ddxf"

// FeeSplitModel ...
type FeeSplitModel struct {
	Percentage uint16 // decimal = 2
}

// Marketplace ...
// this contract is agreed by seller and marketplace
type Marketplace interface {
	// should sign by both seller and mp
	SetFeeSplitModel(sellerAcc ddxf.OntID, model FeeSplitModel)
	GetFeeSplitModel(sellerAcc ddxf.OntID) FeeSplitModel

	Settle(sellerAcc ddxf.OntID)
	TransferAmount(buyerAcc, sellerAcc ddxf.OntID, fee ddxf.Fee, n uint32)

	MPAccount() ddxf.OntID
}

type marketplace struct {
}

// NewMarketplace ...
func NewMarketplace() Marketplace {
	return &marketplace{}
}

func (m *marketplace) SetFeeSplitModel(sellerAcc ddxf.OntID, model FeeSplitModel) {

}

func (m *marketplace) GetFeeSplitModel(sellerAcc ddxf.OntID) FeeSplitModel {
	return FeeSplitModel{}
}

func (m *marketplace) Settle(sellerAcc ddxf.OntID) {

}
func (m *marketplace) TransferAmount(buyerAcc, sellerAcc ddxf.OntID, fee ddxf.Fee, n uint32) {

}

func (m *marketplace) MPAccount() ddxf.OntID {
	return ddxf.OntID("")
}
