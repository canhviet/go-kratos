package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(255);not null"` 
	Position     string  `gorm:"type:varchar(100)"` 
	BaseSalary   float64 `gorm:"type:decimal(15,2);not null"` 
	BankAccount  string  `gorm:"type:varchar(50)"` 
	JoinDate     time.Time `gorm:"type:date"` 
	Dependents   int     `gorm:"default:0"` 
	Timesheets   []Timesheet 
	Payrolls     []Payroll   
}