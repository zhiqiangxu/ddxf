package v2

var mpItemOffChain = map[string]interface{}{
	"title":      "thesis title for mp",
	"body":       "thesis body for mp",
	"mpspecific": "mp specific",
	"fee":        map[string]interface{}{},
	"expireDate": 1590064629,
	"stock":      100,
}

var mpItemOnChain = map[string]interface{}{
	"tokens": []interface{}{
		map[string]interface{}{
			"dataids":   "",
			"tokenHash": "",
			"dataHash":  "",
		},
	},
	"descHash":   "",
	"fee":        map[string]interface{}{},
	"expireDate": 1590064629,
	"stock":      100,
}
