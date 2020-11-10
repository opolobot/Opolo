package models

type Guild struct {
	// gorm.Model
	ID    string
	Code  string
	Price uint
}
