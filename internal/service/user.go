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
	UserClient *biz.UserClient
}

func NewUserService(uc *biz.UserClient) *UserService {
	return &UserService{
		UserClient: uc,
	}
}

func (s *UserService) PassengerUserMobileLogin(ctx context.Context, req *pb.UserLoginRegister) (*pb.PleaseReturn, error) {
	if req.Mobile == "" || req.SendSms == "" {
		return pkg.RequestReturnDynamic(http.StatusBadRequest, "验证失败", nil)
	}
	login, err := s.UserClient.UserMobileLogin(ctx, req)
	if err != nil {
		return pkg.RequestReturnDynamic(http.StatusInternalServerError, err.Error(), nil)
	}
	return pkg.RequestReturnDynamic(http.StatusOK, "登录成功", login)
}
