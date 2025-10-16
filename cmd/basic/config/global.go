package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Ctx = context.Background()
var RDB *redis.Client
var Configs Config
