package redis

import (
	"github.com/go-redis/redis"
	"webapp/models"
)

func GetPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求参数获取设置key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	startIndex := (p.Page - 1) * p.Size
	endIndex := startIndex + p.Size - 1
	//按分数从大到小的查询指定数量的查询
	return rdb.ZRevRange(key, startIndex, endIndex).Result()
}

// 根据Ids查询赞成票数
//func GetPostVoteData(ids []string) {
//	for _, id := range ids {
//		key := getRedisKey(KeyPostVotedZSetPrefix + id)
//		rdb.ZCount(key, "1", "1").Val()
//	}
//}

// 根据Ids查询赞成票数
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data := make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	//查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := rdb.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	//使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
