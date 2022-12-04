package gexpay

import (
	"fmt"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

func ApplyWallet(method gexdb.WalletMethod) (address string, err error) {
	timestamp, sign := RequestSign(method, "/usr/newUserWallet")
	applyURL := fmt.Sprintf(
		"%v/usr/newUserWallet?merch_type=%v&merch_addr=%v&sign=%v&timestamp=%v",
		AgentAddr, MerchType[method], MerchAddr[method], sign, timestamp,
	)
	res, err := xhttp.GetMap("%v", applyURL)
	if err != nil {
		xlog.Errorf("ApplyWallet apply wallet fail with %v by %v", err, applyURL)
		return
	}
	if code := res.IntDef(-1, "code"); code != define.Success {
		err = fmt.Errorf("%v", converter.JSON(res))
		xlog.Errorf("ApplyWallet apply wallet fail with %v by %v", err, applyURL)
		return
	}
	address = res.StrDef("", "/wallet/address")
	if len(address) < 1 {
		err = fmt.Errorf("%v", converter.JSON(res))
		xlog.Errorf("ApplyWallet apply wallet fail with %v by %v", err, applyURL)
		return
	}
	return
}
