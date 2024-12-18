package hashMd5

import (
	"crypto/md5"
	"encoding/hex"
)

// HashMd5 用于文件生成唯一hash值
func HashMd5(src []byte) string {
	m := md5.New()
	m.Write(src)
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
