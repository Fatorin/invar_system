package controllers

import (
	"invar/services"
	"invar/status"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailCodeReq struct {
	Email string `json:"email"`
}

// GetEmailCode godoc
// @Summary      獲得信箱驗證碼
// @Description  獲得信箱驗證碼
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        email  body      string  true  "信箱"
// @Success      200    {object}  status.ResponseWtihData{data=string}
// @Failure      400    {object}  status.Response
// @Failure      500    {object}  status.Response
// @Router       /get_email_code [post]
func GetEmailCode(c *gin.Context) {
	var data EmailCodeReq

	bindErr := c.BindJSON(&data)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	err := services.CheckRepeatUserEmail(data.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.ExistUser),
		})
		return
	}

	preCode, errCode := services.RequestMailCode(data.Email)
	if errCode != status.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(errCode),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   preCode,
	})
}
