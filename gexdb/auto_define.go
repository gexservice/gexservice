//auto gen func by autogen
package gexdb

/**
 * @apiDefine BalanceUpdate
 */
/**
 * @apiDefine BalanceObject
 * @apiSuccess (Balance) {Int64} Balance.tid the primary key
 * @apiSuccess (Balance) {Int64} Balance.user_id the balance user id
 * @apiSuccess (Balance) {BalanceArea} Balance.area the balance area, all suported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiSuccess (Balance) {String} Balance.asset the balance asset key
 * @apiSuccess (Balance) {Decimal} Balance.free the balance free amount
 * @apiSuccess (Balance) {Decimal} Balance.locked the balance locked amount
 * @apiSuccess (Balance) {Decimal} Balance.margin the balance margin value
 * @apiSuccess (Balance) {Time} Balance.update_time the balance last update time
 * @apiSuccess (Balance) {Time} Balance.create_time the balance create time
 * @apiSuccess (Balance) {BalanceStatus} Balance.status the balance status, all suported is <a href="#metadata-Balance">BalanceStatusAll</a>
 */

/**
 * @apiDefine BalanceHistoryUpdate
 */
/**
 * @apiDefine BalanceHistoryObject
 * @apiSuccess (BalanceHistory) {Int64} BalanceHistory.tid the primary key
 * @apiSuccess (BalanceHistory) {Int64} BalanceHistory.user_id the balance user id
 * @apiSuccess (BalanceHistory) {String} BalanceHistory.asset the balance asset key
 * @apiSuccess (BalanceHistory) {Decimal} BalanceHistory.valuation the balance valuation
 * @apiSuccess (BalanceHistory) {Time} BalanceHistory.update_time the balance record update time
 * @apiSuccess (BalanceHistory) {Time} BalanceHistory.create_time the balance record create time,is daily zero time
 * @apiSuccess (BalanceHistory) {BalanceHistoryStatus} BalanceHistory.status the balance record status, all suported is <a href="#metadata-BalanceHistory">BalanceHistoryStatusAll</a>
 */

/**
 * @apiDefine BalanceRecordUpdate
 */
/**
 * @apiDefine BalanceRecordObject
 * @apiSuccess (BalanceRecord) {Int64} BalanceRecord.tid the primary key
 * @apiSuccess (BalanceRecord) {Int64} BalanceRecord.creator the balance creator
 * @apiSuccess (BalanceRecord) {Int64} BalanceRecord.balance_id the balance id
 * @apiSuccess (BalanceRecord) {BalanceRecordType} BalanceRecord.type the balance record type, all suported is <a href="#metadata-BalanceRecord">BalanceRecordTypeAll</a>
 * @apiSuccess (BalanceRecord) {Int} BalanceRecord.target the balance target type
 * @apiSuccess (BalanceRecord) {Decimal} BalanceRecord.changed the balance change value
 * @apiSuccess (BalanceRecord) {Object} BalanceRecord.transaction the balance record transaction info
 * @apiSuccess (BalanceRecord) {Time} BalanceRecord.update_time the balance last update time
 * @apiSuccess (BalanceRecord) {Time} BalanceRecord.create_time the balance create time
 * @apiSuccess (BalanceRecord) {BalanceRecordStatus} BalanceRecord.status the balance status, all suported is <a href="#metadata-BalanceRecord">BalanceRecordStatusAll</a>
 */

/**
 * @apiDefine HoldingUpdate
 */
/**
 * @apiDefine HoldingObject
 * @apiSuccess (Holding) {Int64} Holding.tid the primary key
 * @apiSuccess (Holding) {Int64} Holding.user_id the holding user id
 * @apiSuccess (Holding) {String} Holding.symbol the holding symbol
 * @apiSuccess (Holding) {Decimal} Holding.amount the holding amount
 * @apiSuccess (Holding) {Decimal} Holding.open the holding open price
 * @apiSuccess (Holding) {Decimal} Holding.blowup the holding blowup price
 * @apiSuccess (Holding) {Int} Holding.lever the holding lever
 * @apiSuccess (Holding) {Decimal} Holding.margin_used the holding margin used
 * @apiSuccess (Holding) {Decimal} Holding.margin_added the holding margin added
 * @apiSuccess (Holding) {Time} Holding.update_time the holding last update time
 * @apiSuccess (Holding) {Time} Holding.create_time the holding create time
 * @apiSuccess (Holding) {HoldingStatus} Holding.status the holding status, all suported is <a href="#metadata-Holding">HoldingStatusAll</a>
 */

