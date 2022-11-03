## 名词对照表
|英文名词|中文名词|备注|
|:------|:----|:--|
|User|用户表||

## 关于用户
* 管理后台用户管理：使用<a href="#api-User-SearchUser">列出用户接口</a>
* 黄金管理列出用户：使用<a href="#api-User-SearchUser">列出用户接口</a>，传入`ret_balance=1`返回黄金信息
* 列出经纪人：使用<a href="#api-User-SearchUser">列出用户接口</a>，传入`user_role=200`列出经纪人
* 列出下属用户：使用<a href="#api-User-SearchMyUser">列出下属用户接口</a>
* 更新经纪人信息：使用<a href="#api-User-UpdateUser">更新用户接口</a>，经纪人信息保存在`external`中，字段前端自由定义


## 关于kbz登录
* kbz打开小程序之后调用小程序函数获取authcode，`clientId`需要放在配置文件中
```.js
window.xm.getAuthCode({ clientId: "kpf9cc61b5abdc43938a518c1a0cadb2" }).then(function (token) {
    console.log("token", token); // token
});
```
* 获取token之后，调用<a href="#api-User-Login">登录接口</a>，使用`kbz_token`登录系统


## 关于kbz充值
* 首先使用<a href="#api-Order-CreateTopupOrder">创建充值订单接口</a>创建支付订单，返回prepay信息
* 获取prepay信息后，调用小程序拉起支付api
```.js
var prepayID = result.order.prepay_result.prepay_id;
var orderInfo = result.order.prepay_result.order_info;
var sign = result.order.prepay_result.sign;
var singType = result.order.prepay_result.sign_type;
var tradeType = result.order.prepay_result.trade_type;
window.xm.native("startPay", {
    prepayId: prepayID,
    orderInfo: orderInfo,
    sign: sign,
    signType: singType,
    tradeType: tradeType,
}).then((res) => {
    console.log(res);
});
```
* 拉取支付返回成功后，需要轮训<a href="#api-Order-QueryOrder">订单查询接口</a>，检查是否支付完成`OrderStatusDone`
* 测试环境中可以调用<a href="#api-Order-MockPayTopupOrder">充值支付模拟接口</a>模拟支付通知

## 关于行情数据
* 进入k线图后面，需要通过<a href="#api-Market-ListKLine">列出k线图接口</a>获取历史数据，最新数据在最前面
* 然后使用<a href="#api-Market-WsMarket">Websocket行情推送接口</a>订阅k线图数据来获取最新数据
  * 返回的数据追加到历史数据最前
  * 每次推送需要检测kline.start_time与之前一次推送是否一样，如果不一样需要再次追加
* 深度图数据使用<a href="#api-Market-WsMarket">Websocket行情推送接口</a>订阅深度图数据来获取最新数据
  * 每次推送为全量推送
  * 可以与k线图同用一个websocket连接

## 关于钱包
* 通过<a href="#api-Balance-LoadMyBalance">我的钱包接口</a>获取当前用户的钱包信息，以及当日的估值
* 当用用户停在我的页面时，前端需要定时（5s)刷新我的钱包信息

## 关于交易
* 交易基本要求得有一个基础计价币，一个被交易的计量单位，本系统中，计价币为`MMK`，计量单位为`YWE`
  * 讲价币是用来支出和收入的，所以以它为视图的交易类型一共有三种，充值、提现、修改（即管理员直接修改数量）、交易（买就是支出、卖就是收入）,根据可以根据订单类型与订单方向判断、过滤数据，详情请查看<a href="#api-Order-SearchOrder">列出订单接口</a>的样例说明
  * 计量单位是用来持有和卖出的，所以以它为视图的交易类型一共有三种，提取、修改（即管理员直接修改数量）、交易（买就是持有、卖就是卖出）,根据可以根据订单类型与订单方向判断、过滤数据，详情请查看<a href="#api-Order-SearchOrder">列出订单接口</a>的样例说明
  * 本系统中所有说价格的均以MMK做为单位
  * 交易订单类型有充值、提现、交易、提取、修改等类型，对应`OrderTypeTopup, OrderTypeWithdraw, OrderTypeTrade, OrderTypeGoldbar, OrderTypeChangeYWE, OrderTypeChangeMMK`
  * 交易订单状态有可能存在进行中、部分完成（交易中，已经完成部分）、完成、部分取消（交易完成了，只交易了一部分）、取消等状态，对应`OrderStatusPending, OrderStatusPartialled, OrderStatusDone, OrderStatusPartcanceled, OrderStatusCanceled`
  * 交易如果正进行中，系统会锁住对应的钱包的交易量
