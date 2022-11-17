package gexapi

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/util/xtime"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
	"github.com/gexservice/gexservice/matcher"
)

//ListSymbolH is http handler
/**
 *
 * @api {GET} /pub/listSymbol List Symbol
 * @apiName ListSymbol
 * @apiGroup Market
 *
 * @apiParam  {String} type the symbol type, supported in spot/futures
 * @apiParam  {String} orderby the symbol orderby, supported in +rate/-rate/+volume/-volume
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (SymbolInfo) {Array} symbols the symbol info list
 * @apiUse SymbolInfoObject
 * @apiSuccess (KLine) {Object} days the symbol day change line, mapping by key is symbol
 * @apiUse KLineObject
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "days": {
 *         "spot.YWEUSDT": {
 *             "amount": "0.5",
 *             "close": "100",
 *             "count": 1,
 *             "high": "100",
 *             "interv": "1day",
 *             "low": "100",
 *             "open": "100",
 *             "start_time": 1667404800000,
 *             "symbol": "spot.YWEUSDT",
 *             "update_time": 1667486761280,
 *             "volume": "50"
 *         }
 *     },
 *     "symbols": [
 *         {
 *             "base": "YWE",
 *             "fee": "0.002",
 *             "margin_add": "0.01",
 *             "margin_max": "0.99",
 *             "precision_price": 8,
 *             "precision_quantity": 8,
 *             "quote": "USDT",
 *             "symbol": "spot.YWEUSDT"
 *         },
 *         {
 *             "base": "YWE",
 *             "fee": "0.002",
 *             "margin_add": "0.01",
 *             "margin_max": "0.99",
 *             "precision_price": 8,
 *             "precision_quantity": 8,
 *             "quote": "USDT",
 *             "symbol": "futures.YWEUSDT"
 *         }
 *     ]
 * }
 *
 */
func ListSymbolH(s *web.Session) web.Result {
	symbols, _, days := market.ListSymbol(s.Argument("type"), nil, s.Argument("orderby"))
	return s.SendJSON(xmap.M{
		"code":    0,
		"symbols": symbols,
		"days":    days,
	})
}

//LoadSymbolH is http handler
/**
 *
 * @api {GET} /pub/loadSymbol Load Symbol
 * @apiName LoadSymbol
 * @apiGroup Market
 *
 * @apiParam  {String} symbol the symbol
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (SymbolInfo) {Object} symbol the symbol info
 * @apiUse SymbolInfoObject
 * @apiSuccess (KLine) {Object} day the symbol day change line
 * @apiUse KLineObject
 *
 * @apiParamExample  {Query} QueryOrder:
 * symbol=spot.YWEUSDT
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "day": {
 *         "amount": "0.5",
 *         "close": "100",
 *         "count": 1,
 *         "high": "100",
 *         "interv": "1day",
 *         "low": "100",
 *         "open": "100",
 *         "start_time": 1667404800000,
 *         "symbol": "spot.YWEUSDT",
 *         "update_time": 1667486840555,
 *         "volume": "50"
 *     },
 *     "symbol": {
 *         "base": "YWE",
 *         "fee": "0.002",
 *         "margin_add": "0.01",
 *         "margin_max": "0.99",
 *         "precision_price": 8,
 *         "precision_quantity": 8,
 *         "quote": "USDT",
 *         "symbol": "spot.YWEUSDT"
 *     }
 * }
 *
 */
func LoadSymbolH(s *web.Session) web.Result {
	symbol, day := market.LoadSymbol(s.Argument("symbol"))
	return s.SendJSON(xmap.M{
		"code":   0,
		"symbol": symbol,
		"day":    day,
	})
}

