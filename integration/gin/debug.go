// +build !release

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/kaz/pprotein/integration/mux"
)

func Integrate(r *gin.Engine) {
	EnableDebugHandler(r)
	EnableDebugMode(r)
}

func EnableDebugHandler(r *gin.Engine) {
	r.Any("/debug/*path", gin.WrapH(mux.NewDebugHandler()))
}

func EnableDebugMode(r *gin.Engine) {
	gin.SetMode(gin.DebugMode)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
}
