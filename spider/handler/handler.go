package handler

// handler 用于对请求获取到的响应进行处理
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/shawpo/lagou/params"
	"github.com/shawpo/lagou/spider/request"
	"github.com/shawpo/lagou/spider/types"
)

var mutex sync.Mutex

// GetPosition : 从json api 的相应中解析出岗位列表
func GetPosition(reader io.Reader) ([]types.Task, error) {
	var tasks []types.Task
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return tasks, err
	}

	positionResp := PositionResp{}
	err = json.Unmarshal(body, &positionResp)
	if err != nil {
		return tasks, err
	}

	result := positionResp.Content.PositionResult.Result
	if len(result) > 0 {
		pn := positionResp.Content.PageNo + 1
		values := url.Values{
			"first": {"false"},
			"pn":    {strconv.Itoa(pn)},
			"kd":    {params.KD},
		}
		request, err := request.KdPositionRequest(values)
		if err != nil {
			return tasks, err
		}
		task := types.Task{request, GetPosition}
		tasks = append(tasks, task)
	} else {
		return tasks, types.NoneNewTask{}
	}

	for _, position := range result {
		request, err := request.DetailRequest(position.PositionID)
		if err != nil {
			continue
			log.Printf("error in build DetailRequest: %v", err)
		}
		task := types.Task{request, GetDetail}
		tasks = append(tasks, task)
	}
	return tasks, err
}

// GetDetail 从工作详情页面获取数据
func GetDetail(body io.Reader) (task []types.Task, err error) {
	html, err := ioutil.ReadAll(body)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return
	}
	detail := doc.Find(".job_bt div").Text()
	positionID, _ := doc.Find(".resume-deliver .send-CV-btn").Attr("data-position-id")
	if detail == "" {
		filePath := fmt.Sprintf("error-%s-%d.html", params.KD, positionID)
		file, _ := os.Create(filePath)
		file.Write(html)
		err = fmt.Errorf("PositionId %d: can't get detail, response has write in %s, maybe should update cookies",
			positionID, filePath)
		return
	}
	filePath := params.KD + params.POSITIONSEXT
	go func() {
		mutex.Lock()
		{
			// 写入文件
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := f.Write([]byte(detail)); err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
		mutex.Unlock()
	}()

	if err != nil {
		return
	}
	return
}
