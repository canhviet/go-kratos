package biz

import (
	"context"
	"time"

	"myapp/internal/data/model"
	"myapp/internal/repository"
)

type TimesheetUsecase struct {
	repo *repository.TimesheetRepo
}

func NewTimesheetUsecase(repo *repository.TimesheetRepo) *TimesheetUsecase {
	return &TimesheetUsecase{repo: repo}
}

func (uc *TimesheetUsecase) List(ctx context.Context, pageSize int, pageToken string) ([]*model.Timesheet, string, error) {
	return uc.repo.List(ctx, pageSize, pageToken)
}

func (uc *TimesheetUsecase) Get(ctx context.Context, id uint32) (*model.Timesheet, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *TimesheetUsecase) Create(ctx context.Context, employeeID uint32, monthYear time.Time, workingDays int, overtimeHours float64, leaveDays int) (*model.Timesheet, error) {
	timesheet := &model.Timesheet{
		EmployeeID:    uint(employeeID),
		MonthYear:     monthYear,
		WorkingDays:   workingDays,
		OvertimeHours: overtimeHours,
		LeaveDays:     leaveDays,
	}
	err := uc.repo.Create(ctx, timesheet)
	if err != nil {
		return nil, err
	}
	return timesheet, nil
}

func (uc *TimesheetUsecase) Update(ctx context.Context, id uint32, employeeID uint32, monthYear time.Time, workingDays int, overtimeHours float64, leaveDays int) (*model.Timesheet, error) {
	timesheet, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	timesheet.EmployeeID = uint(employeeID)
	timesheet.MonthYear = monthYear
	timesheet.WorkingDays = workingDays
	timesheet.OvertimeHours = overtimeHours
	timesheet.LeaveDays = leaveDays
	err = uc.repo.Update(ctx, timesheet)
	if err != nil {
		return nil, err
	}
	return timesheet, nil
}

func (uc *TimesheetUsecase) Delete(ctx context.Context, id uint32) error {
	return uc.repo.Delete(ctx, id)
}