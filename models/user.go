package models

type User struct {
	UserID  uint64  `json:"userId" gorm:"primaryKey"`
	Balance float64 `json:"balance"`
}
