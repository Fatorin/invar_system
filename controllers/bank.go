package controllers

import (
	"invar/middlewares"
	"invar/services"
	"invar/status"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetBanksReq struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
}

// GetBanks godoc
// @Summary      獲得全部使用者的財產狀況
// @Description  獲得全部使用者的財產狀況
// @Tags         Bank
// @Accept       json
// @Produce      json
// @Param        email     body      string  false  "信箱"
// @Param        username  body      string  false  "姓名"
// @Success      200       {object}  status.ResponseWtihData{data=[]models.Bank}
// @Failure      400       {object}  status.Response
// @Failure      500       {object}  status.Response
// @Router       /bank [get]
// @Security     BearerAuth
func GetBanks(c *gin.Context) {
	var request GetBanksReq

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	banks, err := services.GetBanks(request.Email, request.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   banks,
	})
}

// GetBank godoc
// @Summary      獲得使用者的財產狀況
// @Description  獲得使用者的財產狀況
// @Tags         Bank
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Bank ID"
// @Success      200  {object}  status.ResponseWtihData{data=models.Bank}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /bank/{id} [get]
// @Router       /admin/bank/{id} [get]
// @Security     BearerAuth
func GetBank(c *gin.Context) {
	roleType := c.GetInt(middlewares.ROLE_TYPE)
	roleID := c.GetInt(middlewares.ROLE_ID)
	bankID, _ := strconv.Atoi(c.Param("id"))

	bank, err := services.GetBank(uint(bankID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	if roleType == middlewares.User && roleID != int(bank.UserID) {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   bank,
	})
}
