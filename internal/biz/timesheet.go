package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	v1 "myapp/api/timesheet/v1"
	"myapp/internal/data/model"
	"myapp/internal/repository"
)

type TimesheetUsecase struct {
	repo repository.TimesheetRepo
}

func NewTimesheetUsecase(repo repository.TimesheetRepo) *TimesheetUsecase {
	return &TimesheetUsecase{repo: repo}
}

func (uc *TimesheetUsecase) Create(ctx context.Context, req *v1.CreateTimesheetRequest) error {
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	workDate := req.WorkDate.AsTime().In(location).Truncate(24 * time.Hour)

	exists, err := uc.repo.ExistsByEmployeeAndDate(ctx, uint(req.EmployeeId), workDate)
	if err != nil {
		return fmt.Errorf("check duplicate attendance: %w", err)
	}
	if exists {
		return errors.New("attendance already recorded for this date")
	}

	ts := &model.Timesheet{
		EmployeeID:    uint(req.EmployeeId),
		WorkDate:      workDate,
		HoursWorked:   req.HoursWorked,
		OvertimeHours: req.OvertimeHours,
		IsLeave:       req.IsLeave,
		LeaveType:     req.LeaveType,
		Note:          req.Note,
	}

	return uc.repo.Create(ctx, ts)
}