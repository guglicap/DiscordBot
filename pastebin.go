package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

func hasPastebinLink(s string) (link, id string) {
	re := regexp.MustCompile("pastebin\\.com\\/(\\w{8})")
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return "", ""
	}
	return matches[0], matches[1]
}

func getPasteLanguage(link string) string {
	resp, err := http.Get("http://" + link)
	check(err)
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	check(err)
	re := regexp.MustCompile("\\/css_lang\\/(\\w+)\\.css")
	matches := re.FindStringSubmatch(string(htmlBytes))
	if matches == nil {
		return ""
	}
	if len(matches) > 1 {
		return matches[1]
	}
	return "text"
}

func getPasteRaw(id string) string {
	resp, err := http.Get("http://pastebin.com/raw/" + id)
	check(err)
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	check(err)
	return string(htmlBytes)
}
