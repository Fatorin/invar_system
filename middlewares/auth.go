package middlewares

import (
	"invar/models"
	"invar/services"
	"invar/status"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleType uint

const (
	Admin = iota + 1
	User
)

const (
	ROLE_TYPE     = "ROLE_TYPE"
	ROLE_ID       = "ROLE_ID"
	REFRESH_TOKEN = "refresh_toekn"
)

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				status.RespStatus: status.NewResponse(status.NoToken),
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				status.RespStatus: status.NewResponse(status.TokenIsInvalid),
			})
			c.Abort()
			return
		}

		roleType, roleID, err := services.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				status.RespStatus: status.NewResponse(status.TokenIsInvalid),
			})
			c.Abort()
			return
		}

		c.Set(ROLE_TYPE, roleType)
		c.Set(ROLE_ID, roleID)

		c.Next()
	}
}

func CheckUserPremission() func(c *gin.Context) {
	return func(c *gin.Context) {
		roleType := c.GetInt(ROLE_TYPE)
		roleID := c.GetInt(ROLE_ID)

		if roleType != int(User) || roleID == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.NotExistUser),
			})
			c.Abort()
			return
		}

		user, err := services.GetUserById(uint(roleID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(status.NotExistUser),
			})
			c.Abort()
			return
		}

		if user.Status == models.Disabled {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.UserDisabled),
			})
			c.Abort()
			return
		}

		if user.Status != models.AuditSuccess {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.UserNoKYC),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckAdminPermission(premissionCode uint) func(c *gin.Context) {
	return func(c *gin.Context) {
		roleType := c.GetInt(ROLE_TYPE)
		roleID := c.GetInt(ROLE_ID)

		if roleType != int(Admin) || roleID == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.NotPermission),
			})
			c.Abort()
			return
		}

		admin, err := services.GetAdminById(uint(roleID))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.NotExistUser),
			})
			c.Abort()
			return
		}

		if !checkAdminPermission(admin.Premissions, int32(premissionCode)) {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.NotPermission),
			})
			c.Abort()
		}

		c.Next()
	}
}

func CheckFolderPermission(premissionCode uint) func(c *gin.Context) {
	return func(c *gin.Context) {
		roleType := c.GetInt(ROLE_TYPE)
		roleID := c.GetInt(ROLE_ID)
		path := c.Param("filepath")

		if len(strings.TrimSpace(path)) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.BadRequest),
			})
		}

		if roleType == 0 || roleID == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				status.RespStatus: status.NewResponse(status.NotPermission),
			})
			c.Abort()
			return
		}

		if roleType == int(Admin) {
			admin, err := services.GetAdminById(uint(roleID))
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					status.RespStatus: status.NewResponse(status.NotExistUser),
				})
				c.Abort()
				return
			}

			if !checkAdminPermission(admin.Premissions, int32(premissionCode)) {
				c.JSON(http.StatusForbidden, gin.H{
					status.RespStatus: status.NewResponse(status.NotPermission),
				})
				c.Abort()
			}
		}

		if roleType == int(User) {
			strs := strings.Split(path, "/")
			userID, _ := strconv.Atoi(strs[0])
			if userID != roleID {
				c.JSON(http.StatusForbidden, gin.H{
					status.RespStatus: status.NewResponse(status.NotPermission),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func checkAdminPermission(permissions []int32, premissionCode int32) bool {
	for _, value := range permissions {
		if value == premissionCode {
			return true
		}
	}

	return false
}
