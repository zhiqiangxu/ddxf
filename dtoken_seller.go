package ddxf

// DTokenSeller def
type DTokenSeller interface {
	Object

	BindDataDDOWithDataSource(dataDDO DataDDO, s StaticDataSource) error
	LookupDataMetaByDataDDO(dataDDO DataDDO, mpID OntID) (string, error)

	BindMPItemWithDataSource(mpItemID string, mpID OntID, synthetic StaticDataSource) error

	MakeTx(item DTokenItem, mpID OntID) (string, error)
}
