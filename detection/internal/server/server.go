package server

import (
	"github.com/gin-gonic/gin"
	"waf-detection/internal/server/controllers"
)

func New(
	controller *controllers.FingerprintController,
) *gin.Engine {
	router := newRouter()

	router.GET("/check/:ip", controller.CheckIP)

	return router
}

func newRouter() *gin.Engine {
	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	return r
}
