// Copyright (C) 2021 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package pakkuconf

import (
	"context"
	"strings"
	"time"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/utypes"

	"github.com/go-redis/redis/v8"
)

// NewRedisCache 事件模块-kafka实现
func NewRedisCache() ipakku.ICache {
	return &redisCache{}
}

type redisConf struct {
	raddrs       string
	username     string
	passwd       string
	dbIndex      int
	readTimeout  int64
	writeTimeout int64
}

// redisCache 基于Kafka的事件
type redisCache struct {
	conf   *redisConf
	clibs  *utypes.SafeMap
	client redis.Cmdable
	ctx    context.Context
}

func (ch *redisCache) initConf(conf ipakku.AppConfig) {
	ch.conf = &redisConf{}
	ch.clibs = utypes.NewSafeMap()
	ch.conf.raddrs = conf.GetConfig("appcache.redis.addrs").ToString("")
	ch.conf.username = conf.GetConfig("appcache.redis.username").ToString("")
	ch.conf.passwd = conf.GetConfig("appcache.redis.passwd").ToString("")
	ch.conf.dbIndex = conf.GetConfig("appcache.redis.dbIndex").ToInt(0)
	ch.conf.readTimeout = conf.GetConfig("appcache.redis.readTimeout").ToInt64(0)
	ch.conf.writeTimeout = conf.GetConfig("appcache.redis.writeTimeout").ToInt64(0)
	if len(ch.conf.raddrs) == 0 {
		logs.Panicln("redis 地址未设置")
	}
	logs.Infoln("redis init ", ch.conf.raddrs)
}

// Init 初始化缓存管理器, 一个对象只能初始化一次
func (ch *redisCache) Init(conf ipakku.AppConfig, appName string) {
	ch.initConf(conf)
	ch.ctx = context.Background()
	addrs := strings.Split(ch.conf.raddrs, ",")
	if len(addrs) > 1 {
		ch.client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        addrs,
			Username:     ch.conf.username,
			Password:     ch.conf.passwd,
			ReadTimeout:  time.Duration(ch.conf.readTimeout),
			WriteTimeout: time.Duration(ch.conf.writeTimeout),
		})
	} else {
		ch.client = redis.NewClient(&redis.Options{
			Addr:         ch.conf.raddrs,
			Username:     ch.conf.username,
			Password:     ch.conf.passwd,
			DB:           ch.conf.dbIndex,
			ReadTimeout:  time.Duration(ch.conf.readTimeout),
			WriteTimeout: time.Duration(ch.conf.writeTimeout),
		})
	}
}

// lib为库名, second:过期时间-1为不过期
func (ch *redisCache) RegLib(clib string, second int64) error {
	if len(clib) == 0 {
		return ipakku.ErrCacheLibIsExist
	}

	if _, ok := ch.clibs.Get(clib); ok {
		return ipakku.ErrCacheLibIsExist
	}
	if second > 0 {
		second = second * 1000000000
	}
	ch.clibs.Put(clib, second)
	return nil
}

// Set 向lib库中设置键为key的值tb
func (ch *redisCache) Set(clib string, key string, tb interface{}) error {
	if exp, ok := ch.clibs.Get(clib); !ok {
		return ipakku.ErrCacheLibNotExist
	} else {
		status := ch.client.Set(ch.ctx, clib+":"+key, tb, time.Duration(exp.(int64)))
		return status.Err()
	}
}

// Del 删除缓存信息
func (ch *redisCache) Del(clib string, key string) error {
	if _, ok := ch.clibs.Get(clib); !ok {
		return ipakku.ErrCacheLibNotExist
	} else {
		status := ch.client.Del(ch.ctx, clib+":"+key)
		return status.Err()
	}
}

// Get 读取缓存信息
func (ch *redisCache) Get(clib string, key string, val interface{}) error {
	if res := ch.client.Get(ch.ctx, clib+":"+key); nil != res.Err() {
		if res.Err().Error() == "redis: nil" {
			return ipakku.ErrNoCacheHit
		}
		return res.Err()
	} else {
		return res.Scan(val)
	}
}

// Keys 获取库的所有key
func (ch *redisCache) Keys(clib string) []string {
	res := ch.client.Keys(ch.ctx, clib+":*")
	return res.Val()
}

// Clear 清空库内容
func (ch *redisCache) Clear(clib string) {
	if keys := ch.Keys(clib); len(keys) > 0 {
		if res := ch.client.Del(ch.ctx, keys...); res.Err() != nil {
			panic(res.Err())
		}
	}
}
