package gexdb

import (
	"context"
	"encoding/json"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/basedb"
)

func LoadCoinRate(ctx context.Context) (rates []xmap.M, err error) {
	var data string
	err = basedb.LoadConf(ctx, ConfigCoinRate, &data)
	if err == nil {
		err = json.Unmarshal([]byte(data), &rates)
	}
	if err == pgx.ErrNoRows {
		err = nil
	}
	return
}

func LoadWithdrawReview(ctx context.Context) (review xmap.M, err error) {
	review, err = LoadWithdrawReviewCall(Pool(), ctx)
	return
}

func LoadWithdrawReviewCall(caller crud.Queryer, ctx context.Context) (review xmap.M, err error) {
	review = xmap.M{}
	err = basedb.LoadConfCall(caller, ctx, ConfigWithdrawReview, &review)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return
}
