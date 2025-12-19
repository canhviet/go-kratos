package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"

	pb_user "myapp/api/user/v1"
	pb_payroll "myapp/api/payroll/v1"
	pb_employee "myapp/api/employee/v1"
	pb_timesheet "myapp/api/timesheet/v1"
	"myapp/internal/conf"
	"myapp/internal/service"
)

func NewHTTPServer(c *conf.Server, user *service.UserService, payroll *service.PayrollService, 
	employee *service.EmployeeService, timesheet *service.TimesheetService) *http.Server {
	srv := http.NewServer(
		http.Address(c.Http.Addr),
	)

	pb_user.RegisterUserHTTPServer(srv, user)
	pb_payroll.RegisterPayrollHTTPServer(srv, payroll)
	pb_employee.RegisterEmployeeHTTPServer(srv, employee)
	pb_timesheet.RegisterTimesheetHTTPServer(srv, timesheet)
	
	return srv
}
