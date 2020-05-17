package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wcc4869/ginessential/common"
)

func main() {
	// 链接数据库
	db := common.InitDB()
	defer db.Close()

	r := gin.Default()

	r = CollectRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
