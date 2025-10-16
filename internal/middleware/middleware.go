package middlewares

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

func Logger() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}
	}
}
