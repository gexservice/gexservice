package gexdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/util/xsql"
)

func TestMessageTemplate(t *testing.T) {
	config := xprop.NewConfig()
	config.LoadPropString(`
[message.withdraw.done.title]
en=Withdraw Success
_=提现成功

[message.withdraw.done.content]
en=you withdraw ${_amount}${_asset} success on ${_time}
_=您于${_time}提现${_amount}${_asset}成功

[message.withdraw.fail.title]
en=Withdraw Fail
_=提现失败

[message.withdraw.fail.content]
en=you withdraw ${_amount}${_asset} Fail on ${_time}, ${_message}
_=您于${_time}提现${_amount}${_asset}失败

[message.topup.title]
en=Topup Success
_=充值成功

[message.topup.content]
en=you topup ${_amount}${_asset} success on ${_time}
_=您于${_time}成功充值${_amount}${_asset}

[message.blowup.title]
en=Blowup Warning
_=爆仓提醒

[message.blowup.content]
en=you position ${_amount} ${_symbol} is blowup on ${_time}, open price is ${_openPrice}, mark price is ${_markPrice}
_=您的仓位${_amount} ${_symbol}于${_time}爆仓，开仓价格为${_openPrice}，标记价格为${_markPrice}
	`)
	template := ReadMessageTemplateByConfig(config)

	values := ParseMessageTemplate(template, xmap.M{}, MessageKeyTopup+".title")
	fmt.Printf("-->%v\n", converter.JSON(values))
}

func TestMessage(t *testing.T) {
	err := AddMultiMessage(ctx, &Message{
		Type:   MessageTypeUser,
		Title:  xsql.M{"test": "abc"},
		Status: MessageStatusNormal,
	})
	if err != nil {
		t.Error(err)
		return
	}

	searcher := &MessageUnifySearcher{}
	searcher.Where.Type = MessageTypeAll
	searcher.Where.Key = "test"
	searcher.Where.Status = MessageStatusAll
	err = searcher.Apply(context.Background())
	if err != nil || len(searcher.Query.Messages) < 1 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}

	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()

	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = AddMultiMessage(ctx, &Message{
			Type:   MessageTypeUser,
			Title:  xsql.M{"test": "abc"},
			Status: MessageStatusNormal,
		})
		return
	})
}
