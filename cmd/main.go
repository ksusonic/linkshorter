package main

import (
	"context"
	"log"

	"github.com/ksusonic/linkshorter/internal/config"
	"github.com/ksusonic/linkshorter/internal/controller"
	"github.com/ksusonic/linkshorter/internal/server"
	"github.com/ksusonic/linkshorter/pkg/ydb"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()

	db := ydb.NewDatabase(ctx, cfg.DatabaseDsn, &log.Logger{})
	defer func() { _ = db.Close(ctx) }()

	srv := server.NewServer()

	{
		api := srv.Group("/api")
		urlController := controller.NewUrlController(db)

		api.POST("/shorten", urlController.Shorten)
		api.GET("/redirect/:id", urlController.Redirect)
	}

	if err := srv.Run(cfg.Address); err != nil {
		return
	}
}
