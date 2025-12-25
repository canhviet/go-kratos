package model

import (
	"time"

	"gorm.io/gorm"
)

type Payroll struct {
	gorm.Model
	EmployeeID    uint      `gorm:"index"`
	MonthYear     time.Time `gorm:"type:date;uniqueIndex:idx_employee_month"` // YYYY-MM-01
	WorkingDays   int       `gorm:"default:0"`
	OvertimeHours float64   `gorm:"type:decimal(8,2);default:0.00"`
	LeaveDays     int       `gorm:"default:0"`
	BasicSalary   float64   `gorm:"type:decimal(15,2)"`
	Allowances    float64   `gorm:"type:decimal(15,2);default:0.00"`
	GrossSalary   float64   `gorm:"type:decimal(15,2)"`
	Deductions    float64   `gorm:"type:decimal(15,2)"`
	NetSalary     float64   `gorm:"type:decimal(15,2)"`
	Status        string    `gorm:"type:varchar(50);default:'calculated'"`
}