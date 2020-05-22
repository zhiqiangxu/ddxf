package v2

var tokenOnChain = map[string]interface{}{
	"dataids":   "", //optional
	"dataHash":  "", //required if dataids empty
	"tokenHash": "",
}

var tokenOffChainWODataID = map[string]interface{}{
	// without dataids
	"dataHash": "xxx",
	"token": map[string]interface{}{
		"title": "xxx",
		"body":  "xxx",
		"id":    1,
		"r":     true,
		"w":     true,
	},
}

var tokenOffChainWithDataID = map[string]interface{}{
	// with dataids
	"dataids": "xxx;yyy",
	"token": map[string]interface{}{
		"r": true,
		"w": true,
	},
}