/**
 * @apiDefine KLineUpdate
 */
/**
 * @apiDefine KLineObject
 * @apiSuccess (KLine) {Int64} KLine.tid the primay key
 * @apiSuccess (KLine) {String} KLine.symbol the kline symbol
 * @apiSuccess (KLine) {String} KLine.interv the kline interval key
 * @apiSuccess (KLine) {Decimal} KLine.amount the kline amount
 * @apiSuccess (KLine) {Int64} KLine.count the kline count
 * @apiSuccess (KLine) {Decimal} KLine.open the kline open price
 * @apiSuccess (KLine) {Decimal} KLine.close the kline close price
 * @apiSuccess (KLine) {Decimal} KLine.low the kline low price
 * @apiSuccess (KLine) {Decimal} KLine.high the kline high price
 * @apiSuccess (KLine) {Decimal} KLine.volume the kline volume price
 * @apiSuccess (KLine) {Time} KLine.start_time the kline start time
 * @apiSuccess (KLine) {Time} KLine.update_time the kline update time
 */

/**
 * @apiDefine MessageUpdate
 * @apiParam (Message) {MessageType} Message.type only required when add, the message type, all suported is <a href="#metadata-Message">MessageTypeAll</a>
 * @apiParam (Message) {Object} Message.title only required when add, the message title
 * @apiParam (Message) {Object} Message.content only required when add, the message content
 * @apiParam (Message) {Int64} Message.to_user_id only required when add,
 */
/**
 * @apiDefine MessageObject
 * @apiSuccess (Message) {Int64} Message.tid the primary key
 * @apiSuccess (Message) {MessageType} Message.type the message type, all suported is <a href="#metadata-Message">MessageTypeAll</a>
 * @apiSuccess (Message) {Object} Message.title the message title
 * @apiSuccess (Message) {Object} Message.content the message content
 * @apiSuccess (Message) {Int64} Message.to_user_id
 * @apiSuccess (Message) {Time} Message.update_time the message update time
 * @apiSuccess (Message) {Time} Message.create_time the message create time
 * @apiSuccess (Message) {MessageStatus} Message.status the message status, all suported is <a href="#metadata-Message">MessageStatusAll</a>
 */

/**
 * @apiDefine OrderUpdate
 * @apiParam (Order) {Int64} [Order.tid] the primary key
 * @apiParam (Order) {Decimal} [Order.quantity] the order expected quantity
 * @apiParam (Order) {Decimal} [Order.price] the order expected price
 * @apiParam (Order) {OrderTriggerType} [Order.trigger_type] the order trigger type, all suported is <a href="#metadata-Order">OrderTriggerTypeAll</a>
 * @apiParam (Order) {Decimal} [Order.trigger_price] the order trigger price
 * @apiParam (Order) {Decimal} [Order.total_price] the order filled total price
 * @apiParam (Order) {OrderStatus} [Order.status] the order status, all suported is <a href="#metadata-Order">OrderStatusAll</a>
 */
