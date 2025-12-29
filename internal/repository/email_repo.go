package repository

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"gopkg.in/gomail.v2"
)

type EmailRepo interface {
	SendPayslip(ctx context.Context, toEmail, employeeName, monthYear string, pdfData []byte) error
}

type emailRepo struct {
	dialer *gomail.Dialer
	from   string
	name   string
}

func NewEmailRepo(host string, port int, username, password, fromEmail, fromName string) EmailRepo {
	d := gomail.NewDialer(host, port, username, password)

	// BỎ QUA KIỂM TRA CERTIFICATE - chỉ dùng cho dev/test
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &emailRepo{
		dialer: d,
		from:   fromEmail,
		name:   fromName,
	}
}

func (r *emailRepo) SendPayslip(ctx context.Context, toEmail, employeeName, monthYear string, pdfData []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(r.from, r.name))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", fmt.Sprintf("Payslip - %s %s", employeeName, monthYear))

	body := fmt.Sprintf(`Dear %s,

Please find your payslip for %s attached to this email.

If you have any questions, please contact the HR department.

Best regards,
HR Team
`, employeeName, monthYear)

	m.SetBody("text/plain", body)

	m.Attach("payslip.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfData)
		return err
	}))

	return r.dialer.DialAndSend(m)
}