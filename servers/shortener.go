package servers

import (
	"fmt"
	"net/http"
	"shortener/domains/entities"
	"shortener/requests"
	"shortener/responses"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	shortToLong = make(map[string]string)
)

func (s *Server) shorten(ctx *gin.Context) {
	var req *requests.Shortener
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	svcEnt := &entities.Shortener{
		LongText: req.Link,
	}
	if err := s.shortSvc.Create(ctx.Request.Context(), svcEnt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// shortToLong[svcEnt.ShortText] = svcEnt.LongText
	if err := s.cacheService.Save(ctx.Request.Context(), svcEnt.ShortText, svcEnt.LongText); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	res := &responses.Shortener{
		Link:  svcEnt.LongText,
		Short: fmt.Sprintf("%s/%s", s.baseURL, svcEnt.ShortText),
	}
	ctx.JSON(http.StatusCreated, res)
}

func (s *Server) visit(ctx *gin.Context) {
	short := ctx.Param("short")
	if short == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   fmt.Errorf("invalid short link"),
		})
		return
	}
	// check for the long version in cache
	rec, ok := shortToLong[short]
	if ok {
		ctx.Redirect(http.StatusMovedPermanently, rec)
		return
	}

	//if not found in cache, check for it in the dbase
	res, err := s.shortSvc.FindByShort(ctx.Request.Context(), short)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// if found, spin up a routine to keep it in the cache
	// shortToLong[short] = res.LongText
	if err := s.cacheService.Save(ctx.Request.Context(), short, res.LongText); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	// redirect to the long version
	ctx.Redirect(http.StatusMovedPermanently, res.LongText)
}
