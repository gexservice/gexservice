//auto gen func by autogen
package gexdb

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/codingeasygo/crud"
)

func TestAutoBalance(t *testing.T) {
	var err error
	for _, value := range BalanceAreaAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if BalanceAreaAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if BalanceAreaAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(BalanceAreaAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(BalanceAreaAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range BalanceStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if BalanceStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if BalanceStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(BalanceStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(BalanceStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithBalance()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	balance := &Balance{}
	balance.Valid()

	table, fields := balance.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(balance.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddBalance(context.Background(), balance)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(balance.TID).IsZero() {
		t.Error("not id")
		return
	}
	balance.Valid()
	err = UpdateBalanceFilter(context.Background(), balance, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateBalanceWheref(context.Background(), balance, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateBalanceFilterWheref(context.Background(), balance, BalanceFilterUpdate, "tid=$%v", balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findBalance, err := FindBalance(context.Background(), balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balance.TID != findBalance.TID {
		t.Error("find id error")
		return
	}
	findBalance, err = FindBalanceWheref(context.Background(), "tid=$%v", balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balance.TID != findBalance.TID {
		t.Error("find id error")
		return
	}
	findBalance, err = FindBalanceFilterWheref(context.Background(), "#all", "tid=$%v", balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balance.TID != findBalance.TID {
		t.Error("find id error")
		return
	}
	findBalance, err = FindBalanceWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{balance.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if balance.TID != findBalance.TID {
		t.Error("find id error")
		return
	}
	findBalance, err = FindBalanceWherefCall(GetQueryer, context.Background(), true, "tid=$%v", balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balance.TID != findBalance.TID {
		t.Error("find id error")
		return
	}
	balanceList, balanceMap, err := ListBalanceByID(context.Background())
	if err != nil || len(balanceList) > 0 || balanceMap == nil || len(balanceMap) > 0 {
		t.Error(err)
		return
	}
	balanceList, balanceMap, err = ListBalanceByID(context.Background(), balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceList) != 1 || balanceList[0].TID != balance.TID || len(balanceMap) != 1 || balanceMap[balance.TID] == nil || balanceMap[balance.TID].TID != balance.TID {
		t.Error("list id error")
		return
	}
	balanceList, balanceMap, err = ListBalanceFilterByID(context.Background(), "#all")
	if err != nil || len(balanceList) > 0 || balanceMap == nil || len(balanceMap) > 0 {
		t.Error(err)
		return
	}
	balanceList, balanceMap, err = ListBalanceFilterByID(context.Background(), "#all", balance.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceList) != 1 || balanceList[0].TID != balance.TID || len(balanceMap) != 1 || balanceMap[balance.TID] == nil || balanceMap[balance.TID].TID != balance.TID {
		t.Error("list id error")
		return
	}
	balanceList = nil
	balanceMap = nil
	err = ScanBalanceByID(context.Background(), []int64{balance.TID}, &balanceList, &balanceMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceList) != 1 || balanceList[0].TID != balance.TID || len(balanceMap) != 1 || balanceMap[balance.TID] == nil || balanceMap[balance.TID].TID != balance.TID {
		t.Error("list id error")
		return
	}
	balanceList = nil
	balanceMap = nil
	err = ScanBalanceFilterByID(context.Background(), "#all", []int64{balance.TID}, &balanceList, &balanceMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceList) != 1 || balanceList[0].TID != balance.TID || len(balanceMap) != 1 || balanceMap[balance.TID] == nil || balanceMap[balance.TID].TID != balance.TID {
		t.Error("list id error")
		return
	}
	balanceList = nil
	balanceMap = nil
	err = ScanBalanceWheref(context.Background(), "tid=$%v", []interface{}{balance.TID}, "", &balanceList, &balanceMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceList) != 1 || balanceList[0].TID != balance.TID || len(balanceMap) != 1 || balanceMap[balance.TID] == nil || balanceMap[balance.TID].TID != balance.TID {
		t.Error("list id error")
		return
	}
	balanceList = nil
	balanceMap = nil
	err = ScanBalanceFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{balance.TID}, "", &balanceList, &balanceMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceList) != 1 || balanceList[0].TID != balance.TID || len(balanceMap) != 1 || balanceMap[balance.TID] == nil || balanceMap[balance.TID].TID != balance.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoBalanceHistory(t *testing.T) {
	var err error
	for _, value := range BalanceHistoryStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if BalanceHistoryStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if BalanceHistoryStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(BalanceHistoryStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(BalanceHistoryStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithBalanceHistory()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	balanceHistory := &BalanceHistory{}
	balanceHistory.Valid()

	table, fields := balanceHistory.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(balanceHistory.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddBalanceHistory(context.Background(), balanceHistory)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(balanceHistory.TID).IsZero() {
		t.Error("not id")
		return
	}
	balanceHistory.Valid()
	err = UpdateBalanceHistoryFilter(context.Background(), balanceHistory, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateBalanceHistoryWheref(context.Background(), balanceHistory, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateBalanceHistoryFilterWheref(context.Background(), balanceHistory, BalanceHistoryFilterUpdate, "tid=$%v", balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findBalanceHistory, err := FindBalanceHistory(context.Background(), balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceHistory.TID != findBalanceHistory.TID {
		t.Error("find id error")
		return
	}
	findBalanceHistory, err = FindBalanceHistoryWheref(context.Background(), "tid=$%v", balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceHistory.TID != findBalanceHistory.TID {
		t.Error("find id error")
		return
	}
	findBalanceHistory, err = FindBalanceHistoryFilterWheref(context.Background(), "#all", "tid=$%v", balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceHistory.TID != findBalanceHistory.TID {
		t.Error("find id error")
		return
	}
	findBalanceHistory, err = FindBalanceHistoryWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{balanceHistory.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if balanceHistory.TID != findBalanceHistory.TID {
		t.Error("find id error")
		return
	}
	findBalanceHistory, err = FindBalanceHistoryWherefCall(GetQueryer, context.Background(), true, "tid=$%v", balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceHistory.TID != findBalanceHistory.TID {
		t.Error("find id error")
		return
	}
	balanceHistoryList, balanceHistoryMap, err := ListBalanceHistoryByID(context.Background())
	if err != nil || len(balanceHistoryList) > 0 || balanceHistoryMap == nil || len(balanceHistoryMap) > 0 {
		t.Error(err)
		return
	}
	balanceHistoryList, balanceHistoryMap, err = ListBalanceHistoryByID(context.Background(), balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceHistoryList) != 1 || balanceHistoryList[0].TID != balanceHistory.TID || len(balanceHistoryMap) != 1 || balanceHistoryMap[balanceHistory.TID] == nil || balanceHistoryMap[balanceHistory.TID].TID != balanceHistory.TID {
		t.Error("list id error")
		return
	}
	balanceHistoryList, balanceHistoryMap, err = ListBalanceHistoryFilterByID(context.Background(), "#all")
	if err != nil || len(balanceHistoryList) > 0 || balanceHistoryMap == nil || len(balanceHistoryMap) > 0 {
		t.Error(err)
		return
	}
	balanceHistoryList, balanceHistoryMap, err = ListBalanceHistoryFilterByID(context.Background(), "#all", balanceHistory.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceHistoryList) != 1 || balanceHistoryList[0].TID != balanceHistory.TID || len(balanceHistoryMap) != 1 || balanceHistoryMap[balanceHistory.TID] == nil || balanceHistoryMap[balanceHistory.TID].TID != balanceHistory.TID {
		t.Error("list id error")
		return
	}
	balanceHistoryList = nil
	balanceHistoryMap = nil
	err = ScanBalanceHistoryByID(context.Background(), []int64{balanceHistory.TID}, &balanceHistoryList, &balanceHistoryMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceHistoryList) != 1 || balanceHistoryList[0].TID != balanceHistory.TID || len(balanceHistoryMap) != 1 || balanceHistoryMap[balanceHistory.TID] == nil || balanceHistoryMap[balanceHistory.TID].TID != balanceHistory.TID {
		t.Error("list id error")
		return
	}
	balanceHistoryList = nil
	balanceHistoryMap = nil
	err = ScanBalanceHistoryFilterByID(context.Background(), "#all", []int64{balanceHistory.TID}, &balanceHistoryList, &balanceHistoryMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceHistoryList) != 1 || balanceHistoryList[0].TID != balanceHistory.TID || len(balanceHistoryMap) != 1 || balanceHistoryMap[balanceHistory.TID] == nil || balanceHistoryMap[balanceHistory.TID].TID != balanceHistory.TID {
		t.Error("list id error")
		return
	}
	balanceHistoryList = nil
	balanceHistoryMap = nil
	err = ScanBalanceHistoryWheref(context.Background(), "tid=$%v", []interface{}{balanceHistory.TID}, "", &balanceHistoryList, &balanceHistoryMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceHistoryList) != 1 || balanceHistoryList[0].TID != balanceHistory.TID || len(balanceHistoryMap) != 1 || balanceHistoryMap[balanceHistory.TID] == nil || balanceHistoryMap[balanceHistory.TID].TID != balanceHistory.TID {
		t.Error("list id error")
		return
	}
	balanceHistoryList = nil
	balanceHistoryMap = nil
	err = ScanBalanceHistoryFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{balanceHistory.TID}, "", &balanceHistoryList, &balanceHistoryMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceHistoryList) != 1 || balanceHistoryList[0].TID != balanceHistory.TID || len(balanceHistoryMap) != 1 || balanceHistoryMap[balanceHistory.TID] == nil || balanceHistoryMap[balanceHistory.TID].TID != balanceHistory.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoBalanceRecord(t *testing.T) {
	var err error
	for _, value := range BalanceRecordTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if BalanceRecordTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if BalanceRecordTypeAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(BalanceRecordTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(BalanceRecordTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range BalanceRecordStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if BalanceRecordStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if BalanceRecordStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(BalanceRecordStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(BalanceRecordStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithBalanceRecord()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	balanceRecord := &BalanceRecord{}
	balanceRecord.Valid()

	table, fields := balanceRecord.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(balanceRecord.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = balanceRecord.Insert(GetQueryer, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(balanceRecord.TID).IsZero() {
		t.Error("not id")
		return
	}
	balanceRecord.Valid()
	err = UpdateBalanceRecordFilter(context.Background(), balanceRecord, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateBalanceRecordWheref(context.Background(), balanceRecord, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateBalanceRecordFilterWheref(context.Background(), balanceRecord, BalanceRecordFilterUpdate, "tid=$%v", balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findBalanceRecord, err := FindBalanceRecord(context.Background(), balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceRecord.TID != findBalanceRecord.TID {
		t.Error("find id error")
		return
	}
	findBalanceRecord, err = FindBalanceRecordWheref(context.Background(), "tid=$%v", balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceRecord.TID != findBalanceRecord.TID {
		t.Error("find id error")
		return
	}
	findBalanceRecord, err = FindBalanceRecordFilterWheref(context.Background(), "#all", "tid=$%v", balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceRecord.TID != findBalanceRecord.TID {
		t.Error("find id error")
		return
	}
	findBalanceRecord, err = FindBalanceRecordWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{balanceRecord.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if balanceRecord.TID != findBalanceRecord.TID {
		t.Error("find id error")
		return
	}
	findBalanceRecord, err = FindBalanceRecordWherefCall(GetQueryer, context.Background(), true, "tid=$%v", balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if balanceRecord.TID != findBalanceRecord.TID {
		t.Error("find id error")
		return
	}
	balanceRecordList, balanceRecordMap, err := ListBalanceRecordByID(context.Background())
	if err != nil || len(balanceRecordList) > 0 || balanceRecordMap == nil || len(balanceRecordMap) > 0 {
		t.Error(err)
		return
	}
	balanceRecordList, balanceRecordMap, err = ListBalanceRecordByID(context.Background(), balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceRecordList) != 1 || balanceRecordList[0].TID != balanceRecord.TID || len(balanceRecordMap) != 1 || balanceRecordMap[balanceRecord.TID] == nil || balanceRecordMap[balanceRecord.TID].TID != balanceRecord.TID {
		t.Error("list id error")
		return
	}
	balanceRecordList, balanceRecordMap, err = ListBalanceRecordFilterByID(context.Background(), "#all")
	if err != nil || len(balanceRecordList) > 0 || balanceRecordMap == nil || len(balanceRecordMap) > 0 {
		t.Error(err)
		return
	}
	balanceRecordList, balanceRecordMap, err = ListBalanceRecordFilterByID(context.Background(), "#all", balanceRecord.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceRecordList) != 1 || balanceRecordList[0].TID != balanceRecord.TID || len(balanceRecordMap) != 1 || balanceRecordMap[balanceRecord.TID] == nil || balanceRecordMap[balanceRecord.TID].TID != balanceRecord.TID {
		t.Error("list id error")
		return
	}
	balanceRecordList = nil
	balanceRecordMap = nil
	err = ScanBalanceRecordByID(context.Background(), []int64{balanceRecord.TID}, &balanceRecordList, &balanceRecordMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceRecordList) != 1 || balanceRecordList[0].TID != balanceRecord.TID || len(balanceRecordMap) != 1 || balanceRecordMap[balanceRecord.TID] == nil || balanceRecordMap[balanceRecord.TID].TID != balanceRecord.TID {
		t.Error("list id error")
		return
	}
	balanceRecordList = nil
	balanceRecordMap = nil
	err = ScanBalanceRecordFilterByID(context.Background(), "#all", []int64{balanceRecord.TID}, &balanceRecordList, &balanceRecordMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceRecordList) != 1 || balanceRecordList[0].TID != balanceRecord.TID || len(balanceRecordMap) != 1 || balanceRecordMap[balanceRecord.TID] == nil || balanceRecordMap[balanceRecord.TID].TID != balanceRecord.TID {
		t.Error("list id error")
		return
	}
	balanceRecordList = nil
	balanceRecordMap = nil
	err = ScanBalanceRecordWheref(context.Background(), "tid=$%v", []interface{}{balanceRecord.TID}, "", &balanceRecordList, &balanceRecordMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceRecordList) != 1 || balanceRecordList[0].TID != balanceRecord.TID || len(balanceRecordMap) != 1 || balanceRecordMap[balanceRecord.TID] == nil || balanceRecordMap[balanceRecord.TID].TID != balanceRecord.TID {
		t.Error("list id error")
		return
	}
	balanceRecordList = nil
	balanceRecordMap = nil
	err = ScanBalanceRecordFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{balanceRecord.TID}, "", &balanceRecordList, &balanceRecordMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(balanceRecordList) != 1 || balanceRecordList[0].TID != balanceRecord.TID || len(balanceRecordMap) != 1 || balanceRecordMap[balanceRecord.TID] == nil || balanceRecordMap[balanceRecord.TID].TID != balanceRecord.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoHolding(t *testing.T) {
	var err error
	for _, value := range HoldingStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if HoldingStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if HoldingStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(HoldingStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(HoldingStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithHolding()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	holding := &Holding{}
	holding.Valid()

	table, fields := holding.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(holding.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = holding.Insert(GetQueryer, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(holding.TID).IsZero() {
		t.Error("not id")
		return
	}
	holding.Valid()
	err = UpdateHoldingFilter(context.Background(), holding, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateHoldingWheref(context.Background(), holding, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateHoldingFilterWheref(context.Background(), holding, HoldingFilterUpdate, "tid=$%v", holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findHolding, err := FindHolding(context.Background(), holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if holding.TID != findHolding.TID {
		t.Error("find id error")
		return
	}
	findHolding, err = FindHoldingWheref(context.Background(), "tid=$%v", holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if holding.TID != findHolding.TID {
		t.Error("find id error")
		return
	}
	findHolding, err = FindHoldingFilterWheref(context.Background(), "#all", "tid=$%v", holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if holding.TID != findHolding.TID {
		t.Error("find id error")
		return
	}
	findHolding, err = FindHoldingWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{holding.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if holding.TID != findHolding.TID {
		t.Error("find id error")
		return
	}
	findHolding, err = FindHoldingWherefCall(GetQueryer, context.Background(), true, "tid=$%v", holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if holding.TID != findHolding.TID {
		t.Error("find id error")
		return
	}
	holdingList, holdingMap, err := ListHoldingByID(context.Background())
	if err != nil || len(holdingList) > 0 || holdingMap == nil || len(holdingMap) > 0 {
		t.Error(err)
		return
	}
	holdingList, holdingMap, err = ListHoldingByID(context.Background(), holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(holdingList) != 1 || holdingList[0].TID != holding.TID || len(holdingMap) != 1 || holdingMap[holding.TID] == nil || holdingMap[holding.TID].TID != holding.TID {
		t.Error("list id error")
		return
	}
	holdingList, holdingMap, err = ListHoldingFilterByID(context.Background(), "#all")
	if err != nil || len(holdingList) > 0 || holdingMap == nil || len(holdingMap) > 0 {
		t.Error(err)
		return
	}
	holdingList, holdingMap, err = ListHoldingFilterByID(context.Background(), "#all", holding.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(holdingList) != 1 || holdingList[0].TID != holding.TID || len(holdingMap) != 1 || holdingMap[holding.TID] == nil || holdingMap[holding.TID].TID != holding.TID {
		t.Error("list id error")
		return
	}
	holdingList = nil
	holdingMap = nil
	err = ScanHoldingByID(context.Background(), []int64{holding.TID}, &holdingList, &holdingMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(holdingList) != 1 || holdingList[0].TID != holding.TID || len(holdingMap) != 1 || holdingMap[holding.TID] == nil || holdingMap[holding.TID].TID != holding.TID {
		t.Error("list id error")
		return
	}
	holdingList = nil
	holdingMap = nil
	err = ScanHoldingFilterByID(context.Background(), "#all", []int64{holding.TID}, &holdingList, &holdingMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(holdingList) != 1 || holdingList[0].TID != holding.TID || len(holdingMap) != 1 || holdingMap[holding.TID] == nil || holdingMap[holding.TID].TID != holding.TID {
		t.Error("list id error")
		return
	}
	holdingList = nil
	holdingMap = nil
	err = ScanHoldingWheref(context.Background(), "tid=$%v", []interface{}{holding.TID}, "", &holdingList, &holdingMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(holdingList) != 1 || holdingList[0].TID != holding.TID || len(holdingMap) != 1 || holdingMap[holding.TID] == nil || holdingMap[holding.TID].TID != holding.TID {
		t.Error("list id error")
		return
	}
	holdingList = nil
	holdingMap = nil
	err = ScanHoldingFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{holding.TID}, "", &holdingList, &holdingMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(holdingList) != 1 || holdingList[0].TID != holding.TID || len(holdingMap) != 1 || holdingMap[holding.TID] == nil || holdingMap[holding.TID].TID != holding.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoKLine(t *testing.T) {
	var err error
	metav := MetaWithKLine()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	kLine := &KLine{}
	kLine.Valid()

	table, fields := kLine.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(kLine.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddKLine(context.Background(), kLine)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(kLine.TID).IsZero() {
		t.Error("not id")
		return
	}
	kLine.Valid()
	err = UpdateKLineFilter(context.Background(), kLine, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateKLineWheref(context.Background(), kLine, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateKLineFilterWheref(context.Background(), kLine, KLineFilterUpdate, "tid=$%v", kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findKLine, err := FindKLine(context.Background(), kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if kLine.TID != findKLine.TID {
		t.Error("find id error")
		return
	}
	findKLine, err = FindKLineWheref(context.Background(), "tid=$%v", kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if kLine.TID != findKLine.TID {
		t.Error("find id error")
		return
	}
	findKLine, err = FindKLineFilterWheref(context.Background(), "#all", "tid=$%v", kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if kLine.TID != findKLine.TID {
		t.Error("find id error")
		return
	}
	findKLine, err = FindKLineWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{kLine.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if kLine.TID != findKLine.TID {
		t.Error("find id error")
		return
	}
	findKLine, err = FindKLineWherefCall(GetQueryer, context.Background(), true, "tid=$%v", kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if kLine.TID != findKLine.TID {
		t.Error("find id error")
		return
	}
	kLineList, kLineMap, err := ListKLineByID(context.Background())
	if err != nil || len(kLineList) > 0 || kLineMap == nil || len(kLineMap) > 0 {
		t.Error(err)
		return
	}
	kLineList, kLineMap, err = ListKLineByID(context.Background(), kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(kLineList) != 1 || kLineList[0].TID != kLine.TID || len(kLineMap) != 1 || kLineMap[kLine.TID] == nil || kLineMap[kLine.TID].TID != kLine.TID {
		t.Error("list id error")
		return
	}
	kLineList, kLineMap, err = ListKLineFilterByID(context.Background(), "#all")
	if err != nil || len(kLineList) > 0 || kLineMap == nil || len(kLineMap) > 0 {
		t.Error(err)
		return
	}
	kLineList, kLineMap, err = ListKLineFilterByID(context.Background(), "#all", kLine.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(kLineList) != 1 || kLineList[0].TID != kLine.TID || len(kLineMap) != 1 || kLineMap[kLine.TID] == nil || kLineMap[kLine.TID].TID != kLine.TID {
		t.Error("list id error")
		return
	}
	kLineList = nil
	kLineMap = nil
	err = ScanKLineByID(context.Background(), []int64{kLine.TID}, &kLineList, &kLineMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(kLineList) != 1 || kLineList[0].TID != kLine.TID || len(kLineMap) != 1 || kLineMap[kLine.TID] == nil || kLineMap[kLine.TID].TID != kLine.TID {
		t.Error("list id error")
		return
	}
	kLineList = nil
	kLineMap = nil
	err = ScanKLineFilterByID(context.Background(), "#all", []int64{kLine.TID}, &kLineList, &kLineMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(kLineList) != 1 || kLineList[0].TID != kLine.TID || len(kLineMap) != 1 || kLineMap[kLine.TID] == nil || kLineMap[kLine.TID].TID != kLine.TID {
		t.Error("list id error")
		return
	}
	kLineList = nil
	kLineMap = nil
	err = ScanKLineWheref(context.Background(), "tid=$%v", []interface{}{kLine.TID}, "", &kLineList, &kLineMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(kLineList) != 1 || kLineList[0].TID != kLine.TID || len(kLineMap) != 1 || kLineMap[kLine.TID] == nil || kLineMap[kLine.TID].TID != kLine.TID {
		t.Error("list id error")
		return
	}
	kLineList = nil
	kLineMap = nil
	err = ScanKLineFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{kLine.TID}, "", &kLineList, &kLineMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(kLineList) != 1 || kLineList[0].TID != kLine.TID || len(kLineMap) != 1 || kLineMap[kLine.TID] == nil || kLineMap[kLine.TID].TID != kLine.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoOrder(t *testing.T) {
	var err error
	for _, value := range OrderTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if OrderTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if OrderTypeAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(OrderTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(OrderTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range OrderSideAll {
		if value.EnumValid(string(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(string("this should invalid")) == nil {
			t.Error("not enum valid")
			return
		}
		if OrderSideAll.EnumValid(string(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if OrderSideAll.EnumValid(string("this should invalid")) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(OrderSideAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(OrderSideAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range OrderTriggerTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if OrderTriggerTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if OrderTriggerTypeAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(OrderTriggerTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(OrderTriggerTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range OrderStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if OrderStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if OrderStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(OrderStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(OrderStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithOrder()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	order := &Order{}
	order.Valid()

	table, fields := order.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(order.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddOrder(context.Background(), order)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(order.TID).IsZero() {
		t.Error("not id")
		return
	}
	order.Valid()
	err = UpdateOrderFilter(context.Background(), order, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateOrderWheref(context.Background(), order, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateOrderFilterWheref(context.Background(), order, OrderFilterUpdate, "tid=$%v", order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findOrder, err := FindOrder(context.Background(), order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if order.TID != findOrder.TID {
		t.Error("find id error")
		return
	}
	findOrder, err = FindOrderWheref(context.Background(), "tid=$%v", order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if order.TID != findOrder.TID {
		t.Error("find id error")
		return
	}
	findOrder, err = FindOrderFilterWheref(context.Background(), "#all", "tid=$%v", order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if order.TID != findOrder.TID {
		t.Error("find id error")
		return
	}
	findOrder, err = FindOrderWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{order.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if order.TID != findOrder.TID {
		t.Error("find id error")
		return
	}
	findOrder, err = FindOrderWherefCall(GetQueryer, context.Background(), true, "tid=$%v", order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if order.TID != findOrder.TID {
		t.Error("find id error")
		return
	}
	orderList, orderMap, err := ListOrderByID(context.Background())
	if err != nil || len(orderList) > 0 || orderMap == nil || len(orderMap) > 0 {
		t.Error(err)
		return
	}
	orderList, orderMap, err = ListOrderByID(context.Background(), order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderList) != 1 || orderList[0].TID != order.TID || len(orderMap) != 1 || orderMap[order.TID] == nil || orderMap[order.TID].TID != order.TID {
		t.Error("list id error")
		return
	}
	orderList, orderMap, err = ListOrderFilterByID(context.Background(), "#all")
	if err != nil || len(orderList) > 0 || orderMap == nil || len(orderMap) > 0 {
		t.Error(err)
		return
	}
	orderList, orderMap, err = ListOrderFilterByID(context.Background(), "#all", order.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderList) != 1 || orderList[0].TID != order.TID || len(orderMap) != 1 || orderMap[order.TID] == nil || orderMap[order.TID].TID != order.TID {
		t.Error("list id error")
		return
	}
	orderList = nil
	orderMap = nil
	err = ScanOrderByID(context.Background(), []int64{order.TID}, &orderList, &orderMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderList) != 1 || orderList[0].TID != order.TID || len(orderMap) != 1 || orderMap[order.TID] == nil || orderMap[order.TID].TID != order.TID {
		t.Error("list id error")
		return
	}
	orderList = nil
	orderMap = nil
	err = ScanOrderFilterByID(context.Background(), "#all", []int64{order.TID}, &orderList, &orderMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderList) != 1 || orderList[0].TID != order.TID || len(orderMap) != 1 || orderMap[order.TID] == nil || orderMap[order.TID].TID != order.TID {
		t.Error("list id error")
		return
	}
	orderList = nil
	orderMap = nil
	err = ScanOrderWheref(context.Background(), "tid=$%v", []interface{}{order.TID}, "", &orderList, &orderMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderList) != 1 || orderList[0].TID != order.TID || len(orderMap) != 1 || orderMap[order.TID] == nil || orderMap[order.TID].TID != order.TID {
		t.Error("list id error")
		return
	}
	orderList = nil
	orderMap = nil
	err = ScanOrderFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{order.TID}, "", &orderList, &orderMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderList) != 1 || orderList[0].TID != order.TID || len(orderMap) != 1 || orderMap[order.TID] == nil || orderMap[order.TID].TID != order.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoOrderComm(t *testing.T) {
	var err error
	for _, value := range OrderCommTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if OrderCommTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if OrderCommTypeAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(OrderCommTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(OrderCommTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range OrderCommStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if OrderCommStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if OrderCommStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(OrderCommStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(OrderCommStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithOrderComm()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	orderComm := &OrderComm{}
	orderComm.Valid()

	table, fields := orderComm.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(orderComm.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddOrderComm(context.Background(), orderComm)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(orderComm.TID).IsZero() {
		t.Error("not id")
		return
	}
	orderComm.Valid()
	err = UpdateOrderCommFilter(context.Background(), orderComm, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateOrderCommWheref(context.Background(), orderComm, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateOrderCommFilterWheref(context.Background(), orderComm, OrderCommFilterUpdate, "tid=$%v", orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findOrderComm, err := FindOrderComm(context.Background(), orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if orderComm.TID != findOrderComm.TID {
		t.Error("find id error")
		return
	}
	findOrderComm, err = FindOrderCommWheref(context.Background(), "tid=$%v", orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if orderComm.TID != findOrderComm.TID {
		t.Error("find id error")
		return
	}
	findOrderComm, err = FindOrderCommFilterWheref(context.Background(), "#all", "tid=$%v", orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if orderComm.TID != findOrderComm.TID {
		t.Error("find id error")
		return
	}
	findOrderComm, err = FindOrderCommWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{orderComm.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if orderComm.TID != findOrderComm.TID {
		t.Error("find id error")
		return
	}
	findOrderComm, err = FindOrderCommWherefCall(GetQueryer, context.Background(), true, "tid=$%v", orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if orderComm.TID != findOrderComm.TID {
		t.Error("find id error")
		return
	}
	orderCommList, orderCommMap, err := ListOrderCommByID(context.Background())
	if err != nil || len(orderCommList) > 0 || orderCommMap == nil || len(orderCommMap) > 0 {
		t.Error(err)
		return
	}
	orderCommList, orderCommMap, err = ListOrderCommByID(context.Background(), orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderCommList) != 1 || orderCommList[0].TID != orderComm.TID || len(orderCommMap) != 1 || orderCommMap[orderComm.TID] == nil || orderCommMap[orderComm.TID].TID != orderComm.TID {
		t.Error("list id error")
		return
	}
	orderCommList, orderCommMap, err = ListOrderCommFilterByID(context.Background(), "#all")
	if err != nil || len(orderCommList) > 0 || orderCommMap == nil || len(orderCommMap) > 0 {
		t.Error(err)
		return
	}
	orderCommList, orderCommMap, err = ListOrderCommFilterByID(context.Background(), "#all", orderComm.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderCommList) != 1 || orderCommList[0].TID != orderComm.TID || len(orderCommMap) != 1 || orderCommMap[orderComm.TID] == nil || orderCommMap[orderComm.TID].TID != orderComm.TID {
		t.Error("list id error")
		return
	}
	orderCommList = nil
	orderCommMap = nil
	err = ScanOrderCommByID(context.Background(), []int64{orderComm.TID}, &orderCommList, &orderCommMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderCommList) != 1 || orderCommList[0].TID != orderComm.TID || len(orderCommMap) != 1 || orderCommMap[orderComm.TID] == nil || orderCommMap[orderComm.TID].TID != orderComm.TID {
		t.Error("list id error")
		return
	}
	orderCommList = nil
	orderCommMap = nil
	err = ScanOrderCommFilterByID(context.Background(), "#all", []int64{orderComm.TID}, &orderCommList, &orderCommMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderCommList) != 1 || orderCommList[0].TID != orderComm.TID || len(orderCommMap) != 1 || orderCommMap[orderComm.TID] == nil || orderCommMap[orderComm.TID].TID != orderComm.TID {
		t.Error("list id error")
		return
	}
	orderCommList = nil
	orderCommMap = nil
	err = ScanOrderCommWheref(context.Background(), "tid=$%v", []interface{}{orderComm.TID}, "", &orderCommList, &orderCommMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderCommList) != 1 || orderCommList[0].TID != orderComm.TID || len(orderCommMap) != 1 || orderCommMap[orderComm.TID] == nil || orderCommMap[orderComm.TID].TID != orderComm.TID {
		t.Error("list id error")
		return
	}
	orderCommList = nil
	orderCommMap = nil
	err = ScanOrderCommFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{orderComm.TID}, "", &orderCommList, &orderCommMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(orderCommList) != 1 || orderCommList[0].TID != orderComm.TID || len(orderCommMap) != 1 || orderCommMap[orderComm.TID] == nil || orderCommMap[orderComm.TID].TID != orderComm.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoUser(t *testing.T) {
	var err error
	for _, value := range UserTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if UserTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if UserTypeAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(UserTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(UserTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range UserRoleAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if UserRoleAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if UserRoleAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(UserRoleAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(UserRoleAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range UserStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if UserStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if UserStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(UserStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(UserStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithUser()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	user := &User{}
	user.Valid()

	table, fields := user.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(user.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddUser(context.Background(), user)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(user.TID).IsZero() {
		t.Error("not id")
		return
	}
	user.Valid()
	err = UpdateUserFilter(context.Background(), user, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateUserWheref(context.Background(), user, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateUserFilterWheref(context.Background(), user, UserFilterUpdate, "tid=$%v", user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findUser, err := FindUser(context.Background(), user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if user.TID != findUser.TID {
		t.Error("find id error")
		return
	}
	findUser, err = FindUserWheref(context.Background(), "tid=$%v", user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if user.TID != findUser.TID {
		t.Error("find id error")
		return
	}
	findUser, err = FindUserFilterWheref(context.Background(), "#all", "tid=$%v", user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if user.TID != findUser.TID {
		t.Error("find id error")
		return
	}
	findUser, err = FindUserWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{user.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if user.TID != findUser.TID {
		t.Error("find id error")
		return
	}
	findUser, err = FindUserWherefCall(GetQueryer, context.Background(), true, "tid=$%v", user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if user.TID != findUser.TID {
		t.Error("find id error")
		return
	}
	userList, userMap, err := ListUserByID(context.Background())
	if err != nil || len(userList) > 0 || userMap == nil || len(userMap) > 0 {
		t.Error(err)
		return
	}
	userList, userMap, err = ListUserByID(context.Background(), user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(userList) != 1 || userList[0].TID != user.TID || len(userMap) != 1 || userMap[user.TID] == nil || userMap[user.TID].TID != user.TID {
		t.Error("list id error")
		return
	}
	userList, userMap, err = ListUserFilterByID(context.Background(), "#all")
	if err != nil || len(userList) > 0 || userMap == nil || len(userMap) > 0 {
		t.Error(err)
		return
	}
	userList, userMap, err = ListUserFilterByID(context.Background(), "#all", user.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(userList) != 1 || userList[0].TID != user.TID || len(userMap) != 1 || userMap[user.TID] == nil || userMap[user.TID].TID != user.TID {
		t.Error("list id error")
		return
	}
	userList = nil
	userMap = nil
	err = ScanUserByID(context.Background(), []int64{user.TID}, &userList, &userMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(userList) != 1 || userList[0].TID != user.TID || len(userMap) != 1 || userMap[user.TID] == nil || userMap[user.TID].TID != user.TID {
		t.Error("list id error")
		return
	}
	userList = nil
	userMap = nil
	err = ScanUserFilterByID(context.Background(), "#all", []int64{user.TID}, &userList, &userMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(userList) != 1 || userList[0].TID != user.TID || len(userMap) != 1 || userMap[user.TID] == nil || userMap[user.TID].TID != user.TID {
		t.Error("list id error")
		return
	}
	userList = nil
	userMap = nil
	err = ScanUserWheref(context.Background(), "tid=$%v", []interface{}{user.TID}, "", &userList, &userMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(userList) != 1 || userList[0].TID != user.TID || len(userMap) != 1 || userMap[user.TID] == nil || userMap[user.TID].TID != user.TID {
		t.Error("list id error")
		return
	}
	userList = nil
	userMap = nil
	err = ScanUserFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{user.TID}, "", &userList, &userMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(userList) != 1 || userList[0].TID != user.TID || len(userMap) != 1 || userMap[user.TID] == nil || userMap[user.TID].TID != user.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoWithdraw(t *testing.T) {
	var err error
	for _, value := range WithdrawTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if WithdrawTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if WithdrawTypeAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(WithdrawTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(WithdrawTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range WithdrawStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
		if WithdrawStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if WithdrawStatusAll.EnumValid(int(-321654)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(WithdrawStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(WithdrawStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithWithdraw()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	withdraw := &Withdraw{}
	withdraw.Valid()

	table, fields := withdraw.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(withdraw.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddWithdraw(context.Background(), withdraw)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(withdraw.TID).IsZero() {
		t.Error("not id")
		return
	}
	withdraw.Valid()
	err = UpdateWithdrawFilter(context.Background(), withdraw, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateWithdrawWheref(context.Background(), withdraw, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateWithdrawFilterWheref(context.Background(), withdraw, WithdrawFilterUpdate, "tid=$%v", withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findWithdraw, err := FindWithdraw(context.Background(), withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if withdraw.TID != findWithdraw.TID {
		t.Error("find id error")
		return
	}
	findWithdraw, err = FindWithdrawWheref(context.Background(), "tid=$%v", withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if withdraw.TID != findWithdraw.TID {
		t.Error("find id error")
		return
	}
	findWithdraw, err = FindWithdrawFilterWheref(context.Background(), "#all", "tid=$%v", withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if withdraw.TID != findWithdraw.TID {
		t.Error("find id error")
		return
	}
	findWithdraw, err = FindWithdrawWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{withdraw.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if withdraw.TID != findWithdraw.TID {
		t.Error("find id error")
		return
	}
	findWithdraw, err = FindWithdrawWherefCall(GetQueryer, context.Background(), true, "tid=$%v", withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if withdraw.TID != findWithdraw.TID {
		t.Error("find id error")
		return
	}
	withdrawList, withdrawMap, err := ListWithdrawByID(context.Background())
	if err != nil || len(withdrawList) > 0 || withdrawMap == nil || len(withdrawMap) > 0 {
		t.Error(err)
		return
	}
	withdrawList, withdrawMap, err = ListWithdrawByID(context.Background(), withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(withdrawList) != 1 || withdrawList[0].TID != withdraw.TID || len(withdrawMap) != 1 || withdrawMap[withdraw.TID] == nil || withdrawMap[withdraw.TID].TID != withdraw.TID {
		t.Error("list id error")
		return
	}
	withdrawList, withdrawMap, err = ListWithdrawFilterByID(context.Background(), "#all")
	if err != nil || len(withdrawList) > 0 || withdrawMap == nil || len(withdrawMap) > 0 {
		t.Error(err)
		return
	}
	withdrawList, withdrawMap, err = ListWithdrawFilterByID(context.Background(), "#all", withdraw.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(withdrawList) != 1 || withdrawList[0].TID != withdraw.TID || len(withdrawMap) != 1 || withdrawMap[withdraw.TID] == nil || withdrawMap[withdraw.TID].TID != withdraw.TID {
		t.Error("list id error")
		return
	}
	withdrawList = nil
	withdrawMap = nil
	err = ScanWithdrawByID(context.Background(), []int64{withdraw.TID}, &withdrawList, &withdrawMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(withdrawList) != 1 || withdrawList[0].TID != withdraw.TID || len(withdrawMap) != 1 || withdrawMap[withdraw.TID] == nil || withdrawMap[withdraw.TID].TID != withdraw.TID {
		t.Error("list id error")
		return
	}
	withdrawList = nil
	withdrawMap = nil
	err = ScanWithdrawFilterByID(context.Background(), "#all", []int64{withdraw.TID}, &withdrawList, &withdrawMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(withdrawList) != 1 || withdrawList[0].TID != withdraw.TID || len(withdrawMap) != 1 || withdrawMap[withdraw.TID] == nil || withdrawMap[withdraw.TID].TID != withdraw.TID {
		t.Error("list id error")
		return
	}
	withdrawList = nil
	withdrawMap = nil
	err = ScanWithdrawWheref(context.Background(), "tid=$%v", []interface{}{withdraw.TID}, "", &withdrawList, &withdrawMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(withdrawList) != 1 || withdrawList[0].TID != withdraw.TID || len(withdrawMap) != 1 || withdrawMap[withdraw.TID] == nil || withdrawMap[withdraw.TID].TID != withdraw.TID {
		t.Error("list id error")
		return
	}
	withdrawList = nil
	withdrawMap = nil
	err = ScanWithdrawFilterWheref(context.Background(), "#all", "tid=$%v", []interface{}{withdraw.TID}, "", &withdrawList, &withdrawMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(withdrawList) != 1 || withdrawList[0].TID != withdraw.TID || len(withdrawMap) != 1 || withdrawMap[withdraw.TID] == nil || withdrawMap[withdraw.TID].TID != withdraw.TID {
		t.Error("list id error")
		return
	}
}
