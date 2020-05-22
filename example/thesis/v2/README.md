# ddxf for thesis

4 core components: **dataid**(optional), **token**, **token item**, **mp item**

## steps

1. define **dataid** (optional), which is the subject of interest.
    1. the offchain part corresponds to the `model` as it is required by business.
    2. the onchain part includes a `descHash`(**descHash=CommonHash(model)**) of the `model`, for static data, `dataHash`(**dataHash=HashOfData(bytes)**) is also required, either in dataid or `tokenHash`(if dataid is not used).

```
var dataOffChain = map[string]interface{}{
	"title": "thesis title for data",
	"body":  "thesis body for data",
	"id":    1,
}

var dataOnChain = map[string]interface{}{
	"dataHash": "",
	"descHash": "",
}
```
    
2. define **token**, the concept of which can vary depending on whether dataid is used.
    1. offchain part is `token` and onchain part is `tokenHash`(**tokenHash=CommonHash(token)**).
    2. if dataid is not used, `token` is actually combination of both the `model` and the `acl`.
    3. if dataid is used, `token` is just the `acl`.

```
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
```

3. define **token item**(can be useful if the seller doesn't use mp), which is composed of multiple tokens and a `desc` of all tokens(**desc can be anything the seller defined**), plus other selling info（`fee`,`expireDate`,`stock`）.

```
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
```

4. define **mp item**, which is composed of multiple tokens and a `desc`(**desc can be anything the mp defined**) of all tokens, plus other selling info（`fee`,`expireDate`,`stock`）.

```
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
```