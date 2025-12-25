package repository

import (
	"context"
	"time"

	"myapp/internal/data"
	"myapp/internal/data/model"
)

type timesheetRepo struct {
	data *data.Data
}


type TimesheetRepo interface {
	Create(ctx context.Context, ts *model.Timesheet) error

	GetMonthlySummary(
		ctx context.Context,
		employeeID uint,
		year int,
		month time.Month,
	) (workingDays int, overtimeHours float64, leaveDays int, err error)

	ExistsByEmployeeAndDate(
		ctx context.Context,
		employeeID uint,
		workDate time.Time,
	) (bool, error)
}

func NewTimesheetRepo(data *data.Data) *timesheetRepo {
	return &timesheetRepo{data: data}
}

func (r *timesheetRepo) Create(ctx context.Context, ts *model.Timesheet) error {
	return r.data.DB.WithContext(ctx).Create(ts).Error
}

func (r *timesheetRepo) ExistsByEmployeeAndDate(
	ctx context.Context,
	employeeID uint,
	workDate time.Time,
) (bool, error) {
	var count int64
	err := r.data.DB.WithContext(ctx).
		Model(&model.Timesheet{}).
		Where("employee_id = ? AND work_date = ?", employeeID, workDate.Truncate(24*time.Hour)).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *timesheetRepo) GetMonthlySummary(
	ctx context.Context,
	employeeID uint,
	year int,
	month time.Month,
) (workingDays int, overtimeHours float64, leaveDays int, err error) {

	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	start := time.Date(year, month, 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var results []struct {
		IsLeave       bool    `gorm:"column:is_leave"`
		OvertimeHours float64 `gorm:"column:overtime_hours"`
	}

	err = r.data.DB.WithContext(ctx).
		Model(&model.Timesheet{}).
		Select("is_leave, overtime_hours").
		Where("employee_id = ? AND work_date BETWEEN ? AND ?", employeeID, start, end).
		Scan(&results).Error
	if err != nil {
		return 0, 0, 0, err
	}

	workingDays = 0
	leaveDays = 0
	overtimeHours = 0.0

	for _, row := range results {
		if row.IsLeave {
			leaveDays++
		} else {
			workingDays++
			overtimeHours += row.OvertimeHours
		}
	}

	return workingDays, overtimeHours, leaveDays, nil
}