package router

import (
	"github.com/barysh-vn/shortener/internal/app"
	"github.com/barysh-vn/shortener/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	linkHandler := handler.LinkHandler{
		LinkService:   app.GetLinkService(),
		RandomService: app.GetRandomService(),
	}

	r.GET("/:id", linkHandler.HandleGet)
	r.POST("/", linkHandler.HandlePost)

	return r
}
