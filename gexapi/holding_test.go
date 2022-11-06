package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/define"
)

func TestHolding(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	//
	listHolding, _ := ts.Should(t, "code", define.Success, "balance/free", xmap.ShouldIsNoZero).GetMap("/usr/listHolding")
	fmt.Printf("listHolding--->%v\n", converter.JSON(listHolding))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()

	pgx.MockerSetCall("Rows.Scan", 1, 2).Should(t, "code", define.ServerError).GetMap("/usr/listHolding")
}
