package biz

import (
	"context"
	"errors"
	pb "kratosItems/api/user"
	"kratosItems/cmd/basic/config"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type UsersTable struct {
	UserId       int64     `gorm:"column:user_id;type:bigint;comment:用户唯一ID;primaryKey;not null;" json:"user_id"`                  // 用户唯一ID
	PhoneNumber  string    `gorm:"column:phone_number;type:varchar(20);comment:手机号;not null;" json:"phone_number"`                 // 手机号
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);comment:加密密码;not null;" json:"password_hash"`             // 加密密码
	AvatarUrl    string    `gorm:"column:avatar_url;type:varchar(500);comment:虚拟头像URL;default:NULL;" json:"avatar_url"`            // 虚拟头像URL
	Nickname     string    `gorm:"column:nickname;type:varchar(100);comment:用户昵称;not null;" json:"nickname"`                       // 用户昵称
	Gender       string    `gorm:"column:gender;type:enum('male', 'female', 'unknown');comment:性别;default:unknown;" json:"gender"` // 性别
	Birthday     time.Time `gorm:"column:birthday;type:date;comment:生日;default:NULL;" json:"birthday"`                             // 生日
	PlanetId     int32     `gorm:"column:planet_id;type:int;comment:所属星球ID;default:NULL;" json:"planet_id"`                        // 所属星球ID
	Intro        string    `gorm:"column:intro;type:text;comment:个人简介;" json:"intro"`                                              // 个人简介
	IsOnline     int8      `gorm:"column:is_online;type:tinyint(1);comment:在线状态(0:离线,1:在线);default:0;" json:"is_online"`           // 在线状态(0:离线,1:在线)
	LastLoginAt  time.Time `gorm:"column:last_login_at;type:datetime;comment:最后登录时间;default:NULL;" json:"last_login_at"`           // 最后登录时间
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;comment:创建时间;default:CURRENT_TIMESTAMP;" json:"created_at"`      // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间;default:CURRENT_TIMESTAMP;" json:"updated_at"`      // 更新时间
}

func (u UsersTable) TableName() string {
	return "users_table"
}

type User interface {
	GetUserMobileInfo(*gorm.DB, *UsersTable) (*UsersTable, error)
}

type UserClient struct {
	repo User
	log  *log.Helper
}

func NewUserClient(repo User, logger log.Logger) *UserClient {
	return &UserClient{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (u *UserClient) UserMobileLogin(ctx context.Context, in *pb.UserLoginRegister) (interface{}, error) {
	info, err := u.repo.GetUserMobileInfo(config.DB, &UsersTable{
		PhoneNumber: in.Mobile,
	})
	if err != nil {
		return nil, errors.New("服务错误")
	}
	if info.UserId == 0 {
		return nil, errors.New("用户不存在")
	}
	return map[string]interface{}{
		"token": "12345",
	}, err
}
