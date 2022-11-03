package gexapi

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
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
	return s.SendJSON(xmap.M{
		"international_prices": InternationalPrice,
	})
}
