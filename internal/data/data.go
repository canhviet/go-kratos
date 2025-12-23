package data

import (
	"myapp/internal/conf"
	"myapp/internal/data/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Data struct {
	DB *gorm.DB
}

func NewData(db *gorm.DB, logger interface{}) (*Data, func(), error) {
	cleanup := func() {}
	return &Data{DB: db}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	dsn := "root:root@tcp(mysql:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Timesheet{})
	db.AutoMigrate(&model.Employee{})
	db.AutoMigrate(&model.Payroll{})
	db.AutoMigrate(&model.User{})

	return db, nil
}
