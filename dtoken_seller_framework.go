package ddxf

// DTokenSellerFramework def
type DTokenSellerFramework interface {
	PublishBatch(synthetic StaticDataSource, ss []StaticDataSource, tokens [][]Token, mpID OntID, fee Fee, opt MarketPlaceItemOption) (err error)

	PublishSingle(s StaticDataSource, tokens []Token, mpID OntID, fee Fee, opt MarketPlaceItemOption) (err error)
}

// NewDTokenSellerFramework is ctor for DTokenSellerFramework
func NewDTokenSellerFramework(dTokenSeller DTokenSeller) DTokenSellerFramework {
	return &dTokenSellerFramework{dTokenSeller: dTokenSeller}
}

type dTokenSellerFramework struct {
	dTokenSeller DTokenSeller
}

func (fw *dTokenSellerFramework) PublishBatch(synthetic StaticDataSource, ss []StaticDataSource, tokens [][]Token, mpID OntID, fee Fee, opt MarketPlaceItemOption) (err error) {

	err = fw.prepareStaticDataSource(synthetic)
	if err != nil {
		return
	}

	for _, s := range ss {
		err = fw.prepareStaticDataSource(s)
		if err != nil {
			return
		}
	}

	assert(len(ss) == len(tokens), "len(ss) != len(tokens)")

	// generate DTokenItem
	templates := make([]TokenTemplate, len(ss))
	for i, s := range ss {
		templates[i] = TokenTemplate{DataID: s.OntID()}
		for _, token := range tokens[i] {
			templates[i].TokenHashes = append(templates[i].TokenHashes, Sha256Bytes(token.AlphabeticalBytes()))
		}
	}
	di := DTokenItem{Fee: fee, ExpiredDate: opt.ExpiredDate, Stocks: opt.Stocks, Templates: templates}

	sellerItem := MarketPlaceSellerItem{Seller: fw.dTokenSeller.OntID(), DataID: synthetic.OntID(), DTokenItem: di}
	mpItemID, err := fw.publish(sellerItem, mpID)
	if err != nil {
		return
	}

	// bind mpItemID with synthetic StaticDataSource locally
	err = fw.dTokenSeller.BindMPItemWithDataSource(mpItemID, mpID, synthetic)
	if err != nil {
		return
	}

	return
}

func (fw *dTokenSellerFramework) publish(item MarketPlaceSellerItem, mpID OntID) (mpItemID string, err error) {

	tx, err := fw.dTokenSeller.MakeTx(item.DTokenItem, mpID)
	if err != nil {
		return
	}

	mpItemID, err = StdInstance().PublishSellerItem(item, mpID, tx)
	return
}

func (fw *dTokenSellerFramework) PublishSingle(s StaticDataSource, tokens []Token, mpID OntID, fee Fee, opt MarketPlaceItemOption) (err error) {
	err = fw.prepareStaticDataSource(s)
	if err != nil {
		return
	}

	template := TokenTemplate{DataID: s.OntID()}
	for _, token := range tokens {
		template.TokenHashes = append(template.TokenHashes, Sha256Bytes(token.AlphabeticalBytes()))
	}
	di := DTokenItem{Fee: fee, ExpiredDate: opt.ExpiredDate, Stocks: opt.Stocks, Templates: []TokenTemplate{template}}

	sellerItem := MarketPlaceSellerItem{Seller: fw.dTokenSeller.OntID(), DataID: s.OntID(), DTokenItem: di}
	mpItemID, err := fw.publish(sellerItem, mpID)
	if err != nil {
		return
	}

	// bind mpItemID with StaticDataSource locally
	err = fw.dTokenSeller.BindMPItemWithDataSource(mpItemID, mpID, s)
	if err != nil {
		return
	}

	return
}

func (fw *dTokenSellerFramework) prepareStaticDataSource(s StaticDataSource) (err error) {

	err = s.EnsureDataDDORegistered()
	if err != nil {
		return
	}

	err = fw.dTokenSeller.BindDataDDOWithDataSource(s.DataDDO(), s)
	return
}
