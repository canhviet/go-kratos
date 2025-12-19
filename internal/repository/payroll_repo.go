package repository

import (
	"context"
	"errors"
	"myapp/internal/data"
	"myapp/internal/data/model"

	"gorm.io/gorm"
)

type payrollRepo struct {
	data *data.Data
}

func NewPayrollRepo(data *data.Data) *payrollRepo {
	return &payrollRepo{data: data}
}

func (r *payrollRepo) GetEmployeeByID(ctx context.Context, id uint) (*model.Employee, error) {
	var emp model.Employee
	if err := r.data.DB.WithContext(ctx).First(&emp, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return &emp, nil
}

func (r *payrollRepo) SavePayroll(ctx context.Context, p *model.Payroll) error {
	return r.data.DB.WithContext(ctx).Create(p).Error
}
