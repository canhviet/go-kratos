package repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/data"
	"myapp/internal/data/model"
	"time"

	"gorm.io/gorm"
)

type PayrollRepo interface {
	SavePayroll(ctx context.Context, p *model.Payroll) error

	GetPayrollByEmployeeAndMonth(
		ctx context.Context,
		employeeID uint,
		monthYear time.Time,
	) (*model.Payroll, error)
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

func (r *payrollRepo) GetPayrollByEmployeeAndMonth(
	ctx context.Context,
	employeeID uint,
	monthYear time.Time,
) (*model.Payroll, error) {
	var payroll model.Payroll
	err := r.data.DB.WithContext(ctx).
		Where("employee_id = ? AND month_year = ?", employeeID, monthYear).
		First(&payroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payroll record not found for this employee and month")
		}
		return nil, fmt.Errorf("query payroll: %w", err)
	}
	return &payroll, nil
}
