package mysql

import (
	"database/sql"
	"studentGrow/models/gorm_model"
	"time"
)

// SelName 根据学号查名字，等一系列数据
func SelName(username string) (string, error) {
	var name string
	err := DB.Model(&gorm_model.User{}).Select("name").Where("username = ?", username).Scan(&name).Error
	if err != nil {
		return "", err
	}
	return name, nil
}

// SelId 根据账号查找id
func SelId(username string) (int, error) {
	var id int
	err := DB.Model(&gorm_model.User{}).Select("id").Where("username = ?", username).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Selfans 根据 id查询粉丝
func Selfans(id int) (int64, error) {
	var fans int64
	err := DB.Table("user_followers").Where("user_id = ?", id).Count(&fans).Error
	if err != nil {
		return 0, err
	}
	return fans, nil
}

// Score 查询积分
func Score(id int) ([]int, error) {
	var score []int
	err := DB.Model(&gorm_model.UserPoint{}).Where("user_id = ?", id).Select("point").Scan(&score).Error
	if err != nil {
		return nil, err
	}
	return score, nil
}

// Frequency 被推举次数
func Frequency(username string) (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.Star{}).Where("username = ?", username).Count(&number).Error
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelHot 查询热度
func SelHot(id int) ([]int, []int, error) {
	var like []int
	var collect []int
	err := DB.Model(&gorm_model.Article{}).Select("like_amount").Where("user_id =?", id).Scan(&like).Error
	err = DB.Model(&gorm_model.Article{}).Select("collect_amount").Where("user_id = ?", id).Scan(&collect).Error
	if err != nil {
		return nil, nil, err
	}
	return like, collect, nil
}

// SelNStarUser 查询未公布且未推举的学号合集(年级管理员搜索表格的数据)
func SelNStarUser(data time.Time, year, page, limit int) ([]string, int64, error) {
	var alluser []string
	var number int64
	_, username, err := SelGradeId(data, year)
	err = DB.Model(&gorm_model.Star{}).Where("type = ?", 1).Where("session = ?", 0).Where("username IN (?)", username).Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&alluser).Error
	err = DB.Model(&gorm_model.Star{}).Where("type = ?", 1).Where("session = ?", 0).Where("username IN (?)", username).Count(&number).Error
	if err != nil {
		return nil, 0, err
	}
	return alluser, number, nil
}

// SelSearchGrade 查询未公布的学号合集
func SelSearchGrade(name string, data time.Time, year int, page, limit int) ([]string, int64, error) {
	var alluser []string
	var number int64
	_, username, err := SelGradeId(data, year)
	err = DB.Model(&gorm_model.Star{}).Where("name LIKE ?", "%"+name+"%").Where("session = ?", 0).Where("username IN (?)", username).Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&alluser).Error
	err = DB.Model(&gorm_model.Star{}).Where("name LIKE ?", "%"+name+"%").Where("session = ?", 0).Where("username IN (?)", username).Count(&number).Error
	if err != nil {
		return nil, 0, err
	}
	return alluser, number, nil
}

// SelPlus 查询入学时间
func SelPlus(username string) (sql.NullTime, error) {
	var plus sql.NullTime
	err := DB.Model(&gorm_model.User{}).Select("plus_time").Where("username = ?", username).Scan(&plus).Error
	if err != nil {
		return sql.NullTime{}, err
	}
	return plus, nil
}

// SelStarColl 院级查询表里学号合集
func SelStarColl(page, limit int) ([]string, int64, error) {
	var alluser []string
	var number int64
	err := DB.Model(&gorm_model.Star{}).Where("type = ?", 2).Where("session = ?", 0).Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&alluser).Error
	err = DB.Model(&gorm_model.Star{}).Where("type = ?", 2).Where("session = ?", 0).Count(&number).Error
	if err != nil {
		return nil, 0, err
	}
	return alluser, number, nil
}

// SelSearchUser 根据名字班级查找学号--班级管理员搜索
func SelSearchUser(name string, class string, page, limit int) ([]string, int64, error) {
	var username []string
	var number int64
	err := DB.Model(&gorm_model.User{}).Where("class = ?", class).Where("name LIKE ?", "%"+name+"%").Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&username).Error
	err = DB.Model(&gorm_model.User{}).Where("class = ?", class).Where("name LIKE ?", "%"+name+"%").Count(&number).Error
	if err != nil {
		return nil, 0, err
	}
	return username, number, nil
}

// SelSearchColl 院级管理员搜索
func SelSearchColl(name string, page, limit int) ([]string, int64, error) {
	var usernamesli []string
	var number int64
	err := DB.Model(&gorm_model.Star{}).Where("name LIKE ?", "%"+name+"%").Where("type = ?", 2).Where("session = ?", 0).Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&usernamesli).Error
	err = DB.Model(&gorm_model.Star{}).Where("name LIKE ?", "%"+name+"%").Where("type = ?", 2).Where("session = ?", 0).Count(&number).Error
	if err != nil {
		return nil, 0, err
	}
	return usernamesli, number, nil
}

