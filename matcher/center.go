package matcher

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/debug"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

/**
 * @apiDefine SymbolInfoObject
 * @apiSuccess (SymbolInfo) {Int64} SymbolInfo.precision_quantity the symbol quantity precision
 * @apiSuccess (SymbolInfo) {Int64} SymbolInfo.precision_price the symbol price precision
 * @apiSuccess (SymbolInfo) {String} SymbolInfo.type the symbol type, all suported is spot/futures
 * @apiSuccess (SymbolInfo) {String} SymbolInfo.symbol the symbol key
 * @apiSuccess (SymbolInfo) {String} SymbolInfo.base the symbol base asset(coin)
 * @apiSuccess (SymbolInfo) {String} SymbolInfo.quote the symbol quote asset(coin)
 * @apiSuccess (SymbolInfo) {Decimal} SymbolInfo.fee the symbol trade fee
 * @apiSuccess (SymbolInfo) {Decimal} SymbolInfo.margin_max the symbol holding margin max to blowup
 * @apiSuccess (SymbolInfo) {Decimal} SymbolInfo.margin_add the symbol holding margin add neary blowup
 */
type SymbolInfo struct {
	PrecisionQuantity int32           `json:"precision_quantity"`
	PrecisionPrice    int32           `json:"precision_price"`
	Type              string          `json:"type"`
	Symbol            string          `json:"symbol"`
	Base              string          `json:"base"`
	Quote             string          `json:"quote"`
	Fee               decimal.Decimal `json:"fee"`
	MarginMax         decimal.Decimal `json:"margin_max"`
	MarginAdd         decimal.Decimal `json:"margin_add"`
}

func (s *SymbolInfo) String() string {
	return s.Symbol
}

type MatcherFeeCache struct {
	Default   map[string]decimal.Decimal
	cacheMax  int
	cacheFee  map[int64]xmap.M
	cacheLast time.Time
	cacheLock sync.RWMutex
}

func NewMatcherFeeCache(cacheMax int) (cache *MatcherFeeCache) {
	cache = &MatcherFeeCache{
		Default:   map[string]decimal.Decimal{},
		cacheMax:  cacheMax,
		cacheFee:  map[int64]xmap.M{},
		cacheLock: sync.RWMutex{},
	}
	return
}

func (m *MatcherFeeCache) readFee(config xmap.M, symbol string) (fee decimal.Decimal) {
	if v, ok := config[symbol]; ok {
		fee = decimal.NewFromFloat(v.(float64))
	} else if v, ok := config["*"]; ok {
		fee = decimal.NewFromFloat(v.(float64))
	} else {
		fee = m.Default[symbol]
	}
	return
}

func (m *MatcherFeeCache) loadCache(userID int64, symbol string) (fee decimal.Decimal, ok bool) {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	if len(m.cacheFee) > m.cacheMax || time.Since(m.cacheLast) > 24*time.Hour {
		m.cacheFee = map[int64]xmap.M{}
		m.cacheLast = time.Now()
	}
	config, ok := m.cacheFee[userID]
	if ok {
		fee = m.readFee(config, symbol)
	}
	return
}

func (m *MatcherFeeCache) updateCache(userID int64, cache xmap.M) {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	m.cacheFee[userID] = cache
}

func (m *MatcherFeeCache) LoadFee(ctx context.Context, userID int64, symbol string) (fee decimal.Decimal, err error) {
	fee, ok := m.loadCache(userID, symbol)
	if ok {
		return
	}
	config, err := gexdb.LoadUserFee(ctx, userID)
	if err != nil {
		return
	}
	m.updateCache(userID, config)
	fee = m.readFee(config, symbol)
	return
}

type MatcherBalancePreparer struct {
	cacheMax     int
	cacheBalance map[string]bool
	cacheLast    time.Time
	cacheLock    sync.RWMutex
}

func NewMatcherBalancePreparer(cacheMax int) (preparer *MatcherBalancePreparer) {
	preparer = &MatcherBalancePreparer{
		cacheMax:     cacheMax,
		cacheBalance: map[string]bool{},
		cacheLock:    sync.RWMutex{},
	}
	return
}

