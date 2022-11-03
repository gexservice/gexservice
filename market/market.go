package market

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/codingeasygo/util/debug"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsort"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/util/xtime"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
	"golang.org/x/net/websocket"
)

var Shared *Market
var Quote string = "USDT"

func Bootstrap() {
	Shared = NewMarket(matcher.Shared.Symbols)
	matcher.Shared.AddMonitor("*", Shared)
	Shared.Start()
}

func ListSymbol() (symbols []*matcher.SymbolInfo, lines map[string]*gexdb.KLine) {
	lines = map[string]*gexdb.KLine{}
	for _, symbol := range Shared.Symbols {
		symbols = append(symbols, symbol)
		line := LoadKLine(symbol.Symbol, "1day")
		if line == nil {
			continue
		}
		lines[symbol.Symbol] = line
	}
	xsort.SortFunc(symbols, func(x, y int) bool {
		linex, liney := lines[symbols[x].Symbol], lines[symbols[y].Symbol]
		if linex == nil || liney == nil {
			return liney == nil
		}
		if linex.Open.Sign() <= 0 || liney.Open.Sign() <= 0 {
			return linex.Open.Sign() >= liney.Open.Sign()
		}
		ratex := linex.Close.Sub(linex.Open).Div(linex.Open)
		ratey := liney.Close.Sub(liney.Open).Div(liney.Open)
		return ratex.GreaterThan(ratey)
	})
	return
}

func LoadSymbol(symbol string) (info *matcher.SymbolInfo, line *gexdb.KLine) {
	info = Shared.Symbols[symbol]
	line = LoadKLine(symbol, "1day")
	return
}

func LoadKLine(symbol, interval string) (line *gexdb.KLine) {
	line = Shared.LoadKLine(symbol, interval)
	return
}

func ListKLine(ctx context.Context, symbol, interval string, startTime, endTime time.Time) (lines []*gexdb.KLine, err error) {
	lines, err = Shared.ListKLine(ctx, symbol, interval, startTime, endTime)
	return
}

func LoadLatestPrice(symbol string) (price decimal.Decimal) {
	price = Shared.LoadLatestPrice(symbol)
	return
}

func ListLatestPrice(symbols ...string) (prices map[string]decimal.Decimal) {
	prices = Shared.ListLatestPrice(symbols...)
	return
}

func LoadDepth(symbol string, max int) (depth *DepthCache) {
	depth = Shared.LoadDepth(symbol, max)
	return
}

func klineKey(symbol, interv string) string {
	return fmt.Sprintf("%v-%v", symbol, interv)
}

type MarketConn struct {
	Conn    *websocket.Conn
	Codec   *websocket.Codec
	KLines  map[string]int
	Depths  map[string]int
	Tickers map[string]int
	Timeout time.Duration
	Ready   bool
	Latest  time.Time
}

func NewMarketConn(conn *websocket.Conn) (mc *MarketConn) {
	mc = &MarketConn{
		Conn:    conn,
		Timeout: 3 * time.Second,
		Latest:  time.Now(),
		KLines:  map[string]int{},
		Depths:  map[string]int{},
	}
	mc.Codec = &websocket.Codec{
		Marshal:   mc.codecMarshal,
		Unmarshal: mc.codecUnmarshal,
	}
	return
}

func (m *MarketConn) codecMarshal(v interface{}) (data []byte, payloadType byte, err error) {
	if s, ok := v.(string); ok {
		switch s {
		case "PING":
			payloadType = websocket.PingFrame
			return
		case "PONG":
			payloadType = websocket.PongFrame
			return
		default:
			payloadType = websocket.TextFrame
			data = []byte(s)
		}
	} else {
		payloadType = websocket.TextFrame
		data, err = json.Marshal(v)
	}
	return
}

func (m *MarketConn) codecUnmarshal(data []byte, payloadType byte, v interface{}) (err error) {
	m.Latest = time.Now()
	switch payloadType {
	case websocket.PingFrame, websocket.PongFrame:
	default:
		err = json.Unmarshal(data, v)
	}
	return
}

