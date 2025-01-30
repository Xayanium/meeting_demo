package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"meeting_demo/internal/models"
	"meeting_demo/utils"
	"net/http"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserLogin(c *gin.Context) {
	// 解析request消息体
	userForm := new(UserLoginRequest)
	err := c.ShouldBindJSON(&userForm)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "System Server Error",
		})
		log.Println("bind user info error: ", err)
		return
	}

	// 前端数据校验，md5加密
	if userForm.Username == "" || userForm.Password == "" {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "Auth Info Is Empty",
		})
		return
	}
	userForm.Password = utils.GetMd5(userForm.Password)

	// 验证用户
	tableData := new(models.UserInfo)
	err = models.DB.Model(&models.UserInfo{}).
		Where("username = ? AND password = ?", userForm.Username, userForm.Password).First(&tableData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, &gin.H{
				"code": -1,
				"msg":  "Auth User Failed",
			})
			return
		}
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "System Server Error",
		})
		log.Println("query userinfo data error: ", err.Error())
		return
	}

	// 生成token
	token, err := utils.GenerateToken(tableData.ID, tableData.Username)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "System Server Error",
		})
		log.Println("generate token error: ", err)
		return
	}

	// 返回token
	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
