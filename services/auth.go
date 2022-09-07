package services

import (
	"errors"
	"invar/database"
	"invar/models"
	"invar/utils"
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/sirupsen/logrus"
)

// Token
var symmetricKey paseto.V4SymmetricKey

func InitSymmetricKey() {
	symmetricKey = paseto.NewV4SymmetricKey()
	logrus.Info("Paseto Key was created.")
}

func ParseToken(token string) (roleType int, roleID int, err error) {
	parser := paseto.NewParser()

	data, err := parser.ParseV4Local(symmetricKey, token, nil)
	if err != nil {
		return roleType, roleID, err
	}

	roleTypeStr, err := data.GetJti()
	if err != nil {
		return roleType, roleID, err
	}

	roleIdStr, err := data.GetSubject()
	if err != nil {
		return roleType, roleID, err
	}

	roleType, err = strconv.Atoi(roleTypeStr)
	if err != nil {
		return roleType, roleID, err
	}

	roleID, err = strconv.Atoi(roleIdStr)
	if err != nil {
		return roleType, roleID, err
	}

	return roleType, roleID, nil
}

func GenerateAccessToken(roleType, roleID uint) string {
	token := paseto.NewToken()
	token.SetJti(strconv.Itoa(int(roleType)))
	token.SetSubject(strconv.Itoa(int(roleID)))
	token.SetExpiration(time.Now().Add(15 * time.Minute))
	encrypted := token.V4Encrypt(symmetricKey, nil)
	return encrypted
}

func GenerateRefreshToken(roleType, roleID uint, ip string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	for {
		token, err := utils.GenerateRefreshToken(64)
		if err != nil {
			return refreshToken, err
		}

		result := database.DB.Where("token = ?", token).First(&refreshToken)
		if result.RowsAffected == 0 {
			refreshToken = models.RefreshToken{
				RoleType:    roleType,
				RoleID:      roleID,
				Token:       token,
				Expires:     time.Now().Add(time.Hour * 24),
				CreatedByIP: ip,
			}

			err = database.DB.Create(&refreshToken).Error
			if err != nil {
				logrus.Error(err)
			}

			return refreshToken, nil
		}
	}
}

func RevokeRefreshToken(token string, ip string) (models.RefreshToken, error) {
	refreshtoken, err := findRefeshToken(token)

	if err != nil {
		return refreshtoken, err
	}

	err = RevokeToken(token, ip)
	if err != nil {
		return refreshtoken, err
	}

	newRefreshToken, err := GenerateRefreshToken(refreshtoken.RoleType, refreshtoken.RoleID, ip)
	if err != nil {
		return refreshtoken, err
	}

	return newRefreshToken, nil
}

func RevokeToken(token string, ip string) error {
	refreshtoken, err := findRefeshToken(token)

	if err != nil {
		return err
	}

	refreshtoken.SetToRevoked(ip)

	err = database.DB.Updates(&refreshtoken).Error
	if err != nil {
		return err
	}

	return nil
}

func findRefeshToken(token string) (models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	result := database.DB.Where("token = ?", token).First(&refreshToken)

	if result.Error != nil {
		return refreshToken, result.Error
	}

	if refreshToken.IsRevoked() {
		return refreshToken, errors.New("was revoked")
	}

	if refreshToken.Expires.Sub(time.Now()) < 0 {
		return refreshToken, errors.New("was expired")
	}

	return refreshToken, nil
}
