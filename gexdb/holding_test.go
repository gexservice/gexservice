package gexdb

import (
	"testing"

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
}
