package mysql

import (
	"fmt"
	_ "gorm.io/gorm/clause"
	"strings"
	"studentGrow/models/gorm_model"
	"studentGrow/models/jrx_model"
	"time"
)

// 将新的用户自述在mysql中进行更行
func UpdateSelfContent(id int, newSelfContent string) error {
	return DB.Table("users").Where("id = ?", id).Update("self_content", newSelfContent).Error
}

// 获取mysql中的用户自述
func GetSelfContent(id int) (string, error) {
	var users gorm_model.User
	err := DB.Where("id = ?", id).First(&users).Error // Unscoped()用于解除搜索时会自动加上deleted_at字段的限制
	return users.SelfContent, err
}

// 根据学号获取id
func GetIdByUsername(username string) (int, error) {
	var users gorm_model.User
	err := DB.Where("username = ?", username).First(&users).Error
	return int(users.ID), err
}

// 根据id获取姓名
func GetNameById(id int) (string, error) {
	var users gorm_model.User
	err := DB.Where("id = ?", id).First(&users).Error
	return users.Name, err
}

// 获取不同的班级
func GetDiffClass() ([]string, error) {
	var diffClassSlice []string
	err := DB.Model(&gorm_model.User{}).Select("class").Distinct("class").Where("LENGTH(class) = 9").Order("class ASC").Scan(&diffClassSlice).Error
	return diffClassSlice, err
}

// 添加单个学生
func AddSingleStudent(users *gorm_model.User) error {
	err := DB.Select("name", "username", "password", "class", "gender", "identity", "head_shot", "plus_time").Create(users).Error
	return err

}

// 添加单个学生记录
func AddSingleStudentRecord(addUserRecord *gorm_model.UserAddRecord) error {
	err := DB.Create(addUserRecord).Error
	return err
}

// 添加单个老师
func AddSingleTeacher(users *gorm_model.User) error {
	err := DB.Select("name", "username", "password", "gender", "identity", "head_shot").Create(users).Error
	return err
}

// 删除单个学生
func DeleteSingleUser(id int) error {
	err := DB.Exec("DELETE FROM users WHERE id = ?", id).Error
	return err
}

// 删除学生记录
func DeleteStudentRecord(deleteUserRecord *gorm_model.UserDeleteRecord) error {
	err := DB.Create(deleteUserRecord).Error
	return err
}

// 删除单个学生的管理员信息
func DeleteSingleUserManager(username string) error {
	err := DB.Unscoped().Model(&gorm_model.UserCasbinRules{}).Where("c_username = ?", username).Delete(nil).Error
	return err
}

// 封禁该用户
func BanStudent(id int) (int, error) {
	var temp int
	var users gorm_model.User
	DB.Take(&users, id)
	var err error
	if users.Ban == false {
		err = DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("ban", 1).Error
		temp = 1
	} else if users.Ban == true {
		err = DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("ban", 0).Error
		temp = 0
	}

	return temp, err
}

// 修改用户信息username
func ChangeStudentMessage(id int, users jrx_model.ChangeStuMesStruct) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Updates(users).Error
	return err
}

// 修改用户信息记录
func EditUserRecord(userEditRecord gorm_model.UserEditRecord) error {
	fmt.Println("userEditRecordStruct ： ", userEditRecord)
	err := DB.Create(&userEditRecord).Error
	return err
}

// 修改老师信息username
func ChangeTeacherMessage(id int, newTeacher jrx_model.ChangeTeacherMesStruct) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Select("name", "username", "password", "gender").Updates(&newTeacher).Error
	return err
}

// 将用户设置为管理员
func SetStuManager(username string, casbinCid string) error {
	// user表设置
	id, err := GetIdByUsername(username)
	if err != nil {
		return err
	}

	err = DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("is_manager", 1).Error
	if err != nil {
		return err
	}

	// casbin_ruler表设置
	casbinUser := gorm_model.UserCasbinRules{
		CUsername: username,
		CasbinCid: casbinCid,
	}
	err = DB.Create(&casbinUser).Error
	fmt.Println("casbinCid:", casbinUser.CasbinCid)
	if err != nil {
		return err
	}

	return err
}

