package controller

import (
	"strconv"
	"webapp/dao/redis"
	"webapp/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//投票功能
/*
/* PostVote 为帖子投票
投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

v=1时，有两种情况
	1.之前没投过票，现在要投赞成票
	2.之前投过反对票，现在要改为赞成票
v=0时，有两种情况
	1.之前投过赞成票，现在要取消
	2.之前投过反对票，现在要取消
v=-1时，有两种情况
	1.之前没投过票，现在要投反对票
	2.之前投过赞成票，现在要改为反对票
*/
//投票限制，每个帖子自发表之日起一个星期内允许用户投票，超过一个星期就不允许用户投票
// 1. 到期之后将redis中保存的赞成票书及反对票数存储到mysql中
// 2. 到期之后删除那个 keyPostVotedPrefix

func PostVote(c *gin.Context) {
	zap.L().Debug("ready to vote")
	p := new(models.VoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParams, errData)
		return
	}
	userId, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	zap.L().Debug(
		"voteForPost",
		zap.Int64("userId", int64(userId)),
		zap.String("postId", p.PostId),
		zap.Int("direction", p.Direction),
	)
	//1.判断投票限制
	if err = redis.VoteForPost(strconv.Itoa(int(userId)), p.PostId, float64(p.Direction)); err != nil {
		zap.L().Error("redis.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
