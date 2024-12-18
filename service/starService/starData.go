package starService

import (
	"fmt"
	"studentGrow/dao/mysql"
	"studentGrow/models"
	"studentGrow/models/constant"
	"time"
)

// StarGrid 查询表格所有数据
func StarGrid(usernameslice []string) ([]models.StarBack, error) {
	var starBack []models.StarBack
	for _, username := range usernameslice {
		//结构体对象存放数据
		//查询 name ，id，userfans，score，hot
		//查询名字
		name, err := mysql.SelName(username)
		if err != nil {
			fmt.Println("StarGridClass err", err)
			return nil, err
		}

		//查询 id
		id, err := mysql.SelId(username)
		if err != nil {
			fmt.Println("StarGridClass err", err)
			return nil, err
		}

		//查询粉丝数
		userFans, err := mysql.Selfans(id)
		if err != nil {
			fmt.Println("StarGridClass Selfans err", err)
			return nil, err
		}
		var score int
		//查询积分
		allScore, err := mysql.Score(id)
		for _, thisScore := range allScore {
			score += thisScore
		}
		if err != nil {
			fmt.Println("StarGridClass Score err", err)
			return nil, err
		}

		//查询被推举次数
		frequency, err := mysql.Frequency(username)
		if err != nil {
			fmt.Println("StarGridClass Frequency err", err)
			return nil, err
		}

		//查询文章数
		article, err := mysql.SelArticleNum(id)
		if err != nil {
			fmt.Println("StarGridClass SelArticleNum err", err)
			return nil, err
		}

		//查询文章质量
		class, err := mysql.SelArticleQua(id, 1)
		grade, err := mysql.SelArticleQua(id, 2)
		college, err := mysql.SelArticleQua(id, 3)
		quality := class*constant.ArticleClass + grade*constant.ArticleGrade + college*constant.ArticleCollege

		var hot int
		//查询热度
		likes, collects, err := mysql.SelHot(id)
		for _, like := range likes {
			hot += like
		}
		for _, collect := range collects {
			hot += collect
		}
		if err != nil {
			fmt.Println("StarGridClass SelHot err", err)
			return nil, err
		}

		star := models.StarBack{
			Username:           username,
			Frequency:          frequency,
			Name:               name,
			User_article_total: article,
			Userfans:           userFans,
			Score:              score,
			Hot:                hot,
			Quality:            int(quality),
			Status:             false,
		}
		starBack = append(starBack, star)
	}
	return starBack, nil
}

