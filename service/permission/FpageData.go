package service

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"studentGrow/dao/mysql"
	"time"
)

// ArticleData 帖子总数
func ArticleData(uidslice []int) (int64, error) {
	var allNumber int64
	//把每个用户发布的帖子个数都查询并把数据相加
	for _, id := range uidslice {
		number, err := mysql.SelArticleNum(id)
		if err != nil {
			fmt.Println("ArticleData mysql.SelArticleNum err", err)
		}
		allNumber += number
	}
	return allNumber, nil
}

// NarticleDataClass 今日新帖数
func NarticleDataClass(uidslice []int) (int64, error) {
	var allNumber int64
	data := time.Now()
	//2.把每个用户发布的帖子个数都查询并把数据相加
	for _, id := range uidslice {
		number, err := mysql.SelArticle(id, data)
		if err != nil {
			fmt.Println("ArticleData mysql.SelArticleNum err", err)
		}
		allNumber += number
	}
	return allNumber, nil
}

// 这个方法用来获取当天日期的前一天
func getYesterdayDate(date time.Time) time.Time {
	// 获取当前日期的年月日
	year, month, day := date.Date()

	// 如果是 1 号,则需要获取上个月的最后一天
	if day == 1 {
		// 减去 1 个月,然后获取那个月的最后一天
		prevMonth := month - 1
		prevYear := year
		if prevMonth == 0 {
			prevMonth = 12
			prevYear--
		}
		return time.Date(prevYear, prevMonth, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	} else {
		// 否则直接减去 1 天即可
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	}
}

// ArticleDataClassRate 昨日新帖数跟今日的相比
func ArticleDataClassRate(uidslice []int, nownumber int64) (float64, error) {
	data := time.Now()
	ydata := getYesterdayDate(data)
	var allNumber int64
	//遍历获取班级成员在指定日期的帖子数
	for _, id := range uidslice {
		number, err := mysql.SelArticle(id, ydata)
		if err != nil {
			fmt.Println("ArticleDataClassRate mysql.SelArticleNum err", err)
			return 0, err
		}
		allNumber += number
	}
	if nownumber == 0 || allNumber == 0 {
		return 0, nil
	}
	//计算比率
	rate, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(nownumber-allNumber)/float64(allNumber)), 64)
	if err != nil {
		return 0, err
	}
	// 将rate乘以100并四舍五入到两位小数
	roundedRate := math.Round(rate*100) / 100
	return roundedRate, nil
}

// TodayVictor 查询今日访客数
func TodayVictor() (int64, error) {
	date := time.Now()
	allNumber, err := mysql.SelLoginNum(date)
	if err != nil {
		return 0, err
	}
	return allNumber, nil
}

// VictorRate 今天和昨天的比率
func VictorRate(todayVictor int64) (float64, error) {
	ydata := getYesterdayDate(time.Now())
	yesdayVictor, err := mysql.SelLoginNum(ydata)
	if todayVictor == 0 || yesdayVictor == 0 {
		return 0, nil
	}
	rate, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(todayVictor-yesdayVictor)/float64(yesdayVictor)), 64)
	if err != nil {
		return 0, err
	}
	// 将rate乘以100并四舍五入到两位小数
	roundedRate := math.Round(rate*100) / 100
	return roundedRate, nil
}

// LikeAmount 获取帖子总赞数
func LikeAmount(uidslice []int) (int, error) {
	var allNumber int
	for _, id := range uidslice {
		number, err := mysql.SelArticleLike(id)
		if err != nil {
			return 0, err
		}
		for _, num := range number {
			allNumber += num
		}
	}
	return allNumber, nil
}

// ReadAmount 获取帖子总阅读数
func ReadAmount(uidslice []int) int {
	var allNumber int
	for _, id := range uidslice {
		number, err := mysql.SelArticleRead(id)
		if err != nil {
			fmt.Println("LikeAmount SelArticleLike err", err)
		}
		for _, num := range number {
			allNumber += num
		}
	}
	return allNumber
}

// PillarData 柱状图
func PillarData() ([]string, []int, error) {
	tagcount, err := mysql.SelTagArticle()
	if err != nil {
		fmt.Println("PillarData SelTagArticle err", err)
		return nil, nil, err
	}

	// 根据人数排序
	sort.Slice(tagcount, func(i, j int) bool {
		return tagcount[i].Count > tagcount[j].Count
	})

	// 只取前7个或更少的元素
	tagName := make([]string, 0, 7)
	count := make([]int, 0, 7)

	for i := 0; i < 7 && i < len(tagcount); i++ {
		tagnames, err := mysql.TagName(tagcount[i].Tag)
		if err != nil {
			fmt.Println("PillarDataTime TagName err", err)
			return nil, nil, err
		}
		tagName = append(tagName, tagnames)
		count = append(count, tagcount[i].Count)
	}

	return tagName, count, nil
}

// PillarDataTime 特定日期的柱状图
func PillarDataTime(date string) ([]string, []int, error) {
	tagcount, err := mysql.SelTagArticleTime(date)
	if err != nil {
		fmt.Println("PillarData SelTagArticle err", err)
		return nil, nil, err
	}

	// 根据人数排序
	sort.Slice(tagcount, func(i, j int) bool {
		return tagcount[i].Count > tagcount[j].Count
	})

	// 只取前7个或更少的元素
	tagName := make([]string, 0, 7)
	count := make([]int, 0, 7)

	for i := 0; i < 7 && i < len(tagcount); i++ {
		tagnames, err := mysql.TagName(tagcount[i].Tag)
		if err != nil {
			fmt.Println("PillarDataTime TagName err", err)
			return nil, nil, err
		}
		tagName = append(tagName, tagnames)
		count = append(count, tagcount[i].Count)
	}

	return tagName, count, nil
}
