package controllers

import (
	"invar/models"
	"invar/services"
	"invar/status"
	"invar/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// GetProducts godoc
// @Summary      獲得全部商品清單
// @Description  獲得全部商品清單
// @Tags         Product
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.ResponseWtihData{data=[]models.Product}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /product [get]
// @Security     BearerAuth
func GetProducts(c *gin.Context) {
	products, err := services.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   products,
	})
}

// GetProductsByAdmin godoc
// @Summary      獲得全部商品清單
// @Description  獲得全部商品清單
// @Tags         Product
// @Accept       json
// @Produce      json
// @Success      200  {object}  status.ResponseWtihData{data=[]models.Product}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /admin/product [get]
// @Security     BearerAuth
func GetProductsByAdmin(c *gin.Context) {
	products, err := services.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   products,
	})
}

// GetProduct godoc
// @Summary      獲得全部商品清單
// @Description  獲得全部商品清單
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "商品ID"
// @Success      200  {object}  status.ResponseWtihData{data=models.Product}
// @Failure      400  {object}  status.Response
// @Failure      500  {object}  status.Response
// @Router       /product/{id} [get]
// @Router       /admin/product/{id} [get]
// @Security     BearerAuth
func GetProduct(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
	}

	product, err := services.GetProduct(uint(id))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistProduct),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   product,
	})
}

type AddProductReq struct {
	Stock           uint            `json:"stock"`
	Status          bool            `json:"status"`
	Price           decimal.Decimal `json:"price"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	ContractChain   string          `json:"contract_chain"`
	ContractType    string          `json:"contract_type"`
	ContractAddress string          `json:"contract_address"`
	BuyPermissions  pq.Int32Array   `json:"buy_remissions" gorm:"type:integer[]" swaggertype:"array,number"`
}

// AddProduct godoc
// @Summary      建立新商品
// @Description  建立新商品
// @Tags         Product
// @Accept       mpfd
// @Produce      json
// @Param        stock             formData  int     true  "庫存數量"  default(1)
// @Param        status            formData  bool    true  "商品狀態"
// @Param        price             formData  int     true  "商品價格"
// @Param        title             formData  string  true  "商品標題"
// @Param        description       formData  string  true  "商品描述"
// @Param        contract_chain    formData  string  true  "區塊鏈名稱"
// @Param        contract_type     formData  string  true  "合約類型"
// @Param        contract_address  formData  string  true  "合約地址"
// @Param        buy_remissions    formData  []int   true  "購買權限"
// @Param        small_image       formData  file    true  "預覽檔案"
// @Param        large_image       formData  file    true  "完整檔案"
// @Success      200               {object}  status.ResponseWtihData{data=models.Product}
// @Failure      400               {object}  status.Response
// @Failure      417               {object}  status.Response
// @Failure      500               {object}  status.Response
// @Router       /admin/product [post]
// @Security     BearerAuth
func AddProduct(c *gin.Context) {
	var data AddProductReq

	bindErr := c.BindWith(&data, binding.Form)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	smallFile, bindErr := c.FormFile("small_image")
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	largeFile, bindErr := c.FormFile("large_image")
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	fileSupportTypes := []int{utils.PNG, utils.JPEG, utils.JPG}
	if utils.CheckMediaType(smallFile.Filename, fileSupportTypes) {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			status.RespStatus: status.NewResponse(status.UnsupportedMediaType),
		})
		return
	}

	if utils.CheckMediaType(largeFile.Filename, fileSupportTypes) {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			status.RespStatus: status.NewResponse(status.UnsupportedMediaType),
		})
		return
	}

	savePath := status.SaveFileFolderPath + status.SaveFilePublicFolderPath + "/"
	smallFileSavePath := utils.FilePathGenerator(smallFile.Filename, savePath)
	err := c.SaveUploadedFile(smallFile, smallFileSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	largeFileSavePath := utils.FilePathGenerator(largeFile.Filename, savePath)
	err = c.SaveUploadedFile(largeFile, largeFileSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	var product = models.Product{
		Stock:           data.Stock,
		Status:          data.Status,
		Price:           data.Price,
		Title:           data.Title,
		Description:     data.Description,
		Image:           largeFileSavePath[1:],
		PreviewImage:    smallFileSavePath[1:],
		ContractChain:   data.ContractChain,
		ContractType:    data.ContractType,
		ContractAddress: data.ContractAddress,
		BuyPermissions:  data.BuyPermissions,
	}

	err = services.AddProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			status.RespStatus: status.NewResponse(status.Unkonwn),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   product,
	})
}

// UpdateProduct godoc
// @Summary      建立新商品
// @Description  建立新商品
// @Tags         Product
// @Accept       mpfd
// @Produce      json
// @Param        stock             formData  int     false  "庫存數量"  default(1)
// @Param        status            formData  bool    false  "商品狀態"
// @Param        price             formData  int     false  "商品價格"
// @Param        title             formData  string  false  "商品標題"
// @Param        description       formData  string  false  "商品描述"
// @Param        contract_chain    formData  string  false  "區塊鏈名稱"
// @Param        contract_type     formData  string  false  "合約類型"
// @Param        contract_address  formData  string  false  "合約地址"
// @Param        buy_remissions    formData  []int   false  "購買權限"
// @Param        small_image       formData  file    false  "預覽檔案"
// @Param        large_image       formData  file    false  "完整檔案"
// @Success      200               {object}  status.ResponseWtihData{data=models.Product}
// @Failure      400               {object}  status.Response
// @Failure      417               {object}  status.Response
// @Failure      500               {object}  status.Response
// @Router       /admin/product [patch]
// @Security     BearerAuth
func UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := services.GetProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.NotExistProduct),
		})
		return
	}

	bindErr := c.BindWith(&product, binding.Form)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.BadRequest),
		})
		return
	}

	savePath := status.SaveFileFolderPath + status.SaveFilePublicFolderPath + "/"
	fileSupportTypes := []int{utils.PNG, utils.JPEG, utils.JPG}
	smallFile, bindErr := c.FormFile("small_image")
	if bindErr == nil {
		if utils.CheckMediaType(smallFile.Filename, fileSupportTypes) {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				status.RespStatus: status.NewResponse(status.UnsupportedMediaType),
			})
			return
		}

		fileSavePath := utils.FilePathGenerator(smallFile.Filename, savePath)
		err := c.SaveUploadedFile(smallFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		err = utils.RemoveFile(product.PreviewImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}
		product.PreviewImage = fileSavePath[1:]
	}

	largeFile, bindErr := c.FormFile("large_image")
	if bindErr == nil {
		if utils.CheckMediaType(largeFile.Filename, fileSupportTypes) {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				status.RespStatus: status.NewResponse(status.UnsupportedMediaType),
			})
			return
		}

		fileSavePath := utils.FilePathGenerator(smallFile.Filename, savePath)
		err := c.SaveUploadedFile(smallFile, fileSavePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}

		err = utils.RemoveFile(product.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				status.RespStatus: status.NewResponse(status.Unkonwn),
			})
			return
		}
		product.Image = fileSavePath[1:]
	}

	err = services.UpdateProduct(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			status.RespStatus: status.NewResponse(status.UpdateFail),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		status.RespStatus: status.NewResponse(status.Success),
		status.RespData:   product,
	})
}
