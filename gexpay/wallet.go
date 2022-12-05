package gexpay

import (
	"fmt"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

func AssignWallet(method gexdb.WalletMethod) (address string, err error) {
	timestamp, sign := RequestSign(method, "/usr/newUserWallet")
	applyURL := fmt.Sprintf(
		"%v/usr/newUserWallet?merch_type=%v&merch_addr=%v&sign=%v&timestamp=%v",
		AgentAddr, MerchType[method], MerchAddr[method], sign, timestamp,
	)
	xlog.Infof("AssignWallet assign wallet is starting by %v", applyURL)
	res, err := xhttp.GetMap("%v", applyURL)
	if err != nil {
		xlog.Errorf("AssignWallet assign wallet fail with %v by %v", err, applyURL)
		return
	}
	if code := res.IntDef(-1, "code"); code != define.Success {
		err = fmt.Errorf("%v", converter.JSON(res))
		xlog.Errorf("AssignWallet assign wallet fail with %v by %v", err, applyURL)
		return
	}
	address = res.StrDef("", "/wallet/address")
	if len(address) < 1 {
		err = fmt.Errorf("%v", converter.JSON(res))
		xlog.Errorf("AssignWallet assign wallet fail with %v by %v", err, applyURL)
		return
	}
	xlog.Infof("AssignWallet assign wallet done with %v by %v", converter.JSON(res), applyURL)
	return
}
