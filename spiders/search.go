package spiders

import (
	"errors"
)

func SearchPageSpider(q string, page string, callback func(Home, error)) {
	data := Home{}
	var e error

	MRTSearchPageSpider(q, page, func(h Home, err error) {
		if err == nil && len(h.Recommends) > 0 {
			data.Recommends = append(data.Recommends, h.Recommends...)
		} else {
			e = errors.Join(err)
		}
	})

	FLTSearchPageSpider(q, page, func(h Home, err error) {
		if err == nil && len(h.Recommends) > 0 {
			data.Recommends = append(data.Recommends, h.Recommends...)
		} else {
			e = errors.Join(err)
		}
	})

	callback(data, e)
}

func MRTSearchPageSpider(q string, page string, callback func(Home, error)) {
	desURL := Meirentu.Doman + "s/" + q + "-" + page + ".html"
	MRTDesURLSpider(desURL, page, Meirentu.Refer, Meirentu.ReferValue, Meirentu_SearchPage_Selector, callback)
}

func FLTSearchPageSpider(q string, page string, callback func(Home, error)) {
	desURL := Fulitu.Doman + "s/" + q + "-" + page + ".html"
	MRTDesURLSpider(desURL, page, "", "", Fulitu_SearchPage_Selector, callback)
}
