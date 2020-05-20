package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/wcc4869/ginessential/model"
)

var DB *gorm.DB

// 初始化数据库
func InitDB() *gorm.DB {
	driverName := viper.GetString("database.drivername")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	charset := viper.GetString("database.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database,err :" + err.Error())
	}

	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
