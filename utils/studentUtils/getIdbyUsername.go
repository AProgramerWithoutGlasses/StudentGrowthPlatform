package studentUtils

import "strings"

// RemoveBlank 用于去除字符串中含有的所有空格
func RemoveBlank(str string) string {
	if strings.Contains(str, " ") {
		str = strings.Replace(str, " ", "", -1)
	}
	return str
}
