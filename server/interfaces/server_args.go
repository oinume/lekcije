package interfaces

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/redis.v4"
)

type ServerArgs struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}
