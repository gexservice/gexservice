//auto gen models by autogen
package gexdb

import (
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

/***** metadata:Balance *****/
type BalanceArea int
type BalanceAreaArray []BalanceArea

const (
	BalanceAreaFunds   BalanceArea = 100 //is funds area
	BalanceAreaSpot    BalanceArea = 200 //is spot area
	BalanceAreaFutures BalanceArea = 300 //is futures area
)

//BalanceAreaAll is the balance area
var BalanceAreaAll = BalanceAreaArray{BalanceAreaFunds, BalanceAreaSpot, BalanceAreaFutures}

//BalanceAreaShow is the balance area
var BalanceAreaShow = BalanceAreaArray{BalanceAreaFunds, BalanceAreaSpot, BalanceAreaFutures}

type BalanceStatus int
type BalanceStatusArray []BalanceStatus

const (
	BalanceStatusNormal BalanceStatus = 100 //is normal
	BalanceStatusLocked BalanceStatus = 200 //is locked
)

//BalanceStatusAll is the balance status
var BalanceStatusAll = BalanceStatusArray{BalanceStatusNormal, BalanceStatusLocked}

//BalanceStatusShow is the balance status
var BalanceStatusShow = BalanceStatusArray{BalanceStatusNormal, BalanceStatusLocked}

/*
 * Balance  represents gex_balance
 * Balance Fields:tid,user_id,area,asset,free,locked,margin,update_time,create_time,status,
 */
type Balance struct {
	T          string          `json:"-" table:"gex_balance"`                              /* the table name tag */
	TID        int64           `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	UserID     int64           `json:"user_id,omitempty" valid:"user_id,r|i,r:0;"`         /* the balance user id */
	Area       BalanceArea     `json:"area,omitempty" valid:"area,r|i,e:0;"`               /* the balance area, Funds=100:is funds area, Spot=200:is spot area, Futures=300:is futures area */
	Asset      string          `json:"asset,omitempty" valid:"asset,r|s,l:0;"`             /* the balance asset key */
	Free       decimal.Decimal `json:"free,omitempty" valid:"free,r|f,r:0;"`               /* the balance free amount */
	Locked     decimal.Decimal `json:"locked,omitempty" valid:"locked,r|f,r:0;"`           /* the balance locked amount */
	Margin     decimal.Decimal `json:"margin,omitempty" valid:"margin,r|f,r:0;"`           /* the balance margin value */
	UpdateTime xsql.Time       `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the balance last update time */
	CreateTime xsql.Time       `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the balance create time */
	Status     BalanceStatus   `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the balance status, Normal=100: is normal, Locked=200: is locked */
}

/***** metadata:BalanceHistory *****/
type BalanceHistoryStatus int
type BalanceHistoryStatusArray []BalanceHistoryStatus

const (
	BalanceHistoryStatusNormal BalanceHistoryStatus = 100 //is normal status
)

//BalanceHistoryStatusAll is the balance record status
var BalanceHistoryStatusAll = BalanceHistoryStatusArray{BalanceHistoryStatusNormal}

//BalanceHistoryStatusShow is the balance record status
var BalanceHistoryStatusShow = BalanceHistoryStatusArray{BalanceHistoryStatusNormal}

/*
 * BalanceHistory  represents gex_balance_history
 * BalanceHistory Fields:tid,user_id,asset,valuation,update_time,create_time,status,
 */
type BalanceHistory struct {
	T          string               `json:"-" table:"gex_balance_history"`                      /* the table name tag */
	TID        int64                `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	UserID     int64                `json:"user_id,omitempty" valid:"user_id,r|i,r:0;"`         /* the balance user id */
	Asset      string               `json:"asset,omitempty" valid:"asset,r|s,l:0;"`             /* the balance asset key */
	Valuation  decimal.Decimal      `json:"valuation,omitempty" valid:"valuation,r|f,r:0;"`     /* the balance valuation */
	UpdateTime xsql.Time            `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the balance record update time */
	CreateTime xsql.Time            `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the balance record create time, is daily zero time */
	Status     BalanceHistoryStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the balance record status, Normal=100: is normal status */
}

/***** metadata:BalanceRecord *****/
type BalanceRecordType int
type BalanceRecordTypeArray []BalanceRecordType

const (
	BalanceRecordTypeTrade    BalanceRecordType = 100 //is trade type
	BalanceRecordTypeTradeFee BalanceRecordType = 110 //is trade fee
	BalanceRecordTypeProfit   BalanceRecordType = 200 //is close profit
	BalanceRecordTypeBlowup   BalanceRecordType = 210 //is blowup
	BalanceRecordTypeTransfer BalanceRecordType = 300 //is transfer
	BalanceRecordTypeChange   BalanceRecordType = 400 //is manual change type
)

//BalanceRecordTypeAll is the balance record type
var BalanceRecordTypeAll = BalanceRecordTypeArray{BalanceRecordTypeTrade, BalanceRecordTypeTradeFee, BalanceRecordTypeProfit, BalanceRecordTypeBlowup, BalanceRecordTypeTransfer, BalanceRecordTypeChange}

//BalanceRecordTypeShow is the balance record type
var BalanceRecordTypeShow = BalanceRecordTypeArray{BalanceRecordTypeTrade, BalanceRecordTypeTradeFee, BalanceRecordTypeProfit, BalanceRecordTypeBlowup, BalanceRecordTypeTransfer, BalanceRecordTypeChange}

type BalanceRecordStatus int
type BalanceRecordStatusArray []BalanceRecordStatus

const (
	BalanceRecordStatusNormal BalanceRecordStatus = 100 //is normal
)

//BalanceRecordStatusAll is the balance status
var BalanceRecordStatusAll = BalanceRecordStatusArray{BalanceRecordStatusNormal}

//BalanceRecordStatusShow is the balance status
var BalanceRecordStatusShow = BalanceRecordStatusArray{BalanceRecordStatusNormal}

/*
 * BalanceRecord  represents gex_balance_record
 * BalanceRecord Fields:tid,creator,balance_id,type,target,changed,update_time,create_time,status,
 */
type BalanceRecord struct {
	T          string              `json:"-" table:"gex_balance_record"`                       /* the table name tag */
	TID        int64               `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	Creator    int64               `json:"creator,omitempty" valid:"creator,r|i,r:0;"`         /* the balance creator */
	BalanceID  int64               `json:"balance_id,omitempty" valid:"balance_id,r|i,r:0;"`   /* the balance id */
	Type       BalanceRecordType   `json:"type,omitempty" valid:"type,r|i,e:0;"`               /* the balance record type, Trade=100: is trade type, TradeFee=110:is trade fee, Profit=200:is close profit, Blowup=210:is blowup, Transfer=300:is transfer, Change=400: is manual change type */
	Target     int                 `json:"target,omitempty" valid:"target,r|i,r:0;"`           /* the balance target type */
	Changed    decimal.Decimal     `json:"changed,omitempty" valid:"changed,r|f,r:0;"`         /* the balance change value */
	UpdateTime xsql.Time           `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the balance last update time */
	CreateTime xsql.Time           `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the balance create time */
	Status     BalanceRecordStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the balance status, Normal=100: is normal */
}

/***** metadata:Holding *****/
type HoldingStatus int
type HoldingStatusArray []HoldingStatus

const (
	HoldingStatusNormal HoldingStatus = 100 //is normal
	HoldingStatusLocked HoldingStatus = 200 //is locked
)

//HoldingStatusAll is the holding status
var HoldingStatusAll = HoldingStatusArray{HoldingStatusNormal, HoldingStatusLocked}

//HoldingStatusShow is the holding status
var HoldingStatusShow = HoldingStatusArray{HoldingStatusNormal, HoldingStatusLocked}

/*
 * Holding  represents gex_holding
 * Holding Fields:tid,user_id,symbol,amount,open,blowup,lever,margin_used,margin_added,update_time,create_time,status,
 */
type Holding struct {
	T           string          `json:"-" table:"gex_holding"`                                /* the table name tag */
	TID         int64           `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                   /* the primary key */
	UserID      int64           `json:"user_id,omitempty" valid:"user_id,r|i,r:0;"`           /* the holding user id */
	Symbol      string          `json:"symbol,omitempty" valid:"symbol,r|s,l:0;"`             /* the holding symbol */
	Amount      decimal.Decimal `json:"amount,omitempty" valid:"amount,r|f,r:0;"`             /* the holding amount */
	Open        decimal.Decimal `json:"open,omitempty" valid:"open,r|f,r:0;"`                 /* the holding open price */
	Blowup      decimal.Decimal `json:"blowup,omitempty" valid:"blowup,r|f,r:0;"`             /* the holding blowup price */
	Lever       int             `json:"lever,omitempty" valid:"lever,r|i,r:0;"`               /* the holding lever */
	MarginUsed  decimal.Decimal `json:"margin_used,omitempty" valid:"margin_used,r|f,r:0;"`   /* the holding margin used */
	MarginAdded decimal.Decimal `json:"margin_added,omitempty" valid:"margin_added,r|f,r:0;"` /* the holding margin added */
	UpdateTime  xsql.Time       `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"`   /* the holding last update time */
	CreateTime  xsql.Time       `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"`   /* the holding create time */
	Status      HoldingStatus   `json:"status,omitempty" valid:"status,r|i,e:0;"`             /* the holding status, Normal=100: is normal, Locked=200: is locked */
}

/***** metadata:KLine *****/

/*
 * KLine  represents gex_kline
 * KLine Fields:tid,symbol,interv,amount,count,open,close,low,high,volume,start_time,update_time,
 */
type KLine struct {
	T          string          `json:"-" table:"gex_kline"`                                /* the table name tag */
	TID        int64           `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primay key */
	Symbol     string          `json:"symbol,omitempty" valid:"symbol,r|s,l:0;"`           /* the kline symbol */
	Interv     string          `json:"interv,omitempty" valid:"interv,r|s,l:0;"`           /* the kline interval key */
	Amount     decimal.Decimal `json:"amount,omitempty" valid:"amount,r|f,r:0;"`           /* the kline amount */
	Count      int64           `json:"count,omitempty" valid:"count,r|i,r:0;"`             /* the kline count */
	Open       decimal.Decimal `json:"open,omitempty" valid:"open,r|f,r:0;"`               /* the kline open price */
	Close      decimal.Decimal `json:"close,omitempty" valid:"close,r|f,r:0;"`             /* the kline close price */
	Low        decimal.Decimal `json:"low,omitempty" valid:"low,r|f,r:0;"`                 /* the kline low price */
	High       decimal.Decimal `json:"high,omitempty" valid:"high,r|f,r:0;"`               /* the kline high price */
	Volume     decimal.Decimal `json:"volume,omitempty" valid:"volume,r|f,r:0;"`           /* the kline volume price */
	StartTime  xsql.Time       `json:"start_time,omitempty" valid:"start_time,r|i,r:1;"`   /* the kline start time */
	UpdateTime xsql.Time       `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the kline update time */
}

/***** metadata:Order *****/
type OrderType int
type OrderTypeArray []OrderType

const (
	OrderTypeTrade   OrderType = 100 //is trade type
	OrderTypeTrigger OrderType = 200 //is trigger trade order
	OrderTypeBlowup  OrderType = 300 //is blow up type
)

//OrderTypeAll is the order type
var OrderTypeAll = OrderTypeArray{OrderTypeTrade, OrderTypeTrigger, OrderTypeBlowup}

//OrderTypeShow is the order type
var OrderTypeShow = OrderTypeArray{OrderTypeTrade, OrderTypeTrigger, OrderTypeBlowup}

type OrderSide string
type OrderSideArray []OrderSide

const (
	OrderSideBuy  OrderSide = "buy"  //is buy side
	OrderSideSell OrderSide = "sell" //is sell side
)

//OrderSideAll is the order side
var OrderSideAll = OrderSideArray{OrderSideBuy, OrderSideSell}

//OrderSideShow is the order side
var OrderSideShow = OrderSideArray{OrderSideBuy, OrderSideSell}

type OrderTriggerType int
type OrderTriggerTypeArray []OrderTriggerType

const (
	OrderTriggerTypeNone       OrderTriggerType = 0   //is none type
	OrderTriggerTypeStopProfit OrderTriggerType = 100 //is stop profit type
	OrderTriggerTypeStopLoss   OrderTriggerType = 200 //is stop loss
)

//OrderTriggerTypeAll is the order trigger type
var OrderTriggerTypeAll = OrderTriggerTypeArray{OrderTriggerTypeNone, OrderTriggerTypeStopProfit, OrderTriggerTypeStopLoss}

//OrderTriggerTypeShow is the order trigger type
var OrderTriggerTypeShow = OrderTriggerTypeArray{OrderTriggerTypeNone, OrderTriggerTypeStopProfit, OrderTriggerTypeStopLoss}

type OrderStatus int
type OrderStatusArray []OrderStatus

const (
	OrderStatusWaiting      OrderStatus = 100 //
	OrderStatusPending      OrderStatus = 200 //is pending
	OrderStatusPartialled   OrderStatus = 300 //is partialled
	OrderStatusDone         OrderStatus = 400 //is done
	OrderStatusPartCanceled OrderStatus = 410 //is partialled canceled
	OrderStatusCanceled     OrderStatus = 420 //is canceled
)

//OrderStatusAll is the order status
var OrderStatusAll = OrderStatusArray{OrderStatusWaiting, OrderStatusPending, OrderStatusPartialled, OrderStatusDone, OrderStatusPartCanceled, OrderStatusCanceled}

//OrderStatusShow is the order status
var OrderStatusShow = OrderStatusArray{OrderStatusWaiting, OrderStatusPending, OrderStatusPartialled, OrderStatusDone, OrderStatusPartCanceled, OrderStatusCanceled}

//OrderOrderbyAll is crud filter
const OrderOrderbyAll = "update_time,create_time"

/*
 * Order  represents gex_order
 * Order Fields:tid,order_id,type,user_id,creator,symbol,side,quantity,filled,price,trigger_type,trigger_price,avg_price,total_price,holding,profit,owned,unhedged,in_balance,in_filled,out_balance,out_filled,fee_balance,fee_filled,fee_rate,transaction,fee_settled_status,fee_settled_next,update_time,create_time,status,
 */
type Order struct {
	T                string           `json:"-" table:"gex_order"`                                              /* the table name tag */
	TID              int64            `json:"tid,omitempty" valid:"tid,o|i,r:0;"`                               /* the primary key */
	OrderID          string           `json:"order_id,omitempty" valid:"order_id,r|s,l:0;"`                     /* the order string id */
	Type             OrderType        `json:"type,omitempty" valid:"type,r|i,e:0;"`                             /* the order type, Trade=100: is trade type, Trigger=200: is trigger trade order, Blowup=300: is blow up type */
	UserID           int64            `json:"user_id,omitempty" valid:"user_id,r|i,r:0;"`                       /* the order user id */
	Creator          int64            `json:"creator,omitempty" valid:"creator,r|i,r:0;"`                       /* the order creator user id */
	Symbol           string           `json:"symbol,omitempty" valid:"symbol,r|s,l:0;"`                         /* the order symbol */
	Side             OrderSide        `json:"side,omitempty" valid:"side,r|s,e:0;"`                             /* the order side, Buy=buy: is buy side, Sell=sell: is sell side */
	Quantity         decimal.Decimal  `json:"quantity,omitempty" valid:"quantity,o|f,r:0;"`                     /* the order expected quantity */
	Filled           decimal.Decimal  `json:"filled,omitempty" valid:"filled,r|f,r:0;"`                         /* the order filled quantity */
	Price            decimal.Decimal  `json:"price,omitempty" valid:"price,o|f,r:0;"`                           /* the order expected price */
	TriggerType      OrderTriggerType `json:"trigger_type,omitempty" valid:"trigger_type,o|i,e:0;"`             /* the order trigger type, None=0:is none type, StopProfit=100: is stop profit type, StopLoss=200: is stop loss */
	TriggerPrice     decimal.Decimal  `json:"trigger_price,omitempty" valid:"trigger_price,o|f,r:0;"`           /* the order trigger price */
	AvgPrice         decimal.Decimal  `json:"avg_price,omitempty" valid:"avg_price,r|f,r:0;"`                   /* the order filled avg price */
	TotalPrice       decimal.Decimal  `json:"total_price,omitempty" valid:"total_price,o|f,r:0;"`               /* the order filled total price */
	Holding          decimal.Decimal  `json:"holding,omitempty" valid:"holding,r|f,r:0;"`                       /* the order holding */
	Profit           decimal.Decimal  `json:"profit,omitempty" valid:"profit,r|f,r:0;"`                         /* the order profit */
	Owned            decimal.Decimal  `json:"owned,omitempty" valid:"owned,r|f,r:0;"`                           /* the order owned count */
	Unhedged         decimal.Decimal  `json:"unhedged,omitempty" valid:"unhedged,r|f,r:0;"`                     /* the order owned is unbalanced */
	InBalance        string           `json:"in_balance,omitempty" valid:"in_balance,r|s,l:0;"`                 /* the in balance asset key */
	InFilled         decimal.Decimal  `json:"in_filled,omitempty" valid:"in_filled,r|f,r:0;"`                   /* the in balance filled amount */
	OutBalance       string           `json:"out_balance,omitempty" valid:"out_balance,r|s,l:0;"`               /* the out balance asset key */
	OutFilled        decimal.Decimal  `json:"out_filled,omitempty" valid:"out_filled,r|f,r:0;"`                 /* the out balance filled amount */
	FeeBalance       string           `json:"fee_balance,omitempty" valid:"fee_balance,r|s,l:0;"`               /* the fee balance asset key */
	FeeFilled        decimal.Decimal  `json:"fee_filled,omitempty" valid:"fee_filled,r|f,r:0;"`                 /* the fee amount */
	FeeRate          decimal.Decimal  `json:"fee_rate,omitempty" valid:"fee_rate,r|f,r:0;"`                     /* the order fee rate */
	Transaction      OrderTransaction `json:"transaction,omitempty" valid:"transaction,r|s,l:0;"`               /* the order transaction info */
	FeeSettledStatus int              `json:"fee_settled_status,omitempty" valid:"fee_settled_status,r|i,r:0;"` /* the order transaction detail */
	FeeSettledNext   xsql.Time        `json:"fee_settled_next,omitempty" valid:"fee_settled_next,r|i,r:1;"`     /* the fee settled time */
	UpdateTime       xsql.Time        `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"`               /* the order update time */
	CreateTime       xsql.Time        `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"`               /* the order create time */
	Status           OrderStatus      `json:"status,omitempty" valid:"status,o|i,e:0;"`                         /* the order status, Waiting=100, Pending=200:is pending, Partialled=300:is partialled, Done=400:is done, PartCanceled=410: is partialled canceled, Canceled=420: is canceled */
}

/***** metadata:OrderComm *****/
type OrderCommType int
type OrderCommTypeArray []OrderCommType

const (
	OrderCommTypeNormal OrderCommType = 100 //is normal type
)

//OrderCommTypeAll is the comm type
var OrderCommTypeAll = OrderCommTypeArray{OrderCommTypeNormal}

//OrderCommTypeShow is the comm type
var OrderCommTypeShow = OrderCommTypeArray{OrderCommTypeNormal}

type OrderCommStatus int
type OrderCommStatusArray []OrderCommStatus

const (
	OrderCommStatusNormal OrderCommStatus = 100 //is normal
)

//OrderCommStatusAll is the comm status
var OrderCommStatusAll = OrderCommStatusArray{OrderCommStatusNormal}

//OrderCommStatusShow is the comm status
var OrderCommStatusShow = OrderCommStatusArray{OrderCommStatusNormal}

/*
 * OrderComm  represents gex_order_comm
 * OrderComm Fields:tid,order_id,user_id,type,in_balance,in_fee,update_time,create_time,status,
 */
type OrderComm struct {
	T          string          `json:"-" table:"gex_order_comm"`                           /* the table name tag */
	TID        int64           `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	OrderID    int64           `json:"order_id,omitempty" valid:"order_id,r|i,r:0;"`       /* the order id */
	UserID     int64           `json:"user_id,omitempty" valid:"user_id,r|i,r:0;"`         /* the user id */
	Type       OrderCommType   `json:"type,omitempty" valid:"type,r|i,e:0;"`               /* the comm type, Normal=100:is normal type */
	InBalance  string          `json:"in_balance,omitempty" valid:"in_balance,r|s,l:0;"`   /* the in balance asset key */
	InFee      decimal.Decimal `json:"in_fee,omitempty" valid:"in_fee,r|f,r:0;"`           /* the in balance fee */
	UpdateTime xsql.Time       `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /*  */
	CreateTime xsql.Time       `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the comm create time */
	Status     OrderCommStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the comm status, Normal=100:is normal */
}

/***** metadata:User *****/
type UserType int
type UserTypeArray []UserType

const (
	UserTypeAdmin  UserType = 10  //is admin user
	UserTypeNormal UserType = 100 //is normal user
)

//UserTypeAll is the user type
var UserTypeAll = UserTypeArray{UserTypeAdmin, UserTypeNormal}

//UserTypeShow is the user type
var UserTypeShow = UserTypeArray{UserTypeAdmin, UserTypeNormal}

type UserRole int
type UserRoleArray []UserRole

const (
	UserRoleNormal UserRole = 100 //is normal
	UserRoleStaff  UserRole = 200 //is staff
	UserRoleMaker  UserRole = 300 //is maker
)

//UserRoleAll is ther user role
var UserRoleAll = UserRoleArray{UserRoleNormal, UserRoleStaff, UserRoleMaker}

//UserRoleShow is ther user role
var UserRoleShow = UserRoleArray{UserRoleNormal, UserRoleStaff, UserRoleMaker}

type UserStatus int
type UserStatusArray []UserStatus

const (
	UserStatusNormal  UserStatus = 100 //is normal
	UserStatusLocked  UserStatus = 200 //is locked
	UserStatusRemoved UserStatus = -1  //is deleted
)

//UserStatusAll is the user status
var UserStatusAll = UserStatusArray{UserStatusNormal, UserStatusLocked, UserStatusRemoved}

//UserStatusShow is the user status
var UserStatusShow = UserStatusArray{UserStatusNormal, UserStatusLocked}

//UserOrderbyAll is crud filter
const UserOrderbyAll = "account,phone,update_time,create_time"

/*
 * User  represents gex_user
 * User Fields:tid,type,role,name,account,phone,password,trade_pass,image,fee,external,favorites,config,update_time,create_time,status,
 */
type User struct {
	T          string        `json:"-" table:"gex_user"`                                 /* the table name tag */
	TID        int64         `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	Type       UserType      `json:"type,omitempty" valid:"type,r|i,e:0;"`               /* the user type,Admin=10:is admin user, Normal=100:is normal user */
	Role       UserRole      `json:"role,omitempty" valid:"role,o|i,e:0;"`               /* ther user role, Normal=100:is normal, Staff=200:is staff, Maker=300:is maker */
	Name       *string       `json:"name,omitempty" valid:"name,o|s,l:0;"`               /* the user name */
	Account    *string       `json:"account,omitempty" valid:"account,o|s,l:0;"`         /* the user account to login */
	Phone      *string       `json:"phone,omitempty" valid:"phone,o|s,p:^\\d{11}$;"`     /* the user phone number to login */
	Password   *string       `json:"password,omitempty" valid:"password,o|s,l:0;"`       /* the user password to login */
	TradePass  *string       `json:"trade_pass,omitempty" valid:"trade_pass,o|s,l:0;"`   /* the user trade password */
	Image      *string       `json:"image,omitempty" valid:"image,o|s,l:0;"`             /* the user image */
	Fee        xsql.M        `json:"fee,omitempty" valid:"fee,r|s,l:0;"`                 /* the user fee */
	External   xsql.M        `json:"external,omitempty" valid:"external,o|s,l:0;"`       /* the user external info */
	Favorites  UserFavorites `json:"favorites,omitempty" valid:"favorites,r|s,l:0;"`     /* the user favorites */
	Config     xsql.M        `json:"config,omitempty" valid:"config,r|s,l:0;"`           /* the user config */
	UpdateTime xsql.Time     `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the last updat time */
	CreateTime xsql.Time     `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the craete time */
	Status     UserStatus    `json:"status,omitempty" valid:"status,o|i,e:0;"`           /* the user status, Normal=100:is normal, Locked=200:is locked, Removed=-1:is deleted */
}

/***** metadata:Withdraw *****/
type WithdrawType int
type WithdrawTypeArray []WithdrawType

const (
	WithdrawTypeWithdraw WithdrawType = 100 //is withdraw type
	WithdrawTypeTopup    WithdrawType = 200 //is topup type
	WithdrawTypeGoldbar  WithdrawType = 300 //is goldbar bar
)

//WithdrawTypeAll is the withdraw order type
var WithdrawTypeAll = WithdrawTypeArray{WithdrawTypeWithdraw, WithdrawTypeTopup, WithdrawTypeGoldbar}

//WithdrawTypeShow is the withdraw order type
var WithdrawTypeShow = WithdrawTypeArray{WithdrawTypeWithdraw, WithdrawTypeTopup, WithdrawTypeGoldbar}

type WithdrawStatus int
type WithdrawStatusArray []WithdrawStatus

const (
	WithdrawStatusPending   WithdrawStatus = 100 //is pending
	WithdrawStatusConfirmed WithdrawStatus = 200 //is confirmed
	WithdrawStatusDone      WithdrawStatus = 300 //is done
	WithdrawStatusCanceled  WithdrawStatus = 320 //is canceled
)

//WithdrawStatusAll is the withdraw order status
var WithdrawStatusAll = WithdrawStatusArray{WithdrawStatusPending, WithdrawStatusConfirmed, WithdrawStatusDone, WithdrawStatusCanceled}

//WithdrawStatusShow is the withdraw order status
var WithdrawStatusShow = WithdrawStatusArray{WithdrawStatusPending, WithdrawStatusConfirmed, WithdrawStatusDone, WithdrawStatusCanceled}

/*
 * Withdraw  represents gex_withdraw
 * Withdraw Fields:tid,order_id,type,user_id,creator,asset,quantity,transaction,update_time,create_time,status,
 */
type Withdraw struct {
	T           string          `json:"-" table:"gex_withdraw"`                             /* the table name tag */
	TID         int64           `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	OrderID     string          `json:"order_id,omitempty" valid:"order_id,r|s,l:0;"`       /* the withdraw order string id */
	Type        WithdrawType    `json:"type,omitempty" valid:"type,r|i,e:0;"`               /* the withdraw order type, Withdraw=100: is withdraw type, Topup=200: is topup type, Goldbar=300: is goldbar bar */
	UserID      int64           `json:"user_id,omitempty" valid:"user_id,r|i,r:0;"`         /* the withdraw order user id */
	Creator     int64           `json:"creator,omitempty" valid:"creator,r|i,r:0;"`         /* the withdraw order creator user id */
	Asset       string          `json:"asset,omitempty" valid:"asset,r|s,l:0;"`             /* the withdraw asset */
	Quantity    decimal.Decimal `json:"quantity,omitempty" valid:"quantity,r|f,r:0;"`       /* the withdraw order quantity */
	Transaction xsql.M          `json:"transaction,omitempty" valid:"transaction,r|s,l:0;"` /* the withdraw order transaction info */
	UpdateTime  xsql.Time       `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the withdraw order update time */
	CreateTime  xsql.Time       `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the withdraw order create time */
	Status      WithdrawStatus  `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the withdraw order status, Pending=100:is pending, Confirmed=200:is confirmed, Done=300:is done, Canceled=320: is canceled */
}
