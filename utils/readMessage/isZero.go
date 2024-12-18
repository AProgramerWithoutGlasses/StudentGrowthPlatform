package readMessage

import "reflect"

// 用于判断一个变量的值是否为零值
func IsZero(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
