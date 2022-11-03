package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/converter"
// 	"github.com/codingeasygo/util/xhash"
// 	"github.com/codingeasygo/util/xsort"
// 	"github.com/codingeasygo/util/xsql"
// 	"github.com/codingeasygo/util/xtime"
// 	"github.com/shopspring/decimal"
// 	"github.com/gexservice/gexservice/base/xlog"
// 	"github.com/gexservice/gexservice/gexapi"
// 	"github.com/gexservice/gexservice/gexdb"
// )

// func GenerateTestUser(staffCount, brokerCount, userCount, orderCount int) (staffs, brokers, users []*gexdb.User, err error) {
// 	caller, err := gexdb.Pool().Begin()
// 	if err != nil {
// 		return
// 	}
// 	defer func() {
// 		if err == nil {
// 			err = caller.Commit()
// 		} else {
// 			caller.Rollback()
// 		}
// 	}()
// 	for i := 0; i < staffCount; i++ {
// 		staff := &gexdb.User{
// 			Type:      gexdb.UserTypeNormal,
// 			Role:      gexdb.UserRoleStaff,
// 			Name:      converter.StringPtr(fmt.Sprintf("Staff%02d", i)),
// 			Account:   converter.StringPtr(fmt.Sprintf("staff%02d", i)),
// 			Password:  converter.StringPtr(xhash.SHA1([]byte("123"))),
// 			TradePass: converter.StringPtr(xhash.SHA1([]byte("123"))),
// 			Status:    gexdb.UserStatusNormal,
// 		}
// 		err = gexdb.AddUserCall(caller, staff)
// 		if err != nil {
// 			break
// 		}
// 		_, err = gexdb.TouchBalanceCall(caller, gexdb.BalanceAssetAll, staff.TID)
// 		if err != nil {
// 			break
// 		}
// 		for j := 0; j < 10; j++ {
// 			_, _, err = gexdb.ChangeBalanceCall(caller, 1000, staff.TID, gexdb.BalanceAssetYWE, decimal.NewFromFloat(1100))
// 			if err != nil {
// 				break
// 			}
// 			_, _, err = gexdb.ChangeBalanceCall(caller, 1000, staff.TID, gexdb.BalanceAssetYWE, decimal.NewFromFloat(-100))
// 			if err != nil {
// 				break
// 			}
// 		}
// 		if err != nil {
// 			break
// 		}
// 		staffs = append(staffs, staff)
// 	}
// 	if err != nil {
// 		return
// 	}
// 	userAll := 0
// 	for i := 0; i < brokerCount; i++ {
// 		broker := &gexdb.User{
// 			Type:      gexdb.UserTypeNormal,
// 			Role:      gexdb.UserRoleBroker,
// 			Name:      converter.StringPtr(fmt.Sprintf("Broker%02d", i)),
// 			Account:   converter.StringPtr(fmt.Sprintf("broker%02d", i)),
// 			Password:  converter.StringPtr(xhash.SHA1([]byte("123"))),
// 			TradePass: converter.StringPtr(xhash.SHA1([]byte("123"))),
// 			KbzOpenid: converter.StringPtr(fmt.Sprintf("broker%02d", i)),
// 			Status:    gexdb.UserStatusNormal,
// 		}
// 		err = gexdb.AddUserCall(caller, broker)
// 		if err != nil {
// 			break
// 		}
// 		_, err = gexdb.TouchBalanceCall(caller, gexdb.BalanceAssetAll, broker.TID)
// 		if err != nil {
// 			break
// 		}
// 		for j := 0; j < userCount; j++ {
// 			user := &gexdb.User{
// 				Type:      gexdb.UserTypeNormal,
// 				Role:      gexdb.UserRoleNormal,
// 				Name:      converter.StringPtr(fmt.Sprintf("User%02d", userAll)),
// 				Account:   converter.StringPtr(fmt.Sprintf("user%02d", userAll)),
// 				Password:  converter.StringPtr(xhash.SHA1([]byte("123"))),
// 				TradePass: converter.StringPtr(xhash.SHA1([]byte("123"))),
// 				BrokerID:  broker.TID,
// 				KbzOpenid: converter.StringPtr(fmt.Sprintf("user%02d", userAll)),
// 				Status:    gexdb.UserStatusNormal,
// 			}
// 			err = gexdb.AddUserCall(caller, user)
// 			if err != nil {
// 				break
// 			}
// 			userAll++
// 			_, err = gexdb.TouchBalanceCall(caller, gexdb.BalanceAssetAll, user.TID)
// 			if err != nil {
// 				break
// 			}
// 			balance := &gexdb.Balance{
// 				UserID: user.TID,
// 				Asset:  gexdb.BalanceAssetMMK,
// 				Free:   decimal.NewFromFloat(10000),
// 				Status: gexdb.BalanceStatusNormal,
// 			}
// 			err = gexdb.IncreaseBalanceCall(caller, balance)
// 			if err != nil {
// 				break
// 			}
// 			for k := 0; k < orderCount; k++ {
// 				order := &gexdb.Order{
// 					Type:      gexdb.OrderTypeTopup,
// 					UserID:    user.TID,
// 					Creator:   user.TID,
// 					Quantity:  decimal.NewFromFloat(1000),
// 					Filled:    decimal.NewFromFloat(1000),
// 					InBalance: gexdb.BalanceAssetMMK,
// 					InFilled:  decimal.NewFromFloat(1000),
// 					Status:    gexdb.OrderStatusDone,
// 				}
// 				err = gexdb.CreateOrderCall(caller, order)
// 				if err != nil {
// 					break
// 				}
// 				order = &gexdb.Order{
// 					Type:       gexdb.OrderTypeWithdraw,
// 					UserID:     user.TID,
// 					Creator:    user.TID,
// 					Quantity:   decimal.NewFromFloat(100),
// 					Filled:     decimal.NewFromFloat(100),
// 					OutBalance: gexdb.BalanceAssetMMK,
// 					OutFilled:  decimal.NewFromFloat(100),
// 					Status:     gexdb.OrderStatusDone,
// 				}
// 				err = gexdb.CreateOrderCall(caller, order)
// 				if err != nil {
// 					break
// 				}
// 			}
// 			if err != nil {
// 				break
// 			}
// 			users = append(users, user)
// 		}
// 		if err != nil {
// 			break
// 		}
// 		brokers = append(brokers, broker)
// 	}
// 	return
// }

