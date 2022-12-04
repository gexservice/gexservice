package gexdb

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
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

func TestLoadWithdrawReview(t *testing.T) {
	clear()
	review, err := LoadWithdrawReview(ctx)
	if err != nil || len(review) > 0 {
		t.Error(err)
		return
	}
	basedb.StoreConf(ctx, ConfigWithdrawReview, converter.JSON(xmap.M{"A": 1}))
	review, err = LoadWithdrawReview(ctx)
	if err != nil || len(review) < 1 {
		t.Error(err)
		return
	}
}
