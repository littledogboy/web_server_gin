package spiders

import (
	"strings"
)

func WebsiteFromURLString(urlStr string) WebSiteType {
	if strings.Contains(urlStr, Meirentu.Doman) {
		return MeirentuSite
	} else if strings.Contains(urlStr, Fulitu.Doman) {
		return FulituSite
	} else if strings.Contains(urlStr, Bestprettygirl.Doman) {
		return BestprettygirlSite
	}
	return UnknownSite
}
