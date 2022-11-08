package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhash"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/maker"
	"github.com/shopspring/decimal"
)

func CheckGenTestUser(ctx context.Context) (err error) {
	if os.Getenv("ENV_GEN_TEST_USER") != "1" {
		return
	}
	userCount := 100
	_, err = gexdb.FindUserByAccount(ctx, "c00")
	if err == nil {
		xlog.Infof("gen test user is skipped, created")
		return
	}
	xlog.Infof("start gen %v test user", userCount)
	caller, err := gexdb.Pool().Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = caller.Commit(ctx)
			xlog.Infof("gen %v test user is done with %v", userCount, err)
		} else {
			caller.Rollback(ctx)
			xlog.Errorf("gen %v test user is fail with %v", userCount, err)
		}
	}()
	for i := 0; i < userCount; i++ {
		user := &gexdb.User{
			Type:      gexdb.UserTypeNormal,
			Role:      gexdb.UserRoleNormal,
			Name:      converter.StringPtr(fmt.Sprintf("c%02d", i)),
			Account:   converter.StringPtr(fmt.Sprintf("c%02d", i)),
			Password:  converter.StringPtr(xhash.SHA1([]byte("123"))),
			TradePass: converter.StringPtr(xhash.SHA1([]byte("123"))),
			Status:    gexdb.UserStatusNormal,
		}
		err = gexdb.AddUserCall(caller, ctx, user)
		if err != nil {
			break
		}
		_, err = gexdb.TouchBalanceCall(caller, ctx, gexdb.BalanceAreaSpot, []string{"YWE", "GCC", "USDT"}, user.TID)
		if err != nil {
			break
		}
		_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaSpot, "YWE", decimal.NewFromFloat(100000))
		if err != nil {
			break
		}
		_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaSpot, "USDT", decimal.NewFromFloat(10000000))
		if err != nil {
			break
		}
		_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaSpot, "GCC", decimal.NewFromFloat(10000000))
		if err != nil {
			break
		}
	}
	return
}

func CheckSetupTestMaker(ctx context.Context) (err error) {
	if os.Getenv("ENV_SETUP_TEST_MAKER") != "1" {
		return
	}
	_, err = gexdb.FindUserByAccount(ctx, "maker")
	if err == nil {
		xlog.Infof("setup test maker is skipped, configured")
		return
	}
	caller, err := gexdb.Pool().Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = caller.Commit(ctx)
			xlog.Infof("setup test maker is done with %v", err)
		} else {
			caller.Rollback(ctx)
			xlog.Errorf("setup test maker is fail with %v", err)
		}
	}()
	user := &gexdb.User{
		Type:      gexdb.UserTypeNormal,
		Role:      gexdb.UserRoleMaker,
		Name:      converter.StringPtr("maker"),
		Account:   converter.StringPtr("maker"),
		Password:  converter.StringPtr(xhash.SHA1([]byte("123"))),
		TradePass: converter.StringPtr(xhash.SHA1([]byte("123"))),
		Status:    gexdb.UserStatusNormal,
	}
	err = gexdb.AddUserCall(caller, ctx, user)
	if err != nil {
		return
	}
	_, err = gexdb.TouchBalanceCall(caller, ctx, gexdb.BalanceAreaSpot, []string{"YWE", "GCC", "USDT"}, user.TID)
	if err != nil {
		return
	}
	_, err = gexdb.TouchBalanceCall(caller, ctx, gexdb.BalanceAreaFunds, []string{"GCC", "USDT"}, user.TID)
	if err != nil {
		return
	}
	_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaSpot, "YWE", decimal.NewFromFloat(100000000))
	if err != nil {
		return
	}
	_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaSpot, "USDT", decimal.NewFromFloat(200000000000))
	if err != nil {
		return
	}
	_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaSpot, "GCC", decimal.NewFromFloat(100000000000))
	if err != nil {
		return
	}
	_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaFutures, "USDT", decimal.NewFromFloat(100000000000))
	if err != nil {
		return
	}
	_, err = gexdb.ChangeBalanceCall(caller, ctx, 0, user.TID, gexdb.BalanceAreaFutures, "GCC", decimal.NewFromFloat(100000000000))
	if err != nil {
		return
	}
	newConfig := func() *maker.Config {
		config := &maker.Config{}
		config.ON = 1
		config.Delay = 500
		config.UserID = user.TID
		config.Open = decimal.NewFromFloat(1000)
		config.Close.Min = decimal.NewFromFloat(-0.01)
		config.Close.Max = decimal.NewFromFloat(0.01)
		config.Vib.Min = decimal.NewFromFloat(-0.03)
		config.Vib.Max = decimal.NewFromFloat(0.03)
		config.Vib.Count = 5
		config.Ticker = decimal.NewFromFloat(0.0001)
		config.Interval = time.Hour.Milliseconds()
		config.Depth.QtyMax = decimal.NewFromFloat(3)
		config.Depth.StepMax = 5
		config.Depth.DiffMax = decimal.NewFromFloat(2)
		config.Depth.DiffMin = decimal.NewFromFloat(0.02)
		config.Depth.Max = 15
		return config
	}
	{
		config := newConfig()
		config.Symbol = "spot.YWEUSDT"
		err = config.Valid()
		if err != nil {
			return
		}
		key := fmt.Sprintf("maker-%v", config.Symbol)
		err = basedb.StoreConfCall(caller, ctx, key, converter.JSON(config))
		if err != nil {
			return
		}
		xlog.Infof("setup test %v maker", config.Symbol)
	}
	{
		config := newConfig()
		config.Symbol = "spot.GCCUSDT"
		config.Delay = 1000
		config.Open = decimal.NewFromFloat(1)
		config.Close.Min = decimal.NewFromFloat(-0.001)
		config.Close.Max = decimal.NewFromFloat(0.001)
		config.Vib.Min = decimal.NewFromFloat(-0.01)
		config.Vib.Max = decimal.NewFromFloat(0.01)
		config.Depth.QtyMax = decimal.NewFromFloat(1000)
		config.Depth.DiffMax = decimal.NewFromFloat(0.05)
		config.Depth.DiffMin = decimal.NewFromFloat(0.01)
		err = config.Valid()
		if err != nil {
			return
		}
		key := fmt.Sprintf("maker-%v", config.Symbol)
		err = basedb.StoreConfCall(caller, ctx, key, converter.JSON(config))
		if err != nil {
			return
		}
		xlog.Infof("setup test %v maker", config.Symbol)
	}
	{
		config := newConfig()
		config.Symbol = "futures.YWEGCC"
		err = config.Valid()
		if err != nil {
			return
		}
		key := fmt.Sprintf("maker-%v", config.Symbol)
		err = basedb.StoreConfCall(caller, ctx, key, converter.JSON(config))
		if err != nil {
			return
		}
		xlog.Infof("setup test %v maker", config.Symbol)
	}
	{
		config := newConfig()
		config.Symbol = "futures.YWEUSDT"
		err = config.Valid()
		if err != nil {
			return
		}
		key := fmt.Sprintf("maker-%v", config.Symbol)
		err = basedb.StoreConfCall(caller, ctx, key, converter.JSON(config))
		if err != nil {
			return
		}
		xlog.Infof("setup test %v maker", config.Symbol)
	}
	return
}

func BootstrapTest(ctx context.Context) (err error) {
	err = CheckSetupTestMaker(ctx)
	if err != nil {
		return
	}
	err = CheckGenTestUser(ctx)
	if err != nil {
		return
	}
	return
}
