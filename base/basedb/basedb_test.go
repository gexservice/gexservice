package basedb

import (
	"context"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/gen"
	"github.com/codingeasygo/crud/pgx"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/base/xlog"
	"go.uber.org/zap"
)

func init() {
	func() {
		defer func() {
			recover()
		}()
		Pool()
	}()
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/base")
	if err != nil {
		panic(err)
	}
	Pool = pgx.Pool
	_, _, err = Pool().Exec(context.Background(), baseupgrade.DROP)
	if err != nil {
		panic(err)
	}
	_, err = CheckDb()
	if err != nil {
		panic(err)
	}
	var l = zap.New(xlog.Core, zap.AddCaller())
	crud.Default.Log = func(caller int, format string, args ...interface{}) {
		l.WithOptions(zap.AddCallerSkip(caller+2)).Sugar().Infof(format, args...)
	}
	crud.Default.ErrNoRows = pgx.ErrNoRows
	crud.Default.NameConv = gen.NameConvPG
	crud.Default.ParmConv = gen.ParmConvPG
}

func clear() {
	_, _, err := Pool().Exec(context.Background(), baseupgrade.CLEAR)
	if err != nil {
		panic(err)
	}
}
