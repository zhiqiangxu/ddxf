package ddxf

// TokenType def
type TokenType byte

const (
	// ONT token
	ONT TokenType = iota
	// ONG token
	ONG
	// OEP4 token
	OEP4
	// OEP5 token
	OEP5
	// OEP8 token
	OEP8
	// OEP68 token
	OEP68
)

// Fee def
type Fee struct {
	ContractAddr string
	Type         TokenType
	Count        uint64
}

// TokenTemplate mirrors tokens
type TokenTemplate struct {
	DataID      OntID    // optional
	ServiceID   OntID    // issuer id
	TokenHashes []string // token hashes from service provider endpoint
}

// MarketPlaceItemOption for marketplace item
type MarketPlaceItemOption struct {
	ExpiredDate int64
	Stocks      uint32
	Auditor     OntID
	OJs         []OntID
}

// DTokenItem is composed from marketplace item options
type DTokenItem struct {
	Fee         Fee
	ExpiredDate int64
	Stocks      uint32
	Templates   []TokenTemplate
}

// Token def
type Token interface {
	AlphabeticalBytes() []byte
}

// CrowdSouring def
type CrowdSouring struct {
	Count       uint32
	ExpiredTime uint32
	UnitFee     Fee
}
