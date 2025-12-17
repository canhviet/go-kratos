package biz

import (
	"context"
	"errors"
	"time"

	v1 "myapp/api/payroll/v1"
	"myapp/internal/data/model"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type PayrollRepo interface {
    GetEmployeeByID(ctx context.Context, id uint) (*model.Employee, error)
	SavePayroll(ctx context.Context, p *model.Payroll) error
}

type PayrollUsecase struct {
	repo PayrollRepo
}

func NewPayrollUsecase(repo PayrollRepo) *PayrollUsecase {
	return &PayrollUsecase{repo: repo}
}

func (uc *PayrollUsecase) CalculatePayroll(ctx context.Context, r *v1.CalculatePayrollRequest) (*v1.CalculatePayrollReply, error) {
	emp, err := uc.repo.GetEmployeeByID(ctx, uint(r.EmployeeId))
	if err != nil {
		return nil, err
	}

	standardDays := 26.0
	basic := emp.BaseSalary * (float64(r.WorkingDays) / standardDays)
	overtimePay := r.OvertimeHours * (emp.BaseSalary / (standardDays * 8)) * 1.5
	gross := basic + overtimePay + r.Allowances

	//deductions
	insuranceDeduction := gross * 0.105
	taxableIncome := gross - insuranceDeduction - 11000000 - float64(emp.Dependents)*4400000
	tax := calculateTax(taxableIncome)
	deductions := insuranceDeduction + tax
	net := gross - deductions

	monthYear, err := time.Parse("2006-01", r.MonthYear)
	if err != nil {
		return nil, errors.New("invalid month_year format")
	}

	payroll := &model.Payroll{
		EmployeeID:  uint(r.EmployeeId),
		MonthYear:   monthYear,
		BasicSalary: basic,
		Allowances:  r.Allowances,
		Deductions:  deductions,
		NetSalary:   net,
		Status:      "paid",  
	}
	if err := uc.repo.SavePayroll(ctx, payroll); err != nil {
		return nil, err
	}

	return &v1.CalculatePayrollReply{
		GrossSalary: gross,
		NetSalary:   net,
		Deductions:  deductions,
	}, nil
}

func calculateTax(income float64) float64 {
	if income <= 0 {
		return 0
	} else if income <= 5000000 {
		return income * 0.05
	} else if income <= 10000000 {
		return 250000 + (income-5000000)*0.10
	} else if income <= 18000000 {
		return 750000 + (income-10000000)*0.15
	} else if income <= 32000000 {
		return 1950000 + (income-18000000)*0.20
	} else if income <= 52000000 {
		return 4750000 + (income-32000000)*0.25
	} else if income <= 80000000 {
		return 9750000 + (income-52000000)*0.30
	}
	return 18150000 + (income-80000000)*0.35
}