package model

type TourPurchaseToken struct {
	Id        string `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID string `gorm:"type:uuid;index"      json:"accountId"`
	TourId    string `gorm:"index" 			   json:"tourId"`
}