func StarGuidGrade(year int, page, limit int) ([]string, int64, error) {
	data := time.Now()
	usename, number, err := mysql.SelNStarUser(data, year, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return usename, number, nil
}

// SearchGrade 年级管理员搜索
func SearchGrade(name string, year int, page, limit int) ([]string, int64, error) {
	data := time.Now()
	gUsername, number, err := mysql.SelSearchGrade(name, data, year, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return gUsername, number, nil
}

// StarClass 返回成长之星的班级之星
func StarClass(session int) ([]models.StarClass, error) {
	var StarClasssli []models.StarClass
	//查询出所有班级之星
	usernamesli, err := mysql.SelStar(session, 1, 0, 0)
	if err != nil {
		return nil, err
	}
	//通过班级进行分组
	starMap, err := GroupByClass(usernamesli)
	if err != nil {
		return nil, err
	}
	//对结构体赋值
	for class, name := range starMap {
		starclass := models.StarClass{
			ClassName: class,
			ClassStar: name,
		}
		StarClasssli = append(StarClasssli, starclass)
	}
	return StarClasssli, nil
}

// GroupByClass 通过班级进行分组
func GroupByClass(usernamesli []string) (map[string][]string, error) {
	starmap := make(map[string][]string)
	for _, username := range usernamesli {
		class, err := mysql.SelClass(username)
		name, err := mysql.SelName(username)
		if err != nil {
			return nil, err
		}
		if _, exists := starmap[class]; !exists {
			starmap[class] = []string{name}
		} else {
			starmap[class] = append(starmap[class], name)
		}
	}
	return starmap, nil
}

// StarGrade 返回年级之星
func StarGrade(session int) ([]models.StarGrade, error) {
	var starGrade []models.StarGrade
	usernamesli, err := mysql.SelStar(session, 2, 0, 0)
	if err != nil {
		return nil, err
	}
	for _, username := range usernamesli {
		name, _ := mysql.SelName(username)
		class, _ := mysql.SelClass(username)
		//赋值结构体
		stargrade := models.StarGrade{
			GradeName:  name,
			GradeClass: class,
		}
		//加入切片
		starGrade = append(starGrade, stargrade)
	}
	return starGrade, nil
}

// StarCollege 返回院级之星
func StarCollege(session int) ([]models.StarGrade, error) {
	var starGrade []models.StarGrade
	usernamesli, err := mysql.SelStar(session, 3, 0, 0)
	if err != nil {
		return nil, err
	}
	for _, username := range usernamesli {
		name, _ := mysql.SelName(username)
		class, _ := mysql.SelClass(username)
		//赋值结构体
		stargrade := models.StarGrade{
			GradeName:  name,
			GradeClass: class,
		}
		//加入切片
		starGrade = append(starGrade, stargrade)
	}
	return starGrade, nil
}

// QStarClass 返回前台成长之星
func QStarClass(starType int, page int, limit int) ([]models.StarStu, error) {
	var starlist []models.StarStu
	session, err := mysql.SelMax()
	if err != nil {
		return nil, err
	}
	usernameslic, err := mysql.SelStar(session, starType, page, limit)
	if err != nil {
		return nil, err
	}
	for _, username := range usernameslic {
		name, err := mysql.SelName(username)
		headshot, err := mysql.SelHead(username)
		starstu := models.StarStu{
			Username:     username,
			Name:         name,
			UserHeadshot: headshot,
		}
		if err != nil {
			return nil, err
		}
		starlist = append(starlist, starstu)
	}
	return starlist, nil
}

// SelNumClass 查询表中跟管理员一个班的有多少人
func SelNumClass(class string) (int, error) {
	classnum, err := mysql.SelStuClass(class)
	if err != nil {
		return 0, err
	}
	return len(classnum), nil
}

// SelNumGrade 查询年级管理员已推选的人数
func SelNumGrade(data time.Time, year int) (int, error) {
	classNum, err := mysql.SelStuGrade(data, year)
	if err != nil {
		return 0, err
	}
	return len(classNum), nil
}

// SelTimeStar 前台通过开始时间和结束时间查询成长之星
func SelTimeStar(starTime, endTime string, starType int, page int, limit int) ([]models.StarStu, error) {
	var starlist []models.StarStu
	usernameslic, err := mysql.SelTimeStar(starTime, endTime, starType, page, limit)
	if err != nil {
		return nil, err
	}
	for _, username := range usernameslic {
		name, err := mysql.SelName(username)
		headshot, err := mysql.SelHead(username)
		starstu := models.StarStu{
			Username:     username,
			Name:         name,
			UserHeadshot: headshot,
		}
		if err != nil {
			return nil, err
		}
		starlist = append(starlist, starstu)
	}
	return starlist, nil
}

func BackNameData(username, role string) ([]string, error) {
	date := time.Now()
	var stuName []string
	var err error
	switch role {
	case "class":
		class, _ := mysql.SelClass(username)
		stuName, err = mysql.SelStuClass(class)
	case "grade1":
		stuName, err = mysql.SelStuGrade(date, -1)
	case "grade2":
		stuName, err = mysql.SelStuGrade(date, -2)
	case "grade3":
		stuName, err = mysql.SelStuGrade(date, -3)
	case "grade4":
		stuName, err = mysql.SelStuGrade(date, -4)
	case "college":
		stuName, err = mysql.SelDataCollege()
	}
	if err != nil {
		return nil, err
	}
	return stuName, nil
}