//Market is struct to market impl
/**
 *
 * @api {WS} /ws/market WS Market
 * @apiName WsMarket
 * @apiGroup Market
 *
 * @apiParam  {String} action subscribe action, supported is "sub.kline"/"sub.depth"/"sub.ticker"
 * @apiParam  {Array} symbols the symbol to sub
 * @apiParam  {String} symbols.symbol the market symbol
 * @apiParam  {Arrasy} [symbols.interval] the kline interval, only for "sub.kline", supported is "5min"/"30min"/"1hour"/"4hour"/"day"/"week"/"mon"
 * @apiParam  {Number} [symbols.max] the depth max size
 *
 * @apiSuccess (Success) {Number} code the response code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {String} action the received action, supported is "sub.kline"/"sub.depth"/"notify.kline"/"notify.depth"
 * @apiSuccess (Ticker) {Object} ticker the received ticker data, only for "notify.ticker"
 * @apiSuccess (Ticker) {Array} ticker.ask the received ticker best ask, the inner data is ["price","quantity"]
 * @apiSuccess (Ticker) {Array} ticker.bid the received ticker best bid, the inner data is ["price","quantity"]
 * @apiSuccess (Ticker) {Decimal} ticker.close the received ticker latest cose price
 * @apiSuccess (Depth) {Object} depth the received depth data, only for "notify.depth"
 * @apiSuccess (Depth) {String} depth.symbol the received depth symbol
 * @apiSuccess (Depth) {Array} depth.bids the received depth bids data, the inner data is ["price","quantity"]
 * @apiSuccess (Depth) {Array} depth.asks the received depth asks data, the inner data is ["price","quantity"]
 * @apiSuccess (KLine) {Object} kline the received kline data, only for "notify.kline"
 * @apiUse KLineObject
 *
 * @apiParamExample  {JSON} Subscribe-KLine:
 * {
 *     "action": "sub.kline",
 *     "symbols": [
 *         {
 *             "symbol": "spot.YWEUSDT",
 *             "interval": "5min"
 *         }
 *     ]
 * }
 * @apiParamExample  {JSON} Subscribe-Depth:
 * {
 *     "action": "sub.depth",
 *     "symbols": [
 *         {
 *             "symbol": "spot.YWEUSDT",
 *             "max": 30
 *         }
 *     ]
 * }
 * @apiParamExample  {JSON} Subscribe-Ticker:
 * {
 *     "action": "sub.ticker",
 *     "symbols": ["spot.YWEUSDT"]
 * }
 *
 * @apiSuccessExample {JSON} Reponse-Depth:
 * {
 *     "action": "sub.kline",
 *     "code": 0
 * }
 *
 * @apiSuccessExample {JSON} Reponse-Depth:
 * {
 *     "action": "sub.depth",
 *     "code": 0
 * }
 * @apiSuccessExample {JSON} Reponse-Ticker:
 * {
 *     "action": "sub.ticker",
 *     "code": 0
 * }
 *
 * @apiSuccessExample {JSON} Notify-Depth:
 * {
 *     "code": 0,
 *     "action": "notify.depth",
 *     "depth": {
 *         "symbol": "spot.YWEUSDT",
 *         "bids": [
 *             [
 *                 "90",
 *                 "2"
 *             ]
 *         ],
 *         "asks": [
 *             [
 *                 "100",
 *                 "2"
 *             ]
 *         ]
 *     }
 * }
 * @apiSuccessExample {JSON} Notify-Ticker:
 * {
 *     "code": 0,
 *     "action": "notify.ticker",
 *     "ticker": {
 *         "symbol": "spot.YWEUSDT",
 *         "bid": [
 *             "90",
 *             "2"
 *         ],
 *         "ask": [
 *             "100",
 *             "2"
 *         ],
 *         "close": "100"
 *     }
 * }
 *
 * @apiSuccessExample {JSON} Notify-KLine:
 * {
 *     "action": "notify.kline",
 *     "code": 0,
 *     "kline": {
 *         "tid": 0,
 *         "symbol": "YWEMMK",
 *         "interval": "5min",
 *         "amount": "1",
 *         "count": 1,
 *         "open": "100",
 *         "close": "100",
 *         "low": "100",
 *         "high": "100",
 *         "volume": "100",
 *         "start_time": 1632577200000,
 *         "update_time": 1632577495577
 *     }
 * }
 *
 */
var MarketOnline *OnlineHander

