package maker

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/centny/orderbook"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xtime"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
)

func init() {
	rand.Seed(xtime.Now())
}

var Verbose = false

var makerAll = map[string]*Maker{}
var makerLock = sync.RWMutex{}

func Bootstrap(ctx context.Context) (err error) {
	var config *Config
	for _, symbol := range matcher.Shared.Symbols {
		config, err = LoadConfig(ctx, symbol.Symbol)
		if err != nil && err != pgx.ErrNoRows {
			break
		}
		if err == pgx.ErrNoRows || config.ON < 1 {
			err = nil
			continue
		}
		err = Start(ctx, config.Symbol)
		if err != nil {
			xlog.Errorf("Bootstrap start %v maker fail with %v", symbol, err)
			break
		}
		xlog.Infof("Bootstrap start %v maker is success", symbol)
	}
	return
}

func Start(ctx context.Context, symbol string) (err error) {
	makerLock.Lock()
	defer makerLock.Unlock()
	_, ok := makerAll[symbol]
	if ok {
		err = fmt.Errorf("running")
		return
	}
	key := fmt.Sprintf("maker-%v", symbol)
	config := &Config{}
	err = basedb.LoadConf(ctx, key, config)
	if err != nil {
		return
	}
	maker := NewMaker(config)
	err = maker.Start(ctx)
	if err == nil {
		makerAll[symbol] = maker
		matcher.Shared.AddMonitor(symbol, maker)
	}
	return
}

func Stop(ctx context.Context, symbol string) (err error) {
	makerLock.Lock()
	defer makerLock.Unlock()
	maker, ok := makerAll[symbol]
	if !ok {
		err = fmt.Errorf("not running")
		return
	}
	maker.Stop()
	delete(makerAll, symbol)
	matcher.Shared.RemoveMonitor(symbol, maker)
	return
}

func LoadConfig(ctx context.Context, symbol string) (config *Config, err error) {
	key := fmt.Sprintf("maker-%v", symbol)
	config = &Config{}
	err = basedb.LoadConf(ctx, key, config)
	return
}

func UpdateConfig(ctx context.Context, config *Config) (err error) {
	key := fmt.Sprintf("maker-%v", config.Symbol)
	err = basedb.StoreConf(ctx, key, converter.JSON(config))
	if err != nil {
		return
	}
	makerLock.Lock()
	defer makerLock.Unlock()
	maker, ok := makerAll[config.Symbol]
	if ok {
		maker.Update(config)
	}
	return
}

func Find(ctx context.Context, symbol string) (maker *Maker) {
	makerLock.Lock()
	defer makerLock.Unlock()
	maker = makerAll[symbol]
	return
}

func List(ctx context.Context) (makerList []*Maker, makerMap map[string]*Maker) {
	makerLock.Lock()
	defer makerLock.Unlock()
	makerMap = map[string]*Maker{}
	for _, maker := range makerAll {
		makerList = append(makerList, maker)
		makerMap[maker.symbol.Symbol] = maker
	}
	return
}

var N1 = decimal.NewFromInt(1)

func randomRate(min float64) (rate decimal.Decimal) {
	if min < 1 {
		n := rand.Int63n(int64(100000 * (1 - min)))
		rate = decimal.NewFromInt(n).Div(decimal.NewFromInt(100000)).Add(decimal.NewFromFloat(min))
	}
	return
}

func randomRateValue(min, max decimal.Decimal) (value decimal.Decimal) {
	value = max.Sub(min).Mul(randomRate(0)).Add(min)
	return
}

func intRate(n, i int) (rate decimal.Decimal) {
	return int64Rate(int64(n), int64(i))
}

func int64Rate(n, i int64) (rate decimal.Decimal) {
	return decimal.NewFromInt(i).Div(decimal.NewFromInt(n))
}

