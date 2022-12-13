package gexdb

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestHolding(t *testing.T) {
	symbol := "spot.YWEUSDT"
	user := testAddUser("TestHolding")
	added, err := TouchHolding(ctx, []string{symbol}, user.TID)
	if err != nil || added != 1 {
		t.Error(err)
		return
	}
	holding, err := FindHoldlingBySymbol(ctx, user.TID, symbol)
	if err != nil || holding.Amount.Sign() != 0 {
		t.Error(err)
		return
	}
	holding.Amount = decimal.NewFromFloat(1)
	holding.Blowup = decimal.NewFromFloat(100)
	holding.MarginAdded = decimal.NewFromFloat(1)
	err = holding.UpdateFilter(Pool, ctx, "")
	if err != nil {
		t.Error(err)
		return
	}
	holdings, err := ListHoldingForBlowupOverCall(Pool(), ctx, symbol, decimal.Zero, decimal.NewFromFloat(100), true)
	if err != nil || len(holdings) < 1 {
		t.Error(err)
		return
	}
	holdings, err = ListHoldingForBlowupFreeCall(Pool(), ctx, symbol, decimal.Zero, decimal.NewFromFloat(110), true)
	if err != nil || len(holdings) < 1 {
		t.Error(err)
		return
	}
	holdings, symbols, err := ListUserHolding(ctx, user.TID, []string{symbol})
	if err != nil || len(holdings) < 1 || len(symbols) < 1 {
		t.Error(err)
		return
	}

	searcher := HoldingUnifySearcher{}
	searcher.Where.Symbol = []string{holding.Symbol}
	searcher.Where.Key = "Test"
	err = searcher.Apply(ctx)
	if err != nil || len(searcher.Query.Holdings) < 1 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}

	holdingAll, symbolAll, err := ListHoldingByUser(ctx, []int64{user.TID}, []string{symbol})
	if err != nil || len(holdingAll) < 1 || len(symbolAll) < 1 {
		t.Error(err)
		return
	}

	symbols, err = ListHoldingSymbol(ctx)
	if err != nil || len(symbols) < 1 {
		t.Error(err)
		return
	}

	holdings, _, err = CountHolding(ctx, 1, time.Time{}, time.Now())
	if err != nil || len(holdings) < 1 {
		t.Error(err)
		return
	}

	holdings, _, err = CountHolding(ctx, -1, time.Time{}, time.Now())
	if err != nil || len(holdings) > 0 {
		t.Error(err)
		return
	}
}
