package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
)

type User struct {
	gorm.Model
	Name     string `gorm :"type:varchar(20);not null"`
	Phone    string `gorm :"type:varchar(110);not null;unique"`
	Password string `gorm:"size:255;not null"`
}

func main() {
	// 链接数据库
	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
		// 获取参数
		name := c.PostForm("name")
		phone := c.PostForm("phone")
		password := c.PostForm("password")

		// 数据验证
		if len(phone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号必须 11 位",
			})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能小于 6 位",
			})
			return
		}
		if len(name) == 0 {
			name = GetRandomString(10)
		}
		log.Println(name, phone, password)
		// 判断手机号是否存在
		if isPhoneExist(db, phone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "该用户已存在",
			})
			return

		}
		// 创建用户
		newUser := User{
			Name:     name,
			Phone:    phone,
			Password: password,
		}
		db.Create(&newUser)
		// 返回结果
		c.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 判断手机号是否存在数据库
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user User
	db.Where("phone=?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

var (
	lowercaseRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	uppercaseRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lettersRunes   = append(lowercaseRunes, uppercaseRunes...)
	digitsRunes    = []rune("0123456789")
	allRunes       = append(lettersRunes, digitsRunes...)
)

// 获取随机字符串，n 是长度
func GetRandomString(n int) string {
	b := make([]rune, n)
	b[0] = lowercaseRunes[rand.Intn(len(lowercaseRunes))]
	for i := 1; i < n; i++ {
		b[i] = allRunes[rand.Intn(len(allRunes))]
	}
	return string(b)
}

// 初始化数据库
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database,err :" + err.Error())
	}

	db.AutoMigrate(&User{})
	return db
}
