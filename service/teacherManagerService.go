package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"gorm.io/gorm"
	"strconv"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/models/jrx_model"
)

// 查询老师
func QueryTeacher(queryTeacherParama jrx_model.QueryTeacherParamStruct) ([]jrx_model.QueryTeacherResStruct, int, error) {
	// 获取老师总数量(长度)
	allTeacherCount64, err := mysql.GetAllUserCount("老师")
	if err != nil {
		return nil, 0, err
	}

	// 总数量类型转换
	allTeacherCount := int(allTeacherCount64)

	// 获取查询老师sql语句
	queryTeacherSql := GetQueryTeacherSql(queryTeacherParama)
	fmt.Println("sql:", queryTeacherSql)

	teacherList, err := mysql.GetTeacherList(queryTeacherSql)
	if err != nil {
		return nil, 0, err
	}

	teacherResList := make([]jrx_model.QueryTeacherResStruct, len(teacherList))

	for i := 0; i < len(teacherList); i++ {
		teacherResList[i].Name = teacherList[i].Name
		teacherResList[i].Username = teacherList[i].Username
		teacherResList[i].Password = teacherList[i].Password
		teacherResList[i].Gender = teacherList[i].Gender
		teacherResList[i].Ban = teacherList[i].Ban

		// 获取管理员等级信息
		if *teacherList[i].IsManager {
			managerType, err := GetManagerType(teacherList[i].Username)
			if err != nil {
				return nil, 0, err
			}
			teacherResList[i].ManagerType = managerType
		} else {
			teacherResList[i].ManagerType = "无"
		}
	}

	return teacherResList, allTeacherCount, err
}

// 获得查询老师的sql语句
func GetQueryTeacherSql(queryTeacherParama jrx_model.QueryTeacherParamStruct) string {
	querySql := `Select name, username, password, gender, is_manager, ban from users where identity = '老师' and deleted_at is NULL`

	if queryTeacherParama.Gender != "" {
		querySql = querySql + " and gender = '" + queryTeacherParama.Gender + "'"
	}

	if queryTeacherParama.Ban != nil {
		querySql = querySql + " and ban = " + strconv.FormatBool(*queryTeacherParama.Ban)
	}

	if queryTeacherParama.IsManager != nil {
		querySql = querySql + " and is_manager = " + strconv.FormatBool(*queryTeacherParama.IsManager)
	}

	if len(queryTeacherParama.SearchSelect) > 0 {
		querySql = querySql + " and " + queryTeacherParama.SearchSelect + " like '%" + queryTeacherParama.SearchMessage + "%'"
	}

	// limit 分页查询语句的拼接
	querySql = querySql + " ORDER BY name ASC" + " limit " + strconv.Itoa(queryTeacherParama.Limit) + " offset " + strconv.Itoa((queryTeacherParama.Page-1)*queryTeacherParama.Limit)

	return querySql

}

func GetManagerType(username string) (string, error) {
	casbinUser, err := mysql.GetManager(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("recored not founf recoery")
			return "无", nil
		} else {
			return "", err
		}
	}

	var managerType string

	switch casbinUser.CasbinCid {
	case "2", "3", "4", "5":
		managerType = "年级管理员"
	case "6":
		managerType = "班级管理员"
	case "1":
		managerType = "院级管理员"
	case "0":
		managerType = "超级管理员"
	case "7":
		managerType = "学生会成员"
	default:
		managerType = "无"
	}

	return managerType, err
}

// 处理管理员
func SetStuManagerService(setUsername string, username string, ManagerType string, year string) error {
	// 判断是否对自己进行设置
	if setUsername == username {
		return errors.New("您不能对自己账号进行设置")
	}

	// 判断这个用户是不是管理员
	isManager, err := mysql.GetIsManagerByUsername(setUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		} else {
			return err
		}
	}

	var casbinCid string
	switch ManagerType {
	case "班级管理员":
		casbinCid = "6"
	case "年级管理员":
		switch year {
		case "2024":
			casbinCid = "2"

		case "2023":
			casbinCid = "3"

		case "2022":
			casbinCid = "4"

		case "2021":
			casbinCid = "5"
		}

	case "院级管理员":
		casbinCid = "1"

	case "学生会成员":
		casbinCid = "7"

	case "取消管理员":
		if isManager {
			err := mysql.CancelStuManager(setUsername, casbinCid)
			return err
		} else {
			return errors.New("该用户不为管理员")
		}

	default:
		return errors.New("非法的管理员类型，请检查请求参数")

	}

	if !isManager { // 不是管理员
		err := mysql.SetStuManager(setUsername, casbinCid)
		if err != nil {
			return err
		}

	} else { // 是管理员
		casbinUser, err := mysql.GetManager(setUsername)
		existedCasbinCid := casbinUser.CasbinCid
		if err != nil {
			return err
		}

		if casbinCid == existedCasbinCid {
			return errors.New("该用户已是" + ManagerType)
		}

		err = mysql.ChangeStuManager(setUsername, casbinCid)
		if err != nil {
			return err
		}
	}

	return err
}