func (m *MatcherBalancePreparer) touchBalance(key string) {
	if len(m.cacheBalance) > m.cacheMax || time.Since(m.cacheLast) > 24*time.Hour {
		m.cacheBalance = map[string]bool{}
		m.cacheLast = time.Now()
	}
	m.cacheBalance[key] = true
}

func (m *MatcherBalancePreparer) PrepareSpotMatcher(ctx context.Context, matcher *SpotMatcher, userID int64) (err error) {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	keyQuote := fmt.Sprintf("%v-%v-%v", matcher.Area, matcher.Quote, userID)
	if m.cacheBalance[keyQuote] {
		return
	}
	_, err = gexdb.TouchBalance(ctx, matcher.Area, []string{matcher.Base, matcher.Quote}, userID)
	if err != nil {
		return
	}
	m.touchBalance(keyQuote)
	return
}

func (m *MatcherBalancePreparer) PrepareFuturesMatcher(ctx context.Context, matcher *FuturesMatcher, userID int64) (err error) {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	keyQuote := fmt.Sprintf("%v-%v-%v", matcher.Area, matcher.Quote, userID)
	if m.cacheBalance[keyQuote] {
		return
	}
	_, err = gexdb.TouchBalance(ctx, matcher.Area, []string{matcher.Quote}, userID)
	if err != nil {
		return
	}
	_, err = gexdb.TouchHolding(ctx, []string{matcher.Symbol}, userID)
	if err != nil {
		return
	}
	m.touchBalance(keyQuote)
	return
}

type MatcherCenter struct {
	Symbols      map[string]*SymbolInfo
	FeeCache     *MatcherFeeCache
	Preparer     *MatcherBalancePreparer
	TriggerDelay time.Duration
	matcherAll   map[string]Matcher
	matcherLock  sync.RWMutex
	monitorAll   map[string]map[string]MatcherMonitor
	monitorLock  sync.RWMutex
	eventRun     int
	eventQueue   chan *MatcherEvent
	exiter       chan int
	waiter       sync.WaitGroup
}

func NewMatcherCenter(eventRun, eventMax, cacheMax int) (center *MatcherCenter) {
	center = &MatcherCenter{
		TriggerDelay: time.Second,
		Symbols:      map[string]*SymbolInfo{},
		FeeCache:     NewMatcherFeeCache(cacheMax),
		Preparer:     NewMatcherBalancePreparer(cacheMax),
		matcherAll:   map[string]Matcher{},
		matcherLock:  sync.RWMutex{},
		monitorAll:   map[string]map[string]MatcherMonitor{},
		monitorLock:  sync.RWMutex{},
		eventQueue:   make(chan *MatcherEvent, eventMax),
		eventRun:     eventRun,
		exiter:       make(chan int, 1),
		waiter:       sync.WaitGroup{},
	}
	return
}

