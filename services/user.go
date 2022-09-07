package services

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"invar/database"
	"invar/models"
	"invar/utils"
	"os"
	"time"

	"github.com/sec51/twofactor"
	"gorm.io/gorm"
)

func RegisterUser(user *models.User) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(user).Error
		if err != nil {
			return err
		}
		code := utils.GenerateReferrerCode(user.ID)
		user.ReferrerCode = code
		err = tx.Updates(&user).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	result := database.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func GetUserById(id uint) (models.User, error) {
	var user models.User

	result := database.DB.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func UpdateUser(user *models.User) error {
	result := database.DB.Updates(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Reset Password
type ResetPasswordCache struct {
	Email        string
	RequestToken string
	AccessToken  string
	ExpireTime   time.Time
}

func GetResetPasswordAccessToekn(email string, reqToken string) (string, error) {
	cache, err := getResetPasswordCache(email)
	if err != nil {
		return "", err
	}

	if cache.ExpireTime.Before(time.Now()) {
		return "", errors.New("invalid")
	}

	if cache.RequestToken != reqToken {
		return "", errors.New("invalid")
	}

	return cache.AccessToken, nil
}

func CheckResetPasswordAccessToekn(email string, accessToken string) bool {
	cache, err := getResetPasswordCache(email)
	if err != nil {
		return false
	}

	if cache.ExpireTime.Before(time.Now()) {
		return false
	}

	if cache.AccessToken != accessToken {
		return false
	}

	return true
}

func SendResetPassword(email string) error {
	var ctx = context.Background()

	baseUrl := os.Getenv("BASE_URL")

	reqToken, err := utils.GenerateRefreshToken(64)
	if err != nil {
		return err
	}

	accessToken, err := utils.GenerateRefreshToken(128)
	if err != nil {
		return err
	}

	encodedEmail := base64.URLEncoding.EncodeToString([]byte(email))

	var resetPasswordCache = ResetPasswordCache{
		Email:        email,
		RequestToken: reqToken,
		AccessToken:  accessToken,
		ExpireTime:   time.Now().Add(time.Minute * 10),
	}

	bytes, err := json.Marshal(resetPasswordCache)

	err = database.RDS.HSet(ctx, database.ResetPasswordCache, email, bytes).Err()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(baseUrl + "api/users/resetpassword?user=" + encodedEmail + "&" + "token=" + reqToken)
	text := "Please follow this link to reset password. It will expire in 15 minutes later. Password link:" + url
	go utils.SendEmail(email, "InVar Reset Password", text)
	return nil
}

func getResetPasswordCache(email string) (ResetPasswordCache, error) {
	var ctx = context.Background()
	cache := ResetPasswordCache{}
	result, err := database.RDS.HGet(ctx, database.ResetPasswordCache, email).Result()
	if err != nil {
		return cache, err
	}

	err = json.Unmarshal([]byte(result), &cache)
	if err != nil {
		return cache, errors.New("Unmarshal fail")
	}

	return cache, nil
}

func CheckRepeatUserEmail(email string) error {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.RowsAffected > 0 {
		return errors.New("email was used")
	}
	return nil
}

func GetUserTFA(user *models.User) ([]byte, error) {
	var otp *twofactor.Totp
	var err error
	if user.TFACode == nil {
		otp, err := twofactor.NewTOTP(user.Email, user.UserName, crypto.SHA1, 6)
		if err != nil {
			return nil, err
		}

		data, err := otp.ToBytes()
		if err != nil {
			return nil, err
		}

		user.TFACode = data
		err = UpdateUser(user)
		if err != nil {
			return nil, err
		}
	} else {
		otp, err = twofactor.TOTPFromBytes(user.TFACode, user.UserName)
		if err != nil {
			return nil, err
		}
	}

	qrCode, err := otp.QR()
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}

func EnableUserTFA(user *models.User, code string) error {
	if user.TFAEnable {
		return errors.New("Enabled")
	}

	otp, err := twofactor.TOTPFromBytes(user.TFACode, user.UserName)
	if err != nil {
		return err
	}

	err = otp.Validate(code)
	if err != nil {
		return err
	}

	user.TFAEnable = true
	err = UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func CheckUserTFA(user *models.User, code string) error {
	otp, err := twofactor.TOTPFromBytes(user.TFACode, user.UserName)
	if err != nil {
		return err
	}

	if !user.TFAEnable {
		return nil
	}

	err = otp.Validate(code)
	if err != nil {
		return err
	}

	return nil
}

func DisableUserTFA(user *models.User, code string, byAdmin bool) error {
	var err error
	if !user.TFAEnable {
		return errors.New("Disabled")
	}

	if !byAdmin {
		err = CheckUserTFA(user, code)
		if err != nil {
			return err
		}
	}

	user.TFAEnable = false
	err = UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
