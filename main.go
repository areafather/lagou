package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	ui "github.com/gizak/termui"
	"github.com/manifoldco/promptui"
	"github.com/shawpo/lagou/analysis"
	"github.com/shawpo/lagou/params"
	"github.com/shawpo/lagou/spider/engine"
	"github.com/shawpo/lagou/spider/handler"
	"github.com/shawpo/lagou/spider/request"
	"github.com/shawpo/lagou/spider/scheduler"
	"github.com/shawpo/lagou/spider/types"
	"github.com/shawpo/lagou/utils"
	"github.com/shawpo/sego"
)

func main() {
	exists, err := utils.ExistPositions(".", params.POSITIONSEXT)
	if len(exists) > 0 && err == nil {
		selects := append(exists, "no")
		prompt := promptui.Select{
			Label: "--已有以下岗位数据，请选择其中一个进行分析，选择no表示获取新的岗位数据",
			Items: selects,
		}

		_, params.KD, err = prompt.Run()

		if err != nil {
			log.Fatalf("获取选择失败 %v\n", err)
		}
	}
	if params.KD != "no" && params.KD != "" {
		analysisKd()
		return
	}

	prompt := promptui.Prompt{
		Label: "输入岗位名称",
	}

	params.KD, err = prompt.Run()

	if err != nil {
		log.Fatalf("获取岗位名称失败：%v\n", err)
	}
	filePath := params.KD + params.POSITIONSEXT
	_, err = os.Create(filePath)
	if err != nil {
		log.Fatalf("创建数据文件失败：%s\n", filePath)
	}
	e := &engine.Engine{
		Scheduler:   scheduler.New(),
		WorkerCount: params.WORKERCOUNT,
	}
	values := url.Values{
		"first": {"true"},
		"pn":    {"1"},
		"kd":    {params.KD},
	}
	request, err := request.KdPositionRequest(values)
	if err != nil {
		log.Fatal(err)
	}
	task := types.Task{
		Request: request,
		Handler: handler.GetPosition,
	}
	e.Run(task)
	analysisKd()
}

func analysisKd() {
	fmt.Println("--开始词频分析：")
	var segment sego.Segmenter
	segment.LoadDictionary(params.ITDIC)
	wordsMap, err := analysis.Analysis(params.KD+params.POSITIONSEXT, segment)
	var sum = len(wordsMap)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("--词频分析已完成！总计%d个关键词！\n--开始进行排序：\n", sum)
	wordRands := utils.RankByWordCount(wordsMap)
	fmt.Println("--排序已完成！")

	displayRank(wordRands)
}

func displayRank(ranks utils.WordCountList) {
	header := [][]string{
		[]string{"关键词", "出现频次"},
	}
	ranks = ([][]string)(ranks)
	var sum = len(ranks)
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	var i, j = 0, 0 + params.TABLEPAGECOUNT
	// 标题
	p := ui.NewPar(fmt.Sprintf("%s职位关键词(总数%d):q退出，w向上翻页，s向下翻页", params.KD, sum))
	p.Height = 3
	p.Width = 60
	p.TextFgColor = ui.ColorRed
	p.BorderLabel = fmt.Sprintf("排名为%d~%d的关键词", i+1, j)
	p.BorderFg = ui.ColorWhite

	table := ui.NewTable()
	if len(ranks) > j {
		table.Rows = append(header, ranks[i:j]...)
	} else {
		table.Rows = append(header, ranks...)
	}
	table.FgColor = ui.ColorYellow
	table.Block.BorderFg = ui.ColorWhite
	table.Separator = false
	table.Width = 60
	table.Height = j - i + 3
	table.TextAlign = ui.AlignCenterHorizontal
	//ui.Render(p, table)
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(8, 2, p)),
		ui.NewRow(
			ui.NewCol(8, 2, table)))

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	// handle key q pressing
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})
	render := func(i, j int) {
		table.Rows = append(header, ranks[i:j]...)
		p.BorderLabel = fmt.Sprintf("排名为%d~%d的关键词", i+1, j)
		ui.Body.Rows[0].Cols[0] = ui.NewCol(8, 2, p)
		ui.Body.Rows[1].Cols[0] = ui.NewCol(8, 2, table)
		//table.Height = j - i + 3
		ui.Render(ui.Body)
	}

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/s", func(e ui.Event) {
		if i+params.TABLEPAGECOUNT > sum {
			render(i, j)
			return
		}

		if j+params.TABLEPAGECOUNT > sum {
			j = sum
			i = i + params.TABLEPAGECOUNT
		} else {
			i = i + params.TABLEPAGECOUNT
			j = j + params.TABLEPAGECOUNT
		}
		render(i, j)

	})
	ui.Handle("/sys/kbd/w", func(e ui.Event) {
		if j-params.TABLEPAGECOUNT < 0 || i-params.TABLEPAGECOUNT < 0 {
			render(i, j)
			return
		}
		j = i
		i = i - params.TABLEPAGECOUNT
		render(i, j)
	})
	ui.Loop()
}
