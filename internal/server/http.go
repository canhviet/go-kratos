package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"

	pb_payroll "myapp/api/payroll/v1"
	pb_employee "myapp/api/employee/v1"
	pb_timesheet "myapp/api/timesheet/v1"
	"myapp/internal/conf"
	"myapp/internal/service"
	"myapp/internal/repository"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
)

func NewHTTPServer(c *conf.Server, auth *conf.Auth, payroll *service.PayrollService, 
	employee *service.EmployeeService, timesheet *service.TimesheetService,
	redisRepo *repository.RedisRepo) *http.Server {
	srv := http.NewServer(
		http.Address(c.Http.Addr),
		http.Middleware(
			recovery.Recovery(),
			AuthMiddleware(auth.JwtSecret, redisRepo),
		),
	)

	pb_payroll.RegisterPayrollHTTPServer(srv, payroll)
	pb_employee.RegisterEmployeeHTTPServer(srv, employee)
	pb_timesheet.RegisterTimesheetHTTPServer(srv, timesheet)
	
	return srv
}