package servers

import (
	"fmt"
	"net/http"
	"shortener/domains/entities"
	"shortener/requests"
	"shortener/responses"

	"github.com/gin-gonic/gin"
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
	// keep this record in memory.
	// ideally should be kept in redis to make sure every instance of this application
	// can access such data centrally.
	shortToLong[svcEnt.ShortText] = svcEnt.LongText
	res := &responses.Shortener{
		Link:  svcEnt.LongText,
		Short: fmt.Sprintf("%s/%s", s.baseURL, svcEnt.ShortText),
	}
	ctx.JSON(http.StatusCreated, res)
}

func (s *Server) visit(ctx *gin.Context) {
	short := ctx.Param("short")
	if short == "" {
		// complain here
	}
	// check for the long version in cache

	//if not found in cache, check for it in the dbase

	// if found, spin up a routine to keep it in the cache
	// redirect to the long version
}
