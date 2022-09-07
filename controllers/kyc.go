package controllers

import (
	"invar/middlewares"
	"invar/models"
	"invar/services"
	"invar/status"
	"invar/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GetKYC godoc
// @Summary      獲得實名資料
// @Description  獲得實名資料
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200                  {object}  status.ResponseWtihData{data=models.UserKYC}
// @Failure      400                  {object}  status.Response
// @Failure      500                  {object}  status.Response
// @Router       /kyc [get]
// @Security     BearerAuth
func GetKYC(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)

	kyc, err := services.GetKYCByUserID(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   kyc,
	})
}

// GetKYCByAdmin godoc
// @Summary      管理者獲得使用者實名資料
// @Description  管理者獲得使用者實名資料
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id                   path      int     true   "User ID"
// @Success      200                  {object}  status.ResponseWtihData{data=models.UserKYC}
// @Failure      400                  {object}  status.Response
// @Failure      500                  {object}  status.Response
// @Router       /admin/kyc/:id [post]
// @Security     BearerAuth
func GetKYCByAdmin(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	kyc, err := services.GetKYCByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   kyc,
	})
}

type UpsertKYCReq struct {
	Name         string `json:"name"`
	IdentityCard string `json:"identity_card"`
}

// AddKYC godoc
// @Summary      新增實名資料
// @Description  新增實名資料
// @Tags         User
// @Accept       mpfd
// @Produce      json
// @Param        name                 formData  string  true  "姓名"
// @Param        identity_card        formData  string  true  "身分證號碼"
// @Param        photo_file           formData  file    true  "大頭貼"
// @Param        id_photo_front_file  formData  file    true  "身分證正面"
// @Param        id_photo_back_file   formData  file    true  "身分證背面"
// @Param        photo_with_id_file   formData  file    true  "與身分證合照"
// @Success      200                  {object}  status.Response
// @Failure      400                  {object}  status.Response
// @Failure      500                  {object}  status.Response
// @Router       /kyc [post]
// @Security     BearerAuth
func AddKYC(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)
	var data UpsertKYCReq

	user, err := services.GetUserById(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	if user.Status == models.AuditSuccess {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.UserHasKYC),
		})
		return
	}

	bindErr := c.BindWith(&data, binding.Form)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	photoFile, photoFileErr := c.FormFile("photo_file")
	idPhotoFrontFile, idPhotoFrontFileErr := c.FormFile("id_photo_front_file")
	idPhotoBackFile, idPhotoBackFileErr := c.FormFile("id_photo_back_file")
	photoWithIDFile, photoWithIDFileErr := c.FormFile("photo_with_id_file")
	if photoFileErr != nil || idPhotoFrontFileErr != nil || idPhotoBackFileErr != nil || photoWithIDFileErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	fileSupportTypes := []int{utils.PNG, utils.JPEG, utils.JPG}
	filesName := []string{photoFile.Filename, idPhotoFrontFile.Filename, idPhotoBackFile.Filename, photoWithIDFile.Filename}
	if utils.CheckMultiFileMediaType(filesName, fileSupportTypes) {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.UnsupportedMediaType),
		})
		return
	}

	savePath := status.SaveFileFolderPath + status.SaveFilePublicFolderPath + "/"
	photoFileSavePath := utils.FilePathGenerator(photoFile.Filename, savePath)
	err = c.SaveUploadedFile(photoFile, photoFileSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	idPhotoFrontFileSavePath := utils.FilePathGenerator(idPhotoFrontFile.Filename, savePath)
	err = c.SaveUploadedFile(idPhotoFrontFile, idPhotoFrontFileSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	idPhotoBackFileSavePath := utils.FilePathGenerator(idPhotoBackFile.Filename, savePath)
	err = c.SaveUploadedFile(idPhotoBackFile, idPhotoBackFileSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	photoWithIDFileSavePath := utils.FilePathGenerator(photoWithIDFile.Filename, savePath)
	err = c.SaveUploadedFile(photoWithIDFile, photoWithIDFileSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	var kyc = models.UserKYC{
		UserID:       uint(roleID),
		Name:         data.Name,
		IdentityCard: data.IdentityCard,
		Photo:        photoFileSavePath[1:],
		IDPhotoFront: idPhotoFrontFileSavePath,
		IDPhotoBack:  idPhotoBackFileSavePath,
		PhotoWithID:  photoWithIDFileSavePath,
	}

	err = services.AddKYC(&kyc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	user.Status = models.Auditing
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
	})
}

// UpdateKYC godoc
// @Summary      更新實名資料
// @Description  更新實名資料
// @Tags         User
// @Accept       mpfd
// @Produce      json
// @Param        name                 formData  string  false  "姓名"
// @Param        identity_card        formData  string  false  "身分證號碼"
// @Param        photo_file           formData  file    false  "大頭貼"
// @Param        id_photo_front_file  formData  file    false  "身分證正面"
// @Param        id_photo_back_file   formData  file    false  "身分證背面"
// @Param        photo_with_id_file   formData  file    false  "與身分證合照"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /kyc [patch]
// @Security     BearerAuth
func UpdateKYC(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)
	var data UpsertKYCReq

	user, err := services.GetUserById(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	kyc, err := services.GetKYCByUserID(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.UserNoKYC),
		})
	}

	if kyc.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	if user.Status == models.AuditSuccess {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.UserHasKYC),
		})
		return
	}

	bindErr := c.BindWith(&data, binding.Form)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	savePath := status.SaveFileFolderPath + status.SaveFilePublicFolderPath + "/"
	fileSupportTypes := []int{utils.PNG, utils.JPEG, utils.JPG}
	photoFile, err := c.FormFile("photo_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(photoFile.Filename, savePath, kyc.Photo, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(photoFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.Photo = fileSavePath[1:]
	}

	idPhotoFrontFile, err := c.FormFile("id_photo_front_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(idPhotoFrontFile.Filename, savePath, kyc.IDPhotoFront, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(idPhotoFrontFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.IDPhotoFront = fileSavePath[1:]
	}

	idPhotoBackFile, err := c.FormFile("id_photo_back_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(idPhotoBackFile.Filename, savePath, kyc.IDPhotoBack, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(idPhotoBackFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.IDPhotoBack = fileSavePath[1:]
	}

	photoWithIDFile, err := c.FormFile("photo_with_id_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(photoWithIDFile.Filename, savePath, kyc.PhotoWithID, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(photoWithIDFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.PhotoWithID = fileSavePath[1:]
	}

	if data.Name != "" {
		kyc.Name = data.Name
	}

	if data.IdentityCard != "" {
		kyc.Name = data.IdentityCard
	}

	user.Status = models.Auditing
	err = services.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}
	//如果KYC審核失敗則寄信給使用者

	err = services.UpdateKYC(&kyc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
	})
}

