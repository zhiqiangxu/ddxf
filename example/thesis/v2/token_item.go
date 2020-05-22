package v2

var tokenItemOnChain = map[string]interface{}{
	"tokens": []interface{}{
		map[string]interface{}{
			"dataids":   "",
			"tokenHash": "",
			"dataHash":  "",
		},
	},
	"descHash":   "", // can be empty if only one token
	"fee":        map[string]interface{}{},
	"expireDate": 1590064629,
	"stock":      100,
}

var tokenItemOffChain = map[string]interface{}{
	"tokens": []interface{}{
		map[string]interface{}{
			"dataids":   "",
			"tokenHash": "",
			"dataHash":  "",
		},
	},
	"desc": map[string]interface{}{
		"anything": "value",
	},
	"fee":        map[string]interface{}{},
	"expireDate": 1590064629,
	"stock":      100,
}
