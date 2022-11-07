package gexapi

import (
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/baseapi"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
)

var SrvAddr = func() string {
	panic("SrvAddr is not initial")
}

//Handle will register all handler
func Handle(pre string, mux *web.SessionMux) {
	baseapi.EditAnnounceAccess = AdminAccess
	baseapi.EditSysConfigAccess = AdminAccess
	baseapi.EditVersionObjectAccess = AdminAccess
	basedb.ConfigAll = gexdb.ConfigAll
	mux.Handle("^"+pre+"/conf/mp(\\?.*)?$", ConfMPH)
	mux.Handle("^"+pre+"/conf/admin(\\?.*)?$", ConfAdminH)
	mux.Handle("^"+pre+"/conf/rule(\\?.*)?$", ConfRuleH)
	mux.Handle("^"+pre+"/conf/goldbar(\\?.*)?$", ConfGoldbarH)
	mux.FilterFunc("^"+pre+"/usr/.*$", LoginAccessF)
	mux.FilterFunc("^"+pre+"/admin/.*$", AdminAccessF)
	baseapi.Handle(pre, mux)
	//index
	mux.HandleFunc("^"+pre+"/pub/index(\\?.*)?$", IndexH)
	//user
	mux.HandleFunc("^"+pre+"/pub/login(\\?.*)?$", LoginH)
	mux.HandleFunc("^"+pre+"/usr/logout(\\?.*)?$", LogoutH)
	mux.HandleFunc("^"+pre+"/usr/userInfo(\\?.*)?$", UserInfoH)
	mux.HandleFunc("^"+pre+"/usr/loadUser(\\?.*)?$", LoadUserH)
	mux.HandleFunc("^"+pre+"/usr/updateUser(\\?.*)?$", UpdateUserH)
	mux.HandleFunc("^"+pre+"/usr/searchUser(\\?.*)?$", SearchUserH)
	// mux.HandleFunc("^"+pre+"/usr/searchMyUser(\\?.*)?$", SearchMyUserH)
	// mux.HandleFunc("^"+pre+"/usr/updateMyUserRemakr(\\?.*)?$", UpdateMyUserRemarkH)
	//balance
	// mux.HandleFunc("^"+pre+"/usr/changeUserBalance(\\?.*)?$", ChangeUserBalanceH)
	mux.HandleFunc("^"+pre+"/usr/loadBalanceOverview(\\?.*)?$", LoadBalanceOverviewH)
	mux.HandleFunc("^"+pre+"/usr/listBalance(\\?.*)?$", ListBalanceH)
	//holding
	mux.HandleFunc("^"+pre+"/usr/listHolding(\\?.*)?$", ListHoldingH)
	mux.HandleFunc("^"+pre+"/usr/changeHoldingLever(\\?.*)?$", ChangeHoldingLeverH)
	//order
	mux.HandleFunc("^"+pre+"/usr/createWithdrawOrder(\\?.*)?$", CreateWithdrawOrderH)
	mux.HandleFunc("^"+pre+"/usr/cancelWithdrawOrder(\\?.*)?$", CancelWithdrawOrderH)
	// mux.HandleFunc("^"+pre+"/usr/createGoldbarOrder(\\?.*)?$", CreateGoldbarOrderH)
	// mux.HandleFunc("^"+pre+"/usr/cancelGoldbarOrder(\\?.*)?$", CancelGoldbarOrderH)
	// mux.HandleFunc("^"+pre+"/usr/verifyGoldbarOrder(\\?.*)?$", VerifyGoldbarOrderH)
	mux.HandleFunc("^"+pre+"/usr/createTopupOrder(\\?.*)?$", CreateTopupOrderH)
	// mux.HandleFunc("^"+pre+"/usr/searchMyUserOrder(\\?.*)?$", SearchMyUserOrderH)
	mux.HandleFunc("^"+pre+"/usr/placeOrder", PlaceOrderH)
	mux.HandleFunc("^"+pre+"/usr/cancelOrder", CancelOrderH)
	mux.HandleFunc("^"+pre+"/usr/searchOrder(\\?.*)?$", SearchOrderH)
	mux.HandleFunc("^"+pre+"/usr/queryOrder(\\?.*)?$", QueryOrderH)
	// mux.HandleFunc("^"+pre+"/usr/countOrderComm(\\?.*)?$", CountOrderCommH)
	//market
	mux.HandleFunc("^"+pre+"/pub/listSymbol(\\?.*)?$", ListSymbolH)
	mux.HandleFunc("^"+pre+"/pub/loadSymbol(\\?.*)?$", LoadSymbolH)
	MarketOnline = NewOnlineHander(mux, market.Shared)
	mux.Handle("^"+pre+"/ws/market(\\?.*)?$", MarketOnline)
	mux.HandleFunc("^"+pre+"/pub/listKLine(\\?.*)?$", ListKLineH)
	mux.HandleFunc("^"+pre+"/pub/loadDepth(\\?.*)?$", LoadDepthH)
	// mux.HandleFunc("^"+pre+"/pub/listMarketOrder(\\?.*)?$", ListMarketOrderH)
	//maker
	mux.HandleFunc("^"+pre+"/admin/loadSymbolMaker", LoadSymbolMakerH)
	mux.HandleFunc("^"+pre+"/admin/updateSymbolMaker", UpdateSymbolMakerH)
	mux.HandleFunc("^"+pre+"/admin/startSymbolMaker", StartSymbolMakerH)
	mux.HandleFunc("^"+pre+"/admin/stopSymbolMaker", StopSymbolMakerH)
}

func RecvValidJSON(s *web.Session, valider gexdb.Validable) (err error) {
	_, err = s.RecvJSON(interface{}(valider))
	if err == nil {
		err = valider.Valid()
	}
	return
}
