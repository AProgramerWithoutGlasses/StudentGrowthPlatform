package mysql

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
	"time"
)

type Pagination struct {
	Page                    int    `json:"page"`
	Label                   string `json:"Label"`
	Limit                   int    `json:"limit"`
	Sort                    string `json:"sort"`
	PersonInCharge          string `json:"person_in_charge"`
	ActivityName            string `json:"activity_name"`
	IsShow                  string `json:"is_show"`
	StartTime               string `json:"start_time"`
	StopTime                string `json:"stop_time"`
	UserClass               string `json:"user_class"`
	Gender                  string `json:"gender"`
	Name                    string `json:"name"`
	Username                string `json:"username"`
	ClassIsPass             string `json:"class_is_pass" `              //班级审核
	RulerIsPass             string `json:"ruler_is_pass" `              //纪权部综测成绩审核结果
	OrganizerMaterialIsPass string `json:"organizer_material_is_pass" ` //组织部材料审核结果
	OrganizerTrainIsPass    string `json:"organizer_train_is_pass" `    //组织部培训审核结果
	All                     bool   `json:"all"`
	ActivityID              uint   `json:"activity_id"`
}

// OpenActivityMsg 查询活动的开放信息
func OpenActivityMsg() (bool, string, gorm_model.JoinAuditDuty) {
	var count int64
	var ActivityMsg gorm_model.JoinAuditDuty
	var FailActivityMsg gorm_model.JoinAuditDuty
	DB.Where("is_show = ?", "true").Find(&ActivityMsg).Count(&count)
	if count < 1 {
		return false, "不存在开放活动", FailActivityMsg
	}
	if count > 1 {
		return false, "开放活动数量异常", FailActivityMsg
	}
	format := "2006-01-02 15:04:05"
	now := time.Now().Format(format)
	startTime := ActivityMsg.StartTime.Format(format)
	stopTime := ActivityMsg.StopTime.Format(format)
	if now <= startTime {
		return false, "活动未到开放时间", ActivityMsg
	}
	if now >= stopTime {
		return false, "活动已结束", ActivityMsg
	}
	return true, "活动已开放", ActivityMsg
}

// 活动信息根据开放状态返回信息，不受时间限制
func OpenActivityStates() (string, string, gorm_model.JoinAuditDuty) {
	var count int64
	var ActivityMsg gorm_model.JoinAuditDuty
	var FailActivityMsg gorm_model.JoinAuditDuty
	DB.Where("is_show = ?", "true").Find(&ActivityMsg).Count(&count)
	if count < 1 {
		return "false", "不存在开放活动", FailActivityMsg
	}
	if count > 1 {
		return "false", "开放活动数量异常", FailActivityMsg
	}
	format := "2006-01-02 15:04:05"
	now := time.Now().Format(format)
	startTime := ActivityMsg.StartTime.Format(format)
	stopTime := ActivityMsg.StopTime.Format(format)
	if now <= startTime {
		return "true", "活动未到开放时间", ActivityMsg
	}
	if now >= stopTime {
		return "true", "活动已结束", ActivityMsg
	}
	return "true", "活动已开放", ActivityMsg
}

// StuFormMsg 查询提交过的信息
func StuFormMsg(username string, activityID uint) (isExist bool, stuMsg gorm_model.JoinAudit) {
	err := DB.Where("username = ? AND join_audit_duty_id = ?", username, activityID).First(&stuMsg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		isExist = false
		return
	}
	isExist = true
	return
}