type Config struct {
	ON       int             `json:"on"`       //if auto start
	UserID   int64           `json:"user_id"`  //maker uesr id
	Delay    int64           `json:"delay"`    //depth change delay
	Interval int64           `json:"interval"` //maker interval by ms
	Symbol   string          `json:"symbol"`   //symbol
	Open     decimal.Decimal `json:"open"`     //open price
	Close    struct {
		Min decimal.Decimal `json:"min"` //close price min change rate
		Max decimal.Decimal `json:"max"` //close price max change rate
	} `json:"close"` //close price
	Vib struct {
		Min   decimal.Decimal `json:"min"`   //price min vib change rate
		Max   decimal.Decimal `json:"max"`   //price max vib change rate
		Count int             `json:"count"` //price vib count
	} `json:"vib"` //vib price
	Ticker decimal.Decimal `json:"ticker"` //ticker max change rate
	Depth  struct {
		QtyMax  decimal.Decimal `json:"qty_max"`  //depth max qty
		StepMax int             `json:"step_max"` //depth step max
		DiffMax decimal.Decimal `json:"diff_max"` //depth ask/bid max diff
		DiffMin decimal.Decimal `json:"diff_min"` //depth ask/bin min diff
		Max     int             `json:"max"`      //depth max count
	} `json:"depth"`
}

func (c *Config) Valid() (err error) {
	if c == nil ||
		c.UserID < 1 || len(c.Symbol) < 1 ||
		c.Delay < 1 || c.Delay > 10*time.Second.Milliseconds() ||
		c.Open.Sign() <= 0 ||
		c.Close.Max.LessThanOrEqual(c.Close.Min) || c.Close.Min.LessThan(decimal.NewFromFloat(-1)) ||
		c.Vib.Max.LessThanOrEqual(c.Vib.Min) || c.Vib.Min.LessThan(decimal.NewFromFloat(-1)) || c.Vib.Count < 1 ||
		c.Ticker.Sign() <= 0 ||
		c.Depth.QtyMax.Sign() <= 0 ||
		c.Depth.StepMax <= 0 ||
		c.Depth.DiffMax.LessThanOrEqual(c.Depth.DiffMin) ||
		c.Depth.Max < 4 {
		err = fmt.Errorf("config must be \n" +
			"user_id/symbol is valid\n" +
			"0<delay<10s\n" +
			"open >0\n" +
			"close.max>close.min>-1\n" +
			"vib.max>vib.min\n" +
			"vib.count>0\n" +
			"ticker>0\n" +
			"interval>1min\n" +
			"depth.qty_max>0\n" +
			"depth.step_max>0\n" +
			"0.1>depth.diff_max>depth.diff_min>-0.1\n" +
			"depth.max>3",
		)
		return
	}
	return
}

func (c *Config) randomClose() (close decimal.Decimal) {
	rate := c.Close.Max.Sub(c.Close.Min).Mul(randomRate(0)).Add(c.Close.Min).Add(N1)
	close = c.Open.Mul(rate)
	return
}

func (c *Config) randomOneVib() (vib decimal.Decimal) {
	rate := c.Vib.Max.Sub(c.Vib.Min).Mul(randomRate(0)).Add(c.Vib.Min).Add(N1)
	vib = c.Open.Mul(rate)
	return
}

func (c *Config) randomVib() (vibs []decimal.Decimal, intervals []time.Duration) {
	totalInterval := c.Interval
	vibs = append(vibs, c.Open)
	intervals = append(intervals, 0)
	for i, n := 0, rand.Intn(c.Vib.Count-3)+4; i < n; i++ {
		vibs = append(vibs, c.randomOneVib())
		interval := intRate(n, n-i).Mul(randomRate(intRate(n, 1).InexactFloat64())).Mul(decimal.NewFromInt(totalInterval)).IntPart()
		totalInterval -= interval
		intervals = append(intervals, time.Duration(interval)*time.Millisecond)
	}
	vibs = append(vibs, c.randomClose())
	intervals = append(intervals, time.Duration(totalInterval)*time.Millisecond)
	return
}

//Scan is sql.Sanner
func (c *Config) Scan(src interface{}) (err error) {
	if src != nil {
		if jsonSrc, ok := src.(string); ok {
			err = json.Unmarshal([]byte(jsonSrc), c)
		} else {
			err = fmt.Errorf("the %v,%v is not string", reflect.TypeOf(src), src)
		}
	}
	return
}

