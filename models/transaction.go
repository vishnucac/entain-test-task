package models

type Transaction struct {
	TransactionID string `gorm:"primaryKey;unique;not null"`
	Amount        string `gorm:"type:varchar(20);not null"`
	State         string `gorm:"type:varchar(10);not null"`
}
