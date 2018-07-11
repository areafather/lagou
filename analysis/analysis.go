package analysis

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/shawpo/lagou/params"
	"github.com/shawpo/sego"
)

// Analysis : do analysis
func Analysis(path string, segment sego.Segmenter) (map[string]int, error) {
	file, err := os.Open(path)

	var wordsMap = make(map[string]int)

	if err != nil {
		return wordsMap, err
	}

	contents, err := ioutil.ReadAll(file)
	// 分词
	segments := segment.Segment(contents)
	// 获取分词结果slice
	words := sego.SegmentsToSlice(segments, true)
	// 获取过滤词典map
	filterMap := FilterMap(params.FILTER)
	// 获取同义替换词典map
	synonymMap := SynonymMap(params.SYNONYM)
	var word string
	for _, word = range words {

		// 过滤
		if _, exist := filterMap[word]; exist {
			continue

		}

		// 忽略非英文，或以数字_开头的多位字符串
		if !isValid(word) {
			continue
		}

		// 同义替换
		if synonym, exist := synonymMap[word]; exist {
			word = synonym
		}

		if _, exist := wordsMap[word]; exist {
			wordsMap[word]++
		} else {
			wordsMap[word] = 1
		}
	}

	return wordsMap, err
}

// SynonymMap : get SynonymMap from file
func SynonymMap(path string) (synonymMap map[string]string) {
	synonymMap = make(map[string]string)
	synonymFile, err := os.Open(path)
	defer synonymFile.Close()
	if err != nil {
		log.Fatalf("无法载入同义替换字典文件 \"%s\" \n", path)
	}

	reader := bufio.NewReader(synonymFile)

	// 逐行读取
	for {
		var word string
		var synonym string
		size, _ := fmt.Fscanln(reader, &word, &synonym)
		if size == 0 {
			// 文件结束
			break
		} else if size < 2 {
			// 无效行
			continue
		}
		synonymMap[word] = synonym
	}
	return
}

// FilterMap : get FilterMap from file
func FilterMap(path string) (filterMap map[string]bool) {
	filterMap = make(map[string]bool)
	filterFile, err := os.Open(path)
	defer filterFile.Close()
	if err != nil {
		log.Fatalf("无法载入过滤字典文件 \"%s\" \n", path)
	}

	reader := bufio.NewReader(filterFile)

	// 逐行读取
	for {
		var word string
		size, _ := fmt.Fscanln(reader, &word)
		if size == 0 {
			// 文件结束
			break
		}
		filterMap[word] = true
	}
	return
}

func isValid(word string) bool {
	length := len([]rune(word))

	isEn, err := regexp.MatchString("\b[a-zA-Z]+\b", word)
	if err != nil {
		return false
	}

	notWord, err := regexp.MatchString(`^(\d|_)`, word)
	if err != nil {
		return false
	}

	if isEn || (length > 1 && !notWord) {
		return true
	} else {
		return false
	}
}
