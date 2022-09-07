package services

import (
	"invar/database"
	"invar/models"

	"github.com/sirupsen/logrus"
)

func GetWhiteLists(userID uint) ([]models.WhiteList, error) {
	var whitelists []models.WhiteList

	result := database.DB.Where("user_id = ?", userID).Find(&whitelists)
	if result.Error != nil {
		logrus.Error("Get WhiteLists fail=", result.Error)
	}

	return whitelists, nil
}

func GetWhiteList(whitelistID uint) (models.WhiteList, error) {
	var whitelist models.WhiteList

	result := database.DB.Where("id = ?", whitelistID).First(&whitelist)
	if result.Error != nil {
		logrus.Error("Get WhiteLists fail=", result.Error)
	}

	return whitelist, nil
}

func AddWhiteList(whitelist *models.WhiteList) error {
	result := database.DB.Create(&whitelist)
	if result.Error != nil {
		logrus.Error("Add WhiteLists fail=", result.Error)
	}

	return nil
}

func UpdateWhiteList(whitelist *models.WhiteList) error {
	result := database.DB.Updates(&whitelist)

	if result.Error != nil {
		logrus.Error("Get WhiteLists fail=", result.Error)
	}

	return nil
}

func DeleteWhiteList(id uint) error {
	var whitelist models.WhiteList
	whitelist.ID = id

	result := database.DB.Delete(&whitelist)

	if result.Error != nil {
		logrus.Error("Delete WhiteLists fail=", result.Error)
	}

	return nil
}
