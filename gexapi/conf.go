package gexapi

import (
	"github.com/codingeasygo/web/handler"
	"github.com/gexservice/gexservice/base/baseapi"
	"github.com/gexservice/gexservice/gexdb"
)

//ConfMPH is http handler to load mp configure
/**
 *
 * @api {GET} /conf/mp MP Config
 * @apiName ConfMP
 * @apiGroup Conf
 *
 *
 * @apiSuccess (Success) {String} xxx_api the xxx api address, split by comma, xxx having binance/okex/huobi/bybit
 * @apiSuccess (Success) {String} xxx_wss the xxx websocket server address, xxx having binance/okex/huobi/bybit
 * @apiSuccess (Success) {String} img_base_url the image base url
 * @apiSuccess (Success) {String} custom_service_url the custom service page web url
 * @apiSuccess (Success) {String} tutorial_url the tutorial page web url
 * @apiSuccess (Success) {String} algorithm_url the algorithm page web url
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "binance_api": "https://api.binance.com,https://api1.binance.com",
 *     "binance_wss": "wss://stream.binance.com:9443,wss://stream1.binance.com:9443",
 *     "img_base_url": "https://xxx.com/"
 * }
 *
 */
var ConfMPH = handler.Map{}

//ConfAdminH is http handler to load admin configure
/**
 *
 * @api {GET} /conf/admin Admin Config
 * @apiName ConfAdmin
 * @apiGroup Conf
 *
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 * }
 *
 */
var ConfAdminH = handler.Map{}

//ConfRuleH is http handler to load rule configure
/**
 *
 * @api {GET} /conf/rule Rule Config
 * @apiName ConfRule
 * @apiGroup Conf
 *
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code":0,
 *     "config":{
 *         "trade_rule":"xxx"
 *     }
 * }
 *
 */
var ConfRuleH = baseapi.ConfigLoader{gexdb.ConfigTradeRule}

//ConfGoldbarH is http handler to load goldbar configure
/**
 *
 * @api {GET} /conf/goldbar Goldbar Config
 * @apiName ConfGoldbar
 * @apiGroup Conf
 *
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "config": {
 *         "goldbar_explain": "xxx",
 *         "goldbar_address": "[{\"city\":\"city1\",\"addresss\":\"address1\"}]",
 *         "goldbar_rate": "1600",
 *         "goldbar_fee": "0.05",
 *         "goldbar_tips": "warning"
 *     }
 * }
 *
 */
var ConfGoldbarH = baseapi.ConfigLoader{gexdb.ConfigGoldbarAddress, gexdb.ConfigGoldbarExplain, gexdb.ConfigGoldbarRate, gexdb.ConfigGoldbarFee, gexdb.ConfigGoldbarTips}
