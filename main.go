package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/wcc4869/ginessential/common"
	"os"
)

func main() {
	InitConfig()
	// 链接数据库
	db := common.InitDB()
	defer db.Close()

	r := gin.Default()

	r = CollectRoute(r)
	port := viper.GetString("server.port")

	fmt.Println(port)

	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("app")               // 设置文件名
	viper.SetConfigType("yml")               // 设置文件格式
	viper.AddConfigPath(workDir + "/config") // 文件路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