//ListKLineH is http handler
/**
 *
 * @api {GET} /pub/listKLine List KLine
 * @apiName ListKLine
 * @apiGroup Market
 *
 * @apiParam  {String} symbol the kline symbol
 * @apiParam  {String} interval the kline interval, supported is "5min"/"30min"/"1hour"/"4hour"/"1day"/"1week"/"1mon"
 * @apiParam  {Number} start_time filter kline kline.start_time>=start_time
 * @apiParam  {Number} end_time filter kline kline.start_time<end_time
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (KLine) {Array} lines the all kline array
 * @apiUse KLineObject
 *
 * @apiParamExample  {Query} QueryOrder:
 * symbol=spot.YWEUSDT&interval=100&start_time=100&end_time=1632578100000
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "lines": [
 *         {
 *             "tid": 0,
 *             "symbol": "YWEMMK",
 *             "interval": "5min",
 *             "amount": "1",
 *             "count": 1,
 *             "open": "100",
 *             "close": "100",
 *             "low": "100",
 *             "high": "100",
 *             "volume": "100",
 *             "start_time": 1632578100000,
 *             "update_time": 1632578330897
 *         }
 *     ]
 * }
 *
 */
func ListKLineH(s *web.Session) web.Result {
	var symbol, interval string
	var startTime, endTime int64
	var err = s.ValidFormat(`
		symbol,R|S,L:0;
		interval,R|S,L:0;
		start_time,O|I,R:0;
		end_time,O|I,R:0;
	`, &symbol, &interval, &startTime, &endTime)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	lines, err := market.ListKLine(s.R.Context(), symbol, interval, xtime.TimeUnix(startTime), xtime.TimeUnix(endTime))
	if err != nil {
		xlog.Warnf("ListKLineH list kline fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":  0,
		"lines": lines,
	})
}

//LoadDepthH is http handler
/**
 *
 * @api {GET} /pub/loadDepth Load Depth
 * @apiName LoadDepth
 * @apiGroup Market
 *
 * @apiParam  {String} symbol the market symbol
 * @apiParam  {Number} [max] max depth
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} depth the depth info
 * @apiSuccess (Success) {Array} depth.bids the depth bid array
 * @apiSuccess (Success) {Array} depth.asks the depth ask array
 *
 * @apiParamExample  {Query} QueryOrder:
 * max=10
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "depth": {
 *         "bids": [
 *             [
 *                 "86.5",
 *                 "0.5"
 *             ]
 *         ],
 *         "asks": [
 *             [
 *                 "125",
 *                 "1"
 *             ]
 *         ]
 *     }
 * }
 *
 */
func LoadDepthH(s *web.Session) web.Result {
	var symbol string
	var max int = 8
	var err = s.ValidFormat(`
		symbol,R|S,L:0;
		max,O|I,R:0;
	`, &symbol, &max)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	depth := market.LoadDepth(symbol, max)
	return s.SendJSON(xmap.M{
		"code":  0,
		"depth": depth,
	})
}

//ListFavoritesSymbolH is http handler
/**
 *
 * @api {GET} /usr/listFavoritesSymbol List Favorites Symbol
 * @apiName ListFavoritesSymbol
 * @apiGroup Market
 *
 * @apiParam  {String} type the symbol type, supported in spot/futures
 * @apiParam  {String} orderby the symbol orderby, supported in +rate/-rate/+volume/-volume
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (SymbolInfo) {Array} symbol the symbol info
 * @apiUse SymbolInfoObject
 * @apiSuccess (KLine) {Object} days the symbol day change line, mapping by key is symbol
 * @apiUse KLineObject
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "days": {
 *         "spot.YWEUSDT": {
 *             "amount": "0.5",
 *             "close": "100",
 *             "count": 1,
 *             "high": "100",
 *             "interv": "1day",
 *             "low": "100",
 *             "open": "100",
 *             "start_time": 1667404800000,
 *             "symbol": "spot.YWEUSDT",
 *             "update_time": 1667486761280,
 *             "volume": "50"
 *         }
 *     },
 *     "symbols": [
 *         {
 *             "base": "YWE",
 *             "fee": "0.002",
 *             "margin_add": "0.01",
 *             "margin_max": "0.99",
 *             "precision_price": 8,
 *             "precision_quantity": 8,
 *             "quote": "USDT",
 *             "symbol": "spot.YWEUSDT"
 *         },
 *         {
 *             "base": "YWE",
 *             "fee": "0.002",
 *             "margin_add": "0.01",
 *             "margin_max": "0.99",
 *             "precision_price": 8,
 *             "precision_quantity": 8,
 *             "quote": "USDT",
 *             "symbol": "futures.YWEUSDT"
 *         }
 *     ]
 * }
 *
 */
