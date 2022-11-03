package basedb

import (
	"context"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
)

func TestConfig(t *testing.T) {
	err := StoreConf(context.Background(), "xx", "100")
	if err != nil {
		t.Error(err)
		return
	}
	var intValPtr *int
	err = LoadConf(context.Background(), "xx", &intValPtr)
	if err != nil || *intValPtr != 100 {
		t.Error(err)
		return
	}
	var intVal int
	err = LoadConf(context.Background(), "xx", &intVal)
	if err != nil || intVal != 100 {
		t.Error(err)
		return
	}
	var floatValPtr *float64
	err = LoadConf(context.Background(), "xx", &floatValPtr)
	if err != nil || *intValPtr != 100 {
		t.Error(err)
		return
	}
	var floatVal float32
	err = LoadConf(context.Background(), "xx", &floatVal)
	if err != nil || intVal != 100 {
		t.Error(err)
		return
	}
	var strValPtr *string
	err = LoadConf(context.Background(), "xx", &strValPtr)
	if err != nil || *strValPtr != "100" {
		t.Error(err)
		return
	}
	var strVal string
	err = LoadConf(context.Background(), "xx", &strVal)
	if err != nil || strVal != "100" {
		t.Error(err)
		return
	}

	//
	err = StoreConf(context.Background(), "xxx", converter.JSON(xmap.M{
		"abc": 1,
	}))
	if err != nil {
		t.Error(err)
		return
	}
	var mapVal xmap.M
	err = LoadConf(context.Background(), "xxx", &mapVal)
	if err != nil || mapVal.Int("/abc") != 1 {
		t.Error(err)
		return
	}
}

func TestUpdateConfigList(t *testing.T) {
	StoreConf(context.Background(), "a0", "1000")
	StoreConf(context.Background(), "a1", "1000")
	err := UpdateConfigList(context.Background(), xmap.M{
		"a0": "10",
		"a1": "11",
	})
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	config, err := LoadConfigList(context.Background(), "a0", "a1")
	if err != nil || config.StrDef("", "a0") != "10" || config.StrDef("", "a1") != "11" {
		t.Errorf("err:%v,%v", err, config.StrDef("", "a0") != "10")
		return
	}
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	//
	pgx.MockerSet("Pool.Begin", 1)
	err = UpdateConfigList(context.Background(), xmap.M{
		"a0": "10",
		"a2": "11",
	})
	if err == nil {
		t.Errorf("err:%v", err)
		return
	}
	pgx.MockerClear()
	//
	pgx.MockerSet("Tx.Exec", 1)
	err = UpdateConfigList(context.Background(), xmap.M{
		"a0": "10",
		"a2": "11",
	})
	if err == nil {
		t.Errorf("err:%v", err)
		return
	}
	pgx.MockerClear()
	//
	pgx.MockerSet("Pool.Query", 1)
	_, err = LoadConfigList(context.Background(), "a0", "a1")
	if err == nil {
		t.Errorf("err:%v", err)
		return
	}
	pgx.MockerClear()
	//
	pgx.MockerSet("Rows.Scan", 1)
	_, err = LoadConfigList(context.Background(), "a0", "a1")
	if err == nil {
		t.Errorf("err:%v", err)
		return
	}
	pgx.MockerClear()
}
