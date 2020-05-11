package contract

import "github.com/zhiqiangxu/ddxf"

// FeeSplitModel ...
type FeeSplitModel struct {
	Percentage uint16 // decimal = 2
}

// Marketplace ...
// this contract is agreed by seller and marketplace
type Marketplace interface {
	SetFeeSplitModel(FeeSplitModel)
	GetFeeSplitModel() FeeSplitModel

	Settle()

	MPAccount() ddxf.OntID
	SellerAccount() ddxf.OntID
}
