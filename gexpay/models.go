package gexpay

import (
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

/***** metadata:Processor *****/
type ProcessorType int
type ProcessorTypeArray []ProcessorType

const (
	ProcessorTypeMerchWithdraw ProcessorType = 100 //is merch withdraw
	ProcessorTypeUserWithdraw  ProcessorType = 200 //is user withdraw
	ProcessorTypeUserTransfer  ProcessorType = 300 //is user transfer
	ProcessorTypeFeeWithdraw   ProcessorType = 500 //is fee withdraw
)

//ProcessorTypeAll is the processor type
var ProcessorTypeAll = ProcessorTypeArray{ProcessorTypeMerchWithdraw, ProcessorTypeUserWithdraw, ProcessorTypeUserTransfer, ProcessorTypeFeeWithdraw}

//ProcessorTypeShow is the processor type
var ProcessorTypeShow = ProcessorTypeArray{ProcessorTypeMerchWithdraw, ProcessorTypeUserWithdraw, ProcessorTypeUserTransfer, ProcessorTypeFeeWithdraw}

type ProcessorStatus int
type ProcessorStatusArray []ProcessorStatus

const (
	ProcessorStatusPending ProcessorStatus = 100 //is pending
	ProcessorStatusDone    ProcessorStatus = 200 //is done
	ProcessorStatusFail    ProcessorStatus = 300 //is fail
)

//ProcessorStatusAll is the processor status
var ProcessorStatusAll = ProcessorStatusArray{ProcessorStatusPending, ProcessorStatusDone, ProcessorStatusFail}

//ProcessorStatusShow is the processor status
var ProcessorStatusShow = ProcessorStatusArray{ProcessorStatusPending, ProcessorStatusDone, ProcessorStatusFail}

/*
 * Processor  represents bc_processor
 * Processor Fields:tid,wallet_id,type,uuid,task_total,task_done,task_fail,user_addr,from_addr,to_addr,asset,amount,try_next,try_count,activated,balanced,approved,transferred,fee_energy_usage,fee_energy_trx,fee_net_usage,fee_net_trx,fee_used,fee_asset,fee_sys,result,synced,notified,update_time,create_time,status,
 */
type Processor struct {
	T              string          `json:"-" table:"bc_processor"`                                       /* the table name tag */
	TID            int64           `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                           /* the primary key */
	WalletID       int64           `json:"wallet_id,omitempty" valid:"wallet_id,r|i,r:0;"`               /* the processor ownere wallet id */
	Type           ProcessorType   `json:"type,omitempty" valid:"type,r|i,e:0;"`                         /* the processor type, MerchWithdraw=100:is merch withdraw, UserWithdraw=200:is user withdraw, UserTransfer=300:is user transfer, FeeWithdraw=500: is fee withdraw */
	UUID           string          `json:"uuid,omitempty" valid:"uuid,r|s,l:0;"`                         /* the processor uuid */
	TaskTotal      int64           `json:"task_total,omitempty" valid:"task_total,r|i,r:0;"`             /* the processor task total count */
	TaskDone       int64           `json:"task_done,omitempty" valid:"task_done,r|i,r:0;"`               /* the processor task done count */
	TaskFail       int64           `json:"task_fail,omitempty" valid:"task_fail,r|i,r:0;"`               /* the processor task total fail */
	UserAddr       string          `json:"user_addr,omitempty" valid:"user_addr,r|s,l:0;"`               /* the processor user address */
	FromAddr       string          `json:"from_addr,omitempty" valid:"from_addr,r|s,l:0;"`               /* the processor from address */
	ToAddr         string          `json:"to_addr,omitempty" valid:"to_addr,r|s,l:0;"`                   /* the processor target address */
	Asset          string          `json:"asset,omitempty" valid:"asset,r|s,l:0;"`                       /* the processor asset key */
	Amount         decimal.Decimal `json:"amount,omitempty" valid:"amount,r|f,r:0;"`                     /* the processor total amount */
	TryNext        xsql.Time       `json:"try_next,omitempty" valid:"try_next,r|i,r:1;"`                 /* the processor try next */
	TryCount       int             `json:"try_count,omitempty" valid:"try_count,r|i,r:0;"`               /* the processor try count */
	Activated      int             `json:"activated,omitempty" valid:"activated,r|i,r:0;"`               /* the processor activated wallet count */
	Balanced       int             `json:"balanced,omitempty" valid:"balanced,r|i,r:0;"`                 /* the processor blanced wallet count */
	Approved       int             `json:"approved,omitempty" valid:"approved,r|i,r:0;"`                 /* the processor approved wallet count */
	Transferred    int             `json:"transferred,omitempty" valid:"transferred,r|i,r:0;"`           /* the processor transferred wallet count */
	FeeEnergyUsage decimal.Decimal `json:"fee_energy_usage,omitempty" valid:"fee_energy_usage,r|f,r:0;"` /* the processor fee energy usage */
	FeeEnergyTRX   decimal.Decimal `json:"fee_energy_trx,omitempty" valid:"fee_energy_trx,r|f,r:0;"`     /* the processor fee energy trx */
	FeeNetUsage    decimal.Decimal `json:"fee_net_usage,omitempty" valid:"fee_net_usage,r|f,r:0;"`       /* the processor fee net usage */
	FeeNetTRX      decimal.Decimal `json:"fee_net_trx,omitempty" valid:"fee_net_trx,r|f,r:0;"`           /* the processor fee net trx */
	FeeUsed        decimal.Decimal `json:"fee_used,omitempty" valid:"fee_used,r|f,r:0;"`                 /* the processor fee used */
	FeeAsset       string          `json:"fee_asset,omitempty" valid:"fee_asset,r|s,l:0;"`               /* the processor fee asset */
	FeeSys         decimal.Decimal `json:"fee_sys,omitempty" valid:"fee_sys,r|f,r:0;"`                   /* the processor fee sys */
	Result         xsql.M          `json:"result,omitempty" valid:"result,r|s,l:0;"`                     /* the processor done message */
	Synced         int             `json:"synced,omitempty" valid:"synced,r|i,r:0;"`                     /* the processor synced */
	Notified       int             `json:"notified,omitempty" valid:"notified,r|i,r:0;"`                 /* the processor notified */
	UpdateTime     xsql.Time       `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"`           /* the processor last update time */
	CreateTime     xsql.Time       `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"`           /* the processor create time */
	Status         ProcessorStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`                     /* the processor status, Pending=100:is pending, Done=200:is done, Fail=300:is fail */
}

/***** metadata:Transaction *****/
type TransactionType int
type TransactionTypeArray []TransactionType

const (
	TransactionTypeNormal    TransactionType = 100 //is normal
	TransactionTypeRecharget TransactionType = 200 //is recharget
	TransactionTypeWithdraw  TransactionType = 300 //is withdraw
	TransactionTypeTransfer  TransactionType = 400 //is transfer
	TransactionTypeActivate  TransactionType = 500 //is activate
)

//TransactionTypeAll is the transaction type
var TransactionTypeAll = TransactionTypeArray{TransactionTypeNormal, TransactionTypeRecharget, TransactionTypeWithdraw, TransactionTypeTransfer, TransactionTypeActivate}

//TransactionTypeShow is the transaction type
var TransactionTypeShow = TransactionTypeArray{TransactionTypeNormal, TransactionTypeRecharget, TransactionTypeWithdraw, TransactionTypeTransfer, TransactionTypeActivate}

type TransactionSynced int
type TransactionSyncedArray []TransactionSynced

const (
	TransactionSyncedSyncing TransactionSynced = 100 //is syncing
	TransactionSyncedSynced  TransactionSynced = 200 //is synced
)

//TransactionSyncedAll is the transaction synced
var TransactionSyncedAll = TransactionSyncedArray{TransactionSyncedSyncing, TransactionSyncedSynced}

//TransactionSyncedShow is the transaction synced
var TransactionSyncedShow = TransactionSyncedArray{TransactionSyncedSyncing, TransactionSyncedSynced}

type TransactionStatus int
type TransactionStatusArray []TransactionStatus

const (
	TransactionStatusPending   TransactionStatus = 100 //is pending
	TransactionStatusConfirmed TransactionStatus = 200 //is confirmed
	TransactionStatusRevert    TransactionStatus = 300 //is revert
)

//TransactionStatusAll is the transaction status
var TransactionStatusAll = TransactionStatusArray{TransactionStatusPending, TransactionStatusConfirmed, TransactionStatusRevert}

//TransactionStatusShow is the transaction status
var TransactionStatusShow = TransactionStatusArray{TransactionStatusPending, TransactionStatusConfirmed, TransactionStatusRevert}

/*
 * Transaction  represents bc_transaction
 * Transaction Fields:tid,type,uuid,txid,contract,from_addr,to_addr,asset,amount,fee_energy_usage,fee_energy_trx,fee_net_usage,fee_net_trx,fee_used,fee_asset,result,synced,notified,update_time,create_time,status,
 */
type Transaction struct {
	T              string            `json:"-" table:"bc_transaction"`                                     /* the table name tag */
	TID            int64             `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                           /* the primary key */
	Type           TransactionType   `json:"type,omitempty" valid:"type,r|i,e:0;"`                         /* the transaction type, Normal=100:is normal, Recharget=200:is recharget, Withdraw=300:is withdraw, Transfer=400:is transfer, Activate=500:is activate */
	UUID           *string           `json:"uuid,omitempty" valid:"uuid,r|s,l:0;"`                         /* the transaction uuid */
	Txid           *string           `json:"txid,omitempty" valid:"txid,r|s,l:0;"`                         /* the transaction id */
	Contract       int               `json:"contract,omitempty" valid:"contract,r|i,r:0;"`                 /* the transaction contract */
	FromAddr       string            `json:"from_addr,omitempty" valid:"from_addr,r|s,l:0;"`               /* the transaction from address */
	ToAddr         string            `json:"to_addr,omitempty" valid:"to_addr,r|s,l:0;"`                   /* the transaction to address */
	Asset          string            `json:"asset,omitempty" valid:"asset,r|s,l:0;"`                       /* the transaction asset key */
	Amount         decimal.Decimal   `json:"amount,omitempty" valid:"amount,r|f,r:0;"`                     /* the transaction amount */
	FeeEnergyUsage decimal.Decimal   `json:"fee_energy_usage,omitempty" valid:"fee_energy_usage,r|f,r:0;"` /* the transaction fee energy usage */
	FeeEnergyTRX   decimal.Decimal   `json:"fee_energy_trx,omitempty" valid:"fee_energy_trx,r|f,r:0;"`     /* the transaction fee energy trx */
	FeeNetUsage    decimal.Decimal   `json:"fee_net_usage,omitempty" valid:"fee_net_usage,r|f,r:0;"`       /* the transaction fee net usage */
	FeeNetTRX      decimal.Decimal   `json:"fee_net_trx,omitempty" valid:"fee_net_trx,r|f,r:0;"`           /* the transaction fee net trx */
	FeeUsed        decimal.Decimal   `json:"fee_used,omitempty" valid:"fee_used,r|f,r:0;"`                 /* the transaction total fee used */
	FeeAsset       string            `json:"fee_asset,omitempty" valid:"fee_asset,r|s,l:0;"`               /* the transaction fee asset */
	Result         xsql.M            `json:"result,omitempty" valid:"result,r|s,l:0;"`                     /* the transaction result */
	Synced         TransactionSynced `json:"synced,omitempty" valid:"synced,r|i,e:0;"`                     /* the transaction synced, Syncing=100:is syncing, Synced=200:is synced */
	Notified       int               `json:"notified,omitempty" valid:"notified,r|i,r:0;"`                 /* the transaction sync status, */
	UpdateTime     xsql.Time         `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"`           /* the transaction update time */
	CreateTime     xsql.Time         `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"`           /* the transaction create time */
	Status         TransactionStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`                     /* the transaction status, Pending=100:is pending, Confirmed=200:is confirmed, Revert=300:is revert */
}
