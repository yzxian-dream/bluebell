package controller

import (
	"fmt"
	"webapp/dao/mysql"
	"webapp/dao/redis"
	"webapp/models"
	"webapp/pkg/snowflake"

	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	// 参数校验

	// 获取作者ID，当前请求的UserID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.AuthorId = userID
	// 生成帖子ID
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	post.PostID = postID
	// 创建帖子
	if err = mysql.CreatePost(&post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
	}
	//community, err := mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	//if err != nil {
	//	zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
	//	ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
	//}

	//if err != nil {
	//	zap.L().Error("logic.CreatePost failed", zap.Error(err))
	//	ResponseError(c, CodeServerBusy)
	//	return
	//}
	//记录帖子创建时间
	err = redis.CreatePost(int64(post.PostID))
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}
func PostDetailHandler(c *gin.Context) {
	postID := c.Param("id")

	//post, err := logic.GetPost(postId)
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed", zap.String("post_id", postID), zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.AuthorName = user.UserName
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.CommunityName = community.CommunityName

	ResponseSuccess(c, post)
}

func GetPostListHandler(c *gin.Context) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		pageNum int64
		size    int64
		err     error
	)
	pageNum, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}

	postList, err := mysql.GetPostList(pageNum, size)

	if err != nil {
		fmt.Println(err)
		return
	}
	data := make([]*models.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
			continue
		}
		post.AuthorName = user.UserName
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		post.CommunityName = community.CommunityName
		data = append(data, post)
	}
	ResponseSuccess(c, data)
}

// 升级版帖子接口
// 获取参数
// 去redis查询id列表
// 根据id去数据库查询帖子详细信息
func GetPostList2Handler(c *gin.Context) {
	type ParamPostList struct {
		Page  int64  `json:"page"`
		Size  int64  `json:"size"`
		Order string `json:"order"`
	}
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	//请求中的query params获取
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostList2Handler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	postList, err := mysql.GetPostList2(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	//提前查询好每篇帖子的投票数
	redis.GetPostVoteData()
	data := make([]*models.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
			continue
		}
		post.AuthorName = user.UserName
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		post.CommunityName = community.CommunityName
		data = append(data, post)
	}
	ResponseSuccess(c, data)
}