// func GenerateTestOrder(staffs, users []*gexdb.User) (err error) {
// 	xlog.Warnf("generate user is done")
// 	for i, staff := range staffs {
// 		_, err = gexapi.MatcherYWEMMK.ProcessLimit(staff.TID, gexdb.OrderSideSell, decimal.NewFromFloat(100), decimal.NewFromFloat(10))
// 		if err != nil {
// 			err = fmt.Errorf("staff sell fail with %v", err)
// 			return
// 		}
// 		for j := 0; j < 10; j++ {
// 			user := users[i*10+j]
// 			_, err = gexapi.MatcherYWEMMK.ProcessLimit(user.TID, gexdb.OrderSideBuy, decimal.NewFromFloat(10), decimal.NewFromFloat(10))
// 			if err != nil {
// 				err = fmt.Errorf("user buy staff fail with %v", err)
// 				return
// 			}
// 		}
// 	}
// 	for {
// 		depth := gexapi.MatcherYWEMMK.Depth(1)
// 		if len(depth.Asks) < 1 && len(depth.Bids) < 1 {
// 			break
// 		}
// 	}
// 	for i := 0; i < 100; i++ {
// 		j := (i + 1) % 100
// 		_, err = gexapi.MatcherYWEMMK.ProcessLimit(users[i].TID, gexdb.OrderSideSell, decimal.NewFromFloat(1.5), decimal.NewFromFloat(3.5))
// 		if err != nil {
// 			err = fmt.Errorf("user sell fail with %v", err)
// 			return
// 		}
// 		_, err = gexapi.MatcherYWEMMK.ProcessLimit(users[j].TID, gexdb.OrderSideBuy, decimal.NewFromFloat(1.5), decimal.NewFromFloat(3.5))
// 		if err != nil {
// 			err = fmt.Errorf("user buy fail with %v", err)
// 			return
// 		}
// 	}
// 	for {
// 		depth := gexapi.MatcherYWEMMK.Depth(1)
// 		if len(depth.Asks) < 1 && len(depth.Bids) < 1 {
// 			break
// 		}
// 	}
// 	for {
// 		err = gexdb.ProcSettleOrderFee()
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				err = nil
// 			}
// 			break
// 		}
// 	}
// 	return
// }

