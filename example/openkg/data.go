package openkg

import (
	"fmt"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
)

// Data ...
type Data struct {
	Name        string
	Download    int
	Maintainers []string
	URL         string
}

// PublishData ...
func PublishData(d Data) (dataID, resourceID string, private keypair.PrivateKey, err error) {

	// 1. generate ont id

	// Generate key pair
	private, public, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	if err != nil {
		return
	}

	dataID = ""

	tx, err := sdk.Native.OntId.NewRegIDWithPublicKeyTransaction(1000, 8000000, dataID, public)
	if err != nil {
		return
	}

	err = sdk.SignToTransaction(tx, signer)
	if err != nil {
		return
	}

	// sdk.PreExecTransaction
	_, err = sdk.SendTransaction(tx)
	if err != nil {
		return
	}

	// 2. add Maintainers as controllers

	// 3. publish to ddxf

	method := "dtokenSellerPublish"
	txHash, err := sdk.WasmVM.InvokeWasmVMSmartContract(1000, 8000000, signer, signer, contractAddress, method, []interface{}{})
	if err != nil {
		return
	}

	timeoutSec := 30 * time.Second

	_, err = sdk.WaitForGenerateBlock(timeoutSec)
	if err != nil {
		return
	}

	fmt.Printf("method:%s, txHash:%s\n", method, txHash.ToHexString())
	event, err := sdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		fmt.Println("GetSmartContractEvent error ", err)
		return
	}
	if event != nil {
		for _, notify := range event.Notify {
			fmt.Printf("%+v\n", notify)
		}
	}
	// sdk.PreExecTransaction

	return
}

var (
	sdk             *osdk.OntologySdk
	signer          *osdk.Account
	walletName      string
	signerAddr      string
	contractAddress common.Address
)

func init() {

	sdk = osdk.NewOntologySdk()
	if true {
		// test
		sdk.NewRpcClient().SetAddress("http://polaris1.ont.io:20336")
	} else {
		// prod
		sdk.NewRpcClient().SetAddress("http://dappnode1.ont.io:20336")
	}

	wallet, err := sdk.OpenWallet(walletName)
	if err != nil {
		panic(fmt.Sprintf("error in OpenWallet:%s\n", err))
	}

	passwd, err := password.GetAccountPassword()
	if err != nil {
		panic(fmt.Sprintf("input password error %s", err))
	}

	signer, err = wallet.GetAccountByAddress(signerAddr, passwd)
	if err != nil {
		panic(fmt.Sprintf("error in GetAccountByAddress:%s\n", err))
	}

	contractAddress, err = common.AddressFromHexString("xxxx")
	if err != nil {
		panic(fmt.Sprintf("error in AddressFromHexString:%s\n", err))
	}
}
