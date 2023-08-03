package ydb

import (
	"context"
	"log"

	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	sdk "github.com/ydb-platform/ydb-go-sdk/v3"
)

type Database struct {
	db     *sdk.Driver
	logger *log.Logger
}

func NewDatabase(ctx context.Context, dsn string, logger *log.Logger) *Database {
	db, err := sdk.Open(ctx, dsn, environ.WithEnvironCredentials(ctx))
	if err != nil {
		panic(err)
	}

	return &Database{
		db:     db,
		logger: logger,
	}
}

func (d *Database) Close(ctx context.Context) error {
	return d.db.Close(ctx)
}
