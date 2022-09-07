package services

import (
	"context"
	"encoding/json"
	"errors"
	"invar/database"
	"invar/status"
	invarTool "invar/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type EmailCode struct {
	Code            string
	ExpireTime      time.Time
	Used            bool
	LastRequestTime time.Time
}

func RequestMailCode(email string) (string, int) {
	var ctx = context.Background()
	currentTime := time.Now()
	var codeData = EmailCode{}

	result, err := database.RDS.HGet(ctx, database.EmailCodeCache, email).Result()
	if err == nil {
		json.Unmarshal([]byte(result), &codeData)
		if currentTime.Sub(codeData.LastRequestTime) < (1 * time.Minute) {
			return "", status.Success
		}
	}

	code := generateMailCode()

	codeData = EmailCode{
		Code:            code,
		ExpireTime:      currentTime.Add(5 * time.Minute),
		Used:            false,
		LastRequestTime: currentTime,
	}

	bytes, err := json.Marshal(codeData)
	if err != nil {
		logrus.Error("Marshal fail, err=", err)
		return "", status.Unkonwn
	}

	err = database.RDS.HSet(ctx, database.EmailCodeCache, email, bytes).Err()
	if err != nil {
		logrus.Error("Redis set fail, err=", err)
		return "", status.Unkonwn
	}

	go invarTool.SendEmail(email, "Invar驗證碼", "您的驗證碼為「"+code+"」")

	return code[0:5], status.Success
}

func ValidEmailCode(email string, code string) error {
	var ctx = context.Background()
	var codeData = EmailCode{}

	result, err := database.RDS.HGet(ctx, database.EmailCodeCache, email).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(result), &codeData)
	if err != nil {
		return errors.New("Unmarshal fail")
	}

	if codeData.Used {
		return errors.New("Used")
	}

	if codeData.ExpireTime.Unix() < time.Now().Unix() {
		return errors.New("Expire")
	}

	if codeData.Code != code {
		return errors.New("Not equal")
	}

	err = database.RDS.HDel(ctx, database.EmailCodeCache, email).Err()
	if err != nil {
		return errors.New("Server err")
	}

	return nil
}

func generateMailCode() string {
	var sb strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 4; i++ {
		temp := 65 + rand.Intn(90-65)
		sb.WriteRune(rune(temp))
	}
	sb.WriteString("-")
	for j := 0; j < 6; j++ {
		temp := rand.Intn(10)
		sb.WriteString(strconv.Itoa(temp))
	}
	return sb.String()
}
