package config

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

// GetOptimizedCacheConfig 获取内存优化的BigCache配置
// 大幅降低内存占用，从默认的~256MB降低到~10MB
func GetOptimizedCacheConfig(ttl time.Duration) bigcache.Config {
	maxMemoryMB := 10 // 默认每个缓存最大10MB（从20降至10）
	if Cache.MaxMemoryMB > 0 {
		maxMemoryMB = Cache.MaxMemoryMB
	}

	shards := 256 // 默认分片数（从1024降至256）
	if Cache.Shards > 0 {
		shards = Cache.Shards
	}

	maxEntries := 500 // 默认每分片最大条目数（从1000降至500）
	if Cache.MaxEntriesPerShard > 0 {
		maxEntries = Cache.MaxEntriesPerShard
	}

	return bigcache.Config{
		// 基础配置
		Shards:             shards,
		LifeWindow:         ttl,
		CleanWindow:        30 * time.Second, // 从1分钟降至30秒，更频繁清理
		MaxEntriesInWindow: maxEntries * shards, // 总条目数 = 分片数 × 每分片条目数

		// 内存优化配置
		MaxEntrySize:     256 * 1024, // 单个条目最大256KB（从512KB降至256KB）
		Verbose:          false,
		HardMaxCacheSize: maxMemoryMB, // 硬限制：最大内存占用MB
		OnRemove:         nil,

		// Logger配置（使用静默模式减少开销）
		Logger: nil,
	}
}

// CreateOptimizedCache 创建优化后的缓存实例
func CreateOptimizedCache(ttl time.Duration) (*bigcache.BigCache, error) {
	config := GetOptimizedCacheConfig(ttl)
	return bigcache.New(context.Background(), config)
}
