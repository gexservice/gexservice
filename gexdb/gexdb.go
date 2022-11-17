package gexdb

import (
	"context"
	"math/rand"
	"os"
	"reflect"
	"strings"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/gen"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xhash"
	"github.com/codingeasygo/util/xtime"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexupgrade"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

func init() {
	rand.Seed(xtime.Now())
	var l = zap.New(xlog.Core, zap.AddCaller())
	crud.Default.Verbose = crud.Default.Verbose || os.Getenv("GEX_DEBUG") == "1"
	crud.Default.Log = func(caller int, format string, args ...interface{}) {
		l.WithOptions(zap.AddCallerSkip(caller+2)).Sugar().Infof(format, args...)
	}
	crud.Default.ErrNoRows = pgx.ErrNoRows
	crud.Default.NameConv = func(on, name string, field reflect.StructField) string {
		if on == "query" && (strings.HasSuffix(field.Type.String(), "OrderTransaction") || strings.HasSuffix(field.Type.String(), "UserFavorites")) {
			return name + "::text"
		}
		return gen.NameConvPG(on, name, field)
	}
	crud.Default.ParmConv = gen.ParmConvPG
}

//Pool will return database connection pool
var Pool = func() *pgx.PgQueryer {
	panic("db is not initial")
}

//Redis will return redis connection
var Redis = func() redis.Conn {
	panic("redis is not initial")
}

//CheckDb will check database if is initial
func CheckDb(ctx context.Context) (created bool, err error) {
	_, _, err = Pool().Exec(ctx, `select tid from gex_user limit 1`)
	if err != nil {
		xlog.Infof("start generate database...")
		_, _, err = Pool().Exec(ctx, gexupgrade.LATEST)
		created = true
	}
	if err == nil {
		_, _, err = Pool().Exec(ctx, gexupgrade.CHECK)
	}
	return
}

func EncryptionUserPassword(pwd string) string {
	return xhash.SHA1([]byte(pwd))
}
