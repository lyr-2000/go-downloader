package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(fName string) string {
	md5 := md5.New()
	md5.Write([]byte(fName))
	md5Data := md5.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}
