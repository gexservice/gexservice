package gexapi

import (
	"strings"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/maker"
	"github.com/gexservice/gexservice/market"
)

//LoadSymbolMakerH is http handler
/**
 *
 * @api {GET} /admin/loadSymbolMaker Load Symbol Maker
 * @apiName LoadSymbolMaker
 * @apiGroup Maker
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Bool} running whether maker is running
 * @apiSuccess (Symbols) {Array} symbol the symbol info list
 * @apiSuccess (Symbols) {String} symbol.base the symbol base asset
 * @apiSuccess (Symbols) {String} symbol.quote the symbol quote asset
 * @apiSuccess (Symbols) {String} symbol.fee the symbol trade fee
 * @apiSuccess (Symbols) {String} symbol.precision_price the symbol price percision
 * @apiSuccess (Symbols) {String} symbol.precision_quantity the symbol quantity percision
 * @apiSuccess (MakerConfig) {Object} maker the symbol maker config, mapping key by symbol
 * @apiSuccess (MakerConfig) {String} maker.on whether auto start when system start
 * @apiSuccess (MakerConfig) {String} maker.user_id the maker user id
 * @apiSuccess (MakerConfig) {String} maker.delay the depth price change delay, must 0<delay<1000, default 20
 * @apiSuccess (MakerConfig) {Int} maker.interval the maker cycle interval(ms), must interval>1min
 * @apiSuccess (MakerConfig) {String} maker.open the open price
 * @apiSuccess (MakerConfig) {String} maker.close.max the close price max change rate after interval, must close.max>close.min>-1
 * @apiSuccess (MakerConfig) {String} maker.close.min the close price min change rate after interval, must close.max>close.min>-1
 * @apiSuccess (MakerConfig) {String} maker.vib.max the vib price max change rate, must vib.max>vib.min>-1
 * @apiSuccess (MakerConfig) {String} maker.vib.min the vib price min change rate, must vib.max>vib.min>-1
 * @apiSuccess (MakerConfig) {Int} maker.vib.count the vib count in interval, must vib.count>0
 * @apiSuccess (MakerConfig) {String} maker.ticker the max ticker change rate, must ticker>0
 * @apiSuccess (MakerConfig) {String} maker.depth.qty_max the depth max qty(per place), must depth.qty_max>0
 * @apiSuccess (MakerConfig) {String} maker.depth.step_max the depth max step(per place), must depth.step_max>0
 * @apiSuccess (MakerConfig) {String} maker.depth.diff_max the depth ask/bid max diff, must depth.diff_max>depth.diff_min
 * @apiSuccess (MakerConfig) {String} maker.depth.diff_min the depth ask/bid max diff, must depth.diff_max>depth.diff_min
 * @apiSuccess (MakerConfig) {String} maker.depth.max the depth max count, must depth.max>3
 * @apiSuccess (BalanceObject) {Object} balances the maker balance info, mapping by asset
 * @apiUse BalanceObject
 * @apiSuccess (Holdingbject) {Object} holding the maker holding info
 * @apiUse Holdingbject
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "balances": {
 *         "USDT": {
 *             "area": 200,
 *             "asset": "USDT",
 *             "create_time": 1667654528608,
 *             "free": "1000",
 *             "locked": "0",
 *             "margin": "0",
 *             "status": 100,
 *             "tid": 1005,
 *             "update_time": 1667654528620,
 *             "user_id": 100002
 *         },
 *         "YWE": {
 *             "area": 200,
 *             "asset": "YWE",
 *             "create_time": 1667654528608,
 *             "free": "999",
 *             "locked": "1",
 *             "margin": "0",
 *             "status": 100,
 *             "tid": 1004,
 *             "update_time": 1667654528628,
 *             "user_id": 100002
 *         }
 *     },
 *     "code": 0,
 *     "holding": null,
 *     "maker": {
 *         "close": {
 *             "max": "0.01",
 *             "min": "-0.01"
 *         },
 *         "delay": 20,
 *         "depth": {
 *             "diff_max": "2",
 *             "diff_min": "0.02",
 *             "max": 15,
 *             "qty_max": "3",
 *             "step_max": 5
 *         },
 *         "interval": 3600000,
 *         "on": 0,
 *         "open": "1000",
 *         "symbol": "spot.YWEUSDT",
 *         "ticker": "0.0001",
 *         "user_id": 100002,
 *         "vib": {
 *             "count": 5,
 *             "max": "0.03",
 *             "min": "-0.03"
 *         }
 *     },
 *     "symbol": {
 *         "area": 0,
 *         "base": "YWE",
 *         "fee": "0.002",
 *         "margin_add": "0.01",
 *         "margin_max": "0.99",
 *         "precision_price": 8,
 *         "precision_quantity": 8,
 *         "quote": "USDT",
 *         "symbol": "spot.YWEUSDT"
 *     },
 *     "user": {
 *         "account": "abc0",
 *         "create_time": 1667654528598,
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1667654528598
 *     }
 * }
 */
