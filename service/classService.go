package service

import (
	"studentGrow/dao/mysql"
	"studentGrow/models/jrx_model"
)

func GetClassListService() ([]jrx_model.Class, error) {
	classNameList, err := mysql.GetClassList()
	if err != nil {
		return nil, err
	}

	// 叫上班级id
	var classList []jrx_model.Class
	i := 1
	for _, c := range classNameList {

		classList = append(classList, jrx_model.Class{
			ClassId:   i,
			ClassName: c,
		})
		i++

	}
	return classList, nil
}

func GetClassByGradeService(grade int) ([]string, error) {
	plusTimeYear := CalculatePlusTimeYearByGrade(grade)

	classNameList, err := mysql.GetClassListByPlusTimeYear(plusTimeYear)
	if err != nil {
		return nil, err
	}

	return classNameList, nil
}

func CalculatePlusTimeYearByGrade(grade int) string {
	var plusTimeYear string
	switch grade {
	case 1:
		plusTimeYear = "2024"
	case 2:
		plusTimeYear = "2023"
	case 3:
		plusTimeYear = "2022"
	case 4:
		plusTimeYear = "2021"
	}

	// 计算入学年份
	return plusTimeYear
}

// 检查班级名称是否符合格式
//func matchFormat(name string) bool {
//	if len(name) != 9 {
//		return false
//	}
//	// 检查前两个字符是否为汉字
//	for i := 0; i < 2; i++ {
//		if !isChinese(rune(name[i])) {
//			return false
//		}
//	}
//	// 检查后三个字符是否为数字
//	for i := 2; i < 5; i++ {
//		if !isDigit(name[i]) {
//			return false
//		}
//	}
//	return true
//}
//
//// 检查字符是否为汉字
//func isChinese(c rune) bool {
//	return unicode.Is(unicode.Scripts["Han"], c)
//}
//
//// 检查字符是否为数字
//func isDigit(c byte) bool {
//	return c >= '0' && c <= '9'
//}
