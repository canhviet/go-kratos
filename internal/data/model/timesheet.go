package model

import (
	"time"

	"gorm.io/gorm"
)

type Timesheet struct {
	gorm.Model
	EmployeeID    uint    `gorm:"index"` 
	MonthYear     time.Time `gorm:"type:date;index"` 
	WorkingDays   int     `gorm:"default:0"` 
	OvertimeHours float64 `gorm:"type:decimal(5,2);default:0"` 
	LeaveDays     int     `gorm:"default:0"` 
	Employee      Employee `gorm:"foreignKey:EmployeeID"`
}