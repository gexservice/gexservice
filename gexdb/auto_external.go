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
	ConfigGoldbarAddress = "goldbar_address"
	ConfigGoldbarExplain = "goldbar_explain"
	ConfigGoldbarRate    = "goldbar_rate"
	ConfigGoldbarFee     = "goldbar_fee"
	ConfigGoldbarTips    = "goldbar_tips"
	ConfigBrokerCommRate = "broker_comm_rate"
	ConfigBrokerDesc     = "broker_desc"
	ConfigTradeRule      = "trade_rule"
)

var ConfigAll = []string{ConfigWelcomeMessage, ConfigWithdrawMax, ConfigGoldbarAddress, ConfigGoldbarExplain, ConfigGoldbarRate, ConfigGoldbarFee, ConfigGoldbarTips, ConfigBrokerCommRate, ConfigTradeRule}

const (
	// BalanceAssetYWE = "YWE"
	// BalanceAssetMMK = "MMK"
	BalanceAssetGoldbar = "YWE"
)

// var (
// 	BalanceAssetAll = []string{BalanceAssetYWE, BalanceAssetMMK}
// )

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

/***** metadata:ExReturnCode *****/

const (
	CodeBalanceNotEnought  = 7100
	CodeBalanceNotFound    = 7110
	CodeOrderNotCancelable = 7200
	CodeOldPasswordInvalid = 7300
)
