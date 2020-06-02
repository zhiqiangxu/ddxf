# meta

```
data_meta = {
    "@context": {
        "xsd":"http://www.w3.org/2001/XMLSchema#",
        "题目":{
            "@id":"http://ont.io/ddxf/thesis/title",
            "@type":"xsd:string"
        },
        "作者":{
            "@id":"http://ont.io/ddxf/thesis/author",
            "@type":"xsd:string"
        },
        "组织":{
            "@id":"http://ont.io/ddxf/thesis/organization",
            "@type":"xsd:string"
        },
        "摘要":{
            "@id":"http://ont.io/ddxf/thesis/description",
            "@type":"xsd:string"
        },
        "关键词":{
            "@id":"http://ont.io/ddxf/thesis/keyword",
            "@type":"xsd:string"
        }
    },
    "题目": "PNT体系建模与可信区块链体系研究",
    "作者": ["did:ont:陈威","did:ont:孙宏伟"],
    "组织": ["did:ont:中国电子科学研究院","did:ont:中国科学院空天信息研究院"],
    "摘要": "PNT...",
    "关键词": ["体系","可信PNT","区块链"],
}

token_meta = ["访问权","引用权"]

item_meta = {
    "@context": {
        "xsd":"http://www.w3.org/2001/XMLSchema#",
        "价格":{
            "@id":"http://ont.io/ddxf/thesis/fee",
            "@type":"xsd:string"
        },
        "库存":{
            "@id":"http://ont.io/ddxf/thesis/fee",
            "@type":"xsd:string"
        },
        "题目":{
            "@id":"http://ont.io/ddxf/thesis/title",
            "@type":"xsd:string"
        },
        "作者":{
            "@id":"http://ont.io/ddxf/thesis/author",
            "@type":"xsd:string"
        },
        "组织":{
            "@id":"http://ont.io/ddxf/thesis/organization",
            "@type":"xsd:string"
        },
        "摘要":{
            "@id":"http://ont.io/ddxf/thesis/description",
            "@type":"xsd:string"
        },
        "关键词":{
            "@id":"http://ont.io/ddxf/thesis/keyword",
            "@type":"xsd:string"
        }
    },
    "题目": "PNT体系建模与可信区块链体系研究",
    "作者": ["陈威","孙宏伟"],
    "组织": ["did:ont:中国电子科学研究院","did:ont:中国科学院空天信息研究院"],
    "摘要": "PNT...",
    "关键词": ["体系","可信PNT","区块链"],
    "fee":"15 ONG",
    "stock":"1000 份",
    "token_templates": [
        {"data_meta":data_meta, "token_meta":token_meta},
        ...
    ],
}
```

