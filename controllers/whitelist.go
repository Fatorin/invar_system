package controllers

import (
	"invar/middlewares"
	"invar/models"
	"invar/services"
	"invar/status"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetWhiteLists godoc
// @Summary      獲得使用者的白名單
// @Description  獲得使用者的白名單
// @Tags         WhiteList
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.ResponseWtihData{data=[]models.WhiteList}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /whitelist [get]
// @Security     BearerAuth
func GetWhiteLists(c *gin.Context) {
	id := c.GetInt(middlewares.ROLE_ID)

	whitelists, err := services.GetWhiteLists(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   whitelists,
	})
}

// GetWhiteLists godoc
// @Summary      獲得使用者的白名單
// @Description  獲得使用者的白名單
// @Tags         WhiteList
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "使用者ID"
// @Success      200  {object}  status.ResponseWtihData{data=[]models.WhiteList}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/whitelist/{id} [get]
// @Security     BearerAuth
func GetWhiteListsByAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	whitelists, err := services.GetWhiteLists(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   whitelists,
	})
}

type AddWhiteListReq struct {
	ChainName string `json:"chain_name"`
	ChainID   uint   `json:"chain_id"`
	Address   string `json:"address"`
	NickName  string `json:"nick_name"`
	Comment   string `json:"comment"`
}

// AddWhiteList godoc
// @Summary      新增使用者的白名單
// @Description  新增使用者的白名單
// @Tags         WhiteList
// @Accept       json
// @Produce      json
// @Param        chain_name     body      string  true  "區塊鏈名稱"
// @Param        chain_id  body      string  true  "區塊鏈ID"
// @Param        address     body      string  true  "地址"
// @Param        nick_name  body      string  false  "暱稱"
// @Param        comment     body      string  false  "備註"
// @Success      200  {object}  status.ResponseWtihData{data=models.WhiteList}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /whitelist [post]
// @Security     BearerAuth
func AddWhiteList(c *gin.Context) {
	var data AddWhiteListReq
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	whitelist := models.WhiteList{
		ChainName: data.ChainName,
		ChainID:   data.ChainID,
		Address:   data.Address,
		NickName:  data.NickName,
		Comment:   data.Comment,
	}

	err = services.AddWhiteList(&whitelist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   whitelist,
	})
}

type UpdateWhiteListReq struct {
	NickName string `json:"nick_name"`
	Comment  string `json:"comment"`
}

// UpdateWhiteList godoc
// @Summary      更新使用者的白名單
// @Description  更新使用者的白名單
// @Tags         WhiteList
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "白名單ID"
// @Param        nick_name  body      string  false  "暱稱"
// @Param        comment     body      string  false  "備註"
// @Success      200  {object}  status.ResponseWtihData{data=models.WhiteList}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /whitelist/{id} [patch]
// @Security     BearerAuth
func UpdateWhiteList(c *gin.Context) {
	var data UpdateWhiteListReq
	roleType := c.GetInt(middlewares.ROLE_TYPE)
	roleID := c.GetInt(middlewares.ROLE_ID)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	whitelist, err := services.GetWhiteList(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistWhiteList),
		})
		return
	}

	if roleType != middlewares.User && whitelist.UserID != uint(roleID) {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	err = c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	whitelist.NickName = data.NickName
	whitelist.Comment = data.Comment

	err = services.UpdateWhiteList(&whitelist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   whitelist,
	})
}

// DeleteWhiteList godoc
// @Summary      刪除使用者的白名單
// @Description  刪除使用者的白名單
// @Tags         WhiteList
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "白名單ID"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /whitelist/{id} [delete]
// @Router       /admin/whitelist/{id} [delete]
// @Security     BearerAuth
func DeleteWhiteList(c *gin.Context) {
	roleType := c.GetInt(middlewares.ROLE_TYPE)
	roleID := c.GetInt(middlewares.ROLE_ID)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.ExistUser),
		})
		return
	}

	if roleType == int(middlewares.User) && roleID != id {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	err = services.DeleteWhiteList(uint(id))
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
