package permission

const (
	QueryAdmin = iota + 1
	ModifyAdmin
	QueryUser
	ModifyUser
	QueryWhiteList
	ModifyWhiteList
	QueryProduct
	ModifyProduct
	QueryOrder
	ModifyOrder
	QueryStack
	ModifyStack
	DeleteStack
	QueryBank
)

func GetDefaultAdminPermission() []int32 {
	return []int32{QueryAdmin, ModifyAdmin, QueryUser, ModifyUser,
		QueryWhiteList, ModifyWhiteList, QueryProduct, ModifyProduct,
		QueryOrder, ModifyOrder, QueryStack, ModifyStack, DeleteStack, QueryBank}
}
