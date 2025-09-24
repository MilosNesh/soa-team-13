package model

type ShoppingCart struct {
	Id        string      `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID string      `gorm:"type:uuid;index"      json:"accountId"`
	Items     []OrderItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}
