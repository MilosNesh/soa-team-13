package model

type OrderItem struct {
	Id             string `gorm:"type:uuid;primaryKey" json:"id"`
	ShoppingCartId string `gorm:"type:uuid;not null;index" json:"shoppingCartId"`
	Name           string `json:"name"`
	Price          int64  `json:"price"`
	TourId         int64  `json:"tourId"`
}
