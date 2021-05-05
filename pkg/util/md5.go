package util

import (
	"crypto/md5"
	"encoding/hex"
)

//返回字符串value经MD5编码后的16进制字符串
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

