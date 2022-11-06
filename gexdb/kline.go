package gexdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/codingeasygo/crud"
	"golang.org/x/net/context"
)

func StringInterv(str string) (interval time.Duration, err error) {
	switch str {
	case "1min":
		interval = time.Minute
	case "5min":
		interval = 5 * time.Minute
	case "30min":
		interval = 30 * time.Minute
	case "1hour":
		interval = 1 * time.Hour
	case "4hour":
		interval = 4 * time.Hour
	case "1day":
		interval = 24 * time.Hour
	case "1week":
		interval = 7 * 24 * time.Hour
	case "1mon":
		interval = 30 * 24 * time.Hour
	default:
		err = fmt.Errorf("%v not supported", str)
	}
	return
}

func IntervString(interval time.Duration) (str string, err error) {
	switch interval {
	case time.Minute:
		str = "1min"
	case 5 * time.Minute:
		str = "5min"
	case 30 * time.Minute:
		str = "30min"
	case 1 * time.Hour:
		str = "1hour"
	case 4 * time.Hour:
		str = "4hour"
	case 24 * time.Hour:
		str = "1day"
	case 7 * 24 * time.Hour:
		str = "1week"
	case 30 * 24 * time.Hour:
		str = "1mon"
	default:
		err = fmt.Errorf("%v not supported", interval)
	}
	return
}

func AddMultiKLine(ctx context.Context, lines ...*KLine) (added int64, err error) {
	if len(lines) < 1 {
		return
	}
	talbe, fileds, _, _ := crud.InsertArgs(lines[0], "^tid#all", nil)
	insertVal := []string{}
	insertArg := []interface{}{}
	for _, line := range lines {
		var param []string
		_, _, param, insertArg = crud.InsertArgs(line, "^tid#all", insertArg)
		insertVal = append(insertVal, "("+strings.Join(param, ",")+")")
	}
	insertSQL := fmt.Sprintf(`insert into %v(%v) values %v`, talbe, strings.Join(fileds, ","), strings.Join(insertVal, ","))
	_, added, err = Pool().Exec(ctx, insertSQL, insertArg...)
	return
}

func ListKLine(ctx context.Context, symbol, interval string, startTime, endTime time.Time) (lines []*KLine, err error) {
	err = crud.QueryWheref(
		Pool, ctx, &KLine{}, "#all",
		"symbol=$%v,interv=$%v,start_time>=$%v,start_time<$%v", []interface{}{symbol, interval, startTime, endTime},
		" order by start_time asc", 0, 0, &lines,
	)
	return
}
