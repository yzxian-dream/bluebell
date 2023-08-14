package models

import "time"

// 定义用户结构体
type Community struct {
	ID    int64  `json:"community_id" db:"community_id"`
	Names string `json:"name" db:"community_name" `
}

type CommunityDetail struct {
	CommunityID   uint64    `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}