func LoadSymbolMakerH(s *web.Session) web.Result {
	var symbol string
	err := s.ValidFormat(`
		symbol,R|S,L:0;
	`, &symbol)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	info, _ := market.LoadSymbol(symbol)
	config, err := maker.LoadConfig(s.R.Context(), symbol)
	if err != nil && err != pgx.ErrNoRows {
		xlog.Errorf("LoadSymbolMakerH load config by %v fail with %v", symbol, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	runner := maker.Find(s.R.Context(), symbol)
	var user *gexdb.User
	var balances map[string]*gexdb.Balance
	var holding *gexdb.Holding
	if config != nil && config.UserID > 0 {
		user, err = gexdb.FindUser(s.R.Context(), config.UserID)
		if err != nil {
			xlog.Errorf("LoadSymbolMakerH find user by %v fail with %v", config.UserID, err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		area := gexdb.BalanceAreaFutures
		if strings.HasPrefix(symbol, "spot.") {
			area = gexdb.BalanceAreaSpot
		}
		_, balances, err = gexdb.ListUserBalance(s.R.Context(), user.TID, area, nil, nil)
		if err != nil {
			xlog.Errorf("LoadSymbolMakerH list balance by %v,%v fail with %v", user.TID, area, err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		holding, err = gexdb.FindHoldlingBySymbol(s.R.Context(), user.TID, symbol)
		if err != nil && err != pgx.ErrNoRows {
			xlog.Errorf("LoadSymbolMakerH load holding by %v,%v fail with %v", user.TID, symbol, err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"symbol":   info,
		"maker":    config,
		"user":     user,
		"balances": balances,
		"holding":  holding,
		"running":  runner != nil,
	})
}

//UpdateSymbolMakerH is http handler
/**
 *
 * @api {GET} /admin/updateSymbolMaker Update Symbol Maker
 * @apiName UpdateSymbolMaker
 * @apiGroup Maker
 *
 * @apiParam  {String} on whether auto start when system start
 * @apiParam  {String} symbol the symbol to maker
 * @apiParam  {String} user_id the maker user id
 * @apiParam  {String} delay the depth price change delay, must 0<delay<1000, default 20
 * @apiParam  {Int} interval the maker cycle interval(ms), must interval>1min
 * @apiParam  {String} open the open price
 * @apiParam  {String} close.max the close price max change rate after interval, must close.max>close.min>-1
 * @apiParam  {String} close.min the close price min change rate after interval, must close.max>close.min>-1
 * @apiParam  {String} vib.max the vib price max change rate, must vib.max>vib.min>-1
 * @apiParam  {String} vib.min the vib price min change rate, must vib.max>vib.min>-1
 * @apiParam  {Int} vib.count the vib count in interval, must vib.count>0
 * @apiParam  {String} ticker the max ticker change rate, must ticker>0
 * @apiParam  {String} depth.qty_max the depth max qty(per place), must depth.qty_max>0
 * @apiParam  {String} depth.step_max the depth max step(per place), must depth.step_max>0
 * @apiParam  {String} depth.diff_max the depth ask/bid max diff, must depth.diff_max>depth.diff_min
 * @apiParam  {String} depth.diff_min the depth ask/bid max diff, must depth.diff_max>depth.diff_min
 * @apiParam  {String} depth.max the depth max count, must depth.max>3
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 * @apiParamExample  {Query} Update
 * {
 *     "close": {
 *         "max": "0.01",
 *         "min": "-0.01"
 *     },
 *     "delay": 20,
 *     "depth": {
 *         "diff_max": "2",
 *         "diff_min": "0.02",
 *         "max": 15,
 *         "qty_max": "3",
 *         "step_max": 5
 *     },
 *     "interval": 3600000,
 *     "on": 0,
 *     "open": "1000",
 *     "symbol": "spot.YWEUSDT",
 *     "ticker": "0.0001",
 *     "user_id": 100002,
 *     "vib": {
 *         "count": 5,
 *         "max": "0.03",
 *         "min": "-0.03"
 *     }
 * }
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 */
func UpdateSymbolMakerH(s *web.Session) web.Result {
	var config maker.Config
	err := RecvValidJSON(s, &config)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	err = maker.UpdateConfig(s.R.Context(), &config)
	if err != nil {
		xlog.Warnf("UpdateSymbolMakerH update maker by %v fail with %v", converter.JSON(config), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//StartSymbolMakerH is http handler
/**
 *
 * @api {GET} /admin/startSymbolMaker Start Symbol Maker
 * @apiName StartSymbolMaker
 * @apiGroup Maker
 *
 * @apiParam  {String} symbol the symbol to start
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 */
func StartSymbolMakerH(s *web.Session) web.Result {
	var symbol string
	err := s.ValidFormat(`
		symbol,R|S,L:0;
	`, &symbol)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	err = maker.Start(s.R.Context(), symbol)
	if err != nil {
		xlog.Warnf("StartSymbolMakerH start maker by %v fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//StopSymbolMakerH is http handler
/**
 *
 * @api {GET} /admin/stopSymbolMaker Stop Symbol Maker
 * @apiName StopSymbolMaker
 * @apiGroup Maker
 *
 * @apiParam  {String} symbol the symbol to stop
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 */
func StopSymbolMakerH(s *web.Session) web.Result {
	var symbol string
	err := s.ValidFormat(`
		symbol,R|S,L:0;
	`, &symbol)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	err = maker.Stop(s.R.Context(), symbol)
	if err != nil {
		xlog.Warnf("StartSymbolMakerH start maker by %v fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}