// 将用户修改管理员等级
func ChangeStuManager(username string, casbinCid string) error {
	// user表设置
	id, err := GetIdByUsername(username)
	if err != nil {
		return err
	}

	err = DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("is_manager", 1).Error
	if err != nil {
		return err
	}

	// casbin_ruler表设置
	casbinUser := gorm_model.UserCasbinRules{
		CUsername: username,
		CasbinCid: casbinCid,
	}
	err = DB.Model(&casbinUser).Where("c_username = ?", username).Update("casbin_cid", casbinCid).Error
	fmt.Println("casbinCid:", casbinUser.CasbinCid)
	if err != nil {
		return err
	}

	return err
}

// 取消用户管理员身份
func CancelStuManager(username string, casbinCid string) error {
	// user表设置
	id, err := GetIdByUsername(username)
	if err != nil {
		return err
	}

	err = DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("is_manager", 0).Error
	if err != nil {
		return err
	}

	// casbin_ruler表设置
	err = DB.Unscoped().Model(&gorm_model.UserCasbinRules{}).Where("c_username = ?", username).Delete(nil).Error

	if err != nil {
		return err
	}

	return err
}

// 判断用户是否存在
func ExistedUsername(username string) error {
	err := DB.Where("username = ?", username).First(&gorm_model.User{}).Error
	return err
}

// 查询选中的用户
func QuerySelectedUser(usernameSlice []string) ([]gorm_model.User, error) {
	var users []gorm_model.User
	err := DB.Where("username IN (?)", usernameSlice).Find(&users).Error
	return users, err
}

func GetAllUserCount(identity string) (int64, error) {
	var user gorm_model.User
	var count int64
	err := DB.Model(&user).Where("identity = ?", identity).Count(&count).Error
	return count, err
}

// GetStuMesList 根据搜索框内容查询学生信息列表
func GetTeacherList(querySql string) ([]jrx_model.QueryTeacherStruct, error) {
	// 从mysql中获取数据到user表中
	var userSlice []jrx_model.QueryTeacherStruct

	err := DB.Raw(querySql).Find(&userSlice).Error
	if err != nil {
		return nil, err
	}

	return userSlice, nil
}

func GetManager(username string) (gorm_model.UserCasbinRules, error) {
	var casbinUser gorm_model.UserCasbinRules
	err := DB.Model(&gorm_model.UserCasbinRules{}).Where("c_username = ?", username).First(&casbinUser).Error
	return casbinUser, err
}

func GetUserListBySql(querySql string) ([]gorm_model.User, error) {
	var userSlice []gorm_model.User
	err := DB.Raw(querySql).Find(&userSlice).Error
	if err != nil {
		return nil, err
	}
	return userSlice, err
}

// 根据学号获取 managerType
func GetIsManagerByUsername(username string) (bool, error) {
	var users gorm_model.User
	err := DB.Where("username = ?", username).First(&users).Error
	return users.IsManager, err
}

func GetHomepageUserMesDao(id int) (*gorm_model.User, error) {
	var userMes gorm_model.User
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).First(&userMes).Error
	if err != nil {
		return nil, err
	}
	return &userMes, err
}

