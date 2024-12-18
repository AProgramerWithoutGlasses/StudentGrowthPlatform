package mysql

import (
	"studentGrow/models"
	"studentGrow/models/gorm_model"
	"time"
)

// SelGradeId 获取大一or大二or大三or大四的用户id 和usernameSlice
func SelGradeId(data time.Time, year int) ([]int, []string, error) {
	var uidslice []int
	var usernameslice []string
	//计算时间的间隔的右端
	// 计算时间间隔的左端
	CurrentYear := data.AddDate(year+1, 9, 1)
	YearAgo := data.AddDate(year, 9, 1)
	err := DB.Model(&gorm_model.User{}).Where("identity = ?", "学生").Where("plus_time >= ?", YearAgo).Where("plus_time <= ?", CurrentYear).Select("id").Scan(&uidslice).Error
	err = DB.Model(&gorm_model.User{}).Where("identity = ?", "学生").Where("plus_time >= ?", YearAgo).Where("plus_time <= ?", CurrentYear).Select("username").Scan(&usernameslice).Error
	if err != nil {
		return nil, nil, err
	}
	return uidslice, usernameslice, nil
}

// SelUid 根据班级查询班成员的id
func SelUid(class string) ([]int, error) {
	//班级成员的id切片
	var uidSlice []int
	err := DB.Model(&gorm_model.User{}).Select("id").Where("deleted_at IS NULL").Where("class = ?", class).Scan(&uidSlice).Error
	// 检查并返回错误
	if err != nil {
		return nil, err
	}
	return uidSlice, nil
}

// SelCollegeId 查询所有人的uid和username
func SelCollegeId() ([]int, []string, error) {
	var uidslice []int
	var usernameslice []string
	err := DB.Model(&gorm_model.User{}).Where("identity <> ?", "开发").Select("id").Scan(&uidslice).Error
	err = DB.Model(&gorm_model.User{}).Where("identity <> ?", "开发").Select("username").Scan(&usernameslice).Error
	if err != nil {
		return nil, nil, err
	}
	return uidslice, usernameslice, nil
}

// SelArticleNum 根据id查询每个用户的帖子数及各级文章个数
func SelArticleNum(id int) (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.Articles{}).Where("user_id = ?", id).Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelArticleQua 根据id查询每个用户的帖子数及各级文章个数
func SelArticleQua(id int, quality int) (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.Articles{}).Where("quality = ?", quality).Where("user_id = ?", id).Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelArticle 根据id查询每个用户目标天数的贴子数
func SelArticle(id int, date time.Time) (int64, error) {
	var number int64
	// 获取当天的开始时间（00:00:00）
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowdate := from.Add(-8 * time.Hour)
	// 获取第二天的开始时间（00:00:00），用于查询截止到当天结束的时间范围
	to := nowdate.Add(24 * time.Hour)

	// 使用 BETWEEN 查询当天的记录数
	err := DB.Model(&gorm_model.Articles{}).Where("user_id = ? ", id).
		Where("deleted_at IS NULL").
		Where("created_at BETWEEN ? AND ?", nowdate, to).
		Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelUsername 根据班级查询班成员的username
func SelUsername(class string, page, limit int) ([]string, int64, error) {
	//班级成员的id切片
	var usernameSlice []string
	var number int64
	err := DB.Model(&gorm_model.User{}).Select("username").Where("class = ?", class).Offset((page - 1) * limit).Limit(limit).Scan(&usernameSlice).Error
	err = DB.Model(&gorm_model.User{}).Select("username").Where("class = ?", class).Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return nil, 0, err
	}
	return usernameSlice, number, nil
}

// SelStudent 查询所有学生人数
func SelStudent() (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.User{}).Where("identity = ?", "学生").Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelTeacher 查询所有教师人数
func SelTeacher() (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.User{}).Where("identity = ?", "老师").Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelArticleLike 查询帖子总赞数
func SelArticleLike(id int) ([]int, error) {
	var number []int
	err := DB.Model(&gorm_model.Articles{}).Select("like_amount").Where("user_id = ?", id).Scan(&number).Error
	// 检查并返回错误
	if err != nil {
		return nil, err
	}
	return number, nil
}

// SelArticleRead 查询帖子总阅读数
func SelArticleRead(id int) ([]int, error) {
	var number []int
	err := DB.Model(&gorm_model.Articles{}).Select("read_amount").Where("user_id = ?", id).Scan(&number).Error
	// 检查并返回错误
	if err != nil {
		return nil, err
	}
	return number, nil
}

// SelTagArticle 查询不同tag下的文章的大小
func SelTagArticle() ([]models.TagAmount, error) {
	var tagcount []models.TagAmount
	err := DB.Model(&gorm_model.ArticleTag{}).Select("tag_id,COUNT(*)as count").Group("tag_id").Scan(&tagcount).Error
	if err != nil {
		return nil, err
	}
	return tagcount, nil
}

// SelTagArticleTime 查询不同tag不同时间下的文章的大小
func SelTagArticleTime(date string) ([]models.TagAmount, error) {
	nowdate, err := time.Parse("2006-01-02", date)
	nowDate := nowdate.Add(-8 * time.Hour)
	//获取第二天的开始时间（00:00:00），用于查询截止到当天结束的时间范围
	to := nowDate.Add(24 * time.Hour)
	var tagcount []models.TagAmount
	err = DB.Model(&gorm_model.ArticleTag{}).Where("created_at BETWEEN ? AND ?", nowDate, to).Select("tag_id,COUNT(*)as count").Group("tag_id").Scan(&tagcount).Error
	if err != nil {
		return nil, err
	}
	return tagcount, nil
}

// TagName 查询标签名
func TagName(id int) (string, error) {
	var tagname string
	err := DB.Model(&gorm_model.Tag{}).Where("id = ?", id).Select("tag_name").Scan(&tagname).Error
	if err != nil {
		return "", err
	}
	return tagname, nil
}

// SelLoginNum 查询指定日期登录的人数
func SelLoginNum(date time.Time) (int64, error) {
	var number int64
	// 获取当天的开始时间（00:00:00）
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowdate := from.Add(-8 * time.Hour)
	// 获取第二天的开始时间（00:00:00），用于查询截止到当天结束的时间范围
	to := nowdate.Add(24 * time.Hour)
	// 使用 BETWEEN 查询当天的记录数
	err := DB.Table("user_login_records").
		Where("created_at BETWEEN ? AND ?", nowdate, to).
		Count(&number).Error
	// 检查并返回错误
	if err != nil {
		return 0, err
	}
	return number, nil

}
