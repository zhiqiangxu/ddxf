package thesis

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ddxf/contract"
	"github.com/zhiqiangxu/util/claim"
)

// Seller ...
type Seller struct {
	descHashMap      map[string][]Thesis
	tokenTemplateMap map[contract.TokenTemplate]Thesis
	jwt              *claim.JWT
}

// NewSeller ...
func NewSeller() *Seller {
	jwt, _ := claim.NewJWT(time.Hour, []byte("secret"))
	return &Seller{descHashMap: make(map[string][]Thesis), tokenTemplateMap: make(map[contract.TokenTemplate]Thesis), jwt: jwt}
}

// Thesis ...
type Thesis struct {
	Title string
	Body  string
}

type (
	// PublishInput ...
	PublishInput struct {
		Theses []Thesis
		MP     contract.Marketplace
	}
	// PublishOutput ...
	PublishOutput struct {
		ResourceID     string
		TokenTemplates contract.TokenTemplates
	}
)

const (
	sellerListenAddr = "0.0.0.0:8081"
	sellerID         = ddxf.OntID("seller xxx")
)

// Publish ...
func (s *Seller) Publish(input PublishInput) (output PublishOutput) {

	resourceID := uuid.New().String()

	descHash, _ := ddxf.HashObject(input.Theses)
	ddo := contract.ResourceDDO{
		Manager:      sellerID,
		ResourceType: contract.RTStaticFile,
		Endpoint:     "http://" + sellerListenAddr,
		DescHash:     descHash,
		MP:           input.MP,
	}
	templates := make(contract.TokenTemplates)
	for _, thesis := range input.Theses {
		tokenHash, _ := ddxf.HashObject(thesis)
		templates[contract.TokenTemplate{TokenHash: tokenHash}] = struct{}{}
	}
	item := contract.DTokenItem{
		Fee: ddxf.Fee{
			ContractAddr: "xxx",
			Type:         ddxf.ONT,
			Count:        1,
		},
		ExpiredDate: time.Now().Add(time.Hour).Unix(),
		Stocks:      100,
		Templates:   templates,
	}

	DDXF().DTokenSellerPublish(resourceID, ddo, item)

	s.descHashMap[descHash] = input.Theses
	for _, thesis := range input.Theses {
		tokenHash, _ := ddxf.HashObject(thesis)
		s.tokenTemplateMap[contract.TokenTemplate{TokenHash: tokenHash}] = thesis
	}

	output.ResourceID = resourceID
	output.TokenTemplates = templates

	return

}

type (
	// UseTokenInput ...
	UseTokenInput struct {
		ResourceID    string
		Buyer         ddxf.OntID
		TokenTemplate contract.TokenTemplate
	}
	// UseTokenOutput ...
	UseTokenOutput struct {
		Thesis Thesis
	}
)

// UseToken ...
func (s *Seller) UseToken(input UseTokenInput) (output UseTokenOutput) {

	DDXF().UseToken(input.ResourceID, input.Buyer, input.TokenTemplate, 1)

	thesis := s.tokenTemplateMap[input.TokenTemplate]

	output.Thesis = thesis

	return
}

type (
	// UseTokenByJWTInput ...
	UseTokenByJWTInput struct {
		ResourceID    string
		Buyer         ddxf.OntID
		TokenTemplate contract.TokenTemplate
	}
	// UseTokenByJWTOutput ...
	UseTokenByJWTOutput struct {
		JWTToken string
	}
)

// UseTokenByJWT ...
func (s *Seller) UseTokenByJWT(input UseTokenByJWTInput) (output UseTokenByJWTOutput) {

	DDXF().UseToken(input.ResourceID, input.Buyer, input.TokenTemplate, 1)

	claim := map[string]interface{}{
		"TokenTemplate": input.TokenTemplate,
	}
	jwtToken, _ := s.jwt.Sign(claim)

	output.JWTToken = jwtToken

	return
}

type (
	// DownloadWithJWTInput ...
	DownloadWithJWTInput struct {
		JWTToken string
	}
	// DownloadWithJWTOutput ...
	DownloadWithJWTOutput struct {
		Thesis Thesis
	}
)

// DownloadWithJWT ...
func (s *Seller) DownloadWithJWT(input DownloadWithJWTInput) (output DownloadWithJWTOutput) {
	ok, claim := s.jwt.Verify(input.JWTToken)
	if !ok {
		return
	}

	tt := claim["TokenTemplate"]
	ttBytes, _ := json.Marshal(tt)
	var tokenTemplate contract.TokenTemplate
	err := json.Unmarshal(ttBytes, &tokenTemplate)
	if err != nil {
		return
	}

	thesis := s.tokenTemplateMap[tokenTemplate]

	output.Thesis = thesis

	return
}

type (
	// LookupByDescHashInput ...
	LookupByDescHashInput struct {
		DescHash string
	}
	// LookupByDescHashOutput ...
	LookupByDescHashOutput struct {
		Theses []Thesis
	}
)

// LookupByDescHash ...
func (s *Seller) LookupByDescHash(input LookupByDescHashInput) (output LookupByDescHashOutput) {
	output.Theses = s.descHashMap[input.DescHash]
	return
}

type (
	// LookupByTokenTemplateInput ...
	LookupByTokenTemplateInput struct {
		TokenTemplate contract.TokenTemplate
	}
	// LookupByTokenTemplateOutput ...
	LookupByTokenTemplateOutput struct {
		Thesis Thesis
	}
)

// LookupByTokenTemplate ...
func (s *Seller) LookupByTokenTemplate(input LookupByTokenTemplateInput) (output LookupByTokenTemplateOutput) {
	output.Thesis = s.tokenTemplateMap[input.TokenTemplate]
	return
}
