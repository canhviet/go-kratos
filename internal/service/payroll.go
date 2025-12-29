package service

import (
	"context"
	"fmt"

	v1 "myapp/api/payroll/v1"
	"myapp/internal/biz"

	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *PayrollService) ExportPayrollPDF(ctx context.Context, req *v1.ExportPayrollPDFRequest) (*v1.ExportPayrollPDFReply, error) {
	pdfData, err := s.uc.ExportPayrollPDF(ctx, req.EmployeeId, req.MonthYear)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot generate PDF: %v", err)
	}

	hctx, ok := ctx.(http.Context)
	if !ok {
		return &v1.ExportPayrollPDFReply{
			PdfData:  pdfData,
			Filename: fmt.Sprintf("payslip_%d_%s.pdf", req.EmployeeId, req.MonthYear),
		}, nil
	}

	w := hctx.Response()
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="payslip_%d_%s.pdf"`, req.EmployeeId, req.MonthYear))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfData)))

	if _, err := w.Write(pdfData); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write PDF")
	}

	return &v1.ExportPayrollPDFReply{}, nil
}

func (s *PayrollService) SendPayslipEmail(ctx context.Context, req *v1.SendPayslipEmailRequest) (*v1.SendPayslipEmailReply, error) {
	err := s.uc.SendPayslipEmail(ctx, req.EmployeeId, req.MonthYear, req.ToEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "send email failed: %v", err)
	}

	return &v1.SendPayslipEmailReply{
		Message: "Payslip sent successfully via email",
	}, nil
}