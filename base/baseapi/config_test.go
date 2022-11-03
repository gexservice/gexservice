package baseapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/codingeasygo/web/httptest"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
)

func TestConfig(t *testing.T) {
	_ = basedb.StoreConf(context.Background(), "a0", "v0")
	_ = basedb.StoreConf(context.Background(), "a1", "v1")
	_ = basedb.StoreConf(context.Background(), "a2", "v100")
	basedb.ConfigAll = []string{"a0", "a1", "a2"}
	updateSysConfig, err := ts.PostJSONMap(
		xmap.M{
			"a0": "vv0",
			"a1": "vv1",
		}, "/usr/updateSysConfig")
	if err != nil || updateSysConfig.IntDef(-1, "code") != 0 {
		t.Errorf("err:%v,updateSysConfig:%v", err, updateSysConfig)
		return
	}
	fmt.Printf("updateSysConfig--->%v\n", converter.JSON(updateSysConfig))
	loadSysConfig, err := ts.GetMap("/usr/loadSysConfig")
	if err != nil || loadSysConfig.IntDef(-1, "code") != 0 || loadSysConfig.Str("/config/a0") != "vv0" {
		t.Errorf("err:%v,loadSysConfig:%v", err, loadSysConfig)
		return
	}
	fmt.Printf("loadSysConfig--->%v\n", converter.JSON(loadSysConfig))

	testConfig := httptest.NewMuxServer()
	testConfig.Mux.Handle("/testConfig", ConfigLoader{"a0", "a1"})
	tesConfig, err := testConfig.GetMap("/testConfig")
	if err != nil || tesConfig.IntDef(-1, "code") != 0 || tesConfig.Str("/config/a0") != "vv0" {
		t.Errorf("err:%v,tesConfig:%v", err, tesConfig)
		return
	}

	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	updateSysConfig, err = ts.PostJSONMap("xx", "/usr/updateSysConfig")
	if err != nil || updateSysConfig.IntDef(-1, "code") != define.ArgsInvalid {
		t.Errorf("err:%v,updateSysConfig:%v", err, updateSysConfig)
		return
	}
	EditSysConfigAccess = func(s *web.Session) bool { return false }
	updateSysConfig, err = ts.PostJSONMap(xmap.New(), "/usr/updateSysConfig")
	if err != nil || updateSysConfig.IntDef(-1, "code") != define.NotAccess {
		t.Errorf("err:%v,updateSysConfig:%v", err, updateSysConfig)
		return
	}
	loadSysConfig, err = ts.GetMap("/usr/loadSysConfig")
	if err != nil || loadSysConfig.IntDef(-1, "code") != define.NotAccess {
		t.Errorf("err:%v,loadSysConfig:%v", err, loadSysConfig)
		return
	}
	pgx.MockerClear()
	EditSysConfigAccess = func(s *web.Session) bool { return true }
	//
	pgx.MockerSet("Pool.Begin", 1)
	updateSysConfig, err = ts.PostJSONMap(xmap.New(), "/usr/updateSysConfig")
	if err != nil || updateSysConfig.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,updateSysConfig:%v", err, updateSysConfig)
		return
	}
	pgx.MockerClear()
	//
	pgx.MockerSet("Pool.Query", 1)
	loadSysConfig, err = ts.GetMap("/usr/loadSysConfig")
	if err != nil || loadSysConfig.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,loadSysConfig:%v", err, loadSysConfig)
		return
	}
	pgx.MockerClear()
	pgx.MockerSet("Pool.Query", 1)
	tesConfig, err = testConfig.GetMap("/testConfig")
	if err != nil || tesConfig.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,tesConfig:%v", err, tesConfig)
		return
	}
	pgx.MockerClear()
}
