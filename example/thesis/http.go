package thesis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Init ...
func Init() {
	seller := NewSeller()
	buyer := NewBuyer()
	mp := NewMP()

	r := gin.Default()
	r.GET("/ddxf/mp/jsonld/types", func(c *gin.Context) {
		types := mp.JSONLDTypes()
		c.JSON(200, types)
	})
	r.GET("/ddxf/mp/jsonld/type", func(c *gin.Context) {
		var input JSONLDInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := mp.JSONLD(input)
		c.JSON(200, output)
	})
	r.GET("/ddxf/lookupByDescHash", func(c *gin.Context) {
		var input LookupByDescHashInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := seller.LookupByDescHash(input)
		c.JSON(200, output)
	})
	r.GET("/ddxf/lookupByTokenTemplate", func(c *gin.Context) {
		var input LookupByTokenTemplateInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := seller.LookupByTokenTemplate(input)
		c.JSON(200, output)
	})

	r.POST("/seller/publish", func(c *gin.Context) {
		var input PublishInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := seller.Publish(input)
		c.JSON(200, output)
	})
	r.POST("/seller/useToken", func(c *gin.Context) {
		var input UseTokenInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := seller.UseToken(input)
		c.JSON(200, output)
	})
	r.POST("/seller/useTokenByJWT", func(c *gin.Context) {
		var input UseTokenByJWTInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := seller.UseTokenByJWT(input)
		c.JSON(200, output)
	})
	r.POST("/seller/downloadWithJWT", func(c *gin.Context) {
		var input DownloadWithJWTInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := seller.DownloadWithJWT(input)
		c.JSON(200, output)
	})

	r.POST("/buyer/buydtoken", func(c *gin.Context) {
		var input BuyDtokenInput
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		output := buyer.BuyDtoken(input)
		c.JSON(200, output)
	})
	r.Run(sellerListenAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
