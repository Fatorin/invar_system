package status

const RespStatus = "status"
const RespData = "data"
const SaveFileFolderPath = "./static"
const SaveFilePublicFolderPath = "/public"

const (
	// Common 0~999
	Unkonwn              = -1
	Success              = 0
	BadRequest           = 1
	UpdateFail           = 2
	NoToken              = 3
	TokenIsInvalid       = 4
	TokenExpired         = 5
	NotPermission        = 6
	ValidCodeIsIncorrect = 7
	ValidCodeIsExpired   = 8
	RequestTooFrequently = 9
	UserDisabled         = 10
	UserNoKYC            = 11
	UserHasKYC           = 12
	TooLarge             = 13
	UnsupportedMediaType = 14
	// Register
	PasswordInvalid        = 1001
	PasswordNotEqual       = 1002
	OldPasswordIsIncorrect = 1003
	NotExistUser           = 1004
	ExistUser              = 1005
	// Login
	IncorrectLoginInfo = 2001
	IncorrectTFA       = 2002
	// Admin
	// User
	// WhiteList
	NotExistWhiteList = 5001
	// Product
	NotExistProduct = 6001
	OutOfStock      = 6002
	HasBeenRemoved  = 6003
	// Order
	NotExistOrder = 7001
	// Stack
	NotExistStack = 8001
)

type Response struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

type ResponseWtihData struct {
	Status Response `json:"status"`
	Data   any      `json:"data"`
}

var errorText = map[int]string{
	// Common
	Unkonwn:              "未知的錯誤",
	Success:              "成功",
	BadRequest:           "請求資料有誤",
	UpdateFail:           "更新失敗",
	NoToken:              "令牌不存在",
	TokenIsInvalid:       "令牌不合法",
	TokenExpired:         "令牌已過期",
	NotPermission:        "沒有權限",
	ValidCodeIsIncorrect: "驗證碼不正確",
	ValidCodeIsExpired:   "驗證碼已過期或已使用",
	RequestTooFrequently: "請求次數過於頻繁",
	UserDisabled:         "帳戶已停用",
	UserNoKYC:            "使用者尚未通過實名認證",
	UserHasKYC:           "使用者已通過實名認證",
	TooLarge:             "請求文件過大",
	UnsupportedMediaType: "不支援的檔案",
	// Register
	PasswordInvalid:        "密碼不合法",
	PasswordNotEqual:       "密碼不一致",
	OldPasswordIsIncorrect: "舊密碼錯誤",
	NotExistUser:           "不存在的使用者",
	ExistUser:              "已存在的使用者",
	// Login
	IncorrectLoginInfo: "帳號或密碼錯誤",
	IncorrectTFA:       "二階段驗證碼錯誤",
	// WhiteList
	NotExistWhiteList: "不存在的白名單",
	// Product
	NotExistProduct: "不存在的商品",
	OutOfStock:      "商品庫存不足",
	HasBeenRemoved:  "已下架",
	// Order
	NotExistOrder: "不存在的訂單",
	// Stack
	NotExistStack: "不存在的質押項目",
}

func ErrorText(code int) string {
	return errorText[code]
}

func NewResponse(code int) Response {
	resp := Response{
		StatusCode:    code,
		StatusMessage: ErrorText(code),
	}

	return resp
}
