package model

import (
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	gorm.Model
	EmployeeID   uint    `gorm:"index"` 
	MonthYear    time.Time `gorm:"type:date;index"` 
	BasicSalary  float64 `gorm:"type:decimal(15,2)"` 
	Allowances   float64 `gorm:"type:decimal(15,2)"` 
	Deductions   float64 `gorm:"type:decimal(15,2)"` 
	NetSalary    float64 `gorm:"type:decimal(15,2)"` 
	Status       string  `gorm:"type:varchar(50);default:'pending'"` 
}