// 处理老师管理员
func SetTeacherManagerService(username string, ManagerType string) error {
	// 判断这个用户是不是管理员
	isManager, err := mysql.GetIsManagerByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		} else {
			return err
		}
	}

	var casbinCid string
	switch ManagerType {
	case "大一管理员":
		casbinCid = "2"

	case "大二管理员":
		casbinCid = "3"

	case "大三管理员":
		casbinCid = "4"

	case "大四管理员":
		casbinCid = "5"

	case "院级管理员":
		casbinCid = "1"

	case "学生会成员":
		casbinCid = "7"

	case "取消管理员":
		if isManager {
			err := mysql.CancelStuManager(username, casbinCid)
			return err
		} else {
			return errors.New("该用户不为管理员")
		}

	default:
		return errors.New("非法的管理员类型，请检查请求参数")
	}

	if !isManager { // 不是管理员
		err := mysql.SetStuManager(username, casbinCid)
		if err != nil {
			return err
		}

	} else { // 是管理员
		casbinUser, err := mysql.GetManager(username)
		existedCasbinCid := casbinUser.CasbinCid
		if err != nil {
			return err
		}

		if casbinCid == existedCasbinCid {
			return errors.New("该用户已是" + ManagerType)
		}

		err = mysql.ChangeStuManager(username, casbinCid)
		if err != nil {
			return err
		}
	}

	return err
}

func AddSingleTeacherService(addSingleTeacherReqStruct gorm_model.User) error {
	addSingleTeacherReqStruct.Identity = "老师"
	addSingleTeacherReqStruct.HeadShot = "https://student-grow.oss-cn-beijing.aliyuncs.com/image/user_headshot/user_headshot_5.png"
	err := mysql.AddSingleTeacher(&addSingleTeacherReqStruct)
	if err != nil {
		return err
	}
	return err
}

// 删除单个老师
func DeleteSingleTeacherService(input gorm_model.User) error {
	id, err := mysql.GetIdByUsername(input.Username)
	if err != nil {
		return err
	}

	err = mysql.DeleteSingleUser(id)
	if err != nil {
		return err
	}
	return err
}

// 编辑老师信息
func EditTeacherService(newTeacher jrx_model.ChangeTeacherMesStruct) error {
	id, err := mysql.GetIdByUsername(newTeacher.OldUsername)
	if err != nil {
		return err
	}

	fmt.Println("newTeacher : ", newTeacher)
	err = mysql.ChangeTeacherMessage(id, newTeacher)
	if err != nil {
		return err
	}

	return err
}

// 获取导出老师信息的 excel表格
func GetSelectedTeacherExcel(selectedTeacher []jrx_model.QueryTeacherStruct) (*bytes.Buffer, error) {
	// 提取处学号数组
	usernameSlice := make([]string, len(selectedTeacher))
	for i, v := range selectedTeacher {
		usernameSlice[i] = v.Username
	}
	fmt.Println(usernameSlice)

	// 从所有用户中查出选中的用户
	users, err := mysql.QuerySelectedUser(usernameSlice)
	if err != nil {
		return nil, err
	}

	// 创建 Excel 文件
	f := excelize.NewFile()

	// 设置表头
	f.SetCellValue("Sheet1", "A1", "姓名")
	f.SetCellValue("Sheet1", "B1", "账号")
	f.SetCellValue("Sheet1", "C1", "性别")

	// 填充数据
	for i, user := range users {
		row := i + 2 // 从第二行开始填充数据
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), user.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), user.Username)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), user.Gender)
	}

	// 将 Excel 文件写入内存
	excelData, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return excelData, err
}
