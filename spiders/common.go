package spiders

import (
	"math/rand"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
