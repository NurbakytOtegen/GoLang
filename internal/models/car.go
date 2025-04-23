package models

type Car struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}
