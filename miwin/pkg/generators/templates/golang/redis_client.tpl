package client

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	conf "github.com/mimokpl/kratos-bootstrap/api/gen/go/conf/v1"
	redisClient "github.com/mimokpl/kratos-bootstrap/cache/redis"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(ctx *bootstrap.Context) (*redis.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil
	}

	cli := redisClient.NewClient(cfg.Data, ctx.NewLoggerHelper("redis/data/{{.Service}}-service"))

    return cli, func() {
        if err := cli.Close(); err != nil {
            l.Error(err)
        }
    }, nil
}
