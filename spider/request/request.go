package request

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/shawpo/lagou/params"
	"github.com/shawpo/lagou/utils"
)

// KdPositionRequest : 返回json api的请求
func KdPositionRequest(values url.Values) (request *http.Request, err error) {
	request, err = http.NewRequest("POST",
		"https://www.lagou.com/jobs/positionAjax.json?needAddtionalResult=false",
		strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	request.Header = http.Header{
		"content-type":       {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Accept-Encoding":    {"gzip, deflate"},
		"Host":               {"www.lagou.com"},
		"Origin":             {"http://www.lagou.com"},
		"User-Agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"},
		"X-Requested-With":   {"XMLHttpRequest"},
		"Referer":            {"https://www.lagou.com/jobs/list_" + values.Get("kd") + "?labelWords=&fromSearch=true&suginput="},
		"Proxy-Connection":   {"keep-alive"},
		"X-Anit-Forge-Code":  {"0"},
		"X-Anit-Forge-Token": {"None"},
	}
	return
}

// DetailRequest : 返回详情页的请求
func DetailRequest(positionID int) (request *http.Request, err error) {
	detailURL := fmt.Sprintf("https://www.lagou.com/jobs/%d.html", positionID)
	request, err = http.NewRequest("GET", detailURL, nil)
	if err != nil {
		return
	}
	return
}

// Fetch : 发起请求并返回请求的响应
func Fetch(request *http.Request) (*http.Response, error) {
	// 添加cookie
	request.Header["Cookie"] = []string{params.COOKIE}
	client := &http.Client{
	//Jar: CookieJar,
	}
	log.Printf("Fetching url: %s", request.URL)
	utils.RandTimeSleep(1, 2)
	return client.Do(request)
}
