package model

import (
	"time"

	"gorm.io/gorm"
)

type Timesheet struct {
	gorm.Model
	EmployeeID    uint      `gorm:"index"`
	WorkDate      time.Time `gorm:"type:date;uniqueIndex:idx_employee_date"`
	HoursWorked   float64   `gorm:"type:decimal(5,2);default:8.00"`           
	OvertimeHours float64   `gorm:"type:decimal(5,2);default:0.00"`           
	IsLeave       bool      `gorm:"default:false"`                           
	LeaveType     string    `gorm:"type:varchar(50)"`                        
	Note          string    `gorm:"type:text"`
}

func (Timesheet) TableName() string {
	return "timesheets"
}