package gexapi

import (
	"regexp"
	"strings"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/xlog"
)

var InternationalURL1 = `https://gold.usd-cny.com/data/jin.js`
var InternationalURL2 = `https://www.ibjarates.com`
var InternationalClient = xhttp.Shared
var InternationalPrice = []xmap.M{}

func ProcRefreshInternationalPrice() (err error) {
	err = pgx.ErrNoRows
	locations := []xmap.M{}
	{
		data, xerr := InternationalClient.GetText("%v", InternationalURL1)
		if xerr != nil {
			xlog.Warnf("ProcRefreshInternationalPrice load international price from %v fail with %v", InternationalURL1, xerr)
			return
		}
		lines := strings.Split(data, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "var")
			line = strings.TrimSpace(line)
			parts := strings.SplitN(line, "=", 2)
			if len(parts) < 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			if !(key == "hq_str_gds_AUTD" || key == "hq_str_hf_XAU" || key == "hq_str_hf_GC") {
				continue
			}
			data := strings.TrimSpace(parts[1])
			data = strings.Trim(data, `"`)
			parts = strings.SplitN(data, ",", 2)
			price, xerr := converter.Float64Val(parts[0])
			if xerr != nil {
				xlog.Warnf("ProcRefreshInternationalPrice load international price fail with %v by %v", xerr, line)
				return
			}
			switch key {
			case "hq_str_gds_AUTD":
				locations = append(locations, xmap.M{
					"location": xmap.M{
						"MM": "တရုတ်",
						"CN": "中国",
						"US": "China",
					},
					"price":    price,
					"currency": "￥",
					"unit":     "g",
				})
			case "hq_str_hf_XAU":
				locations = append(locations, xmap.M{
					"location": xmap.M{
						"MM": "London",
						"CN": "伦敦",
						"US": "London",
					},
					"price":    price,
					"currency": "$",
					"unit":     "oz",
				})
				// case "hq_str_hf_GC":
				// 	locations = append(locations, xmap.M{
				// 		"location": xmap.M{
				// 			"MM": "New York",
				// 			"CN": "纽约",
				// 			"US": "New York",
				// 		},
				// 		"price": price,
				// 	})
			}
		}
	}
	{
		data, xerr := InternationalClient.GetText("%v", InternationalURL2)
		if xerr != nil {
			xlog.Warnf("ProcRefreshInternationalPrice load international price from %v fail with %v", InternationalURL2, xerr)
			return
		}
		matched := regexp.MustCompile(`lblrate24K">[0-9\.]+`).FindString(data)
		if len(matched) < 1 {
			xlog.Warnf("ProcRefreshInternationalPrice load international price fail with %v", xerr)
			return
		}
		parts := strings.SplitN(matched, ">", 2)
		price, _ := converter.Float64Val(parts[1])
		locations = append(locations, xmap.M{
			"location": xmap.M{
				"MM": "India",
				"CN": "印度",
				"US": "India",
			},
			"price":    price,
			"currency": "₹",
			"unit":     "g",
		})
	}
	InternationalPrice = locations
	return
}
