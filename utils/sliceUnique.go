package utils

import "cmp"

// 切片去重
//泛型Ordered
//~int | ~int8 | ~int16 | ~int32 | ~int64 |
//~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
//~float32 | ~float64 |
//~string

func SliceUnique[T cmp.Ordered](tList []T) (tResList []T) {
	if len(tList) < 2 {
		tResList = tList
		return
	}
	m1 := make(map[T]byte)
	for _, v := range tList {
		l := len(m1)
		m1[v] = 0
		if len(m1) != l {
			tResList = append(tResList, v)
		}
	}
	return
}
