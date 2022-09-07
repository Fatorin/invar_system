package services

import (
	"errors"
	"invar/database"
	"invar/models"
	"invar/status"
	"invar/utils"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order

	result := database.DB.Preload("OrderItem").Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		logrus.Error("Get Orders fail=", result.Error)
	}

	return orders, nil
}

func GetOrdersByAdmin(serial, email, username string, createAt time.Time) ([]models.Order, error) {
	var orders []models.Order

	result := database.DB.Joins("User")
	if username != "" {
		result = result.Where(&models.User{UserName: username})
	}

	if email != "" {
		result = result.Where(&models.User{Email: email})
	}

	if serial != "" {
		result = result.Where("serial = ?", serial)
	}

	if !createAt.IsZero() {
		result = result.Where("create_at = ?", createAt)
	}

	result = result.Preload("OrderItems").Find(&orders)

	if result.Error != nil {
		logrus.Error("Get Orders fail=", result.Error)
	}

	return orders, nil
}

func GetOrder(key uint) (models.Order, error) {
	order := models.Order{}
	order.ID = key

	result := database.DB.Preload("OrderItem").Find(&order)
	if result.Error != nil {
		logrus.Error("Get Orders fail=", result.Error)
	}

	return order, nil
}

func AddOrder(userID uint, orderItems []models.OrderItem) error {
	var order = models.Order{
		UserID:           userID,
		Status:           models.Ordered,
		TotalAmount:      clacOrderTotal(orderItems),
		OrderItems:       orderItems,
		PaymentLimitTime: time.Now().Add(time.Hour * 24 * 7),
	}

	var products = make([]models.Product, 0)
	for _, v := range orderItems {
		product, err := GetProduct(v.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < v.Quantity {
			return errors.New(status.ErrorText(status.OutOfStock))
		}
		product.Stock -= v.Quantity
		products = append(products, product)
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&order).Error
		if err != nil {
			return err
		}
		code := utils.GenerateOrderSerialCode(order.ID)
		order.Serial = code
		err = tx.Updates(&order).Error
		if err != nil {
			return err
		}
		return nil
	})

	logrus.Error("Add Orders fail=", err)
	return err
}

func UpdateOrder(order *models.Order) error {
	result := database.DB.Updates(&order)

	if result.Error != nil {
		logrus.Error("Get Orders fail=", result.Error)
	}

	return nil
}

func CancelOrder(order *models.Order) error {
	order.Status = models.Cancel

	var products = make([]models.Product, 0)
	for _, v := range order.OrderItems {
		product, err := GetProduct(v.ProductID)
		if err != nil {
			return err
		}
		product.Stock += v.Quantity
		products = append(products, product)
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Updates(&products).Error
		if err != nil {
			return err
		}

		err = tx.Updates(&order).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func CompletedOrder(order *models.Order) error {
	order.Status = models.Completed

	err := database.DB.Updates(&order).Error
	if err != nil {
		return err
	}

	return nil
}

func clacOrderTotal(orderItems []models.OrderItem) decimal.Decimal {
	total := decimal.NewFromInt(0)
	for _, v := range orderItems {
		total = total.Add(v.Price)
	}
	return total
}
