package service

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/models/jrx_model"
	"studentGrow/utils/fileProcess"
	"time"
)

func GetHomepageMesService(username string) (*jrx_model.HomepageMesStruct, error) {
	homepageMes := &jrx_model.HomepageMesStruct{}

	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	// 从user表中获取数据
	userMes, err := mysql.GetHomepageUserMesDao(id)
	if err != nil {
		return nil, err
	}

	// 从其他表中获取数据
	userFans, err := mysql.GetHomepageFansCountDao(id)
	if err != nil {
		return nil, err
	}

	userConcern, err := mysql.GetHomepageConcernCountDao(id)
	if err != nil {
		return nil, err
	}

	userLike, err := mysql.GetHomepageLikeCountDao(id)
	if err != nil {
		return nil, err
	}

	userTopicPointStruct, err := GetTopicPointsService(username)
	if err != nil {
		return nil, err
	}

	// 将获得的数据存储到 homepageMesStruct中
	homepageMes.Username = userMes.Username
	homepageMes.Ban = userMes.Ban
	homepageMes.Name = userMes.Name
	homepageMes.UserHeadShot = userMes.HeadShot
	homepageMes.UserMotto = userMes.Motto
	homepageMes.UserFans = userFans
	homepageMes.UserConcern = userConcern
	homepageMes.UserLike = userLike
	homepageMes.Point = userTopicPointStruct.TotalPoint
	homepageMes.UserClass = userMes.Class

	fmt.Printf("homepageMes2: %+v\n", homepageMes)

	return homepageMes, err
}

func UpdateHomepageMottoService(username string, motto string) error {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	err = mysql.UpdateHomepageMottoDao(id, motto)
	if err != nil {
		return err
	}

	return nil
}

func UpdateHomepagePhoneNumberService(username string, phone_number string) error {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	err = mysql.UpdateHomepagePhoneNumberDao(id, phone_number)
	if err != nil {
		return err
	}

	return nil
}

func UpdateHomepageEmailService(username string, email string) error {

	// todo 前端传id

	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	err = mysql.UpdateHomepageEmailDao(id, email)
	if err != nil {
		return err
	}

	return nil
}

func GetHomepageUserDataService(username string) (*jrx_model.HomepageDataStruct, error) {
	userData := &jrx_model.HomepageDataStruct{}
	userDataTemp := &gorm_model.User{}

	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	userDataTemp, err = mysql.GetHomepageUserMesDao(id)
	if err != nil {
		return nil, err
	}

	userData.Name = userDataTemp.Name
	userData.UserHeadShot = userDataTemp.HeadShot
	userData.UserGender = userDataTemp.Gender
	userData.UserClass = userDataTemp.Class
	userData.UserMotto = userDataTemp.Motto
	userData.PhoneNumber = userDataTemp.PhoneNumber
	userData.UserEmail = userDataTemp.MailBox
	userData.UserYear = userDataTemp.PlusTime.Format("2006")

	fmt.Printf("userData : %+v\n", userData)

	return userData, err
}

// 更新个人主页头像
func UpdateHeadshotService(file *multipart.FileHeader, username string) (string, error) {

	// todo 图片上传
	fmt.Println("成功进入业务")
	url, err := fileProcess.UploadFile("image", file)
	if err != nil {
		return "", err
	}

	fmt.Println("成功获得url")

	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return "", err
	}

	err = mysql.UpdateHeadshotDao(id, url)
	if err != nil {
		return "", err
	}

	fmt.Println("成功更换头像DAO")

	return url, err
}

func GetFansListService(username string, tokenUsername string) ([]jrx_model.HomepageFanStruct, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	fansIdList, err := mysql.GetFansIdListDao(id)
	if err != nil {
		return nil, err
	}

	fansList, err := mysql.GetFansListDao(fansIdList)
	if err != nil {
		return nil, err
	}

	tokenId, err := mysql.GetIdByUsername(tokenUsername)
	if err != nil {
		return nil, err
	}

	fansList, err = mysql.GetFansListIsConcernDao(fansList, tokenId)
	if err != nil {
		return nil, err
	}

	return fansList, err
}

func GetConcernListService(username string, tokenUsername string) ([]jrx_model.HomepageFanStruct, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	concernIdList, err := mysql.GetConcernIdListDao(id)
	if err != nil {
		return nil, err
	}

	concernList, err := mysql.GetConcernListDao(concernIdList)
	if err != nil {
		return nil, err
	}

	tokenId, err := mysql.GetIdByUsername(tokenUsername)
	if err != nil {
		return nil, err
	}

	concernList, err = mysql.GetFansListIsConcernDao(concernList, tokenId)
	if err != nil {
		return nil, err
	}

	return concernList, err
}

