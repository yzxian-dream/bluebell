package redis

import "webapp/models"

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
func GetPostVoteData(ids []string) {
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		rdb.ZCount(key, "1", "1").Val()
	}
}
