package store

import (
	"context"
	"encoding/json"
	"github.com/r2day/db"
)

const (
	// CacheInfoPrefix 缓存前缀
	CacheInfoPrefix = "store_cache_info"
)

// Cache 缓存器
type Cache struct {
	Ctx context.Context
}

// NewCache 创建新的门店缓存
func NewCache(ctx context.Context) *Cache {
	return &Cache{
		Ctx: ctx,
	}
}

// Dump 将门店信息存入缓存中
func (s *Cache) Dump(m *Model) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// 存入数据库
	err = db.RDB.HSet(s.Ctx, CacheInfoPrefix, m.StoreID, payload).Err()
	if err != nil {
		return err
	}
	return nil
}

// Load 通过门店id 获取到门店信息
func (s *Cache) Load(storeID string) (*Model, error) {
	data := &Model{}
	// 存入数据库
	payload, err := db.RDB.HGet(s.Ctx, CacheInfoPrefix, storeID).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(payload), data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
