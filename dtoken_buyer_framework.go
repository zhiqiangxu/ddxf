package ddxf

// DTokenBuyer def
type DTokenBuyer interface {
	Object
}

// DTokenBuyerFramework def
type DTokenBuyerFramework interface {
	Buy(mpID OntID, mpItemID string) error
}

type dTokenBuyerFramework struct {
	dTokenBuyer DTokenBuyer
}

// NewDTokenBuyerFramework is ctor for DTokenBuyerFramework
func NewDTokenBuyerFramework(dTokenBuyer DTokenBuyer) DTokenBuyerFramework {
	return &dTokenBuyerFramework{dTokenBuyer: dTokenBuyer}
}

func (fw *dTokenBuyerFramework) Buy(mpID OntID, mpItemID string) (err error) {

	return
}
