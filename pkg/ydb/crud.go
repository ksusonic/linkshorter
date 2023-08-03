package ydb

import (
	"context"
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const (
	selectQuery = `DECLARE $id AS String;
		SELECT link FROM links WHERE id = $id LIMIT 1;`
	insertQuery = `DECLARE $id AS String; DECLARE $link AS String;
		REPLACE INTO links (id, link) VALUES ($id, $link);`
)

func (d *Database) GetLink(ctx context.Context, id string) (string, error) {
	readTx := table.TxControl(
		table.BeginTx(
			table.WithSnapshotReadOnly(),
		),
		table.CommitTx(),
	)
	var res result.Result
	err := d.db.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			_, res, err = s.Execute(ctx, readTx, selectQuery,
				table.NewQueryParameters(
					table.ValueParam("$id", types.StringValueFromString(id)),
				),
				options.WithCollectStatsModeBasic(),
			)
			return err
		},
	)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = res.Close()
	}()
	var link string
	for res.NextResultSet(ctx) {
		for res.NextRow() {
			err = res.ScanNamed(
				named.OptionalWithDefault("link", &link),
			)
			return link, err
		}
	}
	return "", fmt.Errorf("link '%s' is not found", link)
}

func (d *Database) Insert(ctx context.Context, id, url string) error {
	writeTx := table.TxControl(
		table.BeginTx(
			table.WithSerializableReadWrite(),
		),
		table.CommitTx(),
	)
	err := d.db.Table().Do(ctx, func(ctx context.Context, s table.Session) (err error) {
		_, _, err = s.Execute(
			ctx,
			writeTx,
			insertQuery,
			table.NewQueryParameters(
				table.ValueParam("$id", types.StringValueFromString(id)),
				table.ValueParam("$link", types.StringValueFromString(url)),
			),
			options.WithCollectStatsModeBasic(),
		)
		return
	})
	return err
}
