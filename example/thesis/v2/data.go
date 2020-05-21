package v2

var data = map[string]interface{}{
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
		"hash": map[string]interface{}{
			"@type": "xsd:string",
			"@id":   "http://ont.io/ddxf#thesis-hash",
		},
	},
	"title": "thesis title for data",
	"body":  "thesis body for data",
	"hash":  "xxxxx",
}
