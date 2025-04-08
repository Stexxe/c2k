package utils

import (
	"log"
	"strings"
)

func UrlProto(url string) string {
	if strings.HasPrefix(url, "http://") {
		return "http"
	} else if strings.HasPrefix(url, "https://") {
		return "https"
	}

	log.Panicf("url protocol: unexpected url: %s", url)
	return ""
}

func UrlHost(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}

	log.Panicf("url protocol: unexpected url: %s", url)
	return ""
}
