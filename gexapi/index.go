package gexapi

import (
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

//IndexH is http handler
/**
 *
 * @api {GET} /pub/index Index
 * @apiName Index
 * @apiGroup Index
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Array} international_prices the all location price
 * @apiSuccess (Success) {String} international_prices.location the location name by multi language
 * @apiSuccess (Success) {String} international_prices.price the location price
 *
 * @apiParamExample  {Query} QueryOrder:
 * interval=100&start_time=100&end_time=1632578100000
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "international_prices": [
 *         {
 *             "location": {
 *                 "CN": "中国",
 *                 "MM": "တရုတ်",
 *                 "US": "China"
 *             },
 *             "price": 370.85
 *         },
 *         {
 *             "location": {
 *                 "CN": "伦敦",
 *                 "MM": "London",
 *                 "US": "London"
 *             },
 *             "price": 1783.38
 *         },
 *         {
 *             "location": {
 *                 "CN": "纽约",
 *                 "MM": "New York",
 *                 "US": "New York"
 *             },
 *             "price": 1783.74
 *         }
 *     ]
 * }
 *
 */
func IndexH(s *web.Session) web.Result {
	config, err := basedb.LoadConfigList(s.R.Context(), gexdb.ConfigBalanceImage, gexdb.ConfigWelcomeMessage)
	if err != nil {
		xlog.Errorf("IndexH load config fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	config["international_prices"] = InternationalPrice
	return s.SendJSON(config)
}
