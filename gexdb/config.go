package gexdb

import (
	"context"

	"github.com/codingeasygo/crud/pgx"
	"github.com/gexservice/gexservice/base/basedb"
)

func LoadCoinRate(ctx context.Context) (rates map[string]float64, err error) {
	rates = map[string]float64{}
	err = basedb.LoadConf(ctx, ConfigCoinRate, &rates)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return
}
