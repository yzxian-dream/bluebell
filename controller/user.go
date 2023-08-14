package controller

import (
	"errors"
	"fmt"
	"strconv"
	"webapp/dao/mysql"
	"webapp/models"
	"webapp/pkg/jwt"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	//1.拿到参数，参数校验
	var p *models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("请求参数错误", zap.Error(err))
		//判断返回的err是不是validator的类型
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	fmt.Println(p)
	//2.业务处理

	err := mysql.Register(&models.User{
		UserName: p.UserName,
		Password: p.Password,
	})
	if errors.Is(err, mysql.ErrorUserExit) {
		ResponseError(c, CodeUserExist)
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1.拿到参数，参数校验
	var u *models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		zap.L().Error("请求参数错误", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	fmt.Println(u)
	//2.业务处理

	err := mysql.Login(u)
	if errors.Is(err, mysql.ErrorPasswordWrong) {
		ResponseError(c, CodeInvalidPassword)
	}
	//登录账号正确，生成token,放进response
	aToken, rToken, err := jwt.GenToken(u.UserID)
	if err != nil {
		zap.L().Error("生成token失败", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidToken, err.Error())
		return
	}
	//3.返回响应(将返回响应包装成函数，减少重复代码)
	ResponseSuccess(c, gin.H{
		"accessToken":  aToken,
		"refreshToken": rToken,
		"userID":       fmt.Sprintf("%d", u.UserID),
		"username":     u.UserName,
	})
}

func CommunityHandler(c *gin.Context) {
	//2.业务处理
	data, err := mysql.GetCommunityList()
	if err != nil {
		zap.L().Error("mysql.GetCommunityList failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	//3.返回响应(将返回响应包装成函数，减少重复代码)
	ResponseSuccess(c, data)
}

// 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//1.获取社区id
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	//2.业务处理
	data, err := mysql.GetCommunityByID(id)
	if err != nil {
		zap.L().Error("mysql.GetCommunityList failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}

	//3.返回响应(将返回响应包装成函数，减少重复代码)
	ResponseSuccess(c, data)
}
