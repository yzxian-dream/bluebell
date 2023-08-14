package redis

// redis key注意使用命名空间方式
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
const (
	Prefix                 = "bluebell"
	KeyPostTimeZSet        = "post:time"   //帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //帖子及帖子分数
	KeyPostVotedZSetPrefix = "post:voted:" //记录用户及投票类型，参数是post id
)

func getRedisKey(key string) string {
	return Prefix + key
}
