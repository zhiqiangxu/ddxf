package thesis

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ddxf/contract"
	"github.com/zhiqiangxu/util/claim"
)

type descHashAndResourceID struct {
	DescHash   string
	ResourceID string
}

type tokenTemplateAndResourceID struct {
	TokenTemplate contract.TokenTemplate
	ResourceID    string
}

// Seller ...
type Seller struct {
	descHashMap      map[descHashAndResourceID]string
	tokenTemplateMap map[tokenTemplateAndResourceID]Thesis
	jwt              *claim.JWT
}

// NewSeller ...
func NewSeller() *Seller {
	jwt, _ := claim.NewJWT(time.Hour, []byte("secret"))
	return &Seller{descHashMap: make(map[descHashAndResourceID]string), tokenTemplateMap: make(map[tokenTemplateAndResourceID]Thesis), jwt: jwt}
}

// Thesis ...
type Thesis struct {
	ID    int
	Title string
	Body  string
}

// DataHash for this thesis
// TODO impl
func (t *Thesis) DataHash() string {
	h := [4]byte{}
	return string(h[:])
}

// ToMPThesis ...
func (t *Thesis) ToMPThesis() MPThesis {
	return MPThesis{}
}

type (
	// PublishInput ...
	PublishInput struct {
		Theses      []Thesis
		Fee         ddxf.Fee
		ExpiredDate int64
		Stocks      uint32
		Desc        string
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

	descHash := ddxf.Sha256Bytes([]byte(input.Desc))
	ddo := contract.ResourceDDO{
		Manager:      sellerID,
		ResourceType: contract.RTStaticFile,
		Endpoint:     "http://" + sellerListenAddr,
		DescHash:     descHash,
	}
	templates := make(contract.TokenTemplates)
	for _, thesis := range input.Theses {
		tokenHash, _ := ddxf.HashObject(thesis)
		templates[contract.TokenTemplate{TokenHash: tokenHash + thesis.DataHash()}] = struct{}{}
	}
	item := contract.DTokenItem{
		Fee:         input.Fee,
		ExpiredDate: input.ExpiredDate,
		Stocks:      input.Stocks,
		Templates:   templates,
	}

	DDXF().DTokenSellerPublish(resourceID, ddo, item)

	s.descHashMap[descHashAndResourceID{DescHash: descHash, ResourceID: resourceID}] = input.Desc
	for _, thesis := range input.Theses {
		tokenHash, _ := ddxf.HashObject(thesis)
		s.tokenTemplateMap[tokenTemplateAndResourceID{
			TokenTemplate: contract.TokenTemplate{TokenHash: tokenHash + thesis.DataHash()},
			ResourceID:    resourceID}] = thesis
	}

	output.ResourceID = resourceID
	output.TokenTemplates = templates

	return

}

type (
	// MPPublishInput ...
	MPPublishInput struct {
		Theses      []Thesis
		Fee         ddxf.Fee
		ExpiredDate int64
		Stocks      uint32
		MPDesc      string
		MP          *MP
	}
	// MPPublishOutput ...
	MPPublishOutput struct {
		ResourceID     string
		TokenTemplates contract.TokenTemplates
	}
)

// MPPublish ...
func (s *Seller) MPPublish(input MPPublishInput) (output MPPublishOutput) {

	mpTheses := []MPThesis{}
	for _, thesis := range input.Theses {
		mpTheses = append(mpTheses, thesis.ToMPThesis())
	}

	mpoutput := input.MP.PublishMP(
		PublishMPInput{
			Fee:         input.Fee,
			ExpiredDate: input.ExpiredDate,
			Stocks:      input.Stocks,
			MPDesc:      input.MPDesc,
			MPTheses:    mpTheses})
	if !mpoutput.OK {
		return
	}

	resourceID := uuid.New().String()

	descHash := ddxf.Sha256Bytes([]byte(input.MPDesc))

	ddo := contract.ResourceDDO{
		Manager:      sellerID,
		ResourceType: contract.RTStaticFile,
		Endpoint:     "http://" + sellerListenAddr,
		DescHash:     descHash,
		MP:           input.MP.MPC,
	}

	templates := make(contract.TokenTemplates)
	for _, thesis := range input.Theses {
		tokenHash, _ := ddxf.HashObject(thesis)
		templates[contract.TokenTemplate{TokenHash: tokenHash + thesis.DataHash()}] = struct{}{}
	}
	item := contract.DTokenItem{
		Fee:         input.Fee,
		ExpiredDate: input.ExpiredDate,
		Stocks:      input.Stocks,
		Templates:   templates,
	}

	DDXF().DTokenSellerPublish(resourceID, ddo, item)

	setoutput := input.MP.SetResourceID(SetResourceIDInput{ItemID: mpoutput.ItemID, ResourceID: resourceID})
	if !setoutput.OK {
		// TODO handle error
		return
	}

	s.descHashMap[descHashAndResourceID{DescHash: descHash, ResourceID: resourceID}] = input.MPDesc
	for _, thesis := range input.Theses {
		tokenHash, _ := ddxf.HashObject(thesis)
		s.tokenTemplateMap[tokenTemplateAndResourceID{
			TokenTemplate: contract.TokenTemplate{TokenHash: tokenHash + thesis.DataHash()},
			ResourceID:    resourceID}] = thesis
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

	thesis, ok := s.tokenTemplateMap[tokenTemplateAndResourceID{
		TokenTemplate: input.TokenTemplate,
		ResourceID:    input.ResourceID,
	}]
	if !ok {
		return
	}

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
		"ResourceID":    input.ResourceID,
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

	thesis, ok := s.tokenTemplateMap[tokenTemplateAndResourceID{
		TokenTemplate: tokenTemplate,
		ResourceID:    claim["ResourceID"].(string),
	}]
	if !ok {
		return
	}

	output.Thesis = thesis

	return
}

type (
	// LookupByDescHashInput ...
	LookupByDescHashInput struct {
		ResourceID string
		DescHash   string
	}
	// LookupByDescHashOutput ...
	LookupByDescHashOutput struct {
		Desc string
	}
)

// LookupByDescHash ...
func (s *Seller) LookupByDescHash(input LookupByDescHashInput) (output LookupByDescHashOutput) {
	desc, ok := s.descHashMap[descHashAndResourceID{DescHash: input.DescHash, ResourceID: input.ResourceID}]
	if !ok {
		return
	}

	output.Desc = desc
	return
}

type (
	// LookupByTokenTemplateInput ...
	LookupByTokenTemplateInput struct {
		ResourceID    string
		TokenTemplate contract.TokenTemplate
	}
	// LookupByTokenTemplateOutput ...
	LookupByTokenTemplateOutput struct {
		Thesis Thesis
	}
)

// LookupByTokenTemplate ...
func (s *Seller) LookupByTokenTemplate(input LookupByTokenTemplateInput) (output LookupByTokenTemplateOutput) {
	thesis, ok := s.tokenTemplateMap[tokenTemplateAndResourceID{
		TokenTemplate: input.TokenTemplate,
		ResourceID:    input.ResourceID,
	}]
	if !ok {
		return
	}

	output.Thesis = thesis
	return
}

// JSONLDTypes ...
func (s *Seller) JSONLDTypes() []string {
	return []string{"thesis"}
}

// JSONLD ...
func (s *Seller) JSONLD(input JSONLDInput) (output JSONLDOutput) {
	if input.Type != "thesis" {
		return
	}

	output.Context = map[string]interface{}{
		"Desc":  "string",
		"Token": "Thesis",
	}
	return
}
