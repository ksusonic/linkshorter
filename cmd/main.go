package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ksusonic/linkshorter/frontend"
	"github.com/ksusonic/linkshorter/internal/config"
	"github.com/ksusonic/linkshorter/internal/controller"
	"github.com/ksusonic/linkshorter/internal/server"
	"github.com/ksusonic/linkshorter/pkg/ydb"
)

func main() {
	cfg := config.NewConfig()

	db := ydb.NewDatabase(cfg.DatabaseDsn)
	defer func() { _ = db.Close() }()

	router := server.NewServer()

	{
		router.GET("/ping", func(c *gin.Context) {
			c.Status(200)
		})

		// backend
		urlController := controller.NewUrlController(db)

		router.POST("/shorten", urlController.Shorten)
		router.GET("/:id", urlController.Redirect)

		// frontend
		router.LoadHTMLFiles(frontend.Templates...)
		router.GET("/", frontend.Index)
	}

	// entrypoint
	runtime(cfg.Address, router)
}
