package services

import (
	"invar/database"
	"invar/models"

	"github.com/sirupsen/logrus"
)

func AddBank(bank models.Bank) error {
	result := database.DB.Create(&bank)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetBanks(email, username string) ([]models.Bank, error) {
	var banks []models.Bank

	result := database.DB.Joins("User")
	if username != "" {
		result = result.Where(&models.User{UserName: username})
	}

	if email != "" {
		result = result.Where(&models.User{Email: email})
	}

	result = result.Preload("Tokens").Find(&banks)

	if result.Error != nil {
		logrus.Error("Get Banks fail=", result.Error)
	}

	return banks, nil
}

func GetBank(userID uint) (models.Bank, error) {
	var bank models.Bank

	result := database.DB.Where("user_id = ?", userID).Find(&bank)

	if result.Error != nil {
		return bank, result.Error
	}

	return bank, nil
}

func UpdateBank(bank models.Bank) error {
	result := database.DB.Updates(&bank)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AddToken(token *models.Token) error {
	result := database.DB.Create(&token)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteToken(token *models.Token) error {
	result := database.DB.Delete(&token)

	if result != nil {
		return result.Error
	}

	return nil
}