// 获取个人主页粉丝个数
func GetHomepageFansCountDao(id int) (int, error) {
	var count int64
	err := DB.Table("user_followers").Where("user_id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), err
}

// 获取个人主页关注个数
func GetHomepageConcernCountDao(id int) (int, error) {
	var count int64
	err := DB.Table("user_followers").Where("follower_id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), err
}

// 获取个人主页获赞个数(文章获赞)
func GetHomepageLikeCountDao(id int) (int, error) {
	var count int
	// Article表 满足 user_id = id 的所有行中 like_amount 列的值的总和
	err := DB.Model(&gorm_model.Article{}).Where("user_id = ?", id).Select("COALESCE(SUM(like_amount), 0)").Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

// 获取个人主页积分值
func GetHomepagePointDao(id int) (int, error) {
	var count int
	err := DB.Model(&gorm_model.UserPoint{}).Where("user_id = ?", id).Select("COALESCE(SUM(point), 0)").Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

// 修改个人主页用户个签
func UpdateHomepageMottoDao(id int, motto string) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("motto", motto).Error
	return err
}

// 修改个人主页电话号
func UpdateHomepagePhoneNumberDao(id int, phoneNumber string) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("phone_number", phoneNumber).Error
	return err
}

// 修改个人主页邮箱
func UpdateHomepageEmailDao(id int, eamil string) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("mail_box", eamil).Error
	return err
}

// 修改个人主页头像
func UpdateHeadshotDao(id int, url string) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("head_shot", url).Error
	return err
}

func GetFansIdListDao(id int) ([]int, error) {
	var fansId []int
	err := DB.Table("user_followers").Where("user_id", id).Pluck("follower_id", &fansId).Error
	return fansId, err
}

func GetFansListDao(fansId []int) ([]jrx_model.HomepageFanStruct, error) {
	var fansList []jrx_model.HomepageFanStruct
	err := DB.Table("users").Where("id IN (?)", fansId).Find(&fansList).Error

	fmt.Println("fanslist : ", fansList)
	return fansList, err
}

func GetConcernIdListDao(id int) ([]int, error) {
	var concernId []int
	err := DB.Table("user_followers").Where("follower_id", id).Pluck("user_id", &concernId).Error
	return concernId, err
}

func GetConcernListDao(concernId []int) ([]jrx_model.HomepageFanStruct, error) {
	var concernList []jrx_model.HomepageFanStruct
	err := DB.Table("users").Where("id IN (?)", concernId).Find(&concernList).Error
	fmt.Println("concernList : ", concernList)
	return concernList, err
}

func ChangeConcernDao(id int, otherId int) error {
	var count int64
	err := DB.Table("user_followers").Where("user_id = ? AND follower_id = ?", otherId, id).Count(&count).Error //我关注了他吗
	if err != nil {
		return err
	}
	if count > 0 {
		// 我关注了他
		err = DB.Table("user_followers").Where("user_id = ? AND follower_id = ?", otherId, id).Delete(nil).Error // 我取消关注他
	} else {
		// 我未关注他
		userFollower := struct {
			UserID     int
			FollowerID int
		}{
			UserID:     otherId,
			FollowerID: id,
		}
		err = DB.Table("user_followers").Create(&userFollower).Error
	}
	return err
}

func GetHistoryByArticleDao(id int, page int, limit int) ([]jrx_model.HomepageArticleHistoryStruct, error) {
	// 获取该用户阅读过的文章的id
	var articleIds []int
	err := DB.Model(&gorm_model.UserReadRecord{}).
		Select("article_id").
		Where("user_id = ?", id).
		Group("article_id").
		Order("MAX(created_at) DESC"). // 根据你的表结构调整排序字段
		Pluck("article_id", &articleIds).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("articleIds: ", articleIds)

	// 将 articleIds 转换为逗号分隔的字符串
	articleIdsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(articleIds)), ","), "[]")

	// 获取这些文章id的文章信息以及文章发布者的信息  	// 多表查询
	var homepageArticleHistoryList []jrx_model.HomepageArticleHistoryStruct
	query := `
    SELECT articles.id, articles.content, article_pics.pic, articles.comment_amount, articles.like_amount, users.head_shot, users.name
    FROM articles
    LEFT JOIN users ON articles.user_id = users.id
    LEFT JOIN article_pics ON article_pics.article_id = articles.id
    WHERE articles.id IN (?)
    ORDER BY FIND_IN_SET(articles.id, ?)
    LIMIT ? OFFSET ?
`
	err = DB.Raw(query, articleIds, articleIdsStr, limit, (page-1)*limit).Scan(&homepageArticleHistoryList).Error

	if err != nil {
		return nil, err
	}

	return homepageArticleHistoryList, err
}

