package server

import (
	v1 "kratosItems/api/helloworld/v1"
	userv1 "kratosItems/api/user"
	"kratosItems/internal/conf"
	middlewares "kratosItems/internal/middleware"
	"kratosItems/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, user *service.UserService, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middlewares.Logger(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(":8080"))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	userv1.RegisterUserHTTPServer(srv, user)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
