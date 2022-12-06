package market

import (
	"context"
	"fmt"

	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
)

func CalcBalanceOverview(ctx context.Context, userID int64) (totalValue decimal.Decimal, areaValues []xmap.M, err error) {
	for _, area := range gexdb.BalanceAreaAll {
		value, winned, _, _, xerr := CalcBalanceTotalValue(ctx, userID, area)
		if xerr != nil {
			err = xerr
			break
		}
		totalValue = totalValue.Add(value)
		areaValues = append(areaValues, xmap.M{"area": area, "value": value, "today_winned": winned})
	}
	return
}

func CalcBalanceTotalValue(ctx context.Context, userID int64, area gexdb.BalanceArea) (totalValue decimal.Decimal, todayWinned decimal.Decimal, balances []*gexdb.Balance, values map[string]decimal.Decimal, err error) {
	balances, _, err = gexdb.ListUserBalance(ctx, userID, area, nil, nil)
	if err != nil {
		return
	}
	symbols := []string{}
	for _, balance := range balances {
		if balance.Asset == matcher.Quote || (balance.Free.Sign() <= 0 && balance.Locked.Sign() <= 0) {
			continue
		}
		symbols = append(symbols, fmt.Sprintf("%v.%v%v", area, balance.Area, matcher.Quote))
	}
	prices := ListLatestPrice(symbols...)
	values = map[string]decimal.Decimal{}
	for _, balance := range balances {
		if balance.Free.Sign() <= 0 && balance.Locked.Sign() <= 0 {
			values[balance.Asset] = decimal.Zero
			continue
		}
		var value decimal.Decimal
		if balance.Asset == matcher.Quote {
			value = balance.Free.Add(balance.Locked)
		} else {
			symbol := fmt.Sprintf("%v%v%v", area.Prefix(), balance.Asset, matcher.Quote)
			value = balance.Free.Add(balance.Locked).Mul(prices[symbol])
			kline := LoadKLine(symbol, "1day")
			if kline != nil && kline.Open.IsPositive() {
				winned := balance.Free.Add(balance.Locked).Mul(kline.Close.Sub(kline.Open))
				todayWinned = todayWinned.Add(winned)
			}
		}
		values[balance.Asset] = value
		totalValue = totalValue.Add(balance.Free.Add(balance.Locked))
	}
	return
}

func CalcHoldingUnprofit(ctx context.Context, holdings ...*gexdb.Holding) (unprofits map[string]decimal.Decimal, tickers map[string]*gexdb.Ticker) {
	symbols := []string{}
	for _, holding := range holdings {
		symbols = append(symbols, holding.Symbol)
	}
	tickers = ListTicker(symbols...)
	unprofits = map[string]decimal.Decimal{}
	totalUnprofit := decimal.Zero
	for _, holding := range holdings {
		ticker := tickers[holding.Symbol]
		if ticker != nil && ticker.Bid != nil && holding.Amount.IsPositive() {
			unprofit := ticker.Bid[0].Sub(holding.Open).Mul(holding.Amount)
			totalUnprofit = totalUnprofit.Add(unprofit)
			unprofits[holding.Symbol] = unprofit
		} else if ticker != nil && ticker.Ask != nil && holding.Amount.IsNegative() {
			unprofit := ticker.Ask[0].Sub(holding.Open).Mul(holding.Amount)
			totalUnprofit = totalUnprofit.Add(unprofit)
			unprofits[holding.Symbol] = unprofit
		}
	}
	unprofits["total"] = totalUnprofit
	return
}

func ListHoldingUnprofit(ctx context.Context, userIDs ...int64) (unprofits map[int64]map[string]decimal.Decimal, err error) {
	holdingAll, _, err := gexdb.ListHoldingByUser(ctx, userIDs, nil)
	if err != nil {
		return
	}
	unprofits = map[int64]map[string]decimal.Decimal{}
	for userID, holdings := range holdingAll {
		unprofits[userID], _ = CalcHoldingUnprofit(ctx, holdings...)
	}
	return
}
