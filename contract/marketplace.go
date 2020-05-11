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
