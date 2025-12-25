package biz

import (
	"context"
	"time"

	"myapp/internal/data/model"
	"myapp/internal/repository"
)

type EmployeeUsecase struct {
	repo repository.EmployeeRepo
}

func NewEmployeeUsecase(repo repository.EmployeeRepo) *EmployeeUsecase {
	return &EmployeeUsecase{repo: repo}
}

func (uc *EmployeeUsecase) List(ctx context.Context, pageSize int, pageToken string) ([]*model.Employee, string, error) {
	return uc.repo.List(ctx, pageSize, pageToken)
}

func (uc *EmployeeUsecase) Get(ctx context.Context, id uint32) (*model.Employee, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *EmployeeUsecase) Create(ctx context.Context, name string, position string, baseSalary float64, bankAccount string, joinDate time.Time, dependents int) (*model.Employee, error) {
	employee := &model.Employee{
		Name:        name,
		Position:    position,
		BaseSalary:  baseSalary,
		BankAccount: bankAccount,
		JoinDate:    joinDate,
		Dependents:  dependents,
	}
	err := uc.repo.Create(ctx, employee)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (uc *EmployeeUsecase) Update(ctx context.Context, id uint32, name string, position string, baseSalary float64, bankAccount string, joinDate time.Time, dependents int) (*model.Employee, error) {
	employee, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	employee.Name = name
	employee.Position = position
	employee.BaseSalary = baseSalary
	employee.BankAccount = bankAccount
	employee.JoinDate = joinDate
	employee.Dependents = dependents
	err = uc.repo.Update(ctx, employee)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (uc *EmployeeUsecase) Delete(ctx context.Context, id uint32) error {
	return uc.repo.Delete(ctx, id)
}