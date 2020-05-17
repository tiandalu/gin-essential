package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wcc4869/ginessential/common"
	"github.com/wcc4869/ginessential/model"
	"github.com/wcc4869/ginessential/util"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.GetDB()
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
		name = util.GetRandomString(10)
	}
	log.Println(name, phone, password)
	// 判断手机号是否存在
	if isPhoneExist(DB, phone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "该用户已存在",
		})
		return

	}
	// 创建用户
	newUser := model.User{
		Name:     name,
		Phone:    phone,
		Password: password,
	}
	DB.Create(&newUser)
	// 返回结果
	c.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

// 判断手机号是否存在数据库
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("phone=?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
