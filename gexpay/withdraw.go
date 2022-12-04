package gexpay

import (
	"fmt"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

func ApplyWithdraw(withdraw *gexdb.Withdraw) (err error) {
	method := gexdb.WalletMethod(withdraw.Method)
	timestamp, sign := RequestSign(method, "/usr/applyUserWithdraw")
	applyURL := fmt.Sprintf(
		"%v/usr/applyUserWithdraw?merch_type=%v&merch_addr=%v&user_addr=%v&to_addr=%v&asset=%v&amount=%v&try_count=1&uuid=%v&sign=%v&timestamp=%v",
		AgentAddr, MerchType[method], MerchAddr[method], MerchAddr[method], withdraw.Receiver, withdraw.Asset, withdraw.Quantity, withdraw.OrderID, sign, timestamp,
	)
	res, err := xhttp.GetMap("%v", applyURL)
	if err != nil {
		xlog.Errorf("ApplyWithdraw apply withdraw fail with %v by %v", err, applyURL)
		return
	}
	if code := res.IntDef(-1, "code"); code != define.Success && code != define.Duplicate {
		err = fmt.Errorf("%v", converter.JSON(res))
		xlog.Errorf("ApplyWithdraw apply withdraw fail with %v by %v", err, applyURL)
		return
	}
	return
}
