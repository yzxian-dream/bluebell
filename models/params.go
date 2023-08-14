package models

// 定义请求参数
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamPostList struct {
	Page  int64  `json:"page"`
	Size  int64  `json:"size"`
	Order string `json:"order"`
}

const (
	OrderTime  = "time"
	OrderScore = "score"
)
