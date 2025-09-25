package purchase_tour

type CartItem struct {
	ItemID string  `json:"item_id"`
	Price  float64 `json:"price"`
	TourID int64   `json:"tour_id"`
}

type PurchaseCartDetails struct {
	PurchaseID string     `json:"purchase_id"`
	CartID     string     `json:"cart_id"`
	UserID     string     `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalPrice float64    `json:"total_price"`
}

type PurchaseCartCommandType int8

const (
	CalculateTotal PurchaseCartCommandType = iota
	CheckBalance
	DeductBalance
	CompletePurchase
	RefundBalance
	AbortPurchase
	UnknownCommand
)

type PurchaseCartCommand struct {
	Purchase PurchaseCartDetails     `json:"purchase"`
	Type     PurchaseCartCommandType `json:"type"`
}

type PurchaseCartReplyType int8

const (
	TotalCalculated PurchaseCartReplyType = iota
	BalanceSufficient
	BalanceInsufficient
	BalanceDeducted
	BalanceDeductionFailed
	PurchaseCompleted
	BalanceRefunded
	PurchaseAborted
	UnknownReply
)

type PurchaseCartReply struct {
	Purchase PurchaseCartDetails   `json:"purchase"`
	Type     PurchaseCartReplyType `json:"type"`
	Message  string                `json:"message,omitempty"`
}
