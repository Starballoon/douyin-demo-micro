package dal

import (
	"douyin-demo-micro/util"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	if DB != nil {
		return errors.New("multiple initialization is not allowed")
	}
	var err error
	DB, err = gorm.Open(mysql.Open(util.DSN),
		&gorm.Config{
			PrepareStmt: true,
		})
	if err != nil {
		return err
	}
	return InitTables()
}

func InitTables() error {
	var err error
	m := DB.Migrator()
	if !m.HasTable(&Comment{}) {
		err = m.CreateTable(&Comment{})
		if err != nil {
			return err
		}
	}
	return err
}
