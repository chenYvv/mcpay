package idcode

import (
	"errors"
	"strconv"
	"sync"

	"github.com/speps/go-hashids"
)

var (
	encoder  *hashids.HashID
	initOnce sync.Once
)

// Init 在应用启动时调用一次。
// salt      —— 私有盐值
// minLength —— 最小长度(<4 时按 4)
func Init(salt string, minLength int) error {
	var err error
	initOnce.Do(func() {
		if salt == "" {
			err = errors.New("idcode: salt cannot be empty")
			return
		}
		if minLength < 4 {
			minLength = 4
		}

		data := hashids.NewData()
		data.Salt = salt
		data.MinLength = minLength

		encoder, err = hashids.NewWithData(data)
	})
	return err
}

// Encrypt 把字符串形式的数字 ID → 编码串
func Encrypt(idStr string) (string, error) {
	if encoder == nil {
		return "", errors.New("idcode: not initialized, call Init first")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return "", errors.New("idcode: invalid id (must be decimal integer)")
	}

	return encoder.EncodeInt64([]int64{id})
}

// Decrypt 把编码串 → 原始字符串 ID
func Decrypt(code string) (string, error) {
	if encoder == nil {
		return "", errors.New("idcode: not initialized, call Init first")
	}

	ids, err := encoder.DecodeInt64WithError(code)
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", errors.New("idcode: decode result is empty")
	}

	return strconv.FormatInt(ids[0], 10), nil
}
