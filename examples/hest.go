package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/server/hest"
	"net/http"
)

func main() {
	h := hest.NewHest("hest.example")
	h.Router(func(router *gin.RouterGroup) {
		router.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ping ok")
		})
	})
	helix.Startup()
}