func (m *MarketConn) Send(v interface{}) (err error) {
	if m.Conn == nil {
		err = fmt.Errorf("%v", "send not connected")
		return
	}
	m.Conn.SetWriteDeadline(time.Now().Add(m.Timeout))
	err = m.Codec.Send(m.Conn, v)
	return
}

func (m *MarketConn) Receive(v interface{}) (err error) {
	if m.Conn == nil {
		err = fmt.Errorf("%v", "recv not connected")
		return
	}
	err = m.Codec.Receive(m.Conn, v)
	return
}

func (m *MarketConn) Close() (err error) {
	if m.Conn == nil {
		err = fmt.Errorf("%v", "close not connected")
		return
	}
	err = m.Conn.Close()
	return
}

func (m *MarketConn) RemoteAddr() string {
	return m.Conn.Request().RemoteAddr
}

type KLineCache struct {
	Latest time.Time
	Oldest time.Time
	Lines  []*gexdb.KLine
}

type DepthCache struct {
	Bids   [][]decimal.Decimal `json:"bids"`
	Asks   [][]decimal.Decimal `json:"asks"`
	Symbol string              `json:"symbol"`
	Time   xsql.Time           `json:"time"`
}

func (d *DepthCache) Slice(max int) (depth *DepthCache) {
	depth = &DepthCache{
		Asks:   d.Asks,
		Bids:   d.Bids,
		Symbol: d.Symbol,
		Time:   xsql.TimeNow(),
	}
	if len(depth.Asks) > max {
		depth.Asks = depth.Asks[0:max]
	}
	if len(depth.Bids) > max {
		depth.Bids = depth.Bids[0:max]
	}
	return
}

type klineQueueItem struct {
	Conn *MarketConn
}
type depthQueueItem struct {
	Conn   *MarketConn
	Symbol string
}

type Market struct {
	Symbols          map[string]*matcher.SymbolInfo
	WaitTimeout      time.Duration
	KLineGenDelay    time.Duration
	KLineNotifyDelay time.Duration
	NotiryRunner     int
	OnConnect        func(conn *websocket.Conn)
	OnDisconnect     func(conn *websocket.Conn)
	eventQueue       chan *matcher.MatcherEvent
	avgPrice         map[string]decimal.Decimal
	klineVal         map[string]*gexdb.KLine
	klineCache       map[string]*KLineCache
	klineQueue       chan *klineQueueItem
	klineLock        sync.RWMutex
	depthVal         map[string]*DepthCache
	depthQueue       chan *depthQueueItem
	depthLock        sync.RWMutex
	wsconn           map[string]*MarketConn
	wslock           sync.RWMutex
	exiter           chan int
	waiter           sync.WaitGroup
}

func NewMarket(symbols map[string]*matcher.SymbolInfo) (market *Market) {
	market = &Market{
		Symbols:          symbols,
		WaitTimeout:      3 * time.Second,
		KLineGenDelay:    time.Second,
		KLineNotifyDelay: time.Second,
		NotiryRunner:     3,
		eventQueue:       make(chan *matcher.MatcherEvent, 1024),
		avgPrice:         map[string]decimal.Decimal{},
		klineVal:         map[string]*gexdb.KLine{},
		klineCache:       map[string]*KLineCache{},
		klineQueue:       make(chan *klineQueueItem, 1024),
		klineLock:        sync.RWMutex{},
		depthVal:         map[string]*DepthCache{},
		depthQueue:       make(chan *depthQueueItem, 1024),
		depthLock:        sync.RWMutex{},
		wsconn:           map[string]*MarketConn{},
		wslock:           sync.RWMutex{},
		exiter:           make(chan int, 1024),
		waiter:           sync.WaitGroup{},
	}
	return
}

func (m *Market) Start() {
	m.waiter.Add(1)
	go m.loopEvent()
	m.waiter.Add(1)
	go m.loopTriggerKLine()
	for i := 0; i < m.NotiryRunner; i++ {
		m.waiter.Add(1)
		go m.loopNotify()
	}
}