func GetStarDao(id int, page int, limit int) ([]jrx_model.HomepageArticleHistoryStruct, error) {
	// 获取该用户收藏的文章的id
	var articleIds []int
	err := DB.Model(&gorm_model.UserCollectRecord{}).
		Where("user_id = ?", id).
		Pluck("article_id", &articleIds).Error
	if err != nil {
		return nil, err
	}

	// 获取这些文章id的文章信息以及文章发布者的信息  	// 多表查询
	var homepageArticleHistoryList []jrx_model.HomepageArticleHistoryStruct
	err = DB.Model(&gorm_model.Article{}).
		Select("articles.id, articles.content, article_pics.pic, articles.comment_amount, articles.like_amount, users.head_shot, users.name").
		Joins("LEFT JOIN users ON articles.user_id = users.id").
		Joins("LEFT JOIN article_pics ON article_pics.article_id = articles.id").
		Where("articles.id IN (?)", articleIds).
		Order("articles.created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&homepageArticleHistoryList).Error

	if err != nil {
		return nil, err
	}

	return homepageArticleHistoryList, err
}

func GetClassById(id int) (string, error) {
	var users gorm_model.User
	err := DB.Where("id = ?", id).First(&users).Error
	return users.Class, err
}

func GetPlusTimeById(id int) (time.Time, error) {
	var users gorm_model.User
	err := DB.Where("id = ?", id).First(&users).Error
	return users.PlusTime, err
}

func GetClassmateList(class string) ([]jrx_model.HomepageClassmateStruct, error) {
	var classmateList []jrx_model.HomepageClassmateStruct
	err := DB.Table("users").Where("class = ?", class).Find(&classmateList).Error
	return classmateList, err
}

func GetArticleDao(id int, page int, limit int) ([]jrx_model.HomepageArticleHistoryStruct, error) {
	// 获取该用户发布的文章的id
	var articles []jrx_model.HomepageArticleHistoryStruct
	err := DB.Model(&gorm_model.Article{}).
		Select("articles.id, articles.content, articles.comment_amount, articles.like_amount, articles.collect_amount, articles.status, articles.topic, articles.created_at, articles.ban").
		Where("articles.user_id = ?", id).
		Order("articles.created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&articles).Error
	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].PostTime = articles[i].CreatedAt.Format("2006-01-02")

		var tagIds []int
		err := DB.Table("article_tags").
			Where("article_id = ?", articles[i].ID). // 根据文章名获取tag_id
			Pluck("tag_id", &tagIds).Error
		if err != nil {
			return nil, err
		}

		var tagNames []string
		err = DB.Table("tags").
			Where("id in (?)", tagIds). // 根据文章名获取tag_id
			Pluck("tag_name", &tagNames).Error
		if err != nil {
			return nil, err
		}

		articles[i].ArticleTags = tagNames
	}

	return articles, err
}

func ChangeArticleStatusDao(articleId int, ArticleStatus bool) error {
	err := DB.Table("articles").Where("id = ?", articleId).Update("status", ArticleStatus).Error
	return err
}

func BanUserDao(banId int, banEndTime time.Time) error {
	var user gorm_model.User
	user.Ban = true
	user.UserBanEndTime = banEndTime
	err := DB.Model(&gorm_model.User{}).Where("id = ?", banId).Updates(user).Error
	return err
}

// 封禁操作记录
func BanUserRecordDao(banId int, userId int, banTime int) error {
	var userBanRecord gorm_model.UserBanRecord
	userBanRecord.UserID = userId
	userBanRecord.BanId = banId
	userBanRecord.BanTime = banTime
	err := DB.Model(&gorm_model.UserBanRecord{}).Create(&userBanRecord).Error
	return err
}

func UnbanUserDao(banId any) error {
	var user gorm_model.User
	user.Ban = false
	user.UserBanEndTime = time.Now()
	err := DB.Model(&gorm_model.User{}).Where("id = ?", banId).Update("ban", false).Update("user_ban_end_time", time.Now()).Error
	return err
}