/**
 * @apiDefine OrderObject
 * @apiSuccess (Order) {Int64} Order.tid the primary key
 * @apiSuccess (Order) {String} Order.order_id the order string id
 * @apiSuccess (Order) {OrderType} Order.type the order type, all suported is <a href="#metadata-Order">OrderTypeAll</a>
 * @apiSuccess (Order) {Int64} Order.user_id the order user id
 * @apiSuccess (Order) {Int64} Order.creator the order creator user id
 * @apiSuccess (Order) {OrderArea} Order.area the order area, all suported is <a href="#metadata-Order">OrderAreaAll</a>
 * @apiSuccess (Order) {String} Order.symbol the order symbol
 * @apiSuccess (Order) {OrderSide} Order.side the order side, all suported is <a href="#metadata-Order">OrderSideAll</a>
 * @apiSuccess (Order) {Decimal} Order.quantity the order expected quantity
 * @apiSuccess (Order) {Decimal} Order.filled the order filled quantity
 * @apiSuccess (Order) {Decimal} Order.price the order expected price
 * @apiSuccess (Order) {OrderTriggerType} Order.trigger_type the order trigger type, all suported is <a href="#metadata-Order">OrderTriggerTypeAll</a>
 * @apiSuccess (Order) {Decimal} Order.trigger_price the order trigger price
 * @apiSuccess (Order) {Decimal} Order.avg_price the order filled avg price
 * @apiSuccess (Order) {Decimal} Order.total_price the order filled total price
 * @apiSuccess (Order) {Decimal} Order.holding the order holding
 * @apiSuccess (Order) {Decimal} Order.profit the order profit
 * @apiSuccess (Order) {Decimal} Order.owned the order owned count
 * @apiSuccess (Order) {Decimal} Order.unhedged the order owned is unbalanced
 * @apiSuccess (Order) {String} Order.in_balance the in balance asset key
 * @apiSuccess (Order) {Decimal} Order.in_filled the in balance filled amount
 * @apiSuccess (Order) {String} Order.out_balance the out balance asset key
 * @apiSuccess (Order) {Decimal} Order.out_filled the out balance filled amount
 * @apiSuccess (Order) {String} Order.fee_balance the fee balance asset key
 * @apiSuccess (Order) {Decimal} Order.fee_filled the fee amount
 * @apiSuccess (Order) {Decimal} Order.fee_rate the order fee rate
 * @apiSuccess (Order) {OrderTransaction} Order.transaction the order transaction info
 * @apiSuccess (Order) {Int} Order.fee_settled_status the order transaction detail
 * @apiSuccess (Order) {Time} Order.fee_settled_next the fee settled time
 * @apiSuccess (Order) {Time} Order.update_time the order update time
 * @apiSuccess (Order) {Time} Order.create_time the order create time
 * @apiSuccess (Order) {OrderStatus} Order.status the order status, all suported is <a href="#metadata-Order">OrderStatusAll</a>
 */

/**
 * @apiDefine OrderCommUpdate
 */
/**
 * @apiDefine OrderCommObject
 * @apiSuccess (OrderComm) {Int64} OrderComm.tid the primary key
 * @apiSuccess (OrderComm) {Int64} OrderComm.order_id the order id
 * @apiSuccess (OrderComm) {Int64} OrderComm.user_id the user id
 * @apiSuccess (OrderComm) {OrderCommType} OrderComm.type the comm type, all suported is <a href="#metadata-OrderComm">OrderCommTypeAll</a>
 * @apiSuccess (OrderComm) {String} OrderComm.in_balance the in balance asset key
 * @apiSuccess (OrderComm) {Decimal} OrderComm.in_fee the in balance fee
 * @apiSuccess (OrderComm) {Time} OrderComm.update_time
 * @apiSuccess (OrderComm) {Time} OrderComm.create_time the comm create time
 * @apiSuccess (OrderComm) {OrderCommStatus} OrderComm.status the comm status, all suported is <a href="#metadata-OrderComm">OrderCommStatusAll</a>
 */

/**
 * @apiDefine UserUpdate
 * @apiParam (User) {UserRole} [User.role] ther user role, all suported is <a href="#metadata-User">UserRoleAll</a>
 * @apiParam (User) {StringPtr} [User.name] the user name
 * @apiParam (User) {StringPtr} [User.account] the user account to login
 * @apiParam (User) {StringPtr} [User.phone] the user phone number to login
 * @apiParam (User) {StringPtr} [User.email] the user email
 * @apiParam (User) {StringPtr} [User.password] the user password to login
 * @apiParam (User) {StringPtr} [User.trade_pass] the user trade password
 * @apiParam (User) {StringPtr} [User.image] the user image
 * @apiParam (User) {Object} [User.external] the user external info
 * @apiParam (User) {UserStatus} [User.status] the user status, all suported is <a href="#metadata-User">UserStatusAll</a>
 */
