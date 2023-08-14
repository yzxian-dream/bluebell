package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"webapp/models"
	"webapp/pkg/snowflake"
)

// 这一层主要把每一步数据库操作封装成函数
// 待logic层根据业务需求调用
const sercret = "eric_yao"

func Register(user *models.User) (err error) {
	//直接包到register里面做
	//判断有无重复,如果存在就没必要往下走
	sqlStr := "select count(user_id) from user where username = ?"

	var count int64

	if err := db.Get(&count, sqlStr, user.UserName); err != nil {
		return err
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		return ErrorUserExit
	}
	//生成uid

	userID, err := snowflake.GetID()
	//密码加密，保存进数据库
	u := &models.User{
		UserName: user.UserName,
		Password: user.Password,
		UserID:   userID,
	}
	InsertUser(u)
	return
}
func Login(user *models.User) (err error) {
	//判断用户输入密码与数据库中数据是否一致
	inputPassword := user.Password
	sqlStr := "select user_id, username, password from user where username = ?"

	err = db.Get(user, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	password := encyptPassword(inputPassword)
	if password == user.Password {
		return ErrorPasswordWrong
	}
	return
}

func InsertUser(user *models.User) (err error) {
	//user.Password = encyptPassword(user.Password)
	sqlStr := "insert into user(user_id, username,password) value(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	if err != nil {
		return err
	}
	return
}

func encyptPassword(opassword string) string {
	h := md5.New()
	h.Write([]byte(sercret))
	return hex.EncodeToString(h.Sum([]byte(opassword)))

}
func GetUserByID(idStr string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, idStr)
	return
}
