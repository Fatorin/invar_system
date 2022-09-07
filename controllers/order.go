package controllers

import (
	"invar/middlewares"
	"invar/models"
	"invar/services"
	"invar/status"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary      使用者獲得自己的全部訂單
// @Description  使用者獲得自己的全部訂單
// @Tags         Order
// @Accept       json
// @Produce      json
// @Success      200                {object}  status.ResponseWtihData{data=[]models.Order}
// @Failure      400                {object}  status.Response
// @Failure      500                {object}  status.Response
// @Router       /order [get]
// @Security     BearerAuth
func GetOrders(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)
	orders, err := services.GetOrders(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   orders,
	})
}

type GetOrdersByAdminReq struct {
	Serial     string    `json:"serial"`
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	CreateAt   time.Time `json:"create_at"`
	PageNumber int       `json:"page_number"`
	PerPage    int       `json:"per_page"`
}

// GetOrders godoc
// @Summary      獲得所有使用者的訂單
// @Description  獲得所有使用者的訂單
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        serial       body      string  false  "序號"
// @Param        user_name    body      string  false  "姓名"
// @Param        user_email   body      string  false  "信箱"
// @Param        create_at    body      string  false  "建立日期"
// @Param        page_number  body      int     false  "分頁號碼"
// @Param        per_page     body      int     false  "每頁數量"
// @Success      200          {object}  status.ResponseWtihData{data=[]models.Order}
// @Failure      400          {object}  status.Response
// @Failure      500          {object}  status.Response
// @Router       /admin/order [get]
// @Security     BearerAuth
func GetOrdersByAdmin(c *gin.Context) {
	var request GetOrdersByAdminReq
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	orders, err := services.GetOrdersByAdmin(request.Serial, request.UserEmail, request.UserName, request.CreateAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   orders,
	})
}

type AddOrderReq struct {
	OrderDetails []OrderDetail `json:"order_details"`
}

type OrderDetail struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

// AddOrder godoc
// @Summary      新增訂單
// @Description  新增訂單
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        order_details  body      AddOrderReq  true  "商品ID與數量"
// @Success      200            {object}  status.Response
// @Failure      400            {object}  status.Response
// @Failure      500            {object}  status.Response
// @Router       /order [post]
// @Security     BearerAuth
func AddOrder(c *gin.Context) {
	roleID := c.GetInt(middlewares.ROLE_ID)
	var request AddOrderReq

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	orderItems := make([]models.OrderItem, 0)
	for _, v := range request.OrderDetails {
		product, err := services.GetProduct(v.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(status.NotExistProduct),
			})
			return
		}

		if !product.Status {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(status.HasBeenRemoved),
			})
			return
		}

		if product.Stock < v.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				status.RespStatus: status.NewResponse(status.OutOfStock),
			})
			return
		}

		orderItem := models.OrderItem{
			ProductID: product.ID,
			Price:     product.Price,
			Quantity:  v.Quantity,
		}

		orderItems = append(orderItems, orderItem)
	}

	err = services.AddOrder(uint(roleID), orderItems)
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

// CancelOrder godoc
// @Summary      取消訂單
// @Description  取消訂單
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        id                 path      int     true  "訂單ID"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /cancel_order/{id} [patch]
// @Router       /admin/cancel_order/{id} [patch]
// @Security     BearerAuth
func CancelOrder(c *gin.Context) {
	roleType := c.GetInt(middlewares.ROLE_TYPE)
	roleID := c.GetInt(middlewares.ROLE_ID)
	orderID, _ := strconv.Atoi(c.Param("id"))

	order, err := services.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistOrder),
		})
		return
	}

	if roleType == middlewares.User && roleID != int(order.UserID) {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	err = services.CancelOrder(&order)
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

type PaymentOrderReq struct {
	TransactionChain string `json:"transaction_chain"`
	TransactionID    string `json:"transaction_id"`
}

// PaymentOrder godoc
// @Summary      付款通知
// @Description  付款通知
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "訂單ID"
// @Param        transaction_chain  body      string  true  "區塊鏈名稱"
// @Param        transaction_id     body      string  true  "交易序號"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /payment_order/{id} [patch]
// @Security     BearerAuth
func PaymentOrder(c *gin.Context) {
	roleType := c.GetInt(middlewares.ROLE_TYPE)
	roleID := c.GetInt(middlewares.ROLE_ID)
	orderID, _ := strconv.Atoi(c.Param("id"))
	var request PaymentOrderReq

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	order, err := services.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistOrder),
		})
		return
	}

	if roleType == middlewares.User && roleID != int(order.UserID) {
		c.JSON(http.StatusForbidden, gin.H{
			status.RespStatus: status.NewResponse(status.NotPermission),
		})
		return
	}

	order.TransactionChain = request.TransactionChain
	order.TransactionID = request.TransactionID
	order.Status = models.WaitConfirmPayment
	//寄通知給管理者
	err = services.UpdateOrder(&order)
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

// CompletedOrder godoc
// @Summary      完成訂單
// @Description  完成訂單
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "訂單ID"
// @Success      200  {object}  status.Response
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/completed_order/{id} [patch]
// @Security     BearerAuth
func CompletedOrder(c *gin.Context) {
	orderID, _ := strconv.Atoi(c.Param("id"))

	order, err := services.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistOrder),
		})
		return
	}

	err = services.CompletedOrder(&order)
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
