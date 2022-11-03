package market

import (
	"context"
	"fmt"

	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

func CalcBalanceOverview(ctx context.Context, userID int64) (totalValue decimal.Decimal, areaValues []xmap.M, err error) {
	for _, area := range gexdb.BalanceAreaAll {
		value, _, _, xerr := CalcBalanceTotalValue(ctx, userID, area)
		if xerr != nil {
			err = xerr
			break
		}
		totalValue = totalValue.Add(value)
		areaValues = append(areaValues, xmap.M{"area": area, "value": value})
	}
	return
}

func CalcBalanceTotalValue(ctx context.Context, userID int64, area gexdb.BalanceArea) (totalValue decimal.Decimal, balances []*gexdb.Balance, values map[string]decimal.Decimal, err error) {
	balances, _, err = gexdb.ListUserBalance(ctx, userID, area, nil, nil)
	if err != nil {
		return
	}
	symbols := []string{}
	for _, balance := range balances {
		if balance.Asset == Quote || (balance.Free.Sign() <= 0 && balance.Locked.Sign() <= 0) {
			continue
		}
		symbols = append(symbols, fmt.Sprintf("%v.%v%v", area, balance.Area, Quote))
	}
	prices := Shared.ListLatestPrice(symbols...)
	values = map[string]decimal.Decimal{}
	for _, balance := range balances {
		if balance.Free.Sign() <= 0 && balance.Locked.Sign() <= 0 {
			values[balance.Asset] = decimal.Zero
			continue
		}
		var value decimal.Decimal
		if balance.Asset == Quote {
			value = balance.Free.Add(balance.Locked)
		} else {
			symbol := fmt.Sprintf("%v.%v%v", area, balance.Area, Quote)
			value = balance.Free.Add(balance.Locked).Mul(prices[symbol])
		}
		values[balance.Asset] = value
		totalValue = totalValue.Add(balance.Free.Add(balance.Locked))
	}
	return
}
