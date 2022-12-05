package gexapi

import (
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/baseapi"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/captcha"
	"github.com/gexservice/gexservice/base/email"
	"github.com/gexservice/gexservice/base/sms"
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
	//captcha
	captcha.Hand(pre, mux)
	//user
	mux.HandleFunc("^"+pre+"/pub/login(\\?.*)?$", LoginH)
	mux.HandleFunc("^"+pre+"/pub/registerUser(\\?.*)?$", RegisterUserH)
	mux.HandleFunc("^"+pre+"/usr/logout(\\?.*)?$", LogoutH)
	mux.HandleFunc("^"+pre+"/usr/userInfo(\\?.*)?$", UserInfoH)
	mux.HandleFunc("^"+pre+"/usr/loadUser(\\?.*)?$", LoadUserH)
	mux.HandleFunc("^"+pre+"/usr/updateUser(\\?.*)?$", UpdateUserH)
	mux.HandleFunc("^"+pre+"/usr/updateUserConfig(\\?.*)?$", UpdateUserConfigH)
	mux.HandleFunc("^"+pre+"/usr/searchUser(\\?.*)?$", SearchUserH)
	// mux.HandleFunc("^"+pre+"/usr/searchMyUser(\\?.*)?$", SearchMyUserH)
	// mux.HandleFunc("^"+pre+"/usr/updateMyUserRemakr(\\?.*)?$", UpdateMyUserRemarkH)
	//balance
	mux.HandleFunc("^"+pre+"/admin/changeUserBalance(\\?.*)?$", ChangeUserBalanceH)
	mux.HandleFunc("^"+pre+"/usr/loadBalanceOverview(\\?.*)?$", LoadBalanceOverviewH)
	mux.HandleFunc("^"+pre+"/usr/listBalance(\\?.*)?$", ListBalanceH)
	mux.HandleFunc("^"+pre+"/usr/transferBalance(\\?.*)?$", TransferBalanceH)
	mux.HandleFunc("^"+pre+"/usr/listBalanceRecord(\\?.*)?$", ListBalanceRecordH)
	//withdraw
	mux.HandleFunc("^"+pre+"/usr/createWithdraw(\\?.*)?$", CreateWithdrawH)
	mux.HandleFunc("^"+pre+"/usr/cancelWithdraw(\\?.*)?$", CancelWithdrawH)
	mux.HandleFunc("^"+pre+"/usr/listWithdraw(\\?.*)?$", ListWithdrawH)
	mux.HandleFunc("^"+pre+"/usr/confirmWithdraw(\\?.*)?$", ConfirmWithdrawH)
	mux.HandleFunc("^"+pre+"/usr/loadTopupAddress(\\?.*)?$", LoadTopupAddressH)
	//holding
	mux.HandleFunc("^"+pre+"/usr/listHolding(\\?.*)?$", ListHoldingH)
	mux.HandleFunc("^"+pre+"/usr/loadHolding(\\?.*)?$", LoadHoldingH)
	mux.HandleFunc("^"+pre+"/usr/changeHoldingLever(\\?.*)?$", ChangeHoldingLeverH)
	//order
	mux.HandleFunc("^"+pre+"/usr/placeOrder", PlaceOrderH)
	mux.HandleFunc("^"+pre+"/usr/cancelOrder", CancelOrderH)
	mux.HandleFunc("^"+pre+"/usr/cancelAllOrder", CancelAllOrderH)
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
	//market favorites
	mux.HandleFunc("^"+pre+"/usr/listFavoritesSymbol(\\?.*)?$", ListFavoritesSymbolH)
	mux.HandleFunc("^"+pre+"/usr/addFavoritesSymbol(\\?.*)?$", AddFavoritesSymbolH)
	mux.HandleFunc("^"+pre+"/usr/removeFavoritesSymbol(\\?.*)?$", RemoveFavoritesSymbolH)
	mux.HandleFunc("^"+pre+"/usr/switchFavoritesSymbol(\\?.*)?$", SwitchFavoritesSymbolH)
	// mux.HandleFunc("^"+pre+"/pub/listMarketOrder(\\?.*)?$", ListMarketOrderH)
	//maker
	mux.HandleFunc("^"+pre+"/admin/loadSymbolMaker", LoadSymbolMakerH)
	mux.HandleFunc("^"+pre+"/admin/updateSymbolMaker", UpdateSymbolMakerH)
	mux.HandleFunc("^"+pre+"/admin/startSymbolMaker", StartSymbolMakerH)
	mux.HandleFunc("^"+pre+"/admin/stopSymbolMaker", StopSymbolMakerH)
	//message
	mux.HandleFunc("^"+pre+"/usr/addMessage(\\?.*)?$", AddMessageH)
	mux.HandleFunc("^"+pre+"/usr/removeMessage(\\?.*)?$", RemoveMessageH)
	mux.HandleFunc("^"+pre+"/usr/searchMessage(\\?.*)?$", SearchMessageH)
	//sms
	sms.Hand(pre, mux)
	//email
	email.Hand(pre, mux)
}

//Handle will register all handler
func HandleDebug(pre string, mux *web.SessionMux) {
	//sms
	sms.HandDebug(pre, mux)
	//email
	email.HandDebug(pre, mux)
}

func RecvValidJSON(s *web.Session, valider gexdb.Validable) (err error) {
	_, err = s.RecvJSON(interface{}(valider))
	if err == nil {
		err = valider.Valid()
	}
	return
}
