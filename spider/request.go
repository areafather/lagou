package spider

import (
	"net/http"
	"net/url"
	"strings"
	"fmt"
)

const COOKIE  = ``
func KdPositionsRequest(values url.Values) (request *http.Request, err error) {
	request, err = http.NewRequest("POST",
		"https://www.lagou.com/jobs/positionAjax.json?needAddtionalResult=false",
			strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	request.Header = http.Header{
		"content-type": {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Accept-Encoding": {"gzip, deflate"},
		"Host": {"www.lagou.com"},
		"Origin": {"http://www.lagou.com"},
		"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"},
		"X-Requested-With": {"XMLHttpRequest"},
		"Referer": {"https://www.lagou.com/jobs/list_"+ values.Get("kd") +"?labelWords=&fromSearch=true&suginput="},
		"Proxy-Connection": {"keep-alive"},
		"X-Anit-Forge-Code": {"0"},
		"X-Anit-Forge-Token": {"None"},
		"Cookie": {COOKIE},
	}
	return
}

func DetailRequest(positionId int) (request *http.Request, err error) {
	detailUrl := fmt.Sprintf("https://www.lagou.com/jobs/%d.html", positionId)
	request, err = http.NewRequest("GET", detailUrl, nil)
	if err != nil {
		return
	}
	request.Header = http.Header{
		"Cookie": {COOKIE},
	}
	return
}
