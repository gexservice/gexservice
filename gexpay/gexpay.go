package gexpay

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/codingeasygo/util/xhash"
	"github.com/codingeasygo/util/xtime"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/gexdb"
)

const (
	WalletTypeTron     int = 100 //is tron type
	WalletTypeEthereum int = 200 //is ethereum type
)

var AccessToken = "abc"
var AgentAddr = "http://pay.loc:7234"
var MerchType = map[gexdb.WalletMethod]int{
	gexdb.WalletMethodTron:     WalletTypeTron,
	gexdb.WalletMethodEthereum: WalletTypeEthereum,
}
var MerchAddr = map[gexdb.WalletMethod]string{
	gexdb.WalletMethodTron:     "TSM7En6HYpyaBY566ZKVyUo5jpXe1FPSYw",
	gexdb.WalletMethodEthereum: "0x22059c11769bd828f8470582394ee0b53090087f",
}

//Handle will register all handler
func Handle(pre string, mux *web.SessionMux) {
	mux.HandleFunc("^"+pre+"/notify/processor(\\?.*)?$", NotifyProcessorH)
	mux.HandleFunc("^"+pre+"/notify/transaction(\\?.*)?$", NotifyTransactionH)
}

func RequestSign(method gexdb.WalletMethod, key string) (timestamp int64, sign string) {
	timestamp = xtime.Now()
	body := NotifyBody{
		Timestamp: timestamp,
		Key:       key,
	}
	body.Merch.Address = MerchAddr[method]
	body.Merch.Type = MerchType[method]
	sign = body.CalcSign(AccessToken)
	return
}

type NotifyBody struct {
	Merch struct {
		Address string `json:"address"`
		Type    int    `json:"type"`
	} `json:"merch"`
	Timestamp int64  `json:"timestamp"`
	Key       string `json:"key"`
	Sign      string `json:"sign"`
}

func (m *NotifyBody) CalcSign(token string) (sign string) {
	args := url.Values{}
	args.Set("merch_addr", m.Merch.Address)
	args.Set("merch_type", fmt.Sprintf("%v", m.Merch.Type))
	args.Set("timestamp", fmt.Sprintf("%d", m.Timestamp))
	args.Set("key", m.Key)
	signData := fmt.Sprintf("%v&access_token=%v", args.Encode(), token)
	sign = xhash.SHA1([]byte(signData))
	return
}

func (m *NotifyBody) VerifySign(token string) (err error) {
	sign := m.CalcSign(token)
	if !strings.EqualFold(strings.ToLower(sign), strings.ToLower(m.Sign)) {
		err = fmt.Errorf("sign error")
	}
	return
}
