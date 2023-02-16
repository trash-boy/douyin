package encryption

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

func GenerateUserFlag(userID uint) string {
	md5 := md5.New()
	md5.Write([]byte(strconv.FormatInt(int64(userID), 10)))
	return hex.EncodeToString(md5.Sum([]byte("")))
}
