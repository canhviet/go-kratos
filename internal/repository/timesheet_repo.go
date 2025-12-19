package repository

import (
	"context"
	"strconv"

	"myapp/internal/data"
	"myapp/internal/data/model"
)

type TimesheetRepo struct {
	data *data.Data
}

func NewTimesheetRepo(data *data.Data) *TimesheetRepo {
	return &TimesheetRepo{data: data}
}

func (r *TimesheetRepo) List(ctx context.Context, pageSize int, pageToken string) ([]*model.Timesheet, string, error) {
	var offset int
	if pageToken != "" {
		offset, _ = strconv.Atoi(pageToken) // Improve error handling in production
	}

	var timesheets []*model.Timesheet
	err := r.data.DB.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&timesheets).Error
	if err != nil {
		return nil, "", err
	}

	nextToken := ""
	if len(timesheets) == pageSize {
		nextToken = strconv.Itoa(offset + pageSize)
	}

	return timesheets, nextToken, nil
}

func (r *TimesheetRepo) Get(ctx context.Context, id uint32) (*model.Timesheet, error) {
	var timesheet model.Timesheet
	err := r.data.DB.WithContext(ctx).First(&timesheet, id).Error
	if err != nil {
		return nil, err
	}
	return &timesheet, nil
}

func (r *TimesheetRepo) Create(ctx context.Context, timesheet *model.Timesheet) error {
	return r.data.DB.WithContext(ctx).Create(timesheet).Error
}

func (r *TimesheetRepo) Update(ctx context.Context, timesheet *model.Timesheet) error {
	return r.data.DB.WithContext(ctx).Save(timesheet).Error
}

func (r *TimesheetRepo) Delete(ctx context.Context, id uint32) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Timesheet{}, id).Error
}