package JoinAudit

import (
	"errors"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/utils"
	"time"
)

// 班级审核列表的返回
type ResStuMsg struct {
	ID                      uint
	Username                string  `json:"username" `
	UserClass               string  `json:"user_class"`
	Name                    string  `json:"name"`
	Gender                  string  `json:"gender"`
	Major                   string  `json:"major"`
	MoralCoin               float64 `json:"moral_coin" `                 //道德币
	ComprehensiveScore      float64 `json:"comprehensive_score" `        //综测成绩
	TrainScore              int     `json:"training_score" `             //培训测试成绩
	ClassIsPass             string  `json:"class_is_pass"`               //班级审核
	RulerIsPass             string  `json:"ruler_is_pass" `              //纪权部综测成绩审核结果
	OrganizerMaterialIsPass string  `json:"organizer_material_is_pass" ` //组织部材料审核结果
	OrganizerTrainIsPass    string  `json:"organizer_train_is_pass" `    //组织部培训审核结果
	Files                   []File  `json:"files"`                       //用户对应的文件
}

type ResActivityMsg struct {
	ID                    uint
	ActivityName          string    `json:"activity_name"`           //期数
	StartTime             time.Time `json:"start_time"`              //开始时间
	StopTime              time.Time `json:"stop_time"`               //结束时间
	PersonInCharge        string    `json:"person_in_charge"`        //负责任人
	RulerName             string    `json:"ruler_name"`              //纪权部综测审核人
	OrganizerMaterialName string    `json:"organizer_material_name"` //组织部材料审核人
	OrganizerTrainName    string    `json:"organizer_train_name"`    //组织部培训审核人
	IsShow                string    `json:"is_show"`
}
type ResList struct {
	List      []ResStuMsg    `json:"list"`
	Count     int64          `json:"count"`
	ClassList []string       `json:"class_list"`
	Activity  ResActivityMsg `json:"activity"`
}

type File struct {
	FileNote string     `json:"file_note"`
	List     []FileList `json:"list"`
}
type FileList struct {
	ID       uint
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}
type userClass struct {
	UserClass string
}

// 根据活动id和username查询用户文件
func resListWithClass(msgList []gorm_model.JoinAudit, ActivityMsg gorm_model.JoinAuditDuty, count int64) (resListWithAll ResList) {
	var ResListWithStuMsg []ResStuMsg
	ResListWithStuMsg = make([]ResStuMsg, 0)
	classList := make([]string, 0)
	userClassList, _ := mysql.GetUserClassMsgWithActivityID(ActivityMsg.ID)
	for _, v := range userClassList {
		classList = append(classList, v.UserClass)
	}
	for _, val := range msgList {
		fileList, _ := mysql.FilesList(val.Username, val.JoinAuditDutyID)
		filesList := make([]File, 0)
		applicationList := make([]FileList, 0)
		materialList := make([]FileList, 0)
		var resFile FileList
		for _, val := range fileList {
			switch val.Note {
			case "application":
				resFile.FilePath = val.FilePath
				resFile.FileName = val.FileName
				resFile.ID = val.ID
				applicationList = append(applicationList, resFile)
			case "material":
				resFile.FilePath = val.FilePath
				resFile.FileName = val.FileName
				resFile.ID = val.ID
				materialList = append(materialList, resFile)
			}
		}
		filesList = append(filesList, File{
			"application",
			applicationList,
		})
		filesList = append(filesList, File{
			"material",
			materialList,
		})
		StuMsg := ResStuMsg{
			ID:                      val.ID,
			Username:                val.Username,
			UserClass:               val.UserClass,
			Name:                    val.Name,
			Gender:                  val.Gender,
			MoralCoin:               val.MoralCoin,
			ComprehensiveScore:      val.ComprehensiveScore,
			ClassIsPass:             val.ClassIsPass,
			TrainScore:              val.TrainScore,
			RulerIsPass:             val.RulerIsPass,
			OrganizerMaterialIsPass: val.OrganizerMaterialIsPass,
			OrganizerTrainIsPass:    val.OrganizerTrainIsPass,
			Files:                   filesList,
		}
		ResListWithStuMsg = append(ResListWithStuMsg, StuMsg)
	}
	ResListWithActivityMsg := ResActivityMsg{
		ID:                    ActivityMsg.ID,
		ActivityName:          ActivityMsg.ActivityName,
		StartTime:             ActivityMsg.StartTime,
		StopTime:              ActivityMsg.StopTime,
		PersonInCharge:        ActivityMsg.PersonInCharge,
		RulerName:             ActivityMsg.RulerName,
		OrganizerMaterialName: ActivityMsg.OrganizerMaterialName,
		OrganizerTrainName:    ActivityMsg.OrganizerTrainName,
		IsShow:                ActivityMsg.IsShow,
	}
	//班级切片去重
	classList = utils.SliceUnique(classList)
	resListWithAll = ResList{
		List:      ResListWithStuMsg,
		Count:     count,
		ClassList: classList,
		Activity:  ResListWithActivityMsg,
	}
	return
}

// 根据json字段和对应
func ResListWithJSON(cr mysql.Pagination) (ResAllMsgList []ResList, err error) {
	var msgList []gorm_model.JoinAudit
	var count int64
	if !cr.All {
		_, _, openActivityMsg := mysql.OpenActivityStates()
		cr.ActivityID = openActivityMsg.ID
		msgList, count, err = mysql.ComList(gorm_model.JoinAudit{}, cr)
		if err != nil {
			err = errors.New("列表查询出错")
			return nil, err
		}
		ResAllMsgList = append(ResAllMsgList, resListWithClass(msgList, openActivityMsg, count))
	}
	if cr.All {
		activityList, err := mysql.AllActivity()
		if err != nil {
			err = errors.New("活动查询失败")
			return nil, err
		}
		if len(activityList) == 0 {
			err = errors.New("活动信息不存在")
			return nil, err
		}
		for _, v := range activityList {
			cr.ActivityID = v.ID
			msgList, count, err = mysql.ComList(gorm_model.JoinAudit{}, cr)
			if err != nil {
				err = errors.New("列表查询出错")
				return nil, err
			}
			ResAllMsgList = append(ResAllMsgList, resListWithClass(msgList, v, count))
		}
	}
	return
}
