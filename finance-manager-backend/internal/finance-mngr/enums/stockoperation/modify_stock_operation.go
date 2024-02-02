package stockoperation

import "finance-manager-backend/internal/finance-mngr/constants"

type ModifyStockOperation string

const (
	Undefined ModifyStockOperation = ""
	Add ModifyStockOperation = constants.ModifyStockOperationAdd
	Remove ModifyStockOperation = constants.ModifyStockOperationRemove
)