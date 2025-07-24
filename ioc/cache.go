package ioc

import (
	cache "github.com/patrickmn/go-cache"
	"time"
)

func InitCache() *cache.Cache {
	return cache.New(10*time.Minute, 20*time.Minute)
}
