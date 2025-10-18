package service

import (
	"context"
	"kratosItems/internal/biz"
	"kratosItems/pkg"
	"net/http"

	pb "kratosItems/api/user"
)

type UserService struct {
	pb.UnimplementedUserServer
	userUseCase *biz.UserUseCase
}

func NewUserService(uc *biz.UserUseCase) *UserService {
	return &UserService{
		userUseCase: uc,
	}
}

// 手机验证码登录
func (s *UserService) PassengerUserMobileLogin(ctx context.Context, req *pb.UserLoginRegister) (*pb.PleaseReturn, error) {
	if req.Mobile == "" || req.SendSms == "" {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "手机号和验证码不能为空", nil)
	}
	
	result, err := s.userUseCase.UserMobileLogin(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, err.Error(), nil)
	}
	
	return pkg.RequestReturnDynamic(http.StatusOK, "登录成功", result)
}

// 手机验证码注册
func (s *UserService) PassengerUserMobileRegister(ctx context.Context, req *pb.UserRegisterRequest) (*pb.PleaseReturn, error) {
	if req.Mobile == "" || req.SendSms == "" || req.Password == "" {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "手机号、验证码和密码不能为空", nil)
	}
	
	result, err := s.userUseCase.UserMobileRegister(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, err.Error(), nil)
	}
	
	return pkg.RequestReturnDynamic(http.StatusOK, "注册成功", result)
}

// 密码登录
func (s *UserService) PassengerUserPasswordLogin(ctx context.Context, req *pb.UserPasswordLoginRequest) (*pb.PleaseReturn, error) {
	if req.Mobile == "" || req.Password == "" {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "手机号和密码不能为空", nil)
	}
	
	result, err := s.userUseCase.UserPasswordLogin(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, err.Error(), nil)
	}
	
	return pkg.RequestReturnDynamic(http.StatusOK, "登录成功", result)
}

// 获取用户信息
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.PleaseReturn, error) {
	if req.UserId <= 0 {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "用户ID不能为空", nil)
	}
	
	result, err := s.userUseCase.GetUserInfo(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, err.Error(), nil)
	}
	
	return pkg.RequestReturnDynamic(http.StatusOK, "获取成功", result)
}

// 更新用户信息
func (s *UserService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.PleaseReturn, error) {
	if req.UserId <= 0 {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "用户ID不能为空", nil)
	}
	
	result, err := s.userUseCase.UpdateUserInfo(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, err.Error(), nil)
	}
	
	return pkg.RequestReturnDynamic(http.StatusOK, "更新成功", result)
}

// 发送验证码
func (s *UserService) SendVerifyCode(ctx context.Context, req *pb.SendVerifyCodeRequest) (*pb.PleaseReturn, error) {
	if req.Mobile == "" || req.Type == "" {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "手机号和类型不能为空", nil)
	}
	
	if req.Type != "login" && req.Type != "register" {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "验证码类型只能是login或register", nil)
	}
	
	result, err := s.userUseCase.SendVerifyCode(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusInternalServerError, err.Error(), nil)
	}
	
	return pkg.RequestReturnDynamic(http.StatusOK, "发送成功", result)
}