* 买入卖出统一调用<a href="#api-Order-PlaceOrder">下单接口</a>，支持市价和限价两种模式
* 价格、数据、金额最大支持2位小数，即前端处理加减时，步进是0.01
* 限价单：数量和价格必传，下单之前需要检测当前用户的钱包是否有足够的余额，即买单要`quantity*price<=mmk`，卖单`quantity<=ywe`
* 市价单：只需要传入数量，买入市价单还支持传买入总价`total`，总价必须小于钱包余额`total<=mmk`
* 列出当前用户的交易订单使用<a href="#api-Order-SearchOrder">列出订单接口</a>，类型传入`OrderTypeTrade`只列出交易订单，通过`side=buy`或`side=sell`列出买卖订单
* 用户在交易界面时，正常进行中的订单有变化时需要刷新钱包信息
* 订单中字段的详细说明
  * `quantity/price` 为本次交易中用户期望的数量为价格，在市价单、充值、提现、提取等订单中，价格都为0
  * `avg_price/total_price` 为本次交易中成交的平均价格和总价格
  * `in_balance` 为本次交易收入的类型，当卖出、充值时为`MMK`，其他情况为`YWE`
  * `in_filled` 为本次交易收入的最终量
  * `out_balance` 为本次交易支出的类型，当买入、提现时为`MMK`，其他情况为`YWE`
  * `out_filled` 为本次交易支出的最终量
  * `fee_balance` 为本次交易系统扣除的手续费单位，与`in_balance`相同
  * `fee_filled` 为本次交易系统扣除的手续费
  * 不同订单类型的数据结构
    * 充值时：只有`in`，没有`out`、`fee`
    * 提现时：只有`out`，没有`in`、`fee`
    * 提取时：只有`out`，没有`in`、`fee`
    * 交易时：都有`in`、`out`、`fee`
    * 修改时：如果是增加只有`in`，如果是减少只`out`

## 关于黄金提取
* 提取申请：使用<a href="#api-Order-CreateGoldbarOrder">创建黄金订单接口</a>申请，获得提取码，申请后可以使用<a href="#api-Order-CancelGoldbarOrder">取消黄金订单接口</a>取消
* 列出申请：使用<a href="#api-Order-SearchOrder">列出订单接口</a>，类型传入`OrderTypeGoldbar`列出黄金订单，
* 提取管理：使用<a href="#api-Order-SearchOrder">列出订单接口</a>，类型传入`OrderTypeGoldbar`列出黄金订单，状态为`OrderStatusPending`
* 黄金对账单：使用<a href="#api-Order-SearchOrder">列出订单接口</a>，类型传入`OrderTypeGoldbar`和`OrderTypeChangeYWE`列出黄金订单，状态为`OrderStatusDone`
* 提取配置：使用<a href="#api-Conf-ConfGoldbar">黄金提取配置接口</a>

## 关于经纪人
* 列出我的用户使用，使用<a href="#api-User-SearchMyUser">列出我的用户接口</a>
* 列出我的用户交易记录，使用<a href="#api-User-SearchMyUserOrder">列出我的用户交易记录接口</a>，手续费使用`order`里面的`fee_filled`
* 列出我的交易费收，使用<a href="#api-User-SearchMyUserOrder">列出我的用户交易记录接口</a>，收入使用`comms`里面的`in_fee`

## 关于公告
* <a href="#api-Announce">公告接口</a>
* 后台编辑使用<a href="#api-Announce-UpdateAnnounce">更新公告接口</a>，传入`marked=100`标记为首页滚动
* 首页中滚动公告标题，使用<a href="#api-Announce-SearchAnnounce">搜索公告接口</a>，传入`marked=100`列出需要滚动的公告

## 关于规则页
* 后台系统配置里面配置规则html内容
* 点击进入规则页后，使用<a href="#api-Conf-ConfRule">列出规则配置接口</a>加载规则说明