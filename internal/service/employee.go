package service

import (
	"context"

	pb "myapp/api/employee/v1"
	"myapp/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type EmployeeService struct {
	pb.UnimplementedEmployeeServer
	uc *biz.EmployeeUsecase
}

func NewEmployeeService(uc *biz.EmployeeUsecase) *EmployeeService {
	return &EmployeeService{uc: uc}
}

func (s *EmployeeService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	employees, nextToken, err := s.uc.List(ctx, int(req.PageSize), req.PageToken)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListReply{
		NextPageToken: nextToken,
	}
	for _, e := range employees {
		resp.Items = append(resp.Items, &pb.EmployeeItem{
			Id:          uint32(e.ID),
			Name:        e.Name,
			Position:    e.Position,
			BaseSalary:  e.BaseSalary,
			BankAccount: e.BankAccount,
			JoinDate:    timestamppb.New(e.JoinDate),
			Dependents:  int32(e.Dependents),
		})
	}
	return resp, nil
}

func (s *EmployeeService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	employee, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetReply{
		Item: &pb.EmployeeItem{
			Id:          uint32(employee.ID),
			Name:        employee.Name,
			Position:    employee.Position,
			BaseSalary:  employee.BaseSalary,
			BankAccount: employee.BankAccount,
			JoinDate:    timestamppb.New(employee.JoinDate),
			Dependents:  int32(employee.Dependents),
		},
	}, nil
}

func (s *EmployeeService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	employee, err := s.uc.Create(ctx, req.Name, req.Position, req.BaseSalary, req.BankAccount, req.JoinDate.AsTime(), int(req.Dependents))
	if err != nil {
		return nil, err
	}
	return &pb.CreateReply{
		Item: &pb.EmployeeItem{
			Id:          uint32(employee.ID),
			Name:        employee.Name,
			Position:    employee.Position,
			BaseSalary:  employee.BaseSalary,
			BankAccount: employee.BankAccount,
			JoinDate:    timestamppb.New(employee.JoinDate),
			Dependents:  int32(employee.Dependents),
		},
	}, nil
}

func (s *EmployeeService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	employee, err := s.uc.Update(ctx, req.Id, req.Name, req.Position, req.BaseSalary, req.BankAccount, req.JoinDate.AsTime(), int(req.Dependents))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateReply{
		Item: &pb.EmployeeItem{
			Id:          uint32(employee.ID),
			Name:        employee.Name,
			Position:    employee.Position,
			BaseSalary:  employee.BaseSalary,
			BankAccount: employee.BankAccount,
			JoinDate:    timestamppb.New(employee.JoinDate),
			Dependents:  int32(employee.Dependents),
		},
	}, nil
}

func (s *EmployeeService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteReply{}, nil
}