func (m *Market) Stop() {
	m.exiter <- 0
	m.exiter <- 0
	for i := 0; i < m.NotiryRunner; i++ {
		m.exiter <- 0
	}
	m.waiter.Wait()
}

func (m *Market) OnMatched(ctx context.Context, event *matcher.MatcherEvent) {
	select {
	case m.eventQueue <- event:
	default:
	}
}

func (m *Market) loopEvent() {
	defer m.waiter.Done()
	ticker := time.NewTicker(m.KLineGenDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-m.exiter:
			running = false
		case event := <-m.eventQueue:
			m.procGenKLine(event)
			m.procTriggerDepth(event)
		case <-ticker.C:
			m.procGenKLine(nil)
		}
	}
	xlog.Infof("Market event loop is stopped")
}

func (m *Market) listCurrentKLine(symbol string) (lines []*gexdb.KLine) {
	now := time.Now()

	startMin5 := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), (now.Minute()/5)*5, 0, 0, now.Location())
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startMin5), Interv: "5min", UpdateTime: xsql.TimeNow()})

	startMin30 := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), (now.Minute()/30)*30, 0, 0, now.Location())
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startMin30), Interv: "30min", UpdateTime: xsql.TimeNow()})

	startHour1 := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startHour1), Interv: "1hour", UpdateTime: xsql.TimeNow()})

	startHour4 := time.Date(now.Year(), now.Month(), now.Day(), (now.Hour()/4)*4, 0, 0, 0, now.Location())
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startHour4), Interv: "4hour", UpdateTime: xsql.TimeNow()})

	startDay1 := xtime.TimeStartOfToday()
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startDay1), Interv: "1day", UpdateTime: xsql.TimeNow()})

	startWeek := xtime.TimeStartOfWeek()
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startWeek), Interv: "1week", UpdateTime: xsql.TimeNow()})

	startMonth := xtime.TimeStartOfMonth()
	lines = append(lines, &gexdb.KLine{Symbol: symbol, StartTime: xsql.Time(startMonth), Interv: "1mon", UpdateTime: xsql.TimeNow()})
	return
}

func (m *Market) procGenKLine(event *matcher.MatcherEvent) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			xlog.Errorf("Market proc kline panic with %v, callstack is \n%v", perr, debug.CallStatck())
			err = fmt.Errorf("%v", err)
		}
	}()

	var saveLine []*gexdb.KLine
	if event != nil {
		if len(event.Orders) < 1 || event.Orders[0].Filled.Sign() <= 0 {
			return
		}
		order := event.Orders[0]
		avgPrice := order.AvgPrice

		m.klineLock.Lock()
		m.avgPrice[event.Symbol] = avgPrice
		for _, line := range m.listCurrentKLine(event.Symbol) {
			key := klineKey(line.Symbol, line.Interv)
			having, ok := m.klineVal[key]
			if !ok {
				m.klineVal[key] = line
				having = line
				having.Open = avgPrice
				having.Low = avgPrice
				having.High = avgPrice
			} else if having.StartTime.Timestamp() != line.StartTime.Timestamp() {
				saveLine = append(saveLine, having)
				m.klineVal[key] = line
				having = line
				having.Open = avgPrice
				having.Low = avgPrice
				having.High = avgPrice
				having.UpdateTime = xsql.TimeNow()
				delete(m.klineCache, key)
			}
			having.Amount = having.Amount.Add(order.Filled)
			having.Volume = having.Volume.Add(order.TotalPrice)
			having.Count++
			having.Close = avgPrice
			if having.Low.Sign() <= 0 || having.Low.GreaterThan(avgPrice) {
				having.Low = avgPrice
			}
			if having.High.LessThan(avgPrice) {
				having.High = avgPrice
			}
			having.UpdateTime = xsql.TimeNow()
		}
		m.klineLock.Unlock()
	} else {
		m.klineLock.Lock()
		for _, symbol := range m.Symbols {
			for _, line := range m.listCurrentKLine(symbol.Symbol) {
				key := klineKey(line.Symbol, line.Interv)
				having, ok := m.klineVal[key]
				if ok && having.StartTime.Timestamp() != line.StartTime.Timestamp() {
					saveLine = append(saveLine, having)
					m.klineVal[key] = line
					old := having
					having = line
					having.Open = old.Close
					having.Close = old.Close
					having.Low = old.Close
					having.High = old.Close
					having.UpdateTime = xsql.TimeNow()
					delete(m.klineCache, key)
				}
			}
		}
		m.klineLock.Unlock()
	}

	if len(saveLine) > 0 {
		added, xerr := gexdb.AddMultiKLine(context.Background(), saveLine...)
		if xerr != nil {
			err = xerr
			xlog.Errorf("Market add %v kline fail with %v", len(saveLine), xerr)
		} else {
			xlog.Infof("Market add %v kline to database success", added)
		}
	}
	return
}