func ChangeConcernService(username string, othername string) error {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	otherId, err := mysql.GetIdByUsername(othername)
	if err != nil {
		return err
	}

	err = mysql.ChangeConcernDao(id, otherId)
	if err != nil {
		return err
	}

	return err
}

func GetHistoryService(page int, limit int, username string) ([]jrx_model.HomepageArticleHistoryStruct, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	homepageArticleHistoryList, err := mysql.GetHistoryByArticleDao(id, page, limit)
	if err != nil {
		return nil, err
	}

	return homepageArticleHistoryList, err
}

func GetStarService(page int, limit int, username string) ([]jrx_model.HomepageArticleHistoryStruct, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	homepageStarList, err := mysql.GetStarDao(id, page, limit)
	if err != nil {
		return nil, err
	}

	return homepageStarList, err
}

func GetClassmateListService(username string) ([]jrx_model.HomepageClassmateStruct, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	class, err := mysql.GetClassById(id)
	if err != nil {
		return nil, err
	}

	classmateList, err := mysql.GetClassmateList(class)
	if err != nil {
		return nil, err
	}

	return classmateList, err
}

func GetArticleService(page int, limit int, username string, tokenUsername string) ([]jrx_model.HomepageArticleHistoryStruct, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	homepageStarList, err := mysql.GetArticleDao(id, page, limit)
	if err != nil {
		return nil, err
	}

	if username != tokenUsername {
		for i, _ := range homepageStarList {
			if homepageStarList[i].Ban == true {
				homepageStarList[i].Content = "该帖子已被封禁，内容不予展示。"
			}
		}
	}

	return homepageStarList, err
}

func ChangeArticleStatusService(articleId int, articleStatus bool) error {
	err := mysql.ChangeArticleStatusDao(articleId, articleStatus)
	if err != nil {
		return err
	}
	return err
}

func BanHomepageUserService(banUsername string, banTime int, username string) error {
	banId, err := mysql.GetIdByUsername(banUsername)
	if err != nil {
		return err
	}

	userId, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	banEndtime := time.Now().Add(time.Duration(banTime) * 24 * time.Hour)

	// 封禁用户
	err = mysql.BanUserDao(banId, banEndtime)
	if err != nil {
		return err
	}

	// 增加封禁记录
	err = mysql.BanUserRecordDao(banId, userId, banTime)
	if err != nil {
		return err
	}

	return err
}

func UnbanHomepageUserService(banUsername string, username string) error {
	banId, err := mysql.GetIdByUsername(banUsername)
	if err != nil {
		return err
	}

	userId, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	// 解封用户
	err = mysql.UnbanUserDao(banId)
	if err != nil {
		return err
	}

	// 解封记录
	err = mysql.BanUserRecordDao(banId, userId, 0)
	if err != nil {
		return err
	}

	return err
}

func GetIsConcernService(username string, otherUsername string) (bool, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return false, err
	}

	otherId, err := mysql.GetIdByUsername(otherUsername)
	if err != nil {
		return false, err
	}

	isConcern, err := mysql.GetIsConcernDao(id, otherId)
	if err != nil {
		return false, err
	}

	return isConcern, err
}

func GetTracksService(page int, limit int, username string) ([]jrx_model.HomepageTrack, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return nil, err
	}

	Tracks, err := mysql.GetTracksDao(id, page, limit)
	if err != nil {
		return nil, err
	}

	return Tracks, nil
}

