package controllers

import (
	"invar/middlewares"
	"invar/models"
	"invar/services"
	"invar/status"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// GetStacks godoc
// @Summary      獲得可質押的項目
// @Description  獲得可質押的項目
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.ResponseWtihData{data=[]models.Stack}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack [get]
// @Security     BearerAuth
func GetStacks(c *gin.Context) {
	stacks, err := services.GetStacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   stacks,
	})
}

// GetStack godoc
// @Summary      獲得可質押的項目
// @Description  獲得可質押的項目
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stack ID"
// @Success      200  {object}  status.ResponseWtihData{data=models.Stack}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack/{id} [get]
// @Security     BearerAuth
func GetStack(c *gin.Context) {
	stackID, _ := strconv.Atoi(c.Param("id"))
	stack, err := services.GetStack(uint(stackID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   stack,
	})
}

// AddStack godoc
// @Summary      新增得可質押的項目
// @Description  新增可質押的項目
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param          product_id      int  true  "可質押的商品ID"
// @Param        profit  body      int  true  "配息利率"
// @Param        contract_time_bound  body      int  true  "綁約時間"
// @Param        profit_interval_month  body      int  true  "配息週期"
// @Param        stack_permissions  body      []int  true  "User ID"
// @Success      200  {object}  status.ResponseWtihData{data=models.Stack}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack [post]
// @Security     BearerAuth
func AddStack(c *gin.Context) {
	var stack models.Stack

	bindErr := c.BindJSON(&stack)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	err := services.AddStack(&stack)
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

// UpdateStack godoc
// @Summary      更新可質押的項目
// @Description  更新可質押的項目
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Stack ID"
// @Param        name                 formData  string  false  "姓名"
// @Success      200  {object}  status.ResponseWtihData{data=models.Stack}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack/{id} [patch]
// @Security     BearerAuth
func UpdateStack(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	stack, err := services.GetStack(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistStack),
		})
		return
	}

	bindErr := c.BindJSON(&stack)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	err = services.UpdateStack(&stack)
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

// GetStacksRecord godoc
// @Summary      獲得所有的質押紀錄
// @Description  獲得所有的質押紀錄
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.ResponseWtihData{data=[]models.StackRecord}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack_record/ [get]
// @Security     BearerAuth
func GetStacksRecord(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)
	stacksrecord, err := services.GetStacksRecord(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   stacksrecord,
	})
}

type GetStacksRecordByAdminReq struct {
	UserID uint `json:"user_id"`
}

// GetStacksRecord godoc
// @Summary      獲得所有的質押紀錄
// @Description  獲得所有的質押紀錄
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        user_id  body      int  true  "User ID"
// @Success      200  {object}  status.ResponseWtihData{data=[]models.StackRecord}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/stack_record/ [get]
// @Security     BearerAuth
func GetStacksRecordByAdmin(c *gin.Context) {
	var req GetStacksRecordByAdminReq
	bindErr := c.BindJSON(&req)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	stacksrecord, err := services.GetStacksRecord(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   stacksrecord,
	})
}

// GetStackRecord godoc
// @Summary      獲得單筆質押紀錄
// @Description  獲得單筆質押紀錄
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "StackRecord ID"
// @Success      200  {object}  status.ResponseWtihData{data=models.StackRecord}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack_record/{id} [get]
// @Router       /admin/stack_record/{id} [get]
// @Security     BearerAuth
func GetStackRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	roleType := c.GetInt(middlewares.ROLE_TYPE)
	roleID := c.GetInt(middlewares.ROLE_ID)
	stacksrecord, err := services.GetStackRecord(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	if roleType != middlewares.Admin && stacksrecord.UserID != uint(roleID) {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   stacksrecord,
	})
}

// GetStackRecord godoc
// @Summary      獲得單筆質押紀錄
// @Description  獲得單筆質押紀錄
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "StackRecord ID"
// @Success      200  {object}  status.ResponseWtihData{data=models.StackRecord}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /stack_record/{id} [get]
// @Router       /admin/stack_record/{id} [get]
// @Security     BearerAuth
func UnlockStackRecord(c *gin.Context) {
	//解除質押，需要確認是否質押結束
}

func AutoRenewStackRecord(c *gin.Context) {
	//當前質押項目期滿後，可自動續約(只修改當前)
}

func ExchangeStack(c *gin.Context) {
	//將質押結束項目兌換回IVT
}

type AddStackProfitRecordReq struct {
	UserID        uint            `json:"user_id"`
	StackRecordID uint            `json:"stack_record_id"`
	Profit        decimal.Decimal `json:"profit"`
	Comment       string          `json:"comment"`
}

// AddStackProfitRecord godoc
// @Summary      獲得單筆質押紀錄
// @Description  獲得單筆質押紀錄
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        user_id  body      int  true  "User ID"
// @Param        stack_record_id  body      int  true  "質押紀錄ID"
// @Param        profit  body      int  true  "金額"
// @Param        comment  body      string  true  "註解"
// @Success      200  {object}  status.ResponseWtihData{data=models.StackProfitRecord}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/stack_profit_record [post]
// @Security     BearerAuth
func AddStackProfitRecord(c *gin.Context) {
	var recordReq AddStackProfitRecordReq

	bindErr := c.BindJSON(&recordReq)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	bank, err := services.GetBank(recordReq.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	var record = models.StackProfitRecord{
		StackRecordID: recordReq.StackRecordID,
		Profit:        recordReq.Profit,
		Comment:       recordReq.Comment,
	}

	err = services.AddStackProfitRecord(&bank, &record)
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

type DeleteStackProfitRecordReq struct {
	UserID              uint `json:"user_id"`
	StackProfitRecordID uint `json:"stack_record_profit_id"`
}

// DeleteStackProfitRecord godoc
// @Summary      刪除單筆質押配息紀錄
// @Description  刪除單筆質押配息紀錄
// @Tags         Stack
// @Accept       json
// @Produce      json
// @Param        user_id  body      int  true  "User ID"
// @Param        stack_record_profit_id  body      int  true  "質押配息紀錄ID"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/stack_profit_record/ [delete]
// @Security     BearerAuth
func DeleteStackProfitRecord(c *gin.Context) {
	var recordReq DeleteStackProfitRecordReq

	bindErr := c.BindJSON(&recordReq)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	bank, err := services.GetBank(recordReq.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	record, err := services.GetStackProfitRecord(recordReq.StackProfitRecordID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistStack),
		})
		return
	}

	err = services.DeleteStackProfitRecord(&bank, &record)
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
