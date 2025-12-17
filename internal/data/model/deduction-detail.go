package model

import "gorm.io/gorm"

type DeductionDetail struct {
	gorm.Model
	PayrollID uint    `gorm:"index"` 
	Type      string  `gorm:"type:varchar(50)"` 
	Amount    float64 `gorm:"type:decimal(15,2)"` 
	Payroll   Payroll `gorm:"foreignKey:PayrollID"` 
}