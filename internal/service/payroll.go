package service

import (
	"context"

	v1"myapp/api/payroll/v1"
	"myapp/internal/biz"
)

type PayrollService struct {
	v1.UnimplementedPayrollServer

	uc *biz.PayrollUsecase
}

func NewPayrollService(uc *biz.PayrollUsecase) *PayrollService {
	return &PayrollService{uc: uc}
}

func (s *PayrollService) CalculatePayroll(ctx context.Context, req *v1.CalculatePayrollRequest) (*v1.CalculatePayrollReply, error) {
	return s.uc.CalculatePayroll(ctx, req)
}