// ComList 分页查询
func ComList[T any](model T, pagMsg Pagination) (list []T, count int64, err error) {
	if pagMsg.Sort == "asc" {
		pagMsg.Sort = "created_at asc"
	} else {
		pagMsg.Sort = "created_at desc"
	}
	db := DB
	//判断来源，匹配特殊的选项
	switch pagMsg.Label {
	case "ActivityList":
		if pagMsg.IsShow != "" {
			db = db.Where("is_show = ?", pagMsg.IsShow)
		}
		if pagMsg.ActivityName != "" {
			db = db.Where("activity_name like ?", "%"+pagMsg.ActivityName+"%")
		}
		if pagMsg.PersonInCharge != "" {
			db = db.Where("person_in_charge like ?", "%"+pagMsg.PersonInCharge+"%")
		}
		//判断时间区间
		if pagMsg.StartTime != "" && pagMsg.StopTime != "" {
			db = db.Where("stop_time <= ?", pagMsg.StopTime)
			db = db.Where("start_time >= ?", pagMsg.StartTime)
		} else if pagMsg.StartTime != "" && pagMsg.StopTime == "" {
			db = db.Where("start_time >= ?", pagMsg.StartTime)
		} else if pagMsg.StartTime == "" && pagMsg.StopTime != "" {
			db = db.Where("stop_time <= ?", pagMsg.StopTime)
		}
	case "ClassApplicationList":
		db = db.Where("join_audit_duty_id = ?", pagMsg.ActivityID)
		if pagMsg.Gender != "" {
			db = db.Where("gender = ? ", pagMsg.Gender)
		}
		if pagMsg.ClassIsPass != "" {
			db = db.Where("class_is_pass = ?", pagMsg.ClassIsPass)
		}
	case "ActivityRulerList":
		db = db.Where("join_audit_duty_id = ?", pagMsg.ActivityID)
		db = db.Where("class_is_pass = ?", "true")
		if pagMsg.RulerIsPass != "" {
			db = db.Where("ruler_is_pass = ?", pagMsg.RulerIsPass)
		}
		if pagMsg.OrganizerMaterialIsPass != "" {
			db = db.Where("organizer_material_is_pass= ?", pagMsg.OrganizerMaterialIsPass)
		}
	case "ActivityOrganizerTrainList":
		db = db.Where("join_audit_duty_id = ?", pagMsg.ActivityID)
		db = db.Where("class_is_pass = ? AND ruler_is_pass = ? AND organizer_material_is_pass = ?", "true", "true", "true")
		if pagMsg.OrganizerTrainIsPass != "" {
			db.Where("organizer_train_is_pass = ?", pagMsg.OrganizerTrainIsPass)
		}
	}

	//模糊搜索名字条件
	if pagMsg.Name != "" {
		db = db.Where("name like ?", "%"+pagMsg.Name+"%")
	}
	if pagMsg.Username != "" {
		db = db.Where("username like ?", "%"+pagMsg.Username+"%")
	}
	if pagMsg.UserClass != "" {
		db = db.Where("user_class like ?", "%"+pagMsg.UserClass+"%")
	}
	offset := (pagMsg.Page - 1) * pagMsg.Limit
	if offset < 0 {
		offset = 0
	}
	//符合条件的总数
	count = db.Find(&list).RowsAffected
	//分页查询
	if pagMsg.All {
		err = db.Order(pagMsg.Sort).Find(&list).Error
		return
	}
	err = db.Limit(pagMsg.Limit).Offset(offset).Order(pagMsg.Sort).Find(&list).Error
	return
}

// 申请人信息保存和更新
func UpdateStuForm(model gorm_model.JoinAudit) (err error) {
	err = DB.Model(&model).Updates(&model).Error
	return
}

// 申请人信息保存
func CreateStuForm(model gorm_model.JoinAudit) (err error) {
	err = DB.Create(&model).Error
	return
}

// 根据用户学号和活动ID获取学生入团表单信息
func GetStuFromMsg(username string, activityMsgID uint) (stuFromMsg gorm_model.JoinAudit, err error) {
	err = DB.Where("username = ? AND join_audit_duty_id = ?", username, activityMsgID).Take(&stuFromMsg).Error
	return
}

// 根据活动ID跟username查找
func FilesList(username string, joinAuditDutyID uint) (fileList []gorm_model.JoinAuditFile, err error) {
	err = DB.Where("username = ? and join_audit_duty_id = ?", username, joinAuditDutyID).Find(&fileList).Error
	return
}

// 根据username和活动ID删除对应文件
func DelUserFile(username string, activityMsgID uint) (err error) {
	var imageList []gorm_model.JoinAuditFile
	count := DB.Find(&imageList, "username = ? AND join_audit_duty_id = ?", username, activityMsgID).RowsAffected
	if count != 0 {
		err = DB.Delete(&imageList).Error
	}
	return
}

// 根据文件id和删除文件
func DelFileWithID(id int) (count int64) {
	count = DB.Delete(&gorm_model.JoinAuditFile{}, "id = ? ", id).RowsAffected
	return
}

// 更改审核结果,返回更新后的数据
func IsPass(id int, column string, isPass string) (updatedJoinAudit gorm_model.JoinAudit) {
	DB.Model(&gorm_model.JoinAudit{}).Where("id = ?", id).Update(column, isPass)
	DB.Select(column).Where("id = ?", id).First(&updatedJoinAudit)
	return
}

// 根据ID删除活动信息
func DelActivityWithID(id int) (err error) {
	err = DB.Delete(&gorm_model.JoinAuditDuty{}, "id =? ", id).Error
	return
}