//Value will parse to json value
func (c *Config) Value() (driver.Value, error) {
	if c == nil {
		return "{}", nil
	}
	bys, err := json.Marshal(c)
	return string(bys), err
}

func (c *Config) Random(past time.Duration, close decimal.Decimal) (next decimal.Decimal) {
	vibs, intervals := c.randomVib()
	totalInterval := time.Duration(0)
	var i int
	var vib decimal.Decimal
	var interval time.Duration
	for i, vib = range vibs {
		interval = intervals[i]
		totalInterval += intervals[i]
		if totalInterval > past {
			break
		}
	}
	if totalInterval <= past {
		next = close
		return
	}
	minFactor := int64Rate(interval.Milliseconds(), (totalInterval - past).Milliseconds()).Mul(decimal.NewFromFloat(0.5)).InexactFloat64()
	factor := randomRate(minFactor)
	rate := c.Ticker.Mul(decimal.NewFromFloat(2)).Mul(factor).Sub(c.Ticker)
	if close.LessThan(vib) {
		next = close.Add(close.Mul(rate))
	} else {
		next = close.Sub(close.Mul(rate))
	}
	return
}

type Maker struct {
	Config       *Config
	Verbose      bool
	Clear        time.Duration
	symbol       *matcher.SymbolInfo
	depth        *orderbook.Depth
	makingAll    map[string]decimal.Decimal
	makingOrder  map[string]*gexdb.Order
	balances     map[string]*gexdb.Balance
	holding      *gexdb.Holding
	startTime    time.Time
	close        decimal.Decimal
	next         decimal.Decimal
	nextAsk      decimal.Decimal
	nextBid      decimal.Decimal
	nextGen      int64
	nextShow     time.Time
	locker       sync.RWMutex
	makerQueue   chan int
	tickerExiter chan int
	tickerWaiter sync.WaitGroup
	clearLast    time.Time
	clearRemoved int64
	clearShow    time.Time
	waiter       sync.WaitGroup
}

func NewMaker(config *Config) (maker *Maker) {
	maker = &Maker{
		Config:       config,
		Clear:        5 * time.Second,
		makingAll:    map[string]decimal.Decimal{},
		makingOrder:  map[string]*gexdb.Order{},
		balances:     map[string]*gexdb.Balance{},
		locker:       sync.RWMutex{},
		makerQueue:   make(chan int, 1),
		tickerExiter: make(chan int, 1),
		tickerWaiter: sync.WaitGroup{},
		waiter:       sync.WaitGroup{},
	}
	return
}

func (m *Maker) Start(ctx context.Context) (err error) {
	m.close = m.Config.Open
	m.symbol = matcher.Shared.Symbols[m.Config.Symbol]
	if m.symbol == nil {
		err = fmt.Errorf("symbol %v is not found on matcher", m.Config.Symbol)
		return
	}
	if strings.HasPrefix(m.symbol.Symbol, "spot.") {
		_, err = gexdb.TouchBalance(ctx, gexdb.BalanceAreaSpot, []string{m.symbol.Base, m.symbol.Quote}, m.Config.UserID)
		if err != nil {
			return
		}
		_, m.balances, err = gexdb.ListUserBalance(ctx, m.Config.UserID, gexdb.BalanceAreaSpot, nil, nil)
		if err != nil {
			return
		}
	} else {
		_, err = gexdb.TouchBalance(ctx, gexdb.BalanceAreaFutures, []string{m.symbol.Quote}, m.Config.UserID)
		if err != nil {
			return
		}
		_, err = gexdb.TouchHolding(ctx, []string{m.symbol.Symbol}, m.Config.UserID)
		if err != nil {
			return
		}
		_, m.balances, err = gexdb.ListUserBalance(ctx, m.Config.UserID, gexdb.BalanceAreaFutures, nil, nil)
		if err != nil {
			return
		}
		m.holding, err = gexdb.FindHoldlingBySymbol(ctx, m.Config.UserID, m.symbol.Symbol)
		if err != nil {
			return
		}
	}
	m.nextShow = time.Now()
	m.clearShow = time.Now()
	m.waiter.Add(1)
	go m.loopMake()
	m.waiter.Add(1)
	m.tickerWaiter.Add(1)
	go m.loopTicker()
	xlog.Infof("Maker(%v) the maker is started", m.symbol)
	return
}

