package routes

import (
	"laba_service/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type LabaRouter interface {
	Mount()
}

type labaRouterImpl struct {
	v       *gin.RouterGroup
	handler handlers.LabaHandler
}

func NewProductRouter(v *gin.RouterGroup, handler handlers.LabaHandler) LabaRouter {
	return &labaRouterImpl{v: v, handler: handler}
}

func (r *labaRouterImpl) Mount() {
	r.v.Use(cors.Default())
	r.v.GET("", r.handler.GetAll)
	r.v.GET("/export", r.handler.ExportExcel)
	r.v.POST("/import", r.handler.ImportExcel)
}