/**
 * @apiDefine UserObject
 * @apiSuccess (User) {Int64} User.tid the primary key
 * @apiSuccess (User) {UserType} User.type the user type, all suported is <a href="#metadata-User">UserTypeAll</a>
 * @apiSuccess (User) {UserRole} User.role ther user role, all suported is <a href="#metadata-User">UserRoleAll</a>
 * @apiSuccess (User) {StringPtr} User.name the user name
 * @apiSuccess (User) {StringPtr} User.account the user account to login
 * @apiSuccess (User) {StringPtr} User.phone the user phone number to login
 * @apiSuccess (User) {StringPtr} User.email the user email
 * @apiSuccess (User) {StringPtr} User.password the user password to login
 * @apiSuccess (User) {StringPtr} User.trade_pass the user trade password
 * @apiSuccess (User) {StringPtr} User.image the user image
 * @apiSuccess (User) {Object} User.fee the user fee
 * @apiSuccess (User) {Object} User.external the user external info
 * @apiSuccess (User) {UserFavorites} User.favorites the user favorites
 * @apiSuccess (User) {Object} User.config the user config
 * @apiSuccess (User) {Time} User.update_time the last updat time
 * @apiSuccess (User) {Time} User.create_time the craete time
 * @apiSuccess (User) {UserStatus} User.status the user status, all suported is <a href="#metadata-User">UserStatusAll</a>
 */

/**
 * @apiDefine WalletUpdate
 */
/**
 * @apiDefine WalletObject
 * @apiSuccess (Wallet) {Int64} Wallet.tid the wallet primary key
 * @apiSuccess (Wallet) {Int64} Wallet.user_id the wallet user id
 * @apiSuccess (Wallet) {WalletMethod} Wallet.method the wallet type, all suported is <a href="#metadata-Wallet">WalletMethodAll</a>
 * @apiSuccess (Wallet) {String} Wallet.address the wallet address
 * @apiSuccess (Wallet) {Time} Wallet.update_time the wallet update time
 * @apiSuccess (Wallet) {Time} Wallet.create_time the wallet create time
 * @apiSuccess (Wallet) {WalletStatus} Wallet.status the wallet status, all suported is <a href="#metadata-Wallet">WalletStatusAll</a>
 */

/**
 * @apiDefine WithdrawUpdate
 * @apiParam (Withdraw) {WithdrawMethod} Withdraw.method only required when add, the withdraw metod, all suported is <a href="#metadata-Withdraw">WithdrawMethodAll</a>
 * @apiParam (Withdraw) {String} Withdraw.asset only required when add, the withdraw asset
 * @apiParam (Withdraw) {Decimal} Withdraw.quantity only required when add, the withdraw order quantity
 * @apiParam (Withdraw) {String} Withdraw.receiver only required when add, the widhdraw receiver
 */
/**
 * @apiDefine WithdrawObject
 * @apiSuccess (Withdraw) {Int64} Withdraw.tid the primary key
 * @apiSuccess (Withdraw) {String} Withdraw.order_id the withdraw order string id
 * @apiSuccess (Withdraw) {WithdrawType} Withdraw.type the withdraw order type, all suported is <a href="#metadata-Withdraw">WithdrawTypeAll</a>
 * @apiSuccess (Withdraw) {Int64} Withdraw.user_id the withdraw order user id
 * @apiSuccess (Withdraw) {Int64} Withdraw.creator the withdraw order creator user id
 * @apiSuccess (Withdraw) {WithdrawMethod} Withdraw.method the withdraw metod, all suported is <a href="#metadata-Withdraw">WithdrawMethodAll</a>
 * @apiSuccess (Withdraw) {String} Withdraw.asset the withdraw asset
 * @apiSuccess (Withdraw) {Decimal} Withdraw.quantity the withdraw order quantity
 * @apiSuccess (Withdraw) {StringPtr} Withdraw.sender
 * @apiSuccess (Withdraw) {String} Withdraw.receiver the widhdraw receiver
 * @apiSuccess (Withdraw) {Int} Withdraw.processed the withdraw if processed
 * @apiSuccess (Withdraw) {Object} Withdraw.result the withdraw order transaction info
 * @apiSuccess (Withdraw) {Time} Withdraw.update_time the withdraw order update time
 * @apiSuccess (Withdraw) {Time} Withdraw.create_time the withdraw order create time
 * @apiSuccess (Withdraw) {WithdrawStatus} Withdraw.status the withdraw order status, all suported is <a href="#metadata-Withdraw">WithdrawStatusAll</a>
 */