func GetClassList() ([]string, error) {
	var classes []string
	err := DB.Model(&gorm_model.User{}).Where("LENGTH(class) = 9").Distinct("class").Order("class ASC").Pluck("class", &classes).Error
	return classes, err
}

func GetFansListIsConcernDao(fansList []jrx_model.HomepageFanStruct, id int) ([]jrx_model.HomepageFanStruct, error) {
	for i, v := range fansList {
		otherId, err := GetIdByUsername(v.Username)
		if err != nil {
			return nil, err
		}

		// 如果列表中显示了自己的账号
		if otherId == id {
			fansList[i].IsConcern = ""
			continue
		}

		isConcern, err := GetIsConcernDao(id, otherId)
		if err != nil {
			return nil, err
		}

		if isConcern {
			// 我关注了他
			fansList[i].IsConcern = "已关注"
		} else {
			// 我未关注他
			fansList[i].IsConcern = "关注"
		}
	}
	return fansList, nil
}

func GetIsConcernDao(id int, otherId int) (bool, error) {
	// 判断我是否关注这名粉丝
	var count int64
	err := DB.Table("user_followers").Where("user_id = ? AND follower_id = ?", otherId, id).Count(&count).Error //我关注了他吗
	if err != nil {
		return false, err
	}

	if count > 0 {
		// 我关注了他
		return true, nil
	} else {
		// 我未关注他
		return false, err
	}

}

func GetTracksDao(id int, page int, limit int) ([]jrx_model.HomepageTrack, error) {
	var tracks []jrx_model.HomepageTrack
	err := DB.Raw(`
		SELECT ID, CONTENT, NAME, like_amount, comment_amount, 'articles' AS I_TYPE, CreatedAt FROM articles where user_id = id
		UNION ALL
		SELECT a.ID, a.CONTENT, a.NAME, a.like_amount, a.comment_amount, c.comment_content, 'comments' AS I_TYPE, c.CreatedAt FROM articles a where user_id = id
		LEFT JOIN comments c ON a.id = c.article_id
	`).Order("CreatedAt DESC").Scan(&tracks).Error

	return tracks, err
}

func GetTotalPointsByUserAndTopic(id int, name string) (int, error) {
	var totalPoints int
	// 在topic表中根据name查到主键t_id，再去point表中满足id=user_id，t_id=topic_id条件的所有point字段值的总和
	err := DB.Model(&gorm_model.UserPoint{}).
		Select("SUM(user_points.point)").
		Joins("JOIN topics ON user_points.topic_id = topics.id").
		Where("topics.topic_name = ? AND user_points.user_id = ?", name, id).
		Row().Scan(&totalPoints)
	if err != nil {
		if strings.Contains(err.Error(), "converting NULL to int is unsupported") {
			fmt.Println("数据库中该记录的point字段为NULL，为避免error，手动返回0")
			return 0, nil
		} else {
			return 0, err
		}
	}

	return totalPoints, nil
}

func GetUser(id int) (gorm_model.User, error) {
	var user gorm_model.User
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetPasswordById(id int) (pwdFromMysql string, err error) {
	var users gorm_model.User
	err = DB.Where("id = ?", id).First(&users).Error
	return users.Password, err
}

func SaveAdviceDao(advice gorm_model.Advice) error {
	err := DB.Create(&advice).Error
	return err
}

func UpdatePassword(id int, newPwd string) error {
	err := DB.Model(&gorm_model.User{}).Where("id = ?", id).Update("password", newPwd).Error
	if err != nil {
		return err
	}
	return nil
}

func GetClassListByPlusTimeYear(plusTimeYear string) ([]string, error) {
	var classes []string
	err := DB.Model(&gorm_model.User{}).Where("LENGTH(class) = ? AND YEAR(plus_time) = ?", 9, plusTimeYear).Distinct("class").Order("class ASC").Pluck("class", &classes).Error
	return classes, err
}

func DeleteStar(username string) error {
	err := DB.Model(&gorm_model.Star{}).Where("username = ?", username).Delete(nil).Error
	return err
}
