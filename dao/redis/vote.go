package redis

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/go-redis/redis"
)

/*
投票分为四种情况：1.投赞成票 2.投反对票 3.取消投票 4.反转投票

记录文章参与投票的人
更新文章分数：赞成票要加分；反对票减分

direction=1时，有两种情况
	1.之前没投过票，现在要投赞成票       差值的绝对值为1 +432
	2.之前投过反对票，现在要改为赞成票    差值的绝对值为2 +2*432
direction=0时，有两种情况
	1.之前投过赞成票，现在要取消
	2.之前投过反对票，现在要取消
direction=-1时，有两种情况
	1.之前没投过票，现在要投反对票
	2.之前投过赞成票，现在要改为反对票
*/

const oneWeekInSecond = 7 * 24 * 3600

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeat     = errors.New("不允许重复投票")
)

func CreatePost(postId int64) error {
	fmt.Println(rdb)
	_, err := rdb.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	}).Result()
	return err
}

func VoteForPost(userId, postId string, direction float64) error {
	//1.判断投票限制
	//从redis取帖子发布时间
	postTime, err := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postId).Result()
	if err != nil {
		return err
	}
	if float64(time.Now().Unix())-postTime > oneWeekInSecond {
		return ErrVoteTimeExpire
	}
	// 2.更新帖子的分数
	//先查看当前贴子的投票记录
	olddirection := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postId), userId).Val()
	//禁止重复投票
	if direction == olddirection {
		return ErrVoteRepeat
	}
	var dir float64
	if direction > olddirection {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(olddirection - direction)
	_, err = rdb.ZIncrBy(getRedisKey(KeyPostScoreZSet), 432*diff*dir, postId).Result()
	if err != nil {
		return err
	}
	//3. 记录用户为该帖子投票的数据
	if direction == 0 {
		rdb.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postId), userId).Result()
	} else {
		rdb.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postId), redis.Z{
			Score:  direction,
			Member: userId,
		}).Result()
	}
	return err

}
