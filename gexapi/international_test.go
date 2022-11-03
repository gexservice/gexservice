package gexapi

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/codingeasygo/web"
	"github.com/codingeasygo/web/httptest"
	"golang.org/x/net/proxy"
)

func TestInternationalPrice(t *testing.T) {
	url1, url2 := InternationalURL1, InternationalURL2
	dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:1105", nil, nil)
	InternationalClient = xhttp.NewClient(&http.Client{Transport: &http.Transport{
		Dial: dialer.Dial,
	}})
	err := ProcRefreshInternationalPrice()
	if err != pgx.ErrNoRows {
		t.Error(err)
		return
	}
	index, err := ts.GetMap("/pub/index")
	if err != nil || index.Int64("code") != 0 {
		t.Errorf("err:%v,index:%v", err, index)
		return
	}
	fmt.Printf("index--->%v\n", converter.JSON(index))

	InternationalClient = xhttp.Shared
	//
	//test error
	ts := httptest.NewMuxServer()
	ts.Mux.HandleFunc("/", func(s *web.Session) web.Result {
		return s.SendPlainText(`
			x
			hq_str_hf_XAU=
		`)
	})
	InternationalURL1 = ts.URL
	InternationalURL2 = url2
	err = ProcRefreshInternationalPrice()
	if err != pgx.ErrNoRows {
		t.Error(err)
		return
	}
	InternationalURL1 = url1
	InternationalURL2 = ts.URL
	err = ProcRefreshInternationalPrice()
	if err != pgx.ErrNoRows {
		t.Error(err)
		return
	}

	InternationalURL1 = "http://127.0.0.1:6"
	InternationalURL2 = url2
	err = ProcRefreshInternationalPrice()
	if err != pgx.ErrNoRows {
		t.Error(err)
		return
	}
	InternationalURL1 = url1
	InternationalURL2 = "http://127.0.0.1:6"
	err = ProcRefreshInternationalPrice()
	if err != pgx.ErrNoRows {
		t.Error(err)
		return
	}
}