func (m *Maker) Stop() {
	m.tickerExiter <- 0
	m.tickerWaiter.Wait()
	m.makerQueue <- 0
	m.waiter.Wait()
}

func (m *Maker) Update(config *Config) {
	m.tickerExiter <- 0
	m.tickerWaiter.Wait()
	m.Config = config
	m.waiter.Add(1)
	m.tickerWaiter.Add(1)
	go m.loopTicker()
}

func (m *Maker) OnMatched(ctx context.Context, event *matcher.MatcherEvent) {
	m.locker.Lock()
	defer m.locker.Unlock()
	for _, order := range event.Orders {
		if order.UserID != m.Config.UserID || order.Price.Sign() <= 0 {
			continue
		}
		key := order.Price.String()
		switch order.Status {
		case gexdb.OrderStatusPending:
			m.makingAll[key] = m.makingAll[key].Add(order.Quantity)
			m.makingOrder[order.OrderID] = order
		case gexdb.OrderStatusPartialled:
			if m.makingOrder[order.OrderID] != nil {
				reduced := order.Filled.Sub(m.makingOrder[order.OrderID].Filled)
				m.makingAll[key] = m.makingAll[key].Sub(reduced)
			}
			m.makingOrder[order.OrderID] = order
		case gexdb.OrderStatusCanceled, gexdb.OrderStatusDone, gexdb.OrderStatusPartCanceled:
			if m.makingOrder[order.OrderID] != nil {
				reduced := order.Filled.Sub(m.makingOrder[order.OrderID].Filled)
				remain := order.Quantity.Sub(order.Filled)
				m.makingAll[key] = m.makingAll[key].Sub(reduced.Add(remain))
			}
			if m.makingAll[key].Sign() <= 0 {
				delete(m.makingAll, key)
			}
			delete(m.makingOrder, order.OrderID)
		}
	}
	for _, balance := range event.Balances {
		if balance.UserID != m.Config.UserID {
			continue
		}
		m.balances[balance.Asset] = balance
	}
	for _, holding := range event.Holdings {
		if holding.UserID != m.Config.UserID {
			continue
		}
		m.holding = holding
	}
	oldDepth := m.depth
	m.depth = event.Depth
	if len(m.depth.Asks) > 0 && len(m.depth.Bids) > 0 {
		m.close = m.depth.Asks[0][0].Add(m.depth.Bids[0][0]).DivRound(decimal.NewFromFloat(2), m.symbol.PrecisionPrice)
	}
	if (oldDepth != nil && m.depth != nil && len(oldDepth.Asks) > 0 && len(m.depth.Asks) > 0 && !oldDepth.Asks[0][0].Equal(m.depth.Asks[0][0]) && !m.depth.Asks[0][0].Equal(m.nextAsk)) ||
		(oldDepth != nil && m.depth != nil && len(oldDepth.Bids) > 0 && len(m.depth.Bids) > 0 && !oldDepth.Bids[0][0].Equal(m.depth.Bids[0][0]) && !m.depth.Bids[0][0].Equal(m.nextBid)) { //only ask/bid changed
		select {
		case m.makerQueue <- 1:
		default:
		}
	}
}

func (m *Maker) loopTicker() {
	delay := time.Duration(m.Config.Delay) * time.Millisecond
	ticker := time.NewTicker(delay)
	defer func() {
		ticker.Stop()
		m.tickerWaiter.Done()
		m.waiter.Done()
	}()
	running := true
	xlog.Infof("Maker(%v) ticker runner is starting by %v", m.symbol, delay)
	for running {
		select {
		case <-m.tickerExiter:
			running = false
		case <-ticker.C:
			m.makerQueue <- 2
		}
	}
	xlog.Infof("Maker(%v) ticker runner is stopped", m.symbol)
}

