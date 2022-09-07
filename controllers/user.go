package controllers

import (
	"invar/database"
	"invar/middlewares"
	"invar/models"
	"invar/services"
	"invar/status"
	"invar/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Authenticate godoc
// @Summary      使用者帳號登入
// @Description  輸入帳號密碼與二階段驗證進行登入
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        account   body      string  true   "帳號"
// @Param        password  body      string  true   "密碼"
// @Param        tfa       body      string  false  "二階段驗證碼"
// @Success      200             {object}  status.ResponseWtihData{data=string}
// @Failure      400             {object}  status.Response
// @Failure      500             {object}  status.Response
// @Router       /auth [post]
func Authenticate(c *gin.Context) {
	var data AuthenticateReq
	roleType := uint(middlewares.User)

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
		status.RespStatus: status.NewResponse(status.Unkonwn),
		status.RespData:   accessToken,
	})
}

// Register
type RegisterReq struct {
	UserName      string `json:"username"`
	Email         string `json:"email"`
	EmailCode     string `json:"email_code"`
	Password      string `json:"password"`
	PasswordCheck string `json:"password_check"`
	ReferrerCode  string `json:"referrer_code"`
}

// RegisterUser godoc
// @Summary      註冊新帳號
// @Description  註冊新帳號
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        username        body      string  true   "名稱"
// @Param        email           body      string  true   "信箱"
// @Param        email_code      body      string  true   "信箱驗證碼"
// @Param        password        body      string  true   "密碼"
// @Param        password_check  body      string  true   "密碼確認"
// @Param        referrer_code   body      string  false  "邀請碼"
// @Success      200       {object}  status.Response
// @Failure      400    {object}  status.Response
// @Failure      500    {object}  status.Response
// @Router       /register [post]
func RegisterUser(c *gin.Context) {
	var data RegisterReq

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	if !utils.CheckPasswordValid(data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.PasswordInvalid),
		})
		return
	}

	if data.Password != data.PasswordCheck {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.PasswordNotEqual),
		})
		return
	}

	err = services.ValidEmailCode(data.Email, data.EmailCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.ValidCodeIsIncorrect),
		})
		return
	}

	err = services.CheckRepeatUserEmail(data.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			status.RespStatus: status.NewResponse(status.ExistUser),
		})
		return
	}

	user := models.User{
		Email:    data.Email,
		UserName: data.UserName,
	}
	user.SetPassword(data.Password)

	err = services.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
	})
	return
}

type ForgetPasswordReq struct {
	Email string `json:"email"`
}

// ForgetPassword godoc
// @Summary      發送忘記密碼信
// @Description  發送忘記密碼信
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        email  body      string  true  "信箱"
// @Success      200       {object}  status.Response
// @Failure      400            {object}  status.Response
// @Failure      500            {object}  status.Response
// @Router       /forget_password [post]
func ForgetPassword(c *gin.Context) {
	var data ForgetPasswordReq

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	user, err := services.GetUserByEmail(data.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	err = services.SendResetPassword(user.Email)
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

type ForgetPasswordAccessReq struct {
	Email        string `json:"email"`
	RequestToken string `json:"request_token"`
}

// ForgetPasswordAccess godoc
// @Summary      獲得忘記密碼令牌
// @Description  獲得忘記密碼令牌
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        email          body      string  true  "信箱"
// @Param        request_token  body      string  true  "請求碼"
// @Success      200            {object}  status.ResponseWtihData{data=string}
// @Failure      400    {object}  status.Response
// @Failure      500    {object}  status.Response
// @Router       /forget_password_access [post]
func ForgetPasswordAccess(c *gin.Context) {
	var data ForgetPasswordAccessReq

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	accessToken, err := services.GetResetPasswordAccessToekn(data.Email, data.RequestToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.TokenIsInvalid),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   accessToken,
	})
}

type ResetPasswordWithTokenReq struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
	Password    string `json:"password"`
}

// ResetPasswordWithToken godoc
// @Summary      發送重設密碼
// @Description  發送重設密碼
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        email  body      string  true  "信箱"
// @Success      200    {object}  status.ResponseWtihData{data=string}
// @Failure      400             {object}  status.Response
// @Failure      500             {object}  status.Response
// @Router       /reset_password_with_token [post]
func ResetPasswordWithToken(c *gin.Context) {
	var data ResetPasswordWithTokenReq

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	if !utils.CheckPasswordValid(data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.PasswordInvalid),
		})
		return
	}

	if services.CheckResetPasswordAccessToekn(data.Email, data.AccessToken) {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.TokenIsInvalid),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
	})
}

type ChangePasswordReq struct {
	OldPassword   string `json:"old_password"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

// ChangeUserPassword godoc
// @Summary      修改密碼
// @Description  修改密碼
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        old_password    body      string  true  "舊密碼"
// @Param        password        body      string  true  "新密碼"
// @Param        check_password  body      string  true  "確認新密碼"
// @Success      200             {object}  status.ResponseWtihData{data=string}
// @Failure      400             {object}  status.Response
// @Failure      500             {object}  status.Response
// @Router       /change_password [post]
// @Security     BearerAuth
func ChangeUserPassword(c *gin.Context) {
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

	userID := c.GetUint(middlewares.ROLE_ID)
	if userID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	user, err := services.GetUserById(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	err = user.ComparePassword(data.OldPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.OldPasswordIsIncorrect),
		})
		return
	}

	user.SetPassword(data.Password)
	result := services.UpdateUser(&user)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success)})
}

type ChangeUserPasswordByAdminReq struct {
	UserID        int    `json:"user_id"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

// ChangeUserPasswordByAdmin godoc
// @Summary      修改使用者密碼
// @Description  修改使用者密碼
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user_id         body      string  true  "用戶ID"
// @Param        password        body      string  true  "新密碼"
// @Param        check_password  body      string  true  "確認新密碼"
// @Success      200             {object}  status.ResponseWtihData{data=string}
// @Failure      400       {object}  status.Response
// @Failure      500       {object}  status.Response
// @Router       /admin/change_password_by_admin [post]
// @Security     BearerAuth
func ChangeUserPasswordByAdmin(c *gin.Context) {
	var data ChangeUserPasswordByAdminReq

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

	user, err := services.GetUserById(uint(data.UserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	user.SetPassword(data.Password)
	result := database.DB.Model(&user).Updates(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success)})
}

func GetUserTFA(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)

	user, err := services.GetUserById(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.ExistUser),
		})
		return
	}

	qrCode, err := services.GetUserTFA(&user)
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

func EnableUserTFA(c *gin.Context) {
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

func DisableUserTFA(c *gin.Context) {
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

	err = services.DisableUserTFA(&user, data.TFA, false)
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

func DisableUserTFAByAdmin(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)

	user, err := services.GetUserById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	err = services.DisableUserTFA(&user, "", false)
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
