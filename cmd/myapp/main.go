package main

import (
	"flag"
	"fmt"
	"myapp/internal/biz"
	"myapp/internal/conf"
	"myapp/internal/data"
	repo "myapp/internal/repository"
	"myapp/internal/server"
	"myapp/internal/service"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "./configs/config.yaml", "config file")
	flag.Parse()

	// load config
	c := config.New(
		config.WithSource(file.NewSource(confPath)),
	)	
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(fmt.Errorf("load config error: %w", err))
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(fmt.Errorf("scan config error: %w", err))
	}

	// logger
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	// database
	db, err := data.NewDB(bc.Data)
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	d, cleanup, err := data.NewData(db, logger)
	if err != nil {
		panic(fmt.Errorf("init data error: %w", err))
	}
	defer cleanup()

	// business layer
	UserUsecase := biz.NewUserUsecase(repo.NewUserRepo(d))
	EmployeeUsecase := biz.NewEmployeeUsecase(repo.NewEmployeeRepo(d))
	PayrollUsecase := biz.NewPayrollUsecase(repo.NewPayrollRepo(d))
	TimesheetUsecase := biz.NewTimesheetUsecase(repo.NewTimesheetRepo(d))


	// service layer
	UserService := service.NewUserService(UserUsecase)
	EmployeeService := service.NewEmployeeService(EmployeeUsecase)
	PayrollService := service.NewPayrollService(PayrollUsecase)
	TimeSheetService := service.NewTimesheetService(TimesheetUsecase)

	// HTTP server
	httpServer := server.NewHTTPServer(bc.Server, UserService, PayrollService, EmployeeService, TimeSheetService)

	// kratos app
	app := kratos.New(
		kratos.Name("myapp"),
		kratos.Server(httpServer),
		kratos.Logger(logger),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
