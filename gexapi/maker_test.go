package gexapi

import (
	"fmt"
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/maker"
	"github.com/shopspring/decimal"
)

func TestMaker(t *testing.T) {
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	config := &maker.Config{}
	config.Symbol = "spot.YWEUSDT"
	config.Delay = 20
	config.UserID = userabc0.TID
	config.Open = decimal.NewFromFloat(1000)
	config.Close.Min = decimal.NewFromFloat(-0.01)
	config.Close.Max = decimal.NewFromFloat(0.01)
	config.Vib.Min = decimal.NewFromFloat(-0.03)
	config.Vib.Max = decimal.NewFromFloat(0.03)
	config.Vib.Count = 5
	config.Ticker = decimal.NewFromFloat(0.0001)
	config.Interval = time.Hour.Milliseconds()
	config.Depth.QtyMax = decimal.NewFromFloat(3)
	config.Depth.StepMax = 5
	config.Depth.DiffMax = decimal.NewFromFloat(2)
	config.Depth.DiffMin = decimal.NewFromFloat(0.02)
	config.Depth.Max = 15
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap("xx", "/admin/updateSymbolMaker")
	ts.Should(t, "code", define.Success).PostJSONMap(config, "/admin/updateSymbolMaker")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/admin/startSymbolMaker?symbol=%v", "")
	ts.Should(t, "code", define.Success).GetMap("/admin/startSymbolMaker?symbol=%v", config.Symbol)
	ts.Should(t, "code", define.ServerError).GetMap("/admin/startSymbolMaker?symbol=%v", config.Symbol)
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/admin/loadSymbolMaker?symbol=%v", "")
	loadSymbolMaker, _ := ts.Should(t, "code", define.Success).GetMap("/admin/loadSymbolMaker?symbol=%v", config.Symbol)
	fmt.Printf("loadSymbolMaker--->%v\n", converter.JSON(loadSymbolMaker))
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/admin/stopSymbolMaker?symbol=%v", "")
	ts.Should(t, "code", define.Success).GetMap("/admin/stopSymbolMaker?symbol=%v", config.Symbol)
	ts.Should(t, "code", define.ServerError).GetMap("/admin/stopSymbolMaker?symbol=%v", config.Symbol)

	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetRangeCall("Rows.Scan", 2, 7).Should(t, "code", define.ServerError).GetMap("/admin/loadSymbolMaker?symbol=%v", config.Symbol)
	pgx.MockerSetCall("Pool.Exec", 1).Should(t, "code", define.ServerError).PostJSONMap(config, "/admin/updateSymbolMaker")
}