// CreatClass 班级管理员推选班级之星
func CreatClass(username string, name string) error {
	stars := gorm_model.Star{
		Username: username,
		Name:     name,
		Type:     1,
		Session:  0,
	}
	err := DB.Create(&stars).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateGrade 年级管理员推选更新数据
func UpdateGrade(username string) error {
	var star gorm_model.Star
	err := DB.Model(gorm_model.Star{}).Where("username = ?", username).Where("session = ?", 0).Where("type = ?", 1).First(&star).Error
	if err != nil {
		return err
	}
	star.Type = 2
	err = DB.Save(&star).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateCollege 院级管理员推选更新数据
func UpdateCollege(username string) error {
	var star gorm_model.Star
	err := DB.Model(gorm_model.Star{}).Where("username = ?", username).Where("session = ?", 0).Where("type = ?", 2).First(&star).Error
	if err != nil {
		return err
	}
	star.Type = 3
	err = DB.Save(&star).Error
	if err != nil {
		return err
	}
	return nil
}

// SelMax 查询session字段最大值
func SelMax() (int, error) {
	var maxnum int
	err := DB.Model(&gorm_model.Star{}).Select("MAX(session)").Scan(&maxnum).Error
	if err != nil {
		return 0, err
	}
	return maxnum, nil
}

// UpdateSession 更新字段
func UpdateSession(session int) error {
	err := DB.Model(&gorm_model.Star{}).Where("session = ? ", 0).Updates(map[string]interface{}{"session": session}).Error
	if err != nil {
		return err
	}
	return nil
}

// SelStar 查找指定届数的班级之星(前台后台公用)
func SelStar(session int, starType int, page int, limit int) ([]string, error) {
	var username []string
	condition := DB.Model(&gorm_model.Star{}).Where("session = ?", session).Where("type = ?", starType)
	if page == 0 && limit == 0 {
		err := condition.Select("username").Scan(&username).Error
		if err != nil {
			return nil, err
		}
	} else {
		if limit <= 0 {
			limit = 10
		}
		err := condition.Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&username).Error
		if err != nil {
			return nil, err
		}
	}

	return username, nil
}

// Selstarexit 查找这条数据是否存在数据库中
func Selstarexit(username string) (int64, error) {
	var number int64
	err := DB.Model(&gorm_model.Star{}).Where("username = ?", username).Where("type = ?", 1).Where("session = ?", 0).Count(&number).Error
	if err != nil {
		return 0, err
	}
	return number, nil
}

// SelStatus 查询管理员是否可以添加数据
func SelStatus(username string) (bool, error) {
	var status bool
	err := DB.Model(&gorm_model.UserCasbinRules{}).Where("c_username = ?", username).Select("status").Scan(&status).Error
	if err != nil {
		return false, err
	}
	return status, nil
}

// UpdateStatus 批量更新管理员的状态字段
func UpdateStatus() error {
	ok := true
	err := DB.Model(&gorm_model.UserCasbinRules{}).Where("status = ? ", ok).Updates(map[string]interface{}{"status": !ok}).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateOne 更新一个管理员的字段
func UpdateOne(username string) error {
	var user gorm_model.UserCasbinRules
	err := DB.Where("c_username = ?", username).First(&user).Error
	if err != nil {
		return err
	}
	user.Status = true
	err = DB.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// SelTimeStar 查询特定时期的成长之星
func SelTimeStar(starTime, endTime string, starType int, page int, limit int) ([]string, error) {
	var username []string
	start, err := time.Parse("2006-01-02", starTime)
	start = start.Add(-8 * time.Hour)
	end, err := time.Parse("2006-01-02", endTime)
	end = end.AddDate(0, 0, 1)
	end = end.Add(-8 * time.Hour)
	if page <= 0 || limit <= 0 {
		page = 1 // 将页码设置为第一页
		limit = 10
	}
	err = DB.Model(&gorm_model.Star{}).Where("session <> ?", 0).Where("type = ?", starType).Where("created_at BETWEEN ? AND ?", start, end).Offset((page - 1) * limit).Limit(limit).Select("username").Scan(&username).Error
	if err != nil {
		return nil, err
	}
	return username, nil
}

// SelTime 查询表的最大时间以及最小时间
func SelTime() (string, string, error) {
	var maxTime time.Time
	var minTime time.Time
	err := DB.Model(&gorm_model.Star{}).Where("session <> ?", 0).Select("MAX(created_at)").Scan(&maxTime).Error
	err = DB.Model(&gorm_model.Star{}).Select("MIN(created_at)").Scan(&minTime).Error
	if err != nil {
		return "", "", err
	}
	maxdate := maxTime.Format("2006-01-02")
	mindate := minTime.Format("2006-01-02")
	return maxdate, mindate, nil
}

// SelStuClass 查询跟管理员一个班级的姓名合集
func SelStuClass(class string) ([]string, error) {
	var stuSli []string
	err := DB.Model(&gorm_model.Star{}).Where("session = ?", 0).Where("username IN (?)", DB.Model(&gorm_model.User{}).Where("class = ?", class).Select("username")).Select("name").Scan(&stuSli).Error
	if err != nil {
		return nil, err
	}
	return stuSli, nil
}

// SelStuGrade 查询跟已推选未发布的年级之星姓名合集
func SelStuGrade(data time.Time, year int) ([]string, error) {
	var stuSli []string
	_, username, err := SelGradeId(data, year)
	err = DB.Model(&gorm_model.Star{}).Where("type = ?", 2).Where("session = ?", 0).Where("username IN (?)", username).Select("name").Scan(&stuSli).Error
	if err != nil {
		return nil, err
	}
	return stuSli, nil
}

// SelDataCollege 查询院级管理员已经推选的数据
func SelDataCollege() ([]string, error) {
	var name []string
	err := DB.Model(&gorm_model.Star{}).Where("session = ? AND type = ?", 0, 3).Select("name").Scan(&name).Error
	if err != nil {
		return nil, err
	}
	return name, nil
}
