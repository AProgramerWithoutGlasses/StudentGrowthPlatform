package fileIsInList

import "strings"

// FileIsInList 根据传入的文件名后缀判断文件是否在文件后缀白名单内
func FileIsInList(fileName string, whiteList []string) bool {
	nameList := strings.Split(fileName, ".")
	suffix := strings.ToLower(nameList[len(nameList)-1])
	for _, s := range whiteList {
		if suffix == s {
			return true
		}
	}
	return false
}
