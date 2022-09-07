package services

import (
	"invar/database"
	"invar/models"
	"invar/utils"
	"time"

	"gorm.io/gorm"
)

func GetStacks() ([]models.Stack, error) {
	var stacks []models.Stack
	result := database.DB.Find(&stacks)

	if result.Error != nil {
		return stacks, result.Error
	}

	return stacks, nil
}

func GetStack(id uint) (models.Stack, error) {
	var stack models.Stack
	result := database.DB.Where("id = ?", id).Find(&stack)

	if result.Error != nil {
		return stack, result.Error
	}

	return stack, nil
}

func AddStack(stack *models.Stack) error {
	result := database.DB.Create(&stack)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateStack(stack *models.Stack) error {
	result := database.DB.Updates(&stack)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetStacksRecord(userID uint) ([]models.StackRecord, error) {
	var stacksRecord []models.StackRecord

	result := database.DB.Preload("StackProfitRecords").Where("user_id = ?", userID).Find(&stacksRecord)

	if result.Error != nil {
		return stacksRecord, nil
	}

	return stacksRecord, nil
}

func GetStackRecord(id uint) (models.StackRecord, error) {
	var stackRecord models.StackRecord

	result := database.DB.Preload("StackProfitRecords").Where("uid = ?", id).Find(&stackRecord)

	if result.Error != nil {
		return stackRecord, nil
	}

	return stackRecord, nil
}

func AddStackRecord(stack *models.Stack, quantity uint) error {
	currentTime := time.Now()
	var record = models.StackRecord{
		Status:           models.Stacking,
		StackID:          stack.ID,
		Quantity:         quantity,
		EndGetProfitTime: currentTime.AddDate(0, int(stack.ContractTimeBound), 0),
	}

	if stack.ProfitIntervalMonth != 0 {
		record.NextGetProfitTime = currentTime.AddDate(0, int(stack.ProfitIntervalMonth), 0)
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&record).Error
		if err != nil {
			return err
		}

		record.Serial = utils.GenerateStackRecordSerialCode(record.ID)
		err = tx.Updates(&record).Error
		if err != nil {
			return err
		}
		return nil
	})

	//檢查有沒有該TOKEN 沒有就取消
	//call redis

	return err
}

func UpdateStackRecord(stackRecord *models.StackRecord) error {
	result := database.DB.Updates(&stackRecord)
	if result.Error != nil {
		return result.Error
	}

	//call redis

	return nil
}

func GetStackProfitRecords(stackRecordID uint) ([]models.StackProfitRecord, error) {
	var records []models.StackProfitRecord
	result := database.DB.Where("stack_record_id = ?", stackRecordID).Find(&records)
	if result.Error != nil {
		return records, result.Error
	}
	return records, nil
}

func GetStackProfitRecord(id uint) (models.StackProfitRecord, error) {
	var record models.StackProfitRecord
	result := database.DB.Where("id = ?", id).Find(&record)
	if result.Error != nil {
		return record, result.Error
	}
	return record, nil
}

func AddStackProfitRecord(bank *models.Bank, record *models.StackProfitRecord) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&record).Error
		if err != nil {
			return err
		}

		bank.InVarCoin = bank.InVarCoin.Add(record.Profit)
		err = tx.Updates(&bank).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func DeleteStackProfitRecord(bank *models.Bank, record *models.StackProfitRecord) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&record).Error
		if err != nil {
			return err
		}

		bank.InVarCoin = bank.InVarCoin.Sub(record.Profit)
		err = tx.Updates(&bank).Error
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
