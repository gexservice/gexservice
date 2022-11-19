package gexdb

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/util/converter"
	"github.com/gexservice/gexservice/base/basedb"
)

func TestLoadCoinRate(t *testing.T) {
	clear()
	_, err := LoadCoinRate(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	err = basedb.StoreConf(ctx, ConfigCoinRate, converter.JSON([]map[string]float64{
		{"xx": 10},
	}))
	if err != nil {
		t.Error(err)
		return
	}
	rates, err := LoadCoinRate(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("-->%v\n", converter.JSON(rates))
}
