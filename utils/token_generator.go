package utils

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func GenerateRefreshToken(createLen int) (string, error) {
	rand.Seed(time.Now().UnixNano())

	bytes := []byte{}
	bytes = make([]byte, createLen)

	_, err := rand.Read(bytes[:])
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes[:]), nil
}
