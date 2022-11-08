## 名词对照表
|英文名词|中文名词|备注|
|:------|:----|:--|
|User|用户表||

## 关于用户
* 管理后台用户管理：使用<a href="#api-User-SearchUser">列出用户接口</a>

## 关于行情数据
* 列出行情，使用<a href="#api-Market-ListSymbol">列出交易对接口</a>
* 进入k线图后面，需要通过<a href="#api-Market-ListKLine">列出k线图接口</a>获取历史数据，最新数据在最前面
* 然后使用<a href="#api-Market-WsMarket">Websocket行情推送接口</a>订阅k线图数据来获取最新数据
  * 返回的数据追加到历史数据最前
  * 每次推送需要检测kline.start_time与之前一次推送是否一样，如果不一样需要再次追加
* 深度图数据使用<a href="#api-Market-WsMarket">Websocket行情推送接口</a>订阅深度图数据来获取最新数据
  * 每次推送为全量推送
  * 可以与k线图同用一个websocket连接
* 行情收藏
  * 列出收藏的行情<a href="#api-Market-ListFavoritesSymbol">列出收藏交易对</a>
  * 添加行情收藏<a href="#api-Market-AddFavoritesSymbol">添加收藏交易对</a>
  * 删除行情收藏<a href="#api-Market-RemoveFavoritesSymbol">删除收藏交易对</a>
  * 排序行情收藏<a href="#api-Market-SwitchFavoritesSymbol">排序收藏交易对</a>

## 关于钱包
* 钱包总览使用<a href="#api-Balance-LoadBalanceOverview">钱包总览</a>
* 列出钱包使用<a href="#api-Balance-ListBalance">列出钱包</a>
* 列出合约钱包使用<a href="#api-Balance-ListHolding">列出持仓</a>
  * 百分计算为持仓未实现盈亏除以总保证金（已使用保证金和追加保证金），`unprofits[symbol]/(holding[symbol].margin_used+holdings[symbol].margin_added)`
  * 标记价格显示，当`holdings.amout`为多仓（正数）时为`tickers[symbol].bid[0]`, 当`holdings.amout`为空仓（负数）时为`tickers[symbol].ask[0]`
  * 保证金率为持仓未实现盈亏除以总保证金加空闲保证金（已使用保证金+追加保证金+账号空闲余额），`unprofits[symbol]/(holding[symbol].margin_used+holdings[symbol].margin_added+balance.free)`，如果为正数显示绿色，负数显示红色
* 更新杠杆使用<a href="#api-Balance-ChangeHoldingLever">更新仓位杠杆</a>, 注意检查返回`code=CodeBalanceNotEnought/CodeOrderPending`提示用户余额不足或有订单未成交
* 列出钱包历史或资金历史使用<a href="#api-Balance-ListBalanceRecord">列出钱包记录</a>
* 管理员修改钱包余额<a href="#api-Balance-ChangeUserBalance">修改钱包余额</a>
* 划转钱包余额<a href="#api-Balance-TransferBalance">钱包划转</a>
* 当用用户停在我的页面时，前端需要定时（5s)刷新我的钱包信息

## 关于交易
* 每个交易对都有一个计价币（基于引用币`quote`)，一个被交易的计量单位(基于基础币`base`)，例如计价币为`USDT`，计量单位为`YWE`
  * 交易对的计价和计量都有对应的精度，小数点位数，在交易对信息接口中返回，<a href="#api-Market-ListSymbol">列出交易对接口</a>或<a href="#api-Market-LoadSymbol">获取交易对信息</a>
  * 提交订单前，前端需要根据交易对的精度、交易的币等信息检查计算下单的价格与数量，价格与数量的步进即为对应精度的最小量
  * 本系统中所有说价格的均以引用币做为单位
  * 交易订单类型有交易、触发（止盈止损）、爆仓等类型，对应`OrderTypeTrade, OrderTypeTrigger, OrderTypeBlowup`
  * 交易订单状态有可能存在进行中、部分完成（交易中，已经完成部分）、完成、部分取消（交易完成了，只交易了一部分）、取消等状态，对应`OrderStatusWaiting, OrderStatusPending, OrderStatusPartialled, OrderStatusDone, OrderStatusPartcanceled, OrderStatusCanceled`
  * 列出订单，根据可以根据订单类型与订单方向判断、过滤数据，详情请查看<a href="#api-Order-SearchOrder">列出订单接口</a>的样例说明
  * 交易如果正进行中，系统会锁住对应的钱包的交易量
* 买入卖出统一调用<a href="#api-Order-PlaceOrder">下单接口</a>，支持市价和限价两种模式
* 限价单：数量和价格必传，下单之前需要检测当前用户的钱包是否有足够的余额，即现货买单要`quantity*price<=usdt`，现货卖单`quantity<=ywe`，合约时有对应数量的反向仓位或`quantity*price/lever<=usdt`
* 市价单：只需要传入数量，买入市价单还支持传买入总价`total_price`，总价必须小于钱包余额`total_price<=usdt`，另外最少总价为最小单位的10倍数量乘以当前最新卖价
* 列出当前用户的正在交易订单使用<a href="#api-Order-SearchOrder">列出订单接口</a>，传入`type=OrderTypeTrade&status=OrderStatusPending,OrderStatusPartialled`，也可以加上`side=OrderSideBuy`或`side=OrderSideSell`只列出买卖订单
* 用户在交易界面时，正常进行中的订单有变化时需要刷新钱包信息
* 订单中字段的详细说明
  * `quantity/price` 为本次交易中用户期望的数量为价格，在市价单时价格都为0
  * `avg_price/total_price` 为本次交易中成交的平均价格和总价格
  * `in_balance` 为本次交易收入的类型，当卖出时为`USDT`，其他情况为`YWE`
  * `in_filled` 为本次交易收入的最终量
  * `out_balance` 为本次交易支出的类型，当买入时为`USDT`，其他情况为`YWE`
  * `out_filled` 为本次交易支出的最终量
  * `fee_balance` 为本次交易系统扣除的手续费单位，现货时与`in_balance`相同，合约时固定为`USDT`相同
  * `fee_filled` 为本次交易系统扣除的手续费

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