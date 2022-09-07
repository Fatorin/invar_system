package controllers

import (
	"invar/middlewares"
	"invar/services"
	"invar/status"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticateReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	TFA      string `json:"tfa"`
}

// RefreshToken godoc
// @Summary      更換令牌
// @Description  更換令牌
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.ResponseWtihData{data=string}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /refresh_token [post]
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(middlewares.REFRESH_TOKEN)

	if refreshToken == "" || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			status.RespStatus: status.NewResponse(status.NoToken),
		})
		return
	}

	newRefreshToken, err := services.RevokeRefreshToken(refreshToken, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			status.RespStatus: status.NewResponse(status.TokenIsInvalid),
		})
		return
	}

	newAccessToken := services.GenerateAccessToken(newRefreshToken.RoleType, newRefreshToken.RoleID)
	c.SetCookie(middlewares.REFRESH_TOKEN, newRefreshToken.Token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   newAccessToken,
	})
}

// RevokeToken godoc
// @Summary      撤銷令牌
// @Description  撤銷令牌
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /revoke_token [post]
func RevokeToken(c *gin.Context) {
	refreshToken := c.Request.Header.Get(middlewares.REFRESH_TOKEN)

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			status.RespStatus: status.NewResponse(status.NoToken),
		})
		return
	}

	err := services.RevokeToken(refreshToken, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			status.RespStatus: status.NewResponse(status.TokenIsInvalid),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
	})
}

func checkLoginInfo(roleType uint, account, password, tfa string) (roleID uint, errCode int) {
	switch roleType {
	case middlewares.User:
		role, err := services.GetUserByEmail(account)
		if err != nil {
			return roleID, status.NotExistUser
		}

		err = role.ComparePassword(password)
		if err != nil {
			return roleID, status.IncorrectLoginInfo
		}

		err = services.CheckUserTFA(&role, tfa)
		if err != nil {
			return roleID, status.IncorrectTFA
		}

		return role.ID, errCode

	case middlewares.Admin:
		role, err := services.GetAdminByAccount(account)
		if err != nil {
			return roleID, status.NotExistUser
		}

		err = role.ComparePassword(password)
		if err != nil {
			return roleID, status.IncorrectLoginInfo
		}

		return role.ID, errCode
	}

	return roleID, status.Unkonwn
}