func (m *Maker) loopMake() {
	defer m.waiter.Done()
	running := true
	xlog.Infof("Maker(%v) maker runner is starting", m.symbol)
	for running {
		v := <-m.makerQueue
		if v < 1 {
			running = false
		} else {
			m.procMake(v > 1)
		}
	}
	xlog.Infof("Maker(%v) maker runner is stopped", m.symbol)
}

func (m *Maker) procMake(next bool) (err error) {
	ctx := context.Background()
	if next || m.next.Sign() <= 0 {
		m.randomNext()
	}

	//cancle first
	m.procCancle(ctx, m.nextAsk, m.nextBid)

	//place new
	depth := rand.Intn(m.Config.Depth.StepMax-3) + 3
	priceStep := decimal.New(1, -m.symbol.PrecisionPrice)
	for i := 0; i < depth; i++ {
		{ //bid
			step := rand.Intn(m.Config.Depth.StepMax-1) + 1
			price := m.nextBid.Sub(priceStep.Mul(decimal.NewFromInt(int64(step)))).RoundDown(m.symbol.PrecisionPrice)
			m.procPlace(ctx, gexdb.OrderSideBuy, price)
		}
		{ //ask
			step := rand.Intn(m.Config.Depth.StepMax-1) + 1
			price := m.nextAsk.Add(priceStep.Mul(decimal.NewFromInt(int64(step)))).RoundUp(m.symbol.PrecisionPrice)
			m.procPlace(ctx, gexdb.OrderSideSell, price)
		}
	}
	m.procPlace(ctx, gexdb.OrderSideBuy, m.nextBid)
	m.procPlace(ctx, gexdb.OrderSideSell, m.nextAsk)

	//clear
	m.procClear(ctx)
	return
}

func (m *Maker) randomNext() {
	if time.Since(m.startTime).Milliseconds() >= m.Config.Interval {
		m.startTime = time.Now().Add(-time.Second)
	}
	m.next = m.Config.Random(time.Since(m.startTime), m.close).Round(m.symbol.PrecisionPrice)
	diff := randomRateValue(m.Config.Depth.DiffMin, m.Config.Depth.DiffMax).Div(decimal.NewFromFloat(2))
	m.nextAsk = m.next.Add(diff).RoundUp(m.symbol.PrecisionPrice)
	m.nextBid = m.next.Sub(diff).RoundDown(m.symbol.PrecisionPrice)
	m.nextGen++
	if time.Since(m.nextShow) > time.Minute {
		xlog.Infof("Maker(%v) random gen %v price in past %v, latest next:%v,ask:%v,bid:%v", m.symbol, m.nextGen, time.Since(m.nextShow), m.next, m.nextAsk, m.nextBid)
		m.nextGen = 0
		m.nextShow = time.Now()
	}
}

func (m *Maker) checkOrder(ask, bid decimal.Decimal) (cancel, all []string) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	for _, order := range m.makingOrder {
		if order.Status == gexdb.OrderStatusCanceled {
			continue
		}
		if order.Side == gexdb.OrderSideSell && order.Price.LessThan(ask) {
			cancel = append(cancel, order.OrderID)
		} else if order.Side == gexdb.OrderSideBuy && order.Price.GreaterThan(bid) {
			cancel = append(cancel, order.OrderID)
		} else {
			all = append(all, order.OrderID)
		}
	}
	return
}

func (m *Maker) markOrderCanceled(orderIDs map[string]gexdb.OrderStatus) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	for orderID, status := range orderIDs {
		order := m.makingOrder[orderID]
		if order != nil {
			order.Status = status
		}
	}
}