func ListFavoritesSymbolH(s *web.Session) web.Result {
	userID := s.Value("user_id").(int64)
	favorites, err := gexdb.LoadUserFavorites(s.R.Context(), userID)
	if err != nil {
		xlog.Errorf("ListFavoritesSymbolH load user favorites fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	var symbols []*matcher.SymbolInfo
	var days map[string]*gexdb.KLine
	if favorites != nil && len(favorites.Symbols) > 0 {
		symbols, _, days = market.ListSymbol(s.Argument("type"), favorites.Symbols, s.Argument("orderby"))
	}
	return s.SendJSON(xmap.M{
		"code":    0,
		"symbols": symbols,
		"days":    days,
	})
}

//AddFavoritesSymbolH is http handler
/**
 *
 * @api {GET} /usr/addFavoritesSymbol Add Favorites Symbol
 * @apiName AddFavoritesSymbol
 * @apiGroup Market
 *
 * @apiParam  {String} symbol the symbol to add
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 *
 */
func AddFavoritesSymbolH(s *web.Session) web.Result {
	var symbol string
	err := s.ValidFormat(`
		symbol,r|s,l:0;
	`, &symbol)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	err = gexdb.UpdateUserFavorites(s.R.Context(), userID, func(favorites *gexdb.UserFavorites) {
		exists := false
		for _, s := range favorites.Symbols {
			if symbol == s {
				exists = true
				break
			}
		}
		if !exists {
			favorites.Symbols = append(favorites.Symbols, symbol)
		}
	})
	if err != nil {
		xlog.Errorf("AddFavoritesSymbolH update user favorites fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//RemoveFavoritesSymbolH is http handler
/**
 *
 * @api {GET} /usr/removeFavoritesSymbol Remove Favorites Symbol
 * @apiName RemoveFavoritesSymbol
 * @apiGroup Market
 *
 * @apiParam  {String} symbol the symbol to remove, multi split by comma
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 *
 */
func RemoveFavoritesSymbolH(s *web.Session) web.Result {
	var symbols []string
	err := s.ValidFormat(`
		symbol,r|s,l:0;
	`, &symbols)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	err = gexdb.UpdateUserFavorites(s.R.Context(), userID, func(favorites *gexdb.UserFavorites) {
		symbolAll := []string{}
		for _, s := range favorites.Symbols {
			if xsql.StringArray(symbols).HavingOne(s) {
				continue
			}
			symbolAll = append(symbolAll, s)
		}
		favorites.Symbols = symbolAll
	})
	if err != nil {
		xlog.Errorf("RemoveFavoritesSymbolH remove user favorites fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//SwitchFavoritesSymbolH is http handler
/**
 *
 * @api {GET} /usr/switchFavoritesSymbol Switch Favorites Symbol
 * @apiName SwitchFavoritesSymbol
 * @apiGroup Market
 *
 * @apiParam  {String} symbol the symbol to add
 * @apiParam  {String} [to] the symbol to swith, if empty will switch to top
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 *
 */
func SwitchFavoritesSymbolH(s *web.Session) web.Result {
	var symbol, to string
	err := s.ValidFormat(`
		symbol,r|s,l:0;
		to,o|s,l:0;
	`, &symbol, &to)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	err = gexdb.UpdateUserFavorites(s.R.Context(), userID, func(favorites *gexdb.UserFavorites) {
		if len(to) > 0 {
			favorites.SwitchSymbol(symbol, to)
		} else {
			favorites.TopSymbol(symbol)
		}
	})
	if err != nil {
		xlog.Errorf("SwitchFavoritesSymbolH switch favorites symbol fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}
