package readMessage

import jsonvalue "github.com/Andrew-M-C/go.jsonvalue"

// GetParams 传入多个切片
// 每个切片的大小为2:
// 0下标存储对应的字符串字段;1下标存储需要存储的变量指针
func GetParams(j *jsonvalue.V, args ...[]*any) error {
	var err error
	for _, arg := range args {
		v := *(arg[1])
		field := *(arg[0])
		switch v.(type) {
		case int:
			*arg[1], err = j.GetInt(field.(string))
		case string:
			*arg[1], err = j.GetString(field.(string))
		case bool:
			*arg[1], err = j.GetBool(field.(string))
		}
	}
	if err != nil {
		return err
	}
	return nil
}