type UpdateKYCByAdminReq struct {
	Name         string `json:"name"`
	IdentityCard string `json:"identity_card"`
	Status       byte   `json:"status"`
}

// UpdateKYCByAdmin godoc
// @Summary      管理者更新實名資料
// @Description  管理者更新實名資料
// @Tags         User
// @Accept       mpfd
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        name                 formData  string  false  "姓名"
// @Param        identity_card        formData  string  false  "身分證號碼"
// @Param        status               formData  string  false  "審核結果"
// @Param        photo_file           formData  file    false  "大頭貼"
// @Param        id_photo_front_file  formData  file    false  "身分證正面"
// @Param        id_photo_back_file   formData  file    false  "身分證背面"
// @Param        photo_with_id_file   formData  file    false  "與身分證合照"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/kyc/:id [patch]
// @Security     BearerAuth
func UpdateKYCByAdmin(c *gin.Context) {
	kycID, _ := strconv.Atoi(c.Param("id"))
	var data UpdateKYCByAdminReq

	kyc, err := services.GetKYC(uint(kycID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.UserNoKYC),
		})
	}

	user, err := services.GetUserById(kyc.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	bindErr := c.BindWith(&data, binding.Form)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	savePath := status.SaveFileFolderPath + status.SaveFilePublicFolderPath + "/"
	fileSupportTypes := []int{utils.PNG, utils.JPEG, utils.JPG}
	photoFile, err := c.FormFile("photo_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(photoFile.Filename, savePath, kyc.Photo, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(photoFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.Photo = fileSavePath[1:]
	}

	idPhotoFrontFile, err := c.FormFile("id_photo_front_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(idPhotoFrontFile.Filename, savePath, kyc.IDPhotoFront, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(idPhotoFrontFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.IDPhotoFront = fileSavePath[1:]
	}

	idPhotoBackFile, err := c.FormFile("id_photo_back_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(idPhotoBackFile.Filename, savePath, kyc.IDPhotoBack, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(idPhotoBackFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.IDPhotoBack = fileSavePath[1:]
	}

	photoWithIDFile, err := c.FormFile("photo_with_id_file")
	if err == nil {
		fileSavePath, errCode := utils.GenerateNewFilePathAndRemoveOld(photoWithIDFile.Filename, savePath, kyc.PhotoWithID, fileSupportTypes)
		if errCode != status.Success {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(errCode),
			})
			return
		}

		err = c.SaveUploadedFile(photoWithIDFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		kyc.PhotoWithID = fileSavePath[1:]
	}

	if data.Name != "" {
		kyc.Name = data.Name
	}

	if data.IdentityCard != "" {
		kyc.Name = data.IdentityCard
	}

	if user.Status != data.Status {
		user.Status = data.Status
		err = services.UpdateUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}
		//如果KYC審核失敗則寄信給使用者
	}

	err = services.UpdateKYC(&kyc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
	})
}
