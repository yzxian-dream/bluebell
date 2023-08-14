package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    MyCode      `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseErrorWithMsg(ctx *gin.Context, code MyCode, msg interface{}) {
	response := &ResponseData{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
	ctx.JSON(http.StatusBadRequest, response)
}

func ResponseError(ctx *gin.Context, code MyCode) {
	response := &ResponseData{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusInternalServerError, response)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	response := &ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, response)
}
