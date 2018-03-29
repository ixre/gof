package algorithm

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

type Item struct {
	Name string
	Prob float32
}

var (
	items = []*Item{
		{"太阳水", 0.7000},   // 70%
		{"千年雪霜", 0.2000},  // 20%
		{"无极棍", 0.0640},   // 6.4%
		{"召唤神兽", 0.0340},  // 3.4%
		{"极品法杖", 0.0015},  // 0.15%
		{"极品屠龙刀", 0.0005}, // 0.05%
	}
)

type ItemList []*Item

func (a ItemList) Len() int {
	return len(a)
}

func (a ItemList) Less(i, j int) bool {
	return a[i].Prob < a[j].Prob
}

func (a ItemList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// 随机抽取奖品，奖品比例总和应小于100%
func getItem(items []*Item) *Item {
	// 建立区间及区间与奖品的映射
	var bitArr []int
	var bitMap = make(map[int]*Item)
	b := 0
	for _, v := range items {
		b += int(v.Prob * 10000)
		bitArr = append(bitArr, b)
		bitMap[b] = v
	}
	// 生成随机数R
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10000) + 1
	// 获取R所在的区间，使用二分算法进行搜索
	i := sort.SearchInts(bitArr, r)
	fmt.Println(bitArr, "i=", i, ";r=", r)
	if i < len(bitArr) {
		return bitMap[bitArr[i]]
	}
	return nil
}

func TestGetItem(t *testing.T) {
	jpCount := 0
	nmCount := 0

	// 按概率排序
	sort.Sort(ItemList(items))
	for i := 0; i < 3; i++ {
		r := getItem(items)
		if r == nil {
			fmt.Println("很遗憾，您什么都没有抽到")
			continue
		}
		if strings.HasPrefix(r.Name, "极品") {
			jpCount += 1
			//fmt.Println("您抽取到了：", r)
		}
		if strings.HasPrefix(r.Name, "太阳水") {
			nmCount += 1
			//fmt.Println("您抽取到了：", r)
		}
		fmt.Println("您抽取到了：", r.Name)
		time.Sleep(time.Second)
	}
	fmt.Println("抽到普将数量：", nmCount)
	fmt.Println("抽到极品数量：", jpCount)
}