func BootstrapMatcherCenterByConfig(ctx context.Context, config *xprop.Config) (center *MatcherCenter, err error) {
	eventRun := config.IntDef(1, "matcher/matcher_event_run")
	eventMax := config.IntDef(4096, "matcher/matcher_event_max")
	cacheMax := config.IntDef(10000, "matcher/balance_cache_max")
	center = NewMatcherCenter(eventRun, eventMax, cacheMax)
	for _, sec := range config.Seces {
		if !strings.HasPrefix(sec, "matcher.") {
			continue
		}
		if config.StrDef("0", sec+"/on") != "1" {
			continue
		}
		info := &SymbolInfo{
			PrecisionQuantity: 2,
			PrecisionPrice:    2,
			Fee:               decimal.NewFromFloat(0.002),
			MarginMax:         decimal.NewFromFloat(0.99),
			MarginAdd:         decimal.NewFromFloat(0.01),
		}
		err = config.ValidFormat(
			strings.ReplaceAll(`
				_S/precision_quantity,o|i,r:0;
				_S/precision_price,o|i,r:0;
				_S/type,r|s,l:0;
				_S/symbol,r|s,l:0;
				_S/base,r|s,l:0;
				_S/quote,r|s,l:0;
				_S/fee,0|f,r:-1~1;
				_S/margin_max,o|f,r:0~1;
				_S/margin_add,o|f,r:0~1;
			`, "_S", sec),
			&info.PrecisionQuantity, &info.PrecisionPrice, &info.Type, &info.Symbol, &info.Base, &info.Quote, &info.Fee, &info.MarginMax, &info.MarginAdd,
		)
		if err != nil {
			break
		}
		if info.Type == "spot" {
			spot := NewSpotMatcher(info.Symbol, info.Base, info.Quote, center)
			spot.Fee = center.FeeCache
			spot.PrecisionPrice = info.PrecisionPrice
			spot.PrecisionQuantity = info.PrecisionQuantity
			spot.PrepareProcess = center.Preparer.PrepareSpotMatcher
			_, err = spot.Bootstrap(ctx)
			if err != nil {
				xlog.Errorf("Bootstrap init spot matcher by symbol %v fail with \n%v", info.Symbol, ErrStack(err))
				break
			}
			center.FeeCache.Default[info.Symbol] = info.Fee
			center.AddMatcher(info, spot)
			xlog.Infof("Bootstrap register spot matcher by symbol %v", info.Symbol)
		} else if strings.HasPrefix(info.Symbol, "futures.") {
			futures := NewFuturesMatcher(info.Symbol, info.Quote, center)
			futures.Fee = center.FeeCache
			futures.PrecisionPrice = info.PrecisionPrice
			futures.PrecisionQuantity = info.PrecisionQuantity
			futures.MarginMax = info.MarginMax
			futures.MarginAdd = info.MarginAdd
			futures.PrepareProcess = center.Preparer.PrepareFuturesMatcher
			_, err = futures.Bootstrap(ctx)
			if err != nil {
				xlog.Errorf("Bootstrap init futures matcher by symbol %v fail with \n%v", info.Symbol, ErrStack(err))
				break
			}
			center.FeeCache.Default[info.Symbol] = info.Fee
			center.AddMatcher(info, futures)
			xlog.Infof("Bootstrap register futures matcher by symbol %v", info.Symbol)
		} else {
			err = fmt.Errorf("symbol %v is not supported, it must be started with spot. or futures. ", info.Symbol)
			break
		}
	}
	return
}

func (m *MatcherCenter) Start() {
	for i := 0; i < m.eventRun; i++ {
		m.waiter.Add(1)
		go m.loopMatcherEvent()
	}
	m.waiter.Add(1)
	go m.loopTriggerOrder(m.TriggerDelay)
}

func (m *MatcherCenter) Stop() {
	for i := 0; i < m.eventRun; i++ {
		m.exiter <- 0
	}
	m.exiter <- 0
	m.waiter.Wait()
}

func (m *MatcherCenter) AddMatcher(symbol *SymbolInfo, matcher Matcher) {
	m.matcherLock.Lock()
	defer m.matcherLock.Unlock()
	if m.matcherAll[symbol.Symbol] != nil {
		panic(fmt.Sprintf("matcher %v exists", symbol.Symbol))
	}
	m.Symbols[symbol.Symbol] = symbol
	m.matcherAll[symbol.Symbol] = matcher
}

func (m *MatcherCenter) FindMatcher(symbol string) (matcher Matcher) {
	m.matcherLock.RLock()
	defer m.matcherLock.RUnlock()
	matcher = m.matcherAll[symbol]
	return
}

func (m *MatcherCenter) AddMonitor(symbol string, monitor MatcherMonitor) {
	m.matcherLock.Lock()
	defer m.matcherLock.Unlock()
	key := fmt.Sprintf("%p", monitor)
	if m.monitorAll[symbol] == nil {
		m.monitorAll[symbol] = map[string]MatcherMonitor{}
	}
	m.monitorAll[symbol][key] = monitor
}

func (m *MatcherCenter) RemoveMonitor(symbol string, monitor MatcherMonitor) {
	m.matcherLock.Lock()
	defer m.matcherLock.Unlock()
	key := fmt.Sprintf("%p", monitor)
	if m.monitorAll[symbol] != nil {
		delete(m.monitorAll[symbol], key)
	}
}

