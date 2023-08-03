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

func NewDatabase(dsn string) *Database {
	ctx := context.Background()
	db, err := sdk.Open(ctx, dsn, environ.WithEnvironCredentials(ctx))
	if err != nil {
		panic(err)
	}

	return &Database{
		db: db,
	}
}

func (d *Database) Close() error {
	return d.db.Close(context.Background())
}
