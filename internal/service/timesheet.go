// internal/service/timesheet.go
package service

import (
	"context"

	v1 "myapp/api/timesheet/v1"
	"myapp/internal/biz"
)

type TimesheetService struct {
	v1.UnimplementedTimesheetServer
	uc *biz.TimesheetUsecase
}

func NewTimesheetService(uc *biz.TimesheetUsecase) *TimesheetService {
	return &TimesheetService{uc: uc}
}

func (s *TimesheetService) Create(ctx context.Context, req *v1.CreateTimesheetRequest) (*v1.CreateTimesheetReply, error) {
	err := s.uc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.CreateTimesheetReply{Message: "attendance recorded successfully"}, nil
}