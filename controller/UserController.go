package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wcc4869/ginessential/common"
	"github.com/wcc4869/ginessential/dto"
	"github.com/wcc4869/ginessential/model"
	"github.com/wcc4869/ginessential/response"
	"github.com/wcc4869/ginessential/util"
	"golang.org/x/crypto/bcrypt"
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
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须 11 位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能小于 6 位")
		return
	}
	if len(name) == 0 {
		name = util.GetRandomString(10)
	}
	log.Println(name, phone, password)
	// 判断手机号是否存在
	if isPhoneExist(DB, phone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该用户已存在")
		return
	}
	// 创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Phone:    phone,
		Password: string(hashedPassword),
	}
	DB.Create(&newUser)
	// 返回结果
	response.Success(c, nil, "注册成功")
}

func Login(c *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	// 数据验证
	if len(phone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须 11 位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能小于 6 位")
		return
	}
	// 判断手机号存在
	var user model.User
	DB.Where("phone=?", phone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	// 判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Error(c, nil, "密码错误")
		return
	}

	// 发放 token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		return
	}

	// 返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")

}

// 用户信息
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	// 这里 DTO操作，只返回 name+phone
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User)),}, "")

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
