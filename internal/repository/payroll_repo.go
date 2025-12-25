package repository

import (
	"context"
	"myapp/internal/data"
	"myapp/internal/data/model"

)

type PayrollRepo interface {
	SavePayroll(ctx context.Context, p *model.Payroll) error
}

type payrollRepo struct {
	data *data.Data
}

func NewPayrollRepo(data *data.Data) *payrollRepo {
	return &payrollRepo{data: data}
}

func (r *payrollRepo) SavePayroll(ctx context.Context, p *model.Payroll) error {
	return r.data.DB.WithContext(ctx).Create(p).Error
}
