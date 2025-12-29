// internal/biz/payroll.go
package biz

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	v1 "myapp/api/payroll/v1"
	"myapp/internal/data/model"
	"myapp/internal/repository"

	"github.com/jung-kurt/gofpdf"
)

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrNoAttendanceThisMonth = errors.New("no attendance records found for this month")
)

type PayrollUsecase struct {
	payrollRepo   repository.PayrollRepo
	employeeRepo  repository.EmployeeRepo 
	timesheetRepo repository.TimesheetRepo
	emailRepo repository.EmailRepo
}

func NewPayrollUsecase(
	payrollRepo repository.PayrollRepo,
	employeeRepo repository.EmployeeRepo,
	timesheetRepo repository.TimesheetRepo,
	emailRepo repository.EmailRepo,
) *PayrollUsecase {
	return &PayrollUsecase{
		payrollRepo:   payrollRepo,
		employeeRepo:  employeeRepo,
		timesheetRepo: timesheetRepo,
		emailRepo: emailRepo,
	}
}

func (uc *PayrollUsecase) CalculatePayroll(ctx context.Context, r *v1.CalculatePayrollRequest) (*v1.CalculatePayrollReply, error) {
	
	monthYear, err := time.Parse("2006-01", r.MonthYear)
	if err != nil {
		return nil, errors.New("invalid month_year format, expected YYYY-MM")
	}

	emp, err := uc.employeeRepo.GetEmployeeByID(ctx, uint(r.EmployeeId))
	if err != nil {
		return nil, fmt.Errorf("get employee: %w", err)
	}

	workingDays, overtimeHours, leaveDays, err := uc.timesheetRepo.GetMonthlySummary(
		ctx, uint(r.EmployeeId), monthYear.Year(), monthYear.Month())
	if err != nil {
		return nil, fmt.Errorf("get timesheet monthly summary: %w", err)
	}

	if workingDays+leaveDays == 0 {
		return nil, ErrNoAttendanceThisMonth
	}

	const standardWorkingDays = 26.0
	basicSalary := emp.BaseSalary * (float64(workingDays) / standardWorkingDays)
	hourlyRate := emp.BaseSalary / (standardWorkingDays * 8)
	overtimePay := overtimeHours * hourlyRate * 1.5
	grossSalary := basicSalary + overtimePay + r.Allowances

	insurance := grossSalary * 0.105

	taxable := grossSalary - insurance - 11_000_000
	if emp.Dependents > 0 {
		taxable -= float64(emp.Dependents) * 4_400_000
	}

	incomeTax := calculateIncomeTax(taxable)
	totalDeductions := insurance + incomeTax
	netSalary := grossSalary - totalDeductions

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

func (uc *PayrollUsecase) ExportPayrollPDF(ctx context.Context, employeeID uint32, monthYearStr string) ([]byte, error) {
	monthYear, err := time.Parse("2006-01", monthYearStr)
	if err != nil {
		return nil, errors.New("invalid month_year format, expected YYYY-MM")
	}

	payroll, err := uc.payrollRepo.GetPayrollByEmployeeAndMonth(ctx, uint(employeeID), monthYear)
	if err != nil {
		return nil, fmt.Errorf("get payroll record: %w", err)
	}

	emp, err := uc.employeeRepo.GetEmployeeByID(ctx, uint(employeeID))
	if err != nil {
		return nil, fmt.Errorf("get employee info: %w", err)
	}

	const standardWorkingDays = 26.0
	hourlyRate := emp.BaseSalary / (standardWorkingDays * 8)
	overtimePay := payroll.OvertimeHours * hourlyRate * 1.5
	basicAndAllowances := payroll.GrossSalary - overtimePay

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(277, 20, fmt.Sprintf("PAYSLIP - %s", monthYear.Format("January 2006")), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(60, 10, "Full Name:", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 10, emp.Name, "", 1, "L", false, 0, "")

	pdf.CellFormat(60, 10, "Employee ID:", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 10, fmt.Sprintf("%d", emp.ID), "", 1, "L", false, 0, "")

	pdf.CellFormat(60, 10, "Position:", "", 0, "L", false, 0, "")
	pdf.CellFormat(100, 10, emp.Position, "", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFillColor(230, 230, 250) 
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(120, 12, "Description", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 12, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(87, 12, "Amount (VND)", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 13)
	pdf.SetFillColor(255, 255, 255)

	pdf.CellFormat(120, 12, "Basic Salary + Allowances", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 12, fmt.Sprintf("%d working days", payroll.WorkingDays), "1", 0, "C", false, 0, "")
	pdf.CellFormat(87, 12, formatCurrency(basicAndAllowances), "1", 1, "R", false, 0, "")

	pdf.CellFormat(120, 12, "Overtime Pay", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 12, fmt.Sprintf("%.1f hours", payroll.OvertimeHours), "1", 0, "C", false, 0, "")
	pdf.CellFormat(87, 12, formatCurrency(overtimePay), "1", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 15)
	pdf.SetFillColor(220, 240, 255)
	pdf.CellFormat(190, 15, "GROSS SALARY", "1", 0, "R", true, 0, "")
	pdf.CellFormat(87, 15, formatCurrency(payroll.GrossSalary), "1", 1, "R", true, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "", 13)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(190, 12, "Deductions", "", 0, "R", false, 0, "")
	pdf.CellFormat(87, 12, formatCurrency(payroll.Deductions), "", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 18)
	pdf.SetFillColor(200, 255, 200)
	pdf.CellFormat(190, 20, "NET SALARY", "1", 0, "R", true, 0, "")
	pdf.CellFormat(87, 20, formatCurrency(payroll.NetSalary), "1", 1, "R", true, 0, "")

	pdf.Ln(5)

	// Footer
	pdf.SetFont("Arial", "I", 11)
	pdf.CellFormat(277, 10, fmt.Sprintf("Generated on: %s", time.Now().Format("02 January 2006, 15:04")), "", 1, "R", false, 0, "")

	// Output PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generate PDF error: %w", err)
	}

	return buf.Bytes(), nil
}

func formatCurrency(amount float64) string {
	amount = math.Round(amount)

	str := strconv.FormatFloat(amount, 'f', 0, 64)

	var result strings.Builder
	length := len(str)
	for i, char := range str {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}

	return result.String() + " VND"
}

func (uc *PayrollUsecase) SendPayslipEmail(ctx context.Context, employeeID uint32, monthYearStr, toEmail string) error {
	pdfData, err := uc.ExportPayrollPDF(ctx, employeeID, monthYearStr)
	if err != nil {
		return fmt.Errorf("generate PDF: %w", err)
	}

	emp, err := uc.employeeRepo.GetEmployeeByID(ctx, uint(employeeID))
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}

	return uc.emailRepo.SendPayslip(ctx, toEmail, emp.Name, monthYearStr, pdfData)
}