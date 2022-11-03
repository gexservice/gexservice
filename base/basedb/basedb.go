package basedb

import (
	"context"
	"math/rand"
	"strings"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xtime"
	"github.com/gexservice/gexservice/base/baseupgrade"
)

func init() {
	rand.Seed(xtime.Now())
}

var SYS = "_sys"

//Pool will return database connection pool
var Pool = func() *pgx.PgQueryer {
	panic("db is not initial")
}

//CheckDb will check database if is initial
func CheckDb() (created bool, err error) {
	_, _, err = Pool().Exec(context.Background(), `select key from `+SYS+`_config limit 1`)
	if err != nil {
		_, _, err = Pool().Exec(context.Background(), strings.ReplaceAll(baseupgrade.LATEST, "_sys_", SYS+"_"))
		created = true
	}
	return
}

type BaseTableName string

func (b BaseTableName) GetTableName(args ...interface{}) string {
	if len(b) > 0 && len(args) < 2 {
		return SYS + "_" + string(b)
	}
	return SYS + "_" + args[1].(string)
}
