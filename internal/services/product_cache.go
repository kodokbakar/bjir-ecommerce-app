package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ProductListCacheTTL    = 5 * time.Minute
	productListCachePrefix = "products:list"
)

type ProductCache interface {
	GetProductList(ctx context.Context, key string) (*ProductListResult, error)
	SetProductList(ctx context.Context, key string, result *ProductListResult, ttl time.Duration) error
	InvalidateProductList(ctx context.Context) error
}

type redisProductCache struct {
	client *redis.Client
}

func NewRedisProductCache(client *redis.Client) ProductCache {
	if client == nil {
		return nil
	}

	return &redisProductCache{client: client}
}

func (c *redisProductCache) GetProductList(ctx context.Context, key string) (*ProductListResult, error) {
	cachedValue, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var result ProductListResult
	if err := json.Unmarshal([]byte(cachedValue), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *redisProductCache) SetProductList(ctx context.Context, key string, result *ProductListResult, ttl time.Duration) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *redisProductCache) InvalidateProductList(ctx context.Context) error {
	var cursor uint64

	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, productListCachePrefix+":*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

func buildProductListCacheKey(input ProductListInput) string {
	return fmt.Sprintf(
		"%s:page=%d:limit=%d:category_id=%s:category=%s:search=%s:sort_by=%s:sort_order=%s",
		productListCachePrefix,
		input.Page,
		input.Limit,
		escapeCacheKeyPart(input.CategoryID),
		escapeCacheKeyPart(input.CategorySlug),
		escapeCacheKeyPart(input.Search),
		escapeCacheKeyPart(input.SortBy),
		escapeCacheKeyPart(input.SortOrder),
	)
}

func escapeCacheKeyPart(value string) string {
	return url.QueryEscape(value)
}
