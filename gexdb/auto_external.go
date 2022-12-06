package gexdb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

const (
	ConfigWelcomeMessage = "welcome_message"
	ConfigWithdrawMax    = "withdraw_max"
	ConfigWithdrawReview = "withdraw_review"
	ConfigGoldbarAddress = "goldbar_address"
	ConfigGoldbarExplain = "goldbar_explain"
	ConfigGoldbarRate    = "goldbar_rate"
	ConfigGoldbarFee     = "goldbar_fee"
	ConfigGoldbarTips    = "goldbar_tips"
	ConfigTradeRule      = "trade_rule"
	ConfigCoinRate       = "coin_rate"
)

var ConfigAll = []string{ConfigWelcomeMessage, ConfigWithdrawMax, ConfigGoldbarAddress, ConfigGoldbarExplain, ConfigGoldbarRate, ConfigGoldbarFee, ConfigGoldbarTips, ConfigTradeRule, ConfigCoinRate}

const (
	// BalanceAssetYWE = "YWE"
	// BalanceAssetMMK = "MMK"
	BalanceAssetGoldbar = "YWE"
)

// var (
// 	BalanceAssetAll = []string{BalanceAssetYWE, BalanceAssetMMK}
// )

func (b BalanceArea) Prefix() string {
	switch b {
	case BalanceAreaSpot:
		return "spot."
	case BalanceAreaFutures:
		return "futures."
	case BalanceAreaFunds:
		return "funds."
	default:
		return fmt.Sprintf("%v.", int(b))
	}
}

func (b BalanceArea) String() string {
	switch b {
	case BalanceAreaSpot:
		return "Spot"
	case BalanceAreaFutures:
		return "Futures"
	case BalanceAreaFunds:
		return "Funds"
	default:
		return fmt.Sprintf("%v", int(b))
	}
}

type ErrBalanceNotEnought string
type ErrBalanceNotFound string

func (e ErrBalanceNotEnought) Error() string {
	return string(e)
}

func (e ErrBalanceNotFound) Error() string {
	return string(e)
}

func IsErrBalanceNotEnought(err error) bool {
	_, ok := err.(ErrBalanceNotEnought)
	return ok
}

func IsErrBalanceNotFound(err error) bool {
	_, ok := err.(ErrBalanceNotFound)
	return ok
}

type UserFavorites struct {
	Symbols []string `json:"symbols,omitempty"`
}

//Scan is sql.Sanner
func (u *UserFavorites) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), u)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (u *UserFavorites) Value() (driver.Value, error) {
	if u == nil {
		return "{}", nil
	}
	bys, err := json.Marshal(u)
	return string(bys), err
}

func (u *UserFavorites) TopSymbol(symbol string) {
	symbols := []string{symbol}
	for _, s := range u.Symbols {
		if s == symbol {
			continue
		}
		symbols = append(symbols, s)
	}
	u.Symbols = symbols
}

func (u *UserFavorites) SwitchSymbol(a, b string) {
	ia, ib := -1, -1
	for i, s := range u.Symbols {
		if s == a {
			ia = i
		} else if s == b {
			ib = i
		}
	}
	if ia > -1 && ib > -1 {
		u.Symbols[ia], u.Symbols[ib] = u.Symbols[ib], u.Symbols[ia]
	}
}

type OrderTransactionItem struct {
	OrderID    string          `json:"order_id,omitempty"`
	Filled     decimal.Decimal `json:"filled,omitempty"`
	Price      decimal.Decimal `json:"price,omitempty"`
	TotalPrice decimal.Decimal `json:"total_price,omitempty"`
	FeeBalance string          `json:"fee_balance,omitempty"`
	FeeFilled  decimal.Decimal `json:"fee_filled,omitempty"`
	CreateTime xsql.Time       `json:"create_time"`
}

type OrderTransactionGoldbar struct {
	Code    string `json:"code,omitempty"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
}

type OrderTransaction struct {
	Trans   []*OrderTransactionItem  `json:"trans,omitempty"`
	Goldbar *OrderTransactionGoldbar `json:"goldbar,omitempty"`
}

//Scan is sql.Sanner
func (o *OrderTransaction) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), o)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (o *OrderTransaction) Value() (driver.Value, error) {
	if o == nil {
		return "{}", nil
	}
	bys, err := json.Marshal(o)
	return string(bys), err
}

func (h *Holding) CalcMargin(precision int32) (margin decimal.Decimal) {
	return h.Amount.Abs().Mul(h.Open).Div(decimal.NewFromInt(int64(h.Lever))).Round(precision)
}

func (h *Holding) CalcBlowup(precision int32, max decimal.Decimal) (blowup decimal.Decimal) {
	if h.Amount.Sign() != 0 {
		blowup = h.MarginUsed.Mul(max).Add(h.MarginAdded).DivRound(h.Amount, precision).Mul(decimal.NewFromInt(-1)).Add(h.Open)
	}
	return
}

func (h *Holding) Copy() (holding *Holding) {
	holding = &Holding{
		TID:         h.TID,
		UserID:      h.UserID,
		Symbol:      h.Symbol,
		Amount:      h.Amount,
		Open:        h.Open,
		Blowup:      h.Blowup,
		Lever:       h.Lever,
		MarginUsed:  h.MarginUsed,
		MarginAdded: h.MarginAdded,
		UpdateTime:  h.UpdateTime,
		CreateTime:  h.CreateTime,
		Status:      h.Status,
	}
	return
}

func (o *Order) Info() string {
	return fmt.Sprintf(
		"tid:%v,order_id:%v,type:%v,side:%v,qty:%v,filled:%v,price:%v,total_price:%v,holding:%v,fee:%v%v,status:%v",
		o.TID, o.OrderID, o.Type, o.Side, o.Quantity, o.Filled, o.Price, o.TotalPrice, o.Holding, o.FeeFilled, o.FeeBalance, o.Status,
	)
}

type Ticker struct {
	Symbol string            `json:"symbol"`
	Ask    []decimal.Decimal `json:"ask"`
	Bid    []decimal.Decimal `json:"bid"`
}

/**
 * @apiDefine BalanceRecordItemObject
 * @apiSuccess (BalanceRecordItem) {Int64} BalanceRecordItem.tid the primary key
 * @apiSuccess (BalanceRecordItem) {String} BalanceRecord.asset the balance asset
 * @apiSuccess (BalanceRecordItem) {BalanceRecordType} BalanceRecordItem.type the balance record type, all suported is <a href="#metadata-BalanceRecord">BalanceRecordTypeAll</a>
 * @apiSuccess (BalanceRecordItem) {Decimal} BalanceRecordItem.changed the balance change value
 * @apiSuccess (BalanceRecordItem) {Time} BalanceRecordItem.update_time the balance last update time
 */

type BalanceRecordItem struct {
	TID        int64             `json:"tid"`
	Asset      string            `json:"asset"`
	Type       BalanceRecordType `json:"type,omitempty"`
	Target     int               `json:"target,omitempty"`
	Changed    decimal.Decimal   `json:"changed,omitempty"`
	UpdateTime xsql.Time         `json:"update_time,omitempty"`
}

/***** metadata:ExReturnCode *****/

const (
	CodeBalanceNotEnought  = 7100
	CodeBalanceNotFound    = 7110
	CodeOrderNotCancelable = 7200
	CodeOrderPending       = 7210
	CodeOldPasswordInvalid = 7300
)
