package main

import (
	"github.com/ksusonic/linkshorter/internal/config"
	"github.com/ksusonic/linkshorter/internal/controller"
	"github.com/ksusonic/linkshorter/internal/server"
	"github.com/ksusonic/linkshorter/pkg/ydb"
)

func main() {
	cfg := config.NewConfig()

	db := ydb.NewDatabase(cfg.DatabaseDsn)
	defer func() { _ = db.Close() }()

	srv := server.NewServer()

	{
		api := srv.Group("/api")
		urlController := controller.NewUrlController(db)

		api.POST("/shorten", urlController.Shorten)
		api.GET("/redirect/:id", urlController.Redirect)
	}

	// entrypoint
	runtime(cfg.Address, srv)
}
