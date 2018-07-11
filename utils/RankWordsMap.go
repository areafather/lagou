package utils

import (
	"sort"
	"strconv"
)

// WordCountList :
// github.com/gizak/termui 的表格组件需要 [][]string类型的数据
// word and word's count
type WordCountList [][]string

func (p WordCountList) Len() int { return len(p) }
func (p WordCountList) Less(i, j int) bool {
	valuei, _ := strconv.Atoi(p[i][1])
	valuej, _ := strconv.Atoi(p[j][1])
	return valuei < valuej
}
func (p WordCountList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// RankByWordCount : 按出现次数排序
func RankByWordCount(wordFrequencies map[string]int) WordCountList {
	pl := make(WordCountList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = []string{k, strconv.Itoa(v)}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
