package ddxf

// OntID is the id for arbitory party
type OntID string

// Object def
type Object interface {
	OntID() OntID
}

// MPDDO is ddo for marketplace
type MPDDO struct {
	Endpoint string // data service provider uri
}

// MarketPlaceItemMeta is meta for marketplace item
type MarketPlaceItemMeta struct {
	Version uint32
}

// MarketPlaceSellerItem def
type MarketPlaceSellerItem struct {
	Seller     OntID
	DataID     OntID
	DTokenItem DTokenItem
}

// MarketPlace def
type MarketPlace interface {
	Object
	ItemMeta() MarketPlaceItemMeta
	ItemTemplate() string
	ValidateAndPublishItem(item MarketPlaceSellerItem, itemInfo []byte) (mpItemID string, err error)

	EnsureMPDDORegistered() error
}

// MarketPlaceFramework def
type MarketPlaceFramework interface {
	HandleDTokenSellerPublish(item MarketPlaceSellerItem, tx string) (mpItemID string, err error)
	MakeTxForDTokenBuyer(mpItemID string) (tx string, err error)
}

type marketPlaceFramework struct {
	mp MarketPlace
}

// NewMarketPlaceFramework is ctor for MarketPlaceFramework
func NewMarketPlaceFramework(mp MarketPlace) MarketPlaceFramework {
	return &marketPlaceFramework{mp: mp}
}

func (fw *marketPlaceFramework) HandleDTokenSellerPublish(item MarketPlaceSellerItem, tx string) (mpItemID string, err error) {

	mpItemInfo, err := StdInstance().GetDataMeta(item.DataID, fw.mp.OntID())
	if err != nil {
		return
	}

	mpItemID, err = fw.mp.ValidateAndPublishItem(item, mpItemInfo)
	if err != nil {
		return
	}

	err = fw.signAndSendToChain(tx)
	return
}

func (fw *marketPlaceFramework) signAndSendToChain(tx string) (err error) {
	return
}

func (fw *marketPlaceFramework) MakeTxForDTokenBuyer(mpItemID string) (tx string, err error) {
	return
}
