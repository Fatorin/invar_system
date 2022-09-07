package services

import (
	"invar/database"
	"invar/models"

	"github.com/sirupsen/logrus"
)

func GetProducts() ([]models.Product, error) {
	var products []models.Product

	result := database.DB.Find(&products)
	if result.Error != nil {
		logrus.Error("Get Products fail=", result.Error)
	}

	return products, nil
}

func AddProduct(product *models.Product) error {
	result := database.DB.Create(&product)
	if result.Error != nil {
		logrus.Error("Add Products fail=", result.Error)
	}

	return nil
}

func GetProduct(key uint) (models.Product, error) {
	product := models.Product{}
	product.ID = key

	result := database.DB.First(&product)
	if result.Error != nil {
		logrus.Error("Get Products fail=", result.Error)
	}

	return product, nil
}

func UpdateProduct(product *models.Product) error {
	result := database.DB.Updates(&product)

	if result.Error != nil {
		logrus.Error("Get Products fail=", result.Error)
	}

	return nil
}
