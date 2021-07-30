package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine, prefix string) *gin.Engine {
	if r == nil {
		r = gin.Default()
	}

	r.Use(Cors())

	// ///////////////////////////////////////////////////////////
	r.POST("/", wrap(api.Handle))
	// ///////////////////////////////////////////////////////////

	return r
}