func (m *Market) procTriggerDepth(event *matcher.MatcherEvent) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			xlog.Errorf("Market proc trigger depth panic with %v, callstack is \n%v", perr, debug.CallStatck())
			err = fmt.Errorf("%v", err)
		}
	}()
	m.depthLock.Lock()
	depth := &DepthCache{
		Symbol: event.Symbol,
		Asks:   event.Depth.Asks,
		Bids:   event.Depth.Bids,
		Time:   xsql.TimeNow(),
	}
	m.depthVal[depth.Symbol] = depth
	m.depthLock.Unlock()
	conns := []*MarketConn{}
	m.wslock.RLock()
	for _, conn := range m.wsconn {
		conns = append(conns, conn)
	}
	m.wslock.RUnlock()
	for _, conn := range conns {
		if conn.Depths[depth.Symbol] < 1 && conn.Tickers[depth.Symbol] < 1 {
			continue
		}
		select {
		case m.depthQueue <- &depthQueueItem{Conn: conn, Symbol: depth.Symbol}:
		default:
		}
	}
	return
}

func (m *Market) loopTriggerKLine() {
	defer m.waiter.Done()
	ticker := time.NewTicker(m.KLineNotifyDelay)
	defer ticker.Stop()
	running := true
	for running {
		select {
		case <-m.exiter:
			running = false
		case <-ticker.C:
			m.procTriggerKLine()
		}
	}
	xlog.Infof("Market notify loop is stopped")
}

func (m *Market) procTriggerKLine() (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			xlog.Errorf("Market proc trigger kline panic with %v, callstack is \n%v", perr, debug.CallStatck())
			err = fmt.Errorf("%v", err)
		}
	}()

	lines := map[string]*gexdb.KLine{}
	m.klineLock.RLock()
	for key, line := range m.klineVal {
		lines[key] = line
	}
	m.klineLock.RUnlock()

	conns := []*MarketConn{}
	m.wslock.RLock()
	for _, conn := range m.wsconn {
		conns = append(conns, conn)
	}
	m.wslock.RUnlock()

	for _, conn := range conns {
		if !conn.Ready && time.Since(conn.Latest) > m.WaitTimeout {
			conn.Close()
			continue
		}
		having := false
		for key := range conn.KLines {
			line := lines[key]
			if line != nil {
				having = true
				break
			}
		}
		if having {
			m.klineQueue <- &klineQueueItem{Conn: conn}
		}
	}
	return
}

func (m *Market) loopNotify() {
	defer m.waiter.Done()
	running := true
	for running {
		select {
		case <-m.exiter:
			running = false
		case item := <-m.klineQueue:
			m.procNotifyKLine(item.Conn)
		case item := <-m.depthQueue:
			m.procNotifyDepth(item.Conn, item.Symbol)
			m.procNotifyTicker(item.Conn, item.Symbol)
		}
	}
}

func (m *Market) procNotifyKLine(conn *MarketConn) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			xlog.Errorf("Market proc notify kline panic with %v, callstack is \n%v", perr, debug.CallStatck())
		}
	}()
	lines := map[string]*gexdb.KLine{}
	m.klineLock.RLock()
	for key := range conn.KLines {
		line := m.klineVal[key]
		if line != nil {
			lines[key] = line
		}
	}
	m.klineLock.RUnlock()
	for _, line := range lines {
		err = conn.Send(xmap.M{
			"action": "notify.kline",
			"kline":  line,
			"code":   define.Success,
		})
		if err != nil {
			conn.Close()
			break
		}
	}
	return
}

