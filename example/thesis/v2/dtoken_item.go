package v2

var dtokenItem = map[string]interface{}{
	"@context": map[string]interface{}{
		"xsd": "http://www.w3.org/2001/XMLSchema#",
		"title": map[string]interface{}{
			"@type": "xsd:string",
			"@id":   "http://ont.io/ddxf#thesis-title",
		},
		"body": map[string]interface{}{
			"@type": "xsd:string",
			"@id":   "http://ont.io/ddxf#thesis-body",
		},
		"fee": map[string]interface{}{
			"@id": "http://ont.io/ddxf#thesis-fee",
		},
		"expireDate": map[string]interface{}{
			"@type": "xsd:int",
			"@id":   "http://ont.io/ddxf#thesis-expireDate",
		},
		"stock": map[string]interface{}{
			"@type": "xsd:int",
			"@id":   "http://ont.io/ddxf#thesis-stock",
		},
	},
	"title":      "thesis title for data",
	"body":       "thesis body for data",
	"fee":        map[string]interface{}{},
	"expireDate": 1590064629,
	"stock":      100,
}
