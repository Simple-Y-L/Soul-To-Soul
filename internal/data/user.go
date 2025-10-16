package data

import (
	"kratosItems/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Mobile string
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.User {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) GetUserMobileInfo(db *gorm.DB, req *biz.UsersTable) (*biz.UsersTable, error) {
	var user biz.UsersTable
	user.UserId = req.UserId
	user.PhoneNumber = req.PhoneNumber
	var userInfo biz.UsersTable
	err := db.Where(&user).Limit(1).Find(&userInfo).Error
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}
