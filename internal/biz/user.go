package biz

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	pb "kratosItems/api/user"
	"kratosItems/cmd/basic/config"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// 用户表
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

// 验证码存储结构
type VerifyCodeCache struct {
	Mobile    string
	Code      string
	Type      string // login/register
	ExpiredAt time.Time
}

type UserRepo interface {
	GetUserByMobile(*gorm.DB, string) (*UsersTable, error)
	GetUserById(*gorm.DB, int64) (*UsersTable, error)
	CreateUser(*gorm.DB, *UsersTable) (*UsersTable, error)
	UpdateUser(*gorm.DB, *UsersTable) error
	SaveVerifyCode(string, string, string, time.Duration) error
	GetVerifyCode(string, string) (string, error)
	DeleteVerifyCode(string, string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// 手机验证码登录
func (u *UserUseCase) UserMobileLogin(ctx context.Context, req *pb.UserLoginRegister) (interface{}, error) {
	// 验证验证码
	storedCode, err := u.repo.GetVerifyCode(req.Mobile, "login")
	if err != nil {
		return nil, errors.New("验证码已过期或不存在")
	}
	if storedCode != req.SendSms {
		return nil, errors.New("验证码错误")
	}

	// 查找用户
	user, err := u.repo.GetUserByMobile(config.DB, req.Mobile)
	if err != nil {
		return nil, errors.New("用户不存在，请先注册")
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	user.IsOnline = 1
	err = u.repo.UpdateUser(config.DB, user)
	if err != nil {
		u.log.Errorf("更新用户登录时间失败: %v", err)
	}

	// 删除验证码
	u.repo.DeleteVerifyCode(req.Mobile, "login")

	// 生成token（这里简化处理，实际应该使用JWT）
	token := u.generateToken(user.UserId)

	return map[string]interface{}{
		"token":    token,
		"user_id":  user.UserId,
		"nickname": user.Nickname,
		"avatar":   user.AvatarUrl,
	}, nil
}

// 用户注册
func (u *UserUseCase) UserMobileRegister(ctx context.Context, req *pb.UserRegisterRequest) (interface{}, error) {
	// 验证验证码
	storedCode, err := u.repo.GetVerifyCode(req.Mobile, "register")
	if err != nil {
		return nil, errors.New("验证码已过期或不存在")
	}
	if storedCode != req.SendSms {
		return nil, errors.New("验证码错误")
	}

	// 检查用户是否已存在
	existUser, _ := u.repo.GetUserByMobile(config.DB, req.Mobile)
	if existUser != nil && existUser.UserId > 0 {
		return nil, errors.New("用户已存在，请直接登录")
	}

	// 创建新用户
	now := time.Now()
	newUser := &UsersTable{
		PhoneNumber:  req.Mobile,
		PasswordHash: u.hashPassword(req.Password),
		Nickname:     req.Nickname,
		Gender:       req.Gender,
		IsOnline:     1,
		LastLoginAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 如果没有提供昵称，生成默认昵称
	if newUser.Nickname == "" {
		newUser.Nickname = "用户" + req.Mobile[7:]
	}

	// 如果没有提供性别，设置为unknown
	if newUser.Gender == "" {
		newUser.Gender = "unknown"
	}

	user, err := u.repo.CreateUser(config.DB, newUser)
	if err != nil {
		return nil, errors.New("注册失败，请稍后重试")
	}

	// 删除验证码
	u.repo.DeleteVerifyCode(req.Mobile, "register")

	// 生成token
	token := u.generateToken(user.UserId)

	return map[string]interface{}{
		"token":    token,
		"user_id":  user.UserId,
		"nickname": user.Nickname,
		"message":  "注册成功",
	}, nil
}

// 密码登录
func (u *UserUseCase) UserPasswordLogin(ctx context.Context, req *pb.UserPasswordLoginRequest) (interface{}, error) {
	// 查找用户
	user, err := u.repo.GetUserByMobile(config.DB, req.Mobile)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if !u.verifyPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	user.IsOnline = 1
	err = u.repo.UpdateUser(config.DB, user)
	if err != nil {
		u.log.Errorf("更新用户登录时间失败: %v", err)
	}

	// 生成token
	token := u.generateToken(user.UserId)

	return map[string]interface{}{
		"token":    token,
		"user_id":  user.UserId,
		"nickname": user.Nickname,
		"avatar":   user.AvatarUrl,
	}, nil
}

// 获取用户信息
func (u *UserUseCase) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (interface{}, error) {
	user, err := u.repo.GetUserById(config.DB, req.UserId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return map[string]interface{}{
		"user_id":      user.UserId,
		"phone_number": user.PhoneNumber,
		"nickname":     user.Nickname,
		"avatar_url":   user.AvatarUrl,
		"gender":       user.Gender,
		"birthday":     user.Birthday.Format("2006-01-02"),
		"intro":        user.Intro,
		"is_online":    user.IsOnline,
		"created_at":   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// 更新用户信息
func (u *UserUseCase) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (interface{}, error) {
	user, err := u.repo.GetUserById(config.DB, req.UserId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarUrl != "" {
		user.AvatarUrl = req.AvatarUrl
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Birthday != "" {
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			user.Birthday = birthday
		}
	}
	if req.Intro != "" {
		user.Intro = req.Intro
	}
	user.UpdatedAt = time.Now()

	err = u.repo.UpdateUser(config.DB, user)
	if err != nil {
		return nil, errors.New("更新失败")
	}

	return map[string]interface{}{
		"message": "更新成功",
	}, nil
}

// 发送验证码
func (u *UserUseCase) SendVerifyCode(ctx context.Context, req *pb.SendVerifyCodeRequest) (interface{}, error) {
	// 生成6位数字验证码
	code := u.generateVerifyCode()

	// 存储验证码，5分钟过期
	err := u.repo.SaveVerifyCode(req.Mobile, code, req.Type, 5*time.Minute)
	if err != nil {
		return nil, errors.New("发送验证码失败")
	}

	// 这里应该调用短信服务发送验证码
	// 暂时返回验证码用于测试
	u.log.Infof("发送验证码到 %s: %s", req.Mobile, code)

	return map[string]interface{}{
		"message": "验证码发送成功",
		"code":    code, // 测试环境返回验证码，生产环境应该删除
	}, nil
}

// 生成验证码
func (u *UserUseCase) generateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}

// 密码加密
func (u *UserUseCase) hashPassword(password string) string {
	hash := md5.Sum([]byte(password + "soul_salt"))
	return fmt.Sprintf("%x", hash)
}

// 密码验证
func (u *UserUseCase) verifyPassword(password, hash string) bool {
	return u.hashPassword(password) == hash
}

// 生成token（简化版本）
func (u *UserUseCase) generateToken(userId int64) string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%d_%d", userId, timestamp)
	hash := md5.Sum([]byte(data + "soul_token_salt"))
	return fmt.Sprintf("%x", hash)
}