func (m *Market) procNotifyDepth(conn *MarketConn, symbol string) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			xlog.Errorf("Market proc notify depth panic with %v, callstack is \n%v", perr, debug.CallStatck())
		}
	}()
	max := conn.Depths[symbol]
	if max < 1 {
		return
	}
	m.depthLock.RLock()
	depth := m.depthVal[symbol]
	m.depthLock.RUnlock()
	if depth == nil {
		return
	}
	err = conn.Send(xmap.M{
		"action": "notify.depth",
		"depth":  depth.Slice(max),
		"code":   define.Success,
	})
	if err != nil {
		conn.Close()
	}
	return
}

func (m *Market) procNotifyTicker(conn *MarketConn, symbol string) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			xlog.Errorf("Market proc notify ticker panic with %v, callstack is \n%v", perr, debug.CallStatck())
		}
	}()
	on := conn.Tickers[symbol]
	if on < 1 {
		return
	}
	m.depthLock.RLock()
	depth := m.depthVal[symbol]
	m.depthLock.RUnlock()
	if depth == nil || len(depth.Asks) < 1 || len(depth.Bids) < 1 {
		return
	}
	err = conn.Send(xmap.M{
		"action": "notify.ticker",
		"ticker": xmap.M{
			"symbol": symbol,
			"ask":    depth.Asks[0],
			"bid":    depth.Bids[0],
			"close":  m.LoadLatestPrice(symbol),
		},
		"code": define.Success,
	})
	if err != nil {
		conn.Close()
	}
	return
}

func (m *Market) SrvHTTP(s *web.Session) web.Result {
	srv := websocket.Server{Handler: m.HandWs}
	srv.ServeHTTP(s.W, s.R)
	return web.Return
}

func (m *Market) HandWs(raw *websocket.Conn) {
	var err error
	conn := NewMarketConn(raw)
	key := fmt.Sprintf("%p", conn)
	m.wslock.Lock()
	m.wsconn[key] = conn
	m.wslock.Unlock()
	defer func() {
		m.wslock.Lock()
		delete(m.wsconn, key)
		m.wslock.Unlock()
		if m.OnDisconnect != nil {
			m.OnDisconnect(raw)
		}
		conn.Close()
	}()
	if m.OnConnect != nil {
		m.OnConnect(raw)
	}
	xlog.Infof("Market accept ws from %v", raw.Request().RemoteAddr)
	for {
		cmd := xmap.M{}
		err = conn.Receive(&cmd)
		if err != nil {
			break
		}
		if len(cmd) < 1 {
			continue
		}
		action := cmd.StrDef("", "action")
		switch action {
		case "sub.kline":
			err = m.handSubKLine(conn, cmd)
		case "sub.depth":
			err = m.handSubDepth(conn, cmd)
		case "sub.ticker":
			err = m.handSubTicker(conn, cmd)
		default:
			conn.Send(xmap.M{
				"action":  "error",
				"code":    define.ArgsInvalid,
				"message": "unknow action",
			})
			err = fmt.Errorf("unknow action")
		}
		if err != nil {
			break
		}
	}
	xlog.Infof("Market ws from %v is closed by %v", raw.Request().RemoteAddr, err)
}

func (m *Market) handSubKLine(conn *MarketConn, cmd xmap.M) (err error) {
	klines := map[string]int{}
	for _, sub := range cmd.ArrayMapDef(nil, "symbols") {
		symbol := sub.StrDef("", "symbol")
		interval := sub.StrDef("", "interval")
		klines[klineKey(symbol, interval)] = 1
	}
	conn.KLines = klines
	conn.Ready = true
	xlog.Infof("Market ws from %v is ready for kline", conn.RemoteAddr())
	err = conn.Send(xmap.M{
		"action": "sub.kline",
		"code":   define.Success,
	})
	return
}

