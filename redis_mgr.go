/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-29 21:07:46
# File Name: redis_mgr.go
# Description:
####################################################################### */

package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
	rds "github.com/gomodule/redigo/redis"
)

var (
	once  sync.Once
	lock  sync.RWMutex
	pools map[string]*rds.Pool
)

func init() {
	pools = map[string]*rds.Pool{}
}

type redisConfig struct {
	Cfgs map[string]*Cfg `toml:"redis"`
}

type Cfg struct {
	// dial
	DialAddr           string        `toml:"addr"`
	DialUsername       string        `toml:"user"`
	DialPassword       string        `toml:"pawd"`
	DialDatabase       int           `toml:"database"`
	DialConnectTimeout time.Duration `toml:"dial_timeout"`
	DialReadTimeout    time.Duration `toml:"read_timeout"`
	DialWriteTimeout   time.Duration `toml:"write_timeout"`

	// pool
	PoolMaxIdle         int           `toml:"pool_max_idle"`           // 最大闲置连接数
	PoolMaxActive       int           `toml:"pool_max_active"`         // 最大活跃连接数
	PoolIdleTimeout     time.Duration `toml:"pool_idle_time"`          // 闲置的过期时间，在Get方法中会对过期的连接删除
	PoolWait            bool          `toml:"pool_wait"`               // 当活跃连接达到上限，Get时是等待还是返回错误。为false时返回错误，为true时阻塞等待
	PoolMaxConnLifetime time.Duration `toml:"pool_max_conn_life_time"` // 连接最长生存时间，当超过时间会被删除
}

// 验证Redis实例的配置正确性与连通性。
// 参数names是实例的名称列表，如果为空则检测所有配置的实例
func Valid(names ...string) (err error) {
	if len(names) == 0 {
		var cfgs map[string]*Cfg
		if cfgs, err = loadCfgs(); err != nil {
			return
		}
		for k, _ := range cfgs {
			names = append(names, k)
		}
	}
	for _, name := range names {
		var cli rds.Conn
		cli, err = SafeClient(name)
		if err == nil {
			defer cli.Close()
			_, err = cli.Do("PING")
		}
		if err != nil {
			err = fmt.Errorf("redis#%s is invalid, %s", name, err)
			return
		}
	}
	return
}

func DefaultClient() (r rds.Conn) {
	return Client("default")
}

func DefaultPool() (r *rds.Pool) {
	return Pool("default")
}

func Client(name string) (r rds.Conn) {
	r = Pool(name).Get()
	return
}

func SafeClient(name string) (r rds.Conn, err error) {
	var pool *rds.Pool
	pool, err = SafePool(name)
	if err != nil {
		return
	}
	r = pool.Get()
	return
}

func Pool(name string) (r *rds.Pool) {
	var err error
	if r, err = getPool(name); err != nil {
		panic(err)
	}
	return
}

func SafePool(name string) (r *rds.Pool, err error) {
	return getPool(name)
}

func getPool(name string) (r *rds.Pool, err error) {
	lock.RLock()
	r = pools[name]
	lock.RUnlock()
	if r == nil {
		r, err = addPool(name)
	}
	return
}

func addPool(name string) (r *rds.Pool, err error) {
	var cfg *Cfg
	if cfg, err = loadCfg(name); err != nil {
		return
	}
	r = NewRedisPool(cfg)

	lock.Lock()
	pools[name] = r
	lock.Unlock()
	return
}

func loadCfg(name string) (r *Cfg, err error) {
	var cfgs map[string]*Cfg
	if cfgs, err = loadCfgs(); err != nil {
		return
	}
	if r = cfgs[name]; r == nil {
		err = fmt.Errorf("redis#%s not configed", name)
		return
	}
	return
}

func loadCfgs() (r map[string]*Cfg, err error) {
	r = map[string]*Cfg{}

	once.Do(func() {
		config.Get(&redisConfig{}, options.WithOpOnChangeFn(func(cfg interface{}) {
			lock.Lock()
			defer lock.Unlock()
			pools = map[string]*rds.Pool{}
		}))
	})

	cfg := config.Get(&redisConfig{}).(*redisConfig)
	if err == nil && (cfg.Cfgs == nil || len(cfg.Cfgs) == 0) {
		err = fmt.Errorf("not configed")
	}
	if err != nil {
		err = fmt.Errorf("redis load cfgs error, %s", err)
		return
	}
	r = cfg.Cfgs
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
