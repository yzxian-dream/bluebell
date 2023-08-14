package models

type VoteData struct {
	PostId    string `json:"postId" binding:"required"`
	Direction int    `json:"direction" binding:"oneof=1 0 -1"` //赞成票1，反对票-1，取消投票0
}
