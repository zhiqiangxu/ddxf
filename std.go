package ddxf

import (
	"sync"
)

// Std for ddxf
type Std struct {
}

var (
	instanceStd *Std
	lockStd     sync.Mutex
)

// StdInstance is singleton for Std
func StdInstance() *Std {
	if instanceStd != nil {
		return instanceStd
	}

	lockStd.Lock()
	defer lockStd.Unlock()
	if instanceStd != nil {
		return instanceStd
	}

	instanceStd = &Std{}

	return instanceStd
}

// PublishSellerItem publishes seller item to marketplace
func (std *Std) PublishSellerItem(item MarketPlaceSellerItem, mpID OntID, tx string) (mpItemID string, err error) {

	// TODO
	return
}

// ResolveMPDDO resolves mp ddo
func (std *Std) ResolveMPDDO(mpID OntID) (mpDDO MPDDO, err error) {
	// TODO
	return
}

// ResolveDataDDO resolves data ddo
func (std *Std) ResolveDataDDO(dataID OntID) (dataDDO DataDDO, err error) {
	// TODO
	return
}

// GetDataMeta retrieves data meta
func (std *Std) GetDataMeta(dataID, mp OntID) (dataMeta []byte, err error) {
	dataDDO, err := std.ResolveDataDDO(dataID)
	if err != nil {
		return
	}

	dataMeta, err = std.resolveDataMeta(dataDDO, mp)
	return
}

func (std *Std) resolveDataMeta(ddo DataDDO, mp OntID) (dataMeta []byte, err error) {

	// reqBytes, err := json.Marshal(struct {
	// 	Manager  HostID
	// 	MetaHash string
	// 	Hash     string
	// 	MP       OntID
	// }{
	// 	Manager:  ddo.Manager,
	// 	MetaHash: ddo.MetaHash,
	// 	Hash:     ddo.Hash,
	// 	MP:       mp,
	// })
	// if err != nil {
	// 	return
	// }

	// code, _, dataMeta, err := forward.PostJSONRequest(ddo.Endpoint, reqBytes)
	// if err != nil {
	// 	return
	// }

	// if code != 200 {
	// 	err = fmt.Errorf("ResolveData: code (%d) != 200", code)
	// 	return
	// }
	return
}
