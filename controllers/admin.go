package controllers

import (
	"invar/middlewares"
	"invar/services"
	"invar/status"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthenticateAdmin godoc
// @Summary      後台管理者登入
// @Description  後台管理者登入
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        account   body      string  true   "帳號"
// @Param        password  body      string  true   "密碼"
// @Param        tfa       body      string  false  "二階段驗證碼"
// @Success      200    {object}  status.ResponseWtihData{data=string}
// @Failure      400             {object}  status.Response
// @Failure      500             {object}  status.Response
// @Router       /admin_auth [post]
func AuthenticateAdmin(c *gin.Context) {
	var data AuthenticateReq
	roleType := uint(middlewares.Admin)

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	roleID, errCode := checkLoginInfo(roleType, data.Account, data.Password, data.TFA)
	if errCode != status.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(errCode),
		})
		return
	}

	refreshToken, err := services.GenerateRefreshToken(roleType, roleID, c.ClientIP())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	accessToken := services.GenerateAccessToken(roleType, roleID)

	c.SetCookie(middlewares.REFRESH_TOKEN, refreshToken.Token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   accessToken,
	})
}

// ChangeAdminPassword godoc
// @Summary      管理者修改密碼
// @Description  管理者修改密碼
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        old_password    body      string  true  "舊密碼"
// @Param        password        body      string  true  "新密碼"
// @Param        check_password  body      string  true  "確認新密碼"
// @Success      200             {object}  status.ResponseWtihData{data=string}
// @Failure      400       {object}  status.Response
// @Failure      500       {object}  status.Response
// @Router       /admin/change_password [post]
// @Security     BearerAuth
func ChangeAdminPassword(c *gin.Context) {
	var data ChangePasswordReq

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	if data.Password != data.CheckPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.PasswordNotEqual),
		})
		return
	}

	adminID := c.GetUint(middlewares.ROLE_ID)
	if adminID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	admin, err := services.GetAdminById(adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	err = admin.ComparePassword(data.OldPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.OldPasswordIsIncorrect),
		})
		return
	}

	admin.SetPassword(data.Password)
	result := services.UpdateAdmin(&admin)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetAdminTFA(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)

	admin, err := services.GetAdminById(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.ExistUser),
		})
		return
	}

	qrCode, err := services.GetAdminTFA(&admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		status.RespStatus: status.NewResponse(status.ExistUser),
		status.RespData:   qrCode,
	})
}

type TFAReq struct {
	TFA string `json:"tfa"`
}

func EnableAdminTFA(c *gin.Context) {
	var data TFAReq

	roleID := c.GetInt(middlewares.ROLE_ID)

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	user, err := services.GetUserById(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	err = services.EnableUserTFA(&user, data.TFA)
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

func DisableAdminTFA(c *gin.Context) {
	var data TFAReq

	roleID := c.GetInt(middlewares.ROLE_ID)

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	admin, err := services.GetAdminById(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	err = services.DisableAdminTFA(&admin, data.TFA)
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
