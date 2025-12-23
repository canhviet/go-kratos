package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	authv1     "myapp/api/auth/v1"
	employeev1 "myapp/api/employee/v1"
	payrollv1  "myapp/api/payroll/v1"
	timesheetv1 "myapp/api/timesheet/v1"

	"myapp/internal/biz"
	"myapp/internal/conf"
	"myapp/internal/data"
	"myapp/internal/repository"
	"myapp/internal/server"
	"myapp/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/redis/go-redis/v9"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "conf", "./configs/config.yaml", "config file path")
	flag.Parse()

	// Load configuration
	c := config.New(config.WithSource(file.NewSource(confPath)))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(fmt.Errorf("failed to scan config: %w", err))
	}

	// Logger
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)

	// Initialize database
	db, err := data.NewDB(bc.Data)
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	d, cleanup, err := data.NewData(db, logger)
	if err != nil {
		panic(fmt.Errorf("failed to initialize data layer: %w", err))
	}
	defer cleanup()

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     bc.Data.Redis.Addr,
		Password: bc.Data.Redis.Password,
		DB:       int(bc.Data.Redis.GetDb()),
	})
	redisRepo := repository.NewRedisRepo(redisClient)

	// Repositories
	employeeRepo := repository.NewEmployeeRepo(d)
	payrollRepo := repository.NewPayrollRepo(d)
	timesheetRepo := repository.NewTimesheetRepo(d)
	userRepo := repository.NewUserRepo(d)

	// Usecases (Biz layer)
	employeeUsecase := biz.NewEmployeeUsecase(employeeRepo)
	payrollUsecase := biz.NewPayrollUsecase(payrollRepo)
	timesheetUsecase := biz.NewTimesheetUsecase(timesheetRepo)
	authUsecase := biz.NewAuthUsecase(
		userRepo,
		redisRepo,
		bc.Auth.GetJwtSecret(),
		int(bc.Auth.GetTokenExp()),
	)

	// Services
	employeeService := service.NewEmployeeService(employeeUsecase)
	payrollService := service.NewPayrollService(payrollUsecase)
	timesheetService := service.NewTimesheetService(timesheetUsecase)
	authService := service.NewAuthService(authUsecase)

	httpSrv := http.NewServer(
		http.Address(bc.Server.Http.Addr),
		http.Timeout(time.Duration(bc.Server.Http.Timeout)*time.Second),
		http.Middleware(
			recovery.Recovery(),
			server.AuthMiddleware(bc.Auth.GetJwtSecret(), redisRepo),
		),
	)

	authv1.RegisterAuthHTTPServer(httpSrv, authService)
	employeev1.RegisterEmployeeHTTPServer(httpSrv, employeeService)
	payrollv1.RegisterPayrollHTTPServer(httpSrv, payrollService)
	timesheetv1.RegisterTimesheetHTTPServer(httpSrv, timesheetService)

	// Kratos application
	app := kratos.New(
		kratos.Name("myapp"),
		kratos.Version("v1.0.0"),
		kratos.Server(httpSrv),
		kratos.Logger(logger),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}