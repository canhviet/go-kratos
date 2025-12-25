// internal/biz/payroll.go
package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	v1 "myapp/api/payroll/v1"
	"myapp/internal/data/model"
	"myapp/internal/repository"
)

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrNoAttendanceThisMonth = errors.New("no attendance records found for this month")
)

type PayrollUsecase struct {
	payrollRepo   repository.PayrollRepo
	employeeRepo  repository.EmployeeRepo 
	timesheetRepo repository.TimesheetRepo
}

func NewPayrollUsecase(
	payrollRepo repository.PayrollRepo,
	employeeRepo repository.EmployeeRepo,
	timesheetRepo repository.TimesheetRepo,
) *PayrollUsecase {
	return &PayrollUsecase{
		payrollRepo:   payrollRepo,
		employeeRepo:  employeeRepo,
		timesheetRepo: timesheetRepo,
	}
}

func (uc *PayrollUsecase) CalculatePayroll(ctx context.Context, r *v1.CalculatePayrollRequest) (*v1.CalculatePayrollReply, error) {
	// Parse tháng năm
	monthYear, err := time.Parse("2006-01", r.MonthYear)
	if err != nil {
		return nil, errors.New("invalid month_year format, expected YYYY-MM")
	}

	// Lấy nhân viên
	emp, err := uc.employeeRepo.GetEmployeeByID(ctx, uint(r.EmployeeId))
	if err != nil {
		return nil, fmt.Errorf("get employee: %w", err)
	}

	// Tổng hợp chấm công theo tháng từ daily timesheet
	workingDays, overtimeHours, leaveDays, err := uc.timesheetRepo.GetMonthlySummary(
		ctx, uint(r.EmployeeId), monthYear.Year(), monthYear.Month())
	if err != nil {
		return nil, fmt.Errorf("get timesheet monthly summary: %w", err)
	}

	if workingDays+leaveDays == 0 {
		return nil, ErrNoAttendanceThisMonth
	}

	// Tính lương
	const standardWorkingDays = 26.0
	basicSalary := emp.BaseSalary * (float64(workingDays) / standardWorkingDays)
	hourlyRate := emp.BaseSalary / (standardWorkingDays * 8)
	overtimePay := overtimeHours * hourlyRate * 1.5
	grossSalary := basicSalary + overtimePay + r.Allowances

	// Khấu trừ BH (10.5%)
	insurance := grossSalary * 0.105

	// Thu nhập chịu thuế
	taxable := grossSalary - insurance - 11_000_000
	if emp.Dependents > 0 {
		taxable -= float64(emp.Dependents) * 4_400_000
	}

	incomeTax := calculateIncomeTax(taxable)
	totalDeductions := insurance + incomeTax
	netSalary := grossSalary - totalDeductions

	// Lưu bảng lương
	payroll := &model.Payroll{
		EmployeeID:    uint(r.EmployeeId),
		MonthYear:     monthYear,
		WorkingDays:   workingDays,
		OvertimeHours: overtimeHours,
		LeaveDays:     leaveDays,
		BasicSalary:   basicSalary,
		Allowances:    r.Allowances,
		GrossSalary:   grossSalary,
		Deductions:    totalDeductions,
		NetSalary:     netSalary,
		Status:        "calculated",
	}

	if err := uc.payrollRepo.SavePayroll(ctx, payroll); err != nil {
		return nil, fmt.Errorf("save payroll: %w", err)
	}

	return &v1.CalculatePayrollReply{
		GrossSalary:   grossSalary,
		NetSalary:     netSalary,
		Deductions:    totalDeductions,
		WorkingDays:   int32(workingDays),
		OvertimeHours: overtimeHours,
		LeaveDays:     int32(leaveDays),
	}, nil
}

func calculateIncomeTax(income float64) float64 {
	if income <= 0 {
		return 0
	}
	switch {
	case income <= 5_000_000:
		return income * 0.05
	case income <= 10_000_000:
		return 250_000 + (income-5_000_000)*0.10
	case income <= 18_000_000:
		return 750_000 + (income-10_000_000)*0.15
	case income <= 32_000_000:
		return 1_950_000 + (income-18_000_000)*0.20
	case income <= 52_000_000:
		return 4_750_000 + (income-32_000_000)*0.25
	case income <= 80_000_000:
		return 9_750_000 + (income-52_000_000)*0.30
	default:
		return 18_150_000 + (income-80_000_000)*0.35
	}
}