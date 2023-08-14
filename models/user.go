package models

// 定义用户结构体
type User struct {
	UserID   uint64 `db:"user_id" `
	UserName string `db:"username" `
	Password string `db:"password" `
}
