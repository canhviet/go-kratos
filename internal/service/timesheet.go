package service

import (
	"context"

	pb "myapp/api/timesheet/v1"
	"myapp/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type TimesheetService struct {
	pb.UnimplementedTimesheetServer
	uc *biz.TimesheetUsecase
}

func NewTimesheetService(uc *biz.TimesheetUsecase) *TimesheetService {
	return &TimesheetService{uc: uc}
}

func (s *TimesheetService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	timesheets, nextToken, err := s.uc.List(ctx, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListReply{
		NextPageToken: nextToken,
	}
	for _, t := range timesheets {
		resp.Items = append(resp.Items, &pb.TimesheetItem{
			Id:            uint32(t.ID),
			EmployeeId:    uint32(t.EmployeeID),
			MonthYear:     timestamppb.New(t.MonthYear),
			WorkingDays:   int32(t.WorkingDays),
			OvertimeHours: t.OvertimeHours,
			LeaveDays:     int32(t.LeaveDays),
		})
	}
	return resp, nil
}

func (s *TimesheetService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	timesheet, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetReply{
		Item: &pb.TimesheetItem{
			Id:            uint32(timesheet.ID),
			EmployeeId:    uint32(timesheet.EmployeeID),
			MonthYear:     timestamppb.New(timesheet.MonthYear),
			WorkingDays:   int32(timesheet.WorkingDays),
			OvertimeHours: timesheet.OvertimeHours,
			LeaveDays:     int32(timesheet.LeaveDays),
		},
	}, nil
}

func (s *TimesheetService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	timesheet, err := s.uc.Create(ctx, req.EmployeeId, req.MonthYear.AsTime(), int(req.WorkingDays), req.OvertimeHours, int(req.LeaveDays))
	if err != nil {
		return nil, err
	}
	return &pb.CreateReply{
		Item: &pb.TimesheetItem{
			Id:            uint32(timesheet.ID),
			EmployeeId:    uint32(timesheet.EmployeeID),
			MonthYear:     timestamppb.New(timesheet.MonthYear),
			WorkingDays:   int32(timesheet.WorkingDays),
			OvertimeHours: timesheet.OvertimeHours,
			LeaveDays:     int32(timesheet.LeaveDays),
		},
	}, nil
}

func (s *TimesheetService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	timesheet, err := s.uc.Update(ctx, req.Id, req.EmployeeId, req.MonthYear.AsTime(), int(req.WorkingDays), req.OvertimeHours, int(req.LeaveDays))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateReply{
		Item: &pb.TimesheetItem{
			Id:            uint32(timesheet.ID),
			EmployeeId:    uint32(timesheet.EmployeeID),
			MonthYear:     timestamppb.New(timesheet.MonthYear),
			WorkingDays:   int32(timesheet.WorkingDays),
			OvertimeHours: timesheet.OvertimeHours,
			LeaveDays:     int32(timesheet.LeaveDays),
		},
	}, nil
}

func (s *TimesheetService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteReply{}, nil
}