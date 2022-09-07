package services

import (
	"invar/database"
	"invar/models"
)

func GetTokens(bankID, productID uint) ([]models.Token, error) {
	var tokens []models.Token

	result := database.DB.Where("bank_id = ?", bankID)

	if productID != 0 {
		result = result.Where("product_id = ?", productID)
	}

	result = result.Find(&tokens)

	if result.Error != nil {
		return tokens, result.Error
	}

	return tokens, nil
}

func UpdateTokens(tokens *[]models.Token) error {
	result := database.DB.Updates(tokens)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteTokens(tokens *[]models.Token) error {
	result := database.DB.Delete(tokens)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