// 根据ID查询活动是否存在
func GetActivityNumberWithID(id int) (count int64) {
	DB.Model(&gorm_model.JoinAuditDuty{}).Where("id = ?", id).Count(&count)
	return
}

// 创建活动
func CreatActivity(activityMsg gorm_model.JoinAuditDuty) (err error) {
	err = DB.Create(&activityMsg).Error
	return
}

// 根据id更新活动
func UpdateActivityWithID(activityMsg gorm_model.JoinAuditDuty, id uint) (err error) {
	err = DB.Where("id =? ", id).Updates(&activityMsg).Error
	return
}

// 更新活动状态前关闭所有活动
func CloseAllActivity() (err error) {
	err = DB.Model(&gorm_model.JoinAuditDuty{}).Where("is_show = ?", "true").Update("is_show", "false").Error
	return
}

// 查询所有活动的的活动信息
func AllActivity() (activityList []gorm_model.JoinAuditDuty, err error) {
	err = DB.Find(&activityList).Error
	return
}

// 根据id修改分数
func UpdateTrainScoreWithID(id int, score float64) (err error) {
	err = DB.Model(&gorm_model.JoinAudit{}).Where("id = ?", id).Update("train_score", score).Error
	return
}

// 根据id查分数
func GetTrainScoreWithID(id int) (trainScore gorm_model.JoinAudit, err error) {
	err = DB.Select("train_score").Where("id= ?", id).Find(&trainScore).Error
	return
}

// 根据活动ID查询班级信息
type userClass struct {
	UserClass string
}

func GetUserClassMsgWithActivityID(JoinAuditDutyID uint) (classList []userClass, err error) {
	classList = make([]userClass, 0)
	err = DB.Model(&gorm_model.JoinAudit{}).Distinct("user_class").Where("join_audit_duty_id=?", JoinAuditDutyID).Scan(&classList).Error
	return
}

// 查询能否导出信息

func UserListWithOrganizer(activityID int, curMenu string) (userMsg []map[string]interface{}, classList []string) {
	userMsg = make([]map[string]interface{}, 0)
	classList = make([]string, 0)
	var rulerNullCount int64
	var organizerNullCount int64
	switch curMenu {
	case "organizer":
		DB.Model(gorm_model.JoinAudit{}).Where("class_is_pass = ? and join_audit_duty_id = ? and ruler_is_pass = ?", "true", activityID, "null").Count(&rulerNullCount)
		DB.Model(gorm_model.JoinAudit{}).Where(" class_is_pass = ? and join_audit_duty_id = ? and organizer_material_is_pass = ?", "true", activityID, "null").Count(&organizerNullCount)
	case "organizerFinish":
		DB.Model(gorm_model.JoinAudit{}).Where("class_is_pass = ? and ruler_is_pass = ? and organizer_material_is_pass = ? and join_audit_duty_id = ? and organizer_train_is_pass = ?", "true", "true", "true", activityID, "null").Count(&organizerNullCount)
	}
	if rulerNullCount > 0 || organizerNullCount > 0 {
		return userMsg, classList
	}
	switch curMenu {
	case "organizer":
		DB.Model(gorm_model.JoinAudit{}).Select("username", "name", "user_class").Where("join_audit_duty_id = ? and class_is_pass = ? and ruler_is_pass = ? and organizer_material_is_pass = ?", activityID, "true", "true", "true").Scan(&userMsg)
		DB.Model(gorm_model.JoinAudit{}).Distinct("user_class").Where("join_audit_duty_id = ? and class_is_pass = ? and ruler_is_pass = ? and organizer_material_is_pass = ?", activityID, "true", "true", "true").Scan(&classList)
	case "organizerFinish":
		DB.Model(gorm_model.JoinAudit{}).Select("username", "name", "user_class").Where("join_audit_duty_id = ? and class_is_pass = ? and ruler_is_pass = ? and organizer_material_is_pass = ? and organizer_train_is_pass = ?", activityID, "true", "true", "true", "true").Scan(&userMsg)
		DB.Model(gorm_model.JoinAudit{}).Distinct("user_class").Where("join_audit_duty_id = ? and class_is_pass = ? and ruler_is_pass = ? and organizer_material_is_pass = ? and organizer_train_is_pass = ?", activityID, "true", "true", "true", "true").Scan(&classList)
	}
	return userMsg, classList
}

// 查询班级列表
func AllClassList() (classList []string) {
	classList = make([]string, 0)
	DB.Model(&gorm_model.User{}).Distinct("class").Where("class <> ?", "").Find(&classList)
	fmt.Println(classList)
	return
}