// func randomKLine(symbol string, interval time.Duration, prevPrice decimal.Decimal, startTime time.Time) (kline *gexdb.KLine) {
// 	kline = &gexdb.KLine{
// 		Symbol:     symbol,
// 		StartTime:  xsql.Time(startTime),
// 		UpdateTime: xsql.Time(startTime),
// 	}
// 	kline.Interv, _ = gexdb.IntervString(interval)
// 	kline.Count = rand.Int63n(100000) + 1
// 	kline.Open = prevPrice
// 	if rand.Int63n(2) == 0 {
// 		kline.Close = prevPrice.Add(decimal.NewFromFloat(0.05).Mul(decimal.NewFromInt(rand.Int63n(100) + 1)))
// 		kline.High = kline.Close.Mul(decimal.NewFromFloat(1.1))
// 		kline.Low = kline.Open.Mul(decimal.NewFromFloat(0.9))
// 	} else {
// 		kline.Close = prevPrice.Sub(decimal.NewFromFloat(0.05).Mul(decimal.NewFromInt(rand.Int63n(100) + 1)))
// 		kline.High = kline.Open.Mul(decimal.NewFromFloat(1.1))
// 		kline.Low = kline.Close.Mul(decimal.NewFromFloat(0.9))
// 	}
// 	kline.Amount = decimal.NewFromFloat(0.05).Mul(decimal.NewFromInt(rand.Int63n(100000) + 1))
// 	kline.Volume = kline.Amount.Mul(kline.Open.Add(kline.Close).Div(decimal.NewFromFloat(2)))
// 	return
// }

// func GenerateTestMarket(startPrice decimal.Decimal) (err error) {
// 	symbol := "YWEMMK"
// 	klines := []*gexdb.KLine{}
// 	{ //month
// 		now := time.Now()
// 		interval := 30 * 24 * time.Hour
// 		prevPrice := startPrice
// 		for i := -2; i <= 0; i++ {
// 			startTime := time.Date(now.Year(), now.Month()-time.Month(i), 0, 0, 0, 0, 0, now.Location())
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	{ //week
// 		start := xtime.TimeStartOfWeek()
// 		interval := 7 * 24 * time.Hour
// 		prevPrice := startPrice
// 		for i := -2; i <= 0; i++ {
// 			startTime := start.Add(time.Duration(i) * interval)
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	{ //day
// 		start := xtime.TimeStartOfToday()
// 		interval := 24 * time.Hour
// 		prevPrice := startPrice
// 		for i := -100; i <= 0; i++ {
// 			startTime := start.Add(time.Duration(i) * interval)
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	{ //4hour
// 		now := time.Now()
// 		start := time.Date(now.Year(), now.Month(), now.Day(), (now.Hour()/4)*4, 0, 0, 0, now.Location())
// 		interval := 4 * time.Hour
// 		prevPrice := startPrice
// 		for i := -100; i <= 0; i++ {
// 			startTime := start.Add(time.Duration(i) * interval)
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	{ //1hour
// 		now := time.Now()
// 		start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
// 		interval := time.Hour
// 		prevPrice := startPrice
// 		for i := -100; i <= 0; i++ {
// 			startTime := start.Add(time.Duration(i) * interval)
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	{ //30min
// 		now := time.Now()
// 		start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), (now.Minute()/30)*30, 0, 0, now.Location())
// 		interval := 30 * time.Minute
// 		prevPrice := startPrice
// 		for i := -100; i <= 0; i++ {
// 			startTime := start.Add(time.Duration(i) * interval)
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	{ //1min
// 		now := time.Now()
// 		start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), (now.Minute()/5)*5, 0, 0, now.Location())
// 		interval := 5 * time.Minute
// 		prevPrice := startPrice
// 		for i := -100; i <= 0; i++ {
// 			startTime := start.Add(time.Duration(i) * interval)
// 			kline := randomKLine(symbol, interval, prevPrice, startTime)
// 			klines = append(klines, kline)
// 			prevPrice = kline.Close
// 		}
// 	}
// 	xsort.SortFunc(klines, func(x, y int) bool {
// 		return klines[x].StartTime.AsTime().Before(klines[y].StartTime.AsTime())
// 	})
// 	_, err = gexdb.AddKLine(klines...)
// 	return
// }

// func GenerateTestData() (err error) {
// 	staffs, _, users, err := GenerateTestUser(10, 10, 10, 3)
// 	if err != nil {
// 		return
// 	}
// 	err = GenerateTestOrder(staffs, users)
// 	if err != nil {
// 		return
// 	}
// 	err = GenerateTestMarket(decimal.NewFromFloat(1000))
// 	return
// }
