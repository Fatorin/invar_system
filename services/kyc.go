package services

import (
	"invar/database"
	"invar/models"

	"github.com/sirupsen/logrus"
)

func GetKYC(key uint) (models.UserKYC, error) {
	kyc := models.UserKYC{}
	kyc.ID = key

	result := database.DB.First(&kyc)
	if result.Error != nil {
		logrus.Error("Get KYC fail=", result.Error)
	}

	return kyc, nil
}

func GetKYCByUserID(userID uint) (models.UserKYC, error) {
	kyc := models.UserKYC{}
	result := database.DB.Where("user_id = ?", userID).Find(&kyc)
	if result.Error != nil {
		logrus.Error("Get KYC fail=", result.Error)
	}

	return kyc, nil
}

func AddKYC(kyc *models.UserKYC) error {
	err := database.DB.Create(&kyc)

	if err != nil {
		logrus.Error("Create kyc fail, err", err)
		return err.Error
	}

	return nil
}

func UpdateKYC(kyc *models.UserKYC) error {
	err := database.DB.Updates(&kyc)

	if err != nil {
		logrus.Error("Update kyc fail, err", err)
		return err.Error
	}

	return nil
}
