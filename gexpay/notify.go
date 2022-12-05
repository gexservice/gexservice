package gexpay

import (
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

func NotifyTransactionH(s *web.Session) web.Result {
	var body = &struct {
		NotifyBody
		Transaction Transaction `json:"transaction"`
	}{}
	_, err := s.RecvJSON(&body)
	if err != nil {
		return s.SendJSON(xmap.M{"code": define.ArgsInvalid, "message": err.Error()})
	}
	err = body.VerifySign(AccessToken)
	if err != nil {
		xlog.Warnf("NotifyTransactionH verify sign fail by %v", converter.JSON(body))
		return s.SendJSON(xmap.M{"code": define.SignInvalid, "message": err.Error()})
	}
	xlog.Infof("NotifyTransactionH receive transaction by %v", converter.JSON(body))
	if body.Transaction.Status != TransactionStatusConfirmed || body.Transaction.Txid == nil {
		return s.SendJSON(xmap.M{"code": define.Success})
	}
	method := gexdb.WalletMethodTron
	if body.Merch.Type != WalletTypeTron {
		method = gexdb.WalletMethodEthereum
	}
	withdraw, skip, err := gexdb.ReceiveTopup(s.R.Context(), method, body.Transaction.ToAddr, *body.Transaction.Txid, body.Transaction.Asset, body.Transaction.Amount, xmap.M{"from_addr": body.Transaction.FromAddr})
	if err != nil {
		xlog.Errorf("NotifyTransactionH receive topup by err:%v,withdraw:%v", err, converter.JSON(withdraw))
	} else if !skip {
		xlog.Infof("NotifyTransactionH receive topup by err:%v,withdraw:%v", err, converter.JSON(withdraw))
	}
	return s.SendJSON(xmap.M{"code": 0})
}

func NotifyProcessorH(s *web.Session) web.Result {
	var body = &struct {
		NotifyBody
		Processor Processor `json:"processor"`
	}{}
	_, err := s.RecvJSON(&body)
	if err != nil {
		return s.SendJSON(xmap.M{"code": define.ArgsInvalid, "message": err.Error()})
	}
	err = body.VerifySign(AccessToken)
	if err != nil {
		xlog.Warnf("NotifyProcessorH verify sign fail by %v", converter.JSON(body))
		return s.SendJSON(xmap.M{"code": define.SignInvalid, "message": err.Error()})
	}
	xlog.Infof("NotifyProcessorH receive processor by %v", converter.JSON(body))
	if (body.Processor.Type != ProcessorTypeMerchWithdraw && body.Processor.Type != ProcessorTypeUserWithdraw) || body.Processor.Status == ProcessorStatusPending {
		return s.SendJSON(xmap.M{"code": define.Success})
	}
	var withdraw *gexdb.Withdraw
	if body.Processor.Status == ProcessorStatusDone {
		withdraw, err = gexdb.DoneWithdraw(s.R.Context(), body.Processor.UUID, true, body.Processor.Result.AsMap())
	} else {
		withdraw, err = gexdb.DoneWithdraw(s.R.Context(), body.Processor.UUID, false, body.Processor.Result.AsMap())
	}
	if err == nil {
		xlog.Infof("NotifyProcessorH done withdraw by err:%v,withdraw:%v", err, converter.JSON(withdraw))
	} else {
		xlog.Errorf("NotifyProcessorH done withdraw by err:%v,withdraw:%v", err, converter.JSON(withdraw))
	}
	return s.SendJSON(xmap.M{"code": define.Success})
}
