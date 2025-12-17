package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"

	pb_user "myapp/api/user/v1"
	pb_payroll "myapp/api/payroll/v1"
	"myapp/internal/conf"
	"myapp/internal/service"
)

func NewHTTPServer(c *conf.Server, user *service.UserService, payroll *service.PayrollService) *http.Server {
	srv := http.NewServer(
		http.Address(c.Http.Addr),
	)

	pb_user.RegisterUserHTTPServer(srv, user)
	pb_payroll.RegisterPayrollHTTPServer(srv, payroll)
	
	return srv
}