func (m *MatcherCenter) OnMatched(ctx context.Context, event *MatcherEvent) {
	select {
	case m.eventQueue <- event:
	default:
		xlog.Warnf("MatcherCenter matcher event queue is full, skip one for %v", event.Symbol)
	}
}

func (m *MatcherCenter) loopMatcherEvent() {
	defer m.waiter.Done()
	xlog.Infof("MatcherCenter matcher event running is starting")
	running := true
	for running {
		select {
		case <-m.exiter:
			running = false
		case event := <-m.eventQueue:
			m.procMatcherEvent(event)
		}
	}
	xlog.Infof("MatcherCenter matcher event running is stopped")
}

func (m *MatcherCenter) procMatcherEvent(event *MatcherEvent) {
	m.monitorLock.RLock()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("MatcherCenter proc matcher event is panic with %v, call stack is \n%v", rerr, debug.CallStatck())
		}
		cancel()
		m.monitorLock.RUnlock()
	}()
	for _, monitor := range m.monitorAll[event.Symbol] {
		monitor.OnMatched(ctx, event)
	}
	for _, monitor := range m.monitorAll["*"] {
		monitor.OnMatched(ctx, event)
	}
}

func (m *MatcherCenter) loopTriggerOrder(delay time.Duration) {
	defer m.waiter.Done()
	ticker := time.NewTicker(delay)
	running := true
	xlog.Infof("MatcherCenter order trigger is starting by %v ticker", delay)
	for running {
		select {
		case <-m.exiter:
			running = false
		case <-ticker.C:
			m.procTriggerOrder()
		}
	}
	xlog.Infof("MatcherCenter order trigger is stopped")
}

func (m *MatcherCenter) procTriggerOrder() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("MatcherCenter proc trigger order is panic with %v, call stack is \n%v", rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		cancel()
	}()
	for _, symbol := range m.Symbols {
		m.procTriggerSybmolOrder(ctx, symbol.Symbol)
	}
	return
}

func (m *MatcherCenter) procTriggerSybmolOrder(ctx context.Context, symbol string) {
	matcher := m.FindMatcher(symbol)
	if matcher == nil {
		err := fmt.Errorf("symbol %v is not supported", symbol)
		xlog.Warnf("MatcherCenter trigger %v order fail with %v", symbol, err)
		return
	}
	depth := matcher.Depth(1)
	if depth == nil || (len(depth.Asks) < 1 && len(depth.Bids) < 1) {
		// xlog.Warnf("MatcherCenter trigger %v order is skipped for not depth", symbol)
		return
	}
	var ask, bid decimal.Decimal
	if len(depth.Asks) > 0 {
		ask = depth.Asks[0][0]
	}
	if len(depth.Bids) > 0 {
		bid = depth.Bids[0][0]
	}
	orders, err := gexdb.ListOrderForTrigger(ctx, symbol, ask, bid)
	if err != nil {
		xlog.Warnf("MatcherCenter list %v trigger order fail with %v", symbol, err)
		return
	}
	if len(orders) > 0 {
		xlog.Infof("MatcherCenter found %v %v trigger order by ask:%v,bid:%v to apply", len(orders), symbol, ask, bid)
	}
	for _, args := range orders {
		canApply := true
		if strings.HasPrefix(args.Symbol, "futures.") { //check futures if close only
			holding, xerr := gexdb.FindHoldlingBySymbol(ctx, args.UserID, args.Symbol)
			if xerr != nil {
				xlog.Errorf("MatcherCenter find holding by %v,%v fail with %v", args.UserID, args.Symbol, err)
				canApply = false
			} else {
				toChange := args.Quantity
				if args.Side == gexdb.OrderSideSell {
					toChange = decimal.Zero.Sub(args.Quantity)
				}
				willResult := holding.Amount.Add(toChange)
				canApply = willResult.Sign() == 0 || willResult.Sign() == holding.Amount.Sign() //close only
			}
		}
		if canApply {
			_, xerr := matcher.ProcessOrder(ctx, args)
			if xerr == nil {
				xlog.Infof("MatcherCenter apply %v trigger order %v is success, args is %v", symbol, args.TID, converter.JSON(args))
			} else {
				xlog.Warnf("MatcherCenter apply %v trigger order fail with %v, args is %v", symbol, xerr, converter.JSON(args))
			}
		}
		updated, xerr := gexdb.CancelTriggerOrder(ctx, args.UserID, args.Symbol, 0)
		if xerr != nil {
			xlog.Errorf("MatcherCenter cancel trigger order by %v,%v fail with %v", args.UserID, args.Symbol, err)
		}
		if updated > 0 {
			xlog.Infof("MatcherCenter cancel %v,%v other same symbol trigger succss", args.UserID, args.Symbol)
		}
	}
}

