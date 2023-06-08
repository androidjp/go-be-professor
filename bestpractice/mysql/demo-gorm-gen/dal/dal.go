package dal

import (
	"demo-gorm-gen/dal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var DB *gorm.DB
var once sync.Once

func init() {
	once.Do(func() {
		DB = ConnectDB().Debug()
		_ = DB.AutoMigrate(&model.User{}, &model.Passport{})
	})
}

func ConnectDB() (conn *gorm.DB) {
	// 连接mysql8.0 docker环境
	conn, err := gorm.Open(mysql.Open("user01:Password1!@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	return conn
}
