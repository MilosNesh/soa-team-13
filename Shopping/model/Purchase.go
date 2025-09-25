package model

type Purchase struct {
	ID         int         `json:"id"`
	TouristID  string      `json:"tourist_id"`
	TotalPrice int64       `json:"total_price"`
	Items      []OrderItem `json:"items"`
}