func GetTopicPointsService(username string) (jrx_model.HomepageTopicPoint, error) {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return jrx_model.HomepageTopicPoint{}, err
	}
	// todo 数组
	// 1.学习成绩 2.获奖荣誉 3.工作履历 4.社会实践 5.志愿公益 6.文体活动 7.生活日常
	everyPointMap := make(map[string]int)
	everyPointMap["学习成绩"] = 0
	everyPointMap["获奖荣誉"] = 0
	everyPointMap["工作履历"] = 0
	everyPointMap["社会实践"] = 0
	everyPointMap["志愿公益"] = 0
	everyPointMap["文体活动"] = 0
	everyPointMap["生活日常"] = 0

	for k, _ := range everyPointMap {
		everyPointMap[k], err = mysql.GetTotalPointsByUserAndTopic(id, k)
		if everyPointMap[k] > 200 {
			everyPointMap[k] = 200
		}
		if everyPointMap[k] < 0 {
			everyPointMap[k] = 0
		}
		if err != nil {
			return jrx_model.HomepageTopicPoint{}, err
		}
		fmt.Println(k, everyPointMap[k])
	}

	fmt.Println(everyPointMap)

	//studyPoint, err := mysql.GetTotalPointsByUserAndTopic(id, "学习成绩")
	//fmt.Println("studyPoint", studyPoint)
	//if studyPoint > 200 {
	//	studyPoint = 200
	//}
	//if studyPoint < 0 {
	//	studyPoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}
	//
	//honorPoint, err := mysql.GetTotalPointsByUserAndTopic(id, "获奖荣誉")
	//fmt.Println("honorPoint", honorPoint)
	//if honorPoint > 200 {
	//	honorPoint = 200
	//}
	//if honorPoint < 0 {
	//	honorPoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}
	//
	//workPoint, err := mysql.GetTotalPointsByUserAndTopic(id, "工作履历")
	//fmt.Println("workPoint", workPoint)
	//if workPoint > 200 {
	//	workPoint = 200
	//}
	//if workPoint < 0 {
	//	workPoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}
	//
	//socialPoint, err := mysql.GetTotalPointsByUserAndTopic(id, "社会实践")
	//fmt.Println("socialPoint", socialPoint)
	//if socialPoint > 200 {
	//	socialPoint = 200
	//}
	//if socialPoint < 0 {
	//	socialPoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}
	//
	//volunteerPoint, err := mysql.GetTotalPointsByUserAndTopic(id, "志愿公益")
	//fmt.Println("volunteerPoint", volunteerPoint)
	//if volunteerPoint > 200 {
	//	volunteerPoint = 200
	//}
	//if volunteerPoint < 0 {
	//	volunteerPoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}
	//
	//sportPoint, err := mysql.GetTotalPointsByUserAndTopic(id, "文体活动")
	//fmt.Println("sportPoint", sportPoint)
	//if sportPoint > 200 {
	//	sportPoint = 200
	//}
	//if sportPoint < 0 {
	//	sportPoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}
	//
	//lifePoint, err := mysql.GetTotalPointsByUserAndTopic(id, "生活日常")
	//fmt.Println("lifePoint", lifePoint)
	//if lifePoint > 200 {
	//	lifePoint = 200
	//}
	//if lifePoint < 0 {
	//	lifePoint = 0
	//}
	//if err != nil {
	//	return jrx_model.HomepageTopicPoint{}, err
	//}

	totalPointFloat := 0.25*float64(everyPointMap["学习成绩"]) + 0.2*float64(everyPointMap["获奖荣誉"]+everyPointMap["工作履历"]) + 0.1*float64(everyPointMap["社会实践"]+everyPointMap["志愿公益"]+everyPointMap["文体活动"]) + 0.05*float64(everyPointMap["生活日常"])
	totalPointFloat = math.Round(totalPointFloat*10) / 10

	var homepageTopicPoint jrx_model.HomepageTopicPoint
	homepageTopicPoint = jrx_model.HomepageTopicPoint{
		StudyPoint:     everyPointMap["学习成绩"],
		HonorPoint:     everyPointMap["获奖荣誉"],
		WorkPoint:      everyPointMap["工作履历"],
		SocialPoint:    everyPointMap["社会实践"],
		VolunteerPoint: everyPointMap["志愿公益"],
		SportPoint:     everyPointMap["文体活动"],
		LifePoint:      everyPointMap["生活日常"],
		TotalPoint:     totalPointFloat,
	}

	return homepageTopicPoint, nil
}

// 检查输入的原密码与数据库中的原密码是否一致，是则返回ture
func CheckPassword(id int, oldPwd string) (bool, error) {
	PwdFromMysql, err := mysql.GetPasswordById(id)
	if err != nil {
		return false, err
	}

	if PwdFromMysql != oldPwd {
		return false, nil
	}

	return true, nil
}

func ChangePasswordService(username string, oldPwd string, newPwd string) error {
	id, err := mysql.GetIdByUsername(username)
	if err != nil {
		return err
	}

	pwdOk, err := CheckPassword(id, oldPwd)
	if err != nil {
		return err
	}

	if pwdOk {
		err = mysql.UpdatePassword(id, newPwd)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("输入密码与实际密码不符")
}

// 存储用户反馈到数据库中
func SaveAdviceService(username string, advice string) error {
	adviceStruct := gorm_model.Advice{
		Username: username,
		Advice:   advice,
	}

	err := mysql.SaveAdviceDao(adviceStruct)
	if err != nil {
		return err
	}

	return nil
}