func (m *Maker) willPlace(side gexdb.OrderSide, qty, price decimal.Decimal) bool {
	m.locker.RLock()
	defer m.locker.RUnlock()
	making := m.makingAll[price.String()]
	if making.Add(qty).GreaterThan(m.Config.Depth.QtyMax) {
		return false
	}
	//check can place
	if strings.HasPrefix(m.symbol.Symbol, "spot.") {
		if side == gexdb.OrderSideBuy {
			balance := m.balances[m.symbol.Quote]
			if balance != nil && balance.Free.GreaterThanOrEqual(qty.Mul(price)) {
				balance.Free = balance.Free.Sub(qty.Mul(price))
				return true
			}
		} else {
			balance := m.balances[m.symbol.Base]
			if balance != nil && balance.Free.GreaterThanOrEqual(qty) {
				balance.Free = balance.Free.Sub(qty)
				return true
			}
		}
		return false
	}
	if m.holding != nil && side == gexdb.OrderSideBuy && m.holding.Amount.LessThanOrEqual(decimal.Zero.Sub(qty)) {
		m.holding.Amount = m.holding.Amount.Add(qty)
		return true
	} else if m.holding != nil && side == gexdb.OrderSideSell && m.holding.Amount.GreaterThanOrEqual(qty) {
		m.holding.Amount = m.holding.Amount.Sub(qty)
		return true
	} else if balance := m.balances[m.symbol.Quote]; balance != nil && balance.Free.GreaterThan(qty.Mul(price)) {
		balance.Free = balance.Free.Sub(qty.Mul(price))
		return true
	}
	return false
}

func (m *Maker) procCancle(ctx context.Context, ask, bid decimal.Decimal) {
	cancelIDs, allIDs := m.checkOrder(ask, bid)
	if len(allIDs) > 0 {
		picked := map[int]bool{}
		for i, n := 0, rand.Intn(5); i < n; i++ {
			x := rand.Intn(len(allIDs))
			if !picked[x] {
				cancelIDs = append(cancelIDs, allIDs[x])
				picked[x] = true
			}
		}
	}
	if m.Verbose && len(cancelIDs) > 0 {
		xlog.Infof("Maker start cancle %v/%v order by new ask:%v,bid:%v", len(cancelIDs), len(allIDs), ask, bid)
	}
	canceledIDs := map[string]gexdb.OrderStatus{}
	for _, orderID := range cancelIDs {
		order, err := matcher.ProcessCancel(ctx, m.Config.UserID, m.symbol.Symbol, orderID)
		if err != nil {
			xlog.Warnf("Maker cancle order %v fail with\n%v", orderID, matcher.ErrStack(err))
			continue
		}
		if m.Verbose {
			xlog.Infof("Maker cancle order is done with %v", order.Info())
		}
		if order.Status != gexdb.OrderStatusPartialled && order.Status != gexdb.OrderStatusPending {
			canceledIDs[orderID] = order.Status
		}
	}
	if len(canceledIDs) > 0 {
		m.markOrderCanceled(canceledIDs)
	}
}

func (m *Maker) procPlace(ctx context.Context, side gexdb.OrderSide, price decimal.Decimal) {
	qtyStep := decimal.New(1, -m.symbol.PrecisionQuantity)
	qty := randomRateValue(qtyStep, m.Config.Depth.QtyMax.Sub(qtyStep)).RoundUp(m.symbol.PrecisionQuantity)
	if !m.willPlace(side, qty, price) {
		return
	}
	order, err := matcher.ProcessLimit(ctx, m.Config.UserID, m.symbol.Symbol, side, qty, price)
	if err != nil {
		xlog.Warnf("Maker place order by symbol:%v,side:%v,qty:%v,price:%v fail with\n%v", m.symbol.Symbol, side, qty, price, matcher.ErrStack(err))
		return
	}
	if m.Verbose {
		xlog.Infof("Maker place order is done with %v", order.Info())
	}
}

func (m *Maker) procClear(ctx context.Context) {
	if time.Since(m.clearLast) < m.Clear {
		return
	}
	m.clearLast = time.Now()
	cleared, err := gexdb.ClearCanceledOrder(ctx, m.Config.UserID, m.symbol.Symbol, time.Now())
	if err != nil {
		xlog.Warnf("Maker(%v) clear canceled order fail with %v", m.symbol, err)
		return
	}
	m.clearRemoved = cleared
	if time.Since(m.clearShow) > time.Minute {
		xlog.Infof("Maker(%v) clear %v canceled order in past %v", m.symbol, m.clearRemoved, time.Since(m.clearShow))
		m.clearRemoved = 0
		m.clearShow = time.Now()
	}
}
