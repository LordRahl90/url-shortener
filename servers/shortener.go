package servers

import "github.com/gin-gonic/gin"

func (s *Server) shorten(ctx *gin.Context) {

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