func (m *MatcherCenter) ProcessCancel(ctx context.Context, userID int64, symbol string, orderID string) (order *gexdb.Order, err error) {
	matcher := m.FindMatcher(symbol)
	if matcher == nil {
		err = fmt.Errorf("symbol %v is not supported", symbol)
		return
	}
	order, err = matcher.ProcessCancel(ctx, userID, orderID)
	return
}

func (m *MatcherCenter) ProcessMarket(ctx context.Context, userID int64, symbol string, side gexdb.OrderSide, total, quantity decimal.Decimal) (order *gexdb.Order, err error) {
	matcher := m.FindMatcher(symbol)
	if matcher == nil {
		err = fmt.Errorf("symbol %v is not supported", symbol)
		return
	}
	order, err = matcher.ProcessMarket(ctx, userID, side, total, quantity)
	return
}

func (m *MatcherCenter) ProcessLimit(ctx context.Context, userID int64, symbol string, side gexdb.OrderSide, quantity, price decimal.Decimal) (order *gexdb.Order, err error) {
	matcher := m.FindMatcher(symbol)
	if matcher == nil {
		err = fmt.Errorf("symbol %v is not supported", symbol)
		return
	}
	order, err = matcher.ProcessLimit(ctx, userID, side, quantity, price)
	return
}

func (m *MatcherCenter) ProcessOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	if args.Type != gexdb.OrderTypeTrade && args.Type != gexdb.OrderTypeTrigger {
		err = fmt.Errorf("process type must by %d or %d", gexdb.OrderTypeTrade, gexdb.OrderTypeTrigger)
		err = NewErrMatcher(err, "[ProcessOrder] args invalid")
		return
	}
	matcher := m.FindMatcher(args.Symbol)
	if matcher == nil {
		err = fmt.Errorf("symbol %v is not supported", args.Symbol)
		return
	}
	if args.TID < 1 && args.Type == gexdb.OrderTypeTrigger {
		if args.UserID <= 0 || args.Quantity.Sign() <= 0 || args.TriggerPrice.Sign() <= 0 {
			err = fmt.Errorf("process trigger userID/quantity/trigger_price is required or too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.TriggerType != gexdb.OrderTriggerTypeStopProfit && args.TriggerType != gexdb.OrderTriggerTypeStopLoss {
			err = fmt.Errorf("process trigger type must by %v or %v", gexdb.OrderTriggerTypeStopProfit, gexdb.OrderTriggerTypeStopLoss)
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		order = &gexdb.Order{
			UserID:       args.UserID,
			Creator:      args.Creator,
			Type:         gexdb.OrderTypeTrigger,
			OrderID:      gexdb.NewOrderID(),
			Symbol:       args.Symbol,
			Side:         args.Side,
			Quantity:     args.Quantity,
			Price:        args.Price,
			TriggerType:  args.TriggerType,
			TriggerPrice: args.TriggerPrice,
			Status:       gexdb.OrderStatusWaiting,
		}
		err = gexdb.AddOrder(ctx, order)
		return
	}
	order, err = matcher.ProcessOrder(ctx, args)
	return
}

func (m *MatcherCenter) ChangeLever(ctx context.Context, userID int64, symbol string, lever int) (err error) {
	matcher := m.FindMatcher(symbol)
	if matcher == nil {
		err = fmt.Errorf("symbol %v is not supported", symbol)
		return
	}
	err = matcher.ChangeLever(ctx, userID, lever)
	return
}
