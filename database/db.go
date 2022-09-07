package database

import (
	"fmt"
	"invar/models"
	"invar/permission"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var RDS *redis.Client

func Connect(dsn string) {
	var err error
	var errTime = 0
	var reconnectTime = 10 * time.Second
	for {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			fmt.Println("Connect Success.")
			break
		}

		errTime += 1
		fmt.Printf("Try connect fail, time = %v.\n", errTime)
		fmt.Printf("After %v seconds will reconnect.\n", reconnectTime)

		if errTime >= 5 {
			panic("Could not connect with the database!")
		}

		time.Sleep(reconnectTime)
		continue
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.Admin{}, models.User{}, models.UserKYC{},
		models.Token{}, models.Bank{}, models.Withdrawal{},
		models.RefreshToken{}, models.WhiteList{},
		models.Product{}, models.Order{}, models.OrderItem{},
		models.Stack{}, models.StackRecord{}, models.StackProfitRecord{})
}

func InitDefaultAdmin(account, password string) {
	var admin models.Admin

	result := DB.Find(&admin)

	if result.RowsAffected != 0 {
		fmt.Println("Admin created.")
		return
	}

	permissions := permission.GetDefaultAdminPermission()
	admin = models.Admin{
		Account:     account,
		Premissions: permissions,
	}

	admin.SetPassword(password)
	DB.Create(&admin)
	DB.Save(&admin)
}

func SetupRedis(password string) {
	RDS = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		DB:       0,
		Password: password,
	})
}
