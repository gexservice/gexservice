package baseapi

import (
	"context"
	"net/http"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/gen"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xhttp"
	"github.com/codingeasygo/web/httptest"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/base/xlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ts *httptest.Server

func init() {
	var l = zap.New(xlog.Core, zap.AddCaller())
	crud.Default.Log = func(caller int, format string, args ...interface{}) {
		l.WithOptions(zap.AddCallerSkip(caller+2)).Sugar().Infof(format, args...)
	}
	crud.Default.ErrNoRows = pgx.ErrNoRows
	crud.Default.NameConv = gen.NameConvPG
	crud.Default.ParmConv = gen.ParmConvPG
	func() {
		defer func() {
			recover()
		}()
		SrvAddr()
	}()
	xlog.AtomicLevel.Enabled(zapcore.DebugLevel)
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/base")
	if err != nil {
		panic(err)
	}
	basedb.Pool = pgx.Pool
	basedb.Pool().Exec(context.Background(), baseupgrade.DROP)

	_, err = basedb.CheckDb()
	if err != nil {
		panic(err)
	}
	initdata()
	ts = httptest.NewMuxServer()
	// ts.Mux.HandleFunc("^/pub/mlogin(\\?.*)?$", MockLoginH)
	// EnterIntentionVerifyPhoneH = NewVerifyPhone(PhoneCodeTypeVerify, "user", -1)
	Handle("", ts.Mux)
	ts.Mux.HandleNormal("^.*$", http.FileServer(http.Dir("www")))
	SrvAddr = func() string {
		return ts.URL
	}
	xhttp.EnableCookie()
}

func initdata() {

}

func clearCookie() {
	xhttp.ClearCookie()
}
