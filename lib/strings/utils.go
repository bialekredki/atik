package strings

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Encode(value string) string {
	hashedValue := md5.Sum([]byte(value))
	return hex.EncodeToString(hashedValue[:])
}
