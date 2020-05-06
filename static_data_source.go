package ddxf

// DataDDO is ddo for data
type DataDDO struct {
	Manager  HostID // data owner id
	Endpoint string // data service provider uri
	MetaHash string // hash of meta description of data service, hash of json string in key alphabetical sequence
	Hash     string // hash for data source, optional
}

// HostID is owner for data or service
type HostID OntID

// StaticDataSource def
type StaticDataSource interface {
	Object
	DataDDO() DataDDO
	DataMeta(mp OntID) string

	EnsureDataDDORegistered() error
}
