package repository

import (
	"context"
	"errors"
	"strconv"

	"myapp/internal/data"
	"myapp/internal/data/model"

	"gorm.io/gorm"
)

type employeeRepo struct {
	data *data.Data
}

type EmployeeRepo interface {
	List(ctx context.Context, pageSize int, pageToken string) ([]*model.Employee, string, error)
	Get(ctx context.Context, id uint32) (*model.Employee, error)
	Create(ctx context.Context, employee *model.Employee) error
	Update(ctx context.Context, employee *model.Employee) error
	Delete(ctx context.Context, id uint32) error
	GetEmployeeByID(ctx context.Context, id uint) (*model.Employee, error)
}

func NewEmployeeRepo(data *data.Data) *employeeRepo {
	return &employeeRepo{data: data}
}

func (r *employeeRepo) GetEmployeeByID(ctx context.Context, id uint) (*model.Employee, error) {
	var emp model.Employee
	if err := r.data.DB.WithContext(ctx).First(&emp, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepo) List(ctx context.Context, pageSize int, pageToken string) ([]*model.Employee, string, error) {
	var offset int
	if pageToken != "" {
		offset, _ = strconv.Atoi(pageToken) // Simple error handling; improve in production
	}

	var employees []*model.Employee
	err := r.data.DB.WithContext(ctx).Limit(pageSize).Offset(offset).Find(&employees).Error
	if err != nil {
		return nil, "", err
	}

	nextToken := ""
	if len(employees) == pageSize {
		nextToken = strconv.Itoa(offset + pageSize)
	}

	return employees, nextToken, nil
}

func (r *employeeRepo) Get(ctx context.Context, id uint32) (*model.Employee, error) {
	var employee model.Employee
	err := r.data.DB.WithContext(ctx).First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepo) Create(ctx context.Context, employee *model.Employee) error {
	return r.data.DB.WithContext(ctx).Create(employee).Error
}

func (r *employeeRepo) Update(ctx context.Context, employee *model.Employee) error {
	return r.data.DB.WithContext(ctx).Save(employee).Error
}

func (r *employeeRepo) Delete(ctx context.Context, id uint32) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Employee{}, id).Error
}
