package services

import (
	"crypto"
	"errors"
	"invar/database"
	"invar/models"

	"github.com/sec51/twofactor"
)

func RegisterAdmin(admin *models.Admin) error {
	err := database.DB.Create(&admin).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateAdmin(admin *models.Admin) error {
	err := database.DB.Updates(&admin).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckRepeatAdminEmail(email string) error {
	var admin models.Admin
	result := database.DB.Where("email = ?", email).First(&admin)
	if result.RowsAffected > 0 {
		return errors.New("email was used")
	}
	return nil
}

func GetAdminByAccount(account string) (models.Admin, error) {
	var admin models.Admin

	result := database.DB.Where("account = ?", account).First(&admin)

	if result.Error != nil {
		return admin, result.Error
	}

	return admin, nil
}

func GetAdminById(id uint) (models.Admin, error) {
	var admin models.Admin

	result := database.DB.Where("id = ?", id).First(&admin)

	if result.Error != nil {
		return admin, result.Error
	}

	return admin, nil
}

func GetAdminTFA(admin *models.Admin) ([]byte, error) {
	var otp *twofactor.Totp
	var err error
	if admin.TFACode == nil {
		otp, err = twofactor.NewTOTP(admin.Account, admin.Account, crypto.SHA1, 6)
		if err != nil {
			return nil, err
		}

		data, err := otp.ToBytes()
		if err != nil {
			return nil, err
		}

		admin.TFACode = data
		err = UpdateAdmin(admin)
		if err != nil {
			return nil, err
		}
	} else {
		otp, err = twofactor.TOTPFromBytes(admin.TFACode, admin.Account)
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

func EnableAdminTFA(admin *models.Admin, code string) error {
	if admin.TFAEnable {
		return errors.New("Enabled")
	}

	otp, err := twofactor.TOTPFromBytes(admin.TFACode, admin.Account)
	if err != nil {
		return err
	}

	err = otp.Validate(code)
	if err != nil {
		return err
	}

	admin.TFAEnable = true
	err = UpdateAdmin(admin)
	if err != nil {
		return err
	}

	return nil
}

func CheckAdminTFA(admin *models.Admin, code string) error {
	otp, err := twofactor.TOTPFromBytes(admin.TFACode, admin.Account)
	if err != nil {
		return err
	}

	if !admin.TFAEnable {
		return errors.New("Not enable tfa.")
	}

	err = otp.Validate(code)
	if err != nil {
		return err
	}

	return nil
}

func DisableAdminTFA(admin *models.Admin, code string) error {
	if !admin.TFAEnable {
		return errors.New("Disabled")
	}

	err := CheckAdminTFA(admin, code)
	if err != nil {
		return err
	}

	admin.TFAEnable = false
	err = UpdateAdmin(admin)
	if err != nil {
		return err
	}

	return nil
}