func (m *Market) handSubDepth(conn *MarketConn, cmd xmap.M) (err error) {
	depths := map[string]int{}
	for _, sub := range cmd.ArrayMapDef(nil, "symbols") {
		symbol := sub.StrDef("", "symbol")
		max := sub.IntDef(5, "max")
		depths[symbol] = max
	}
	conn.Depths = depths
	conn.Ready = true
	xlog.Infof("Market ws from %v is ready for depth", conn.RemoteAddr())
	err = conn.Send(xmap.M{
		"action": "sub.depth",
		"code":   define.Success,
	})
	return
}

func (m *Market) handSubTicker(conn *MarketConn, cmd xmap.M) (err error) {
	tickers := map[string]int{}
	for _, symbol := range cmd.ArrayStrDef(nil, "symbols") {
		tickers[symbol] = 1
	}
	conn.Tickers = tickers
	conn.Ready = true
	xlog.Infof("Market ws from %v is ready for ticker", conn.RemoteAddr())
	err = conn.Send(xmap.M{
		"action": "sub.ticker",
		"code":   define.Success,
	})
	return
}

func (m *Market) LoadKLine(symbol, interval string) (line *gexdb.KLine) {
	m.klineLock.RLock()
	defer m.klineLock.RUnlock()
	line = m.klineVal[klineKey(symbol, interval)]
	return
}

func (m *Market) ListKLine(ctx context.Context, symbol, interval string, startTime, endTime time.Time) (lines []*gexdb.KLine, err error) {
	key := klineKey(symbol, interval)
	interv, err := gexdb.StringInterv(interval)
	if err != nil {
		return
	}
	startTime = xtime.TimeUnix((xtime.Timestamp(startTime) / interv.Milliseconds()) * interv.Milliseconds())
	m.klineLock.Lock()
	defer m.klineLock.Unlock()
	latest := m.klineVal[key]
	cache := m.klineCache[key]
	if cache == nil {
		cache = &KLineCache{}
		cache.Lines, err = gexdb.ListKLine(ctx, symbol, interval, startTime, time.Now())
		if err != nil {
			return
		}
		if len(cache.Lines) > 0 {
			cache.Latest = cache.Lines[0].StartTime.AsTime()
			cache.Oldest = cache.Lines[len(cache.Lines)-1].StartTime.AsTime()
			m.klineCache[key] = cache
		}
		if latest != nil && latest.StartTime.AsTime().Before(endTime) {
			lines = append(lines, latest)
		}
		lines = append(lines, cache.Lines...)
		return
	}
	if cache.Oldest.After(startTime) {
		var oldLines []*gexdb.KLine
		oldLines, err = gexdb.ListKLine(ctx, symbol, interval, startTime, cache.Oldest)
		if err != nil {
			return
		}
		cache.Lines = append(cache.Lines, oldLines...)
	}
	if latest != nil && latest.StartTime.AsTime().Before(endTime) {
		lines = append(lines, latest)
	}
	for _, line := range cache.Lines {
		if (line.StartTime.AsTime().After(startTime) || line.StartTime.AsTime().Equal(startTime)) && line.StartTime.AsTime().Before(endTime) {
			lines = append(lines, line)
		}
	}
	return
}

func (m *Market) LoadLatestPrice(symbol string) decimal.Decimal {
	m.klineLock.RLock()
	defer m.klineLock.RUnlock()
	return m.avgPrice[symbol]
}

func (m *Market) ListLatestPrice(symbols ...string) (prices map[string]decimal.Decimal) {
	m.klineLock.RLock()
	defer m.klineLock.RUnlock()
	prices = map[string]decimal.Decimal{}
	for _, symbol := range symbols {
		prices[symbol] = m.avgPrice[symbol]
	}
	return
}

func (m *Market) LoadDepth(symbol string, max int) (depth *DepthCache) {
	m.depthLock.RLock()
	defer m.depthLock.RUnlock()
	depth = m.depthVal[symbol]
	if depth != nil {
		depth = depth.Slice(max)
	}
	return
}
