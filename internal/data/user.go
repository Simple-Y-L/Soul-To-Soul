package data

import (
	"fmt"
	"kratosItems/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 根据手机号获取用户信息
func (r *userRepo) GetUserByMobile(db *gorm.DB, mobile string) (*biz.UsersTable, error) {
	var user biz.UsersTable
	err := db.Where("phone_number = ?", mobile).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// 根据用户ID获取用户信息
func (r *userRepo) GetUserById(db *gorm.DB, userId int64) (*biz.UsersTable, error) {
	var user biz.UsersTable
	err := db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// 创建用户
func (r *userRepo) CreateUser(db *gorm.DB, user *biz.UsersTable) (*biz.UsersTable, error) {
	err := db.Create(user).Error
	if err != nil {
		r.log.Errorf("创建用户失败: %v", err)
		return nil, err
	}
	return user, nil
}

// 更新用户信息
func (r *userRepo) UpdateUser(db *gorm.DB, user *biz.UsersTable) error {
	err := db.Save(user).Error
	if err != nil {
		r.log.Errorf("更新用户失败: %v", err)
		return err
	}
	return nil
}

// 保存验证码到缓存（这里使用内存存储，生产环境应该使用Redis）
var verifyCodeCache = make(map[string]biz.VerifyCodeCache)

func (r *userRepo) SaveVerifyCode(mobile, code, codeType string, expiration time.Duration) error {
	key := fmt.Sprintf("%s_%s", mobile, codeType)
	verifyCodeCache[key] = biz.VerifyCodeCache{
		Mobile:    mobile,
		Code:      code,
		Type:      codeType,
		ExpiredAt: time.Now().Add(expiration),
	}
	r.log.Infof("保存验证码: %s -> %s", key, code)
	return nil
}

// 获取验证码
func (r *userRepo) GetVerifyCode(mobile, codeType string) (string, error) {
	key := fmt.Sprintf("%s_%s", mobile, codeType)
	cache, exists := verifyCodeCache[key]
	if !exists {
		return "", fmt.Errorf("验证码不存在")
	}
	
	if time.Now().After(cache.ExpiredAt) {
		delete(verifyCodeCache, key)
		return "", fmt.Errorf("验证码已过期")
	}
	
	return cache.Code, nil
}

// 删除验证码
func (r *userRepo) DeleteVerifyCode(mobile, codeType string) error {
	key := fmt.Sprintf("%s_%s", mobile, codeType)
	delete(verifyCodeCache, key)
	r.log.Infof("删除验证码: %s", key)
	return nil
}
