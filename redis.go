/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-13 11:18:19
# File Name: redis.go
# Description:
####################################################################### */

package redis

import (
	"fmt"
	"time"

	rds "github.com/gomodule/redigo/redis"
)

func NewRedisPool(cfg *Cfg) *rds.Pool {
	pool := &rds.Pool{
		MaxIdle:     20,
		IdleTimeout: 60 * time.Second,
		Wait:        cfg.PoolWait,
		Dial: func() (r rds.Conn, err error) {
			return rds.Dial("tcp", fmt.Sprintf("%s:%d", cfg.DialHost, cfg.DialPort), buildDialOptions(cfg)...)
		},
		TestOnBorrow: func(c rds.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	if cfg.PoolMaxIdle > 0 {
		pool.MaxIdle = cfg.PoolMaxIdle
	}
	if cfg.PoolMaxActive > 0 {
		pool.MaxActive = cfg.PoolMaxActive
	}
	if cfg.PoolIdleTimeout > 0 {
		pool.IdleTimeout = cfg.PoolIdleTimeout * time.Millisecond
	}
	if cfg.PoolMaxConnLifetime > 0 {
		pool.MaxConnLifetime = cfg.PoolMaxConnLifetime * time.Millisecond
	}
	return pool
}

func buildDialOptions(cfg *Cfg) (r []rds.DialOption) {
	if len(cfg.DialUsername) > 0 {
		r = append(r, rds.DialUsername(cfg.DialUsername))
	}
	if len(cfg.DialPassword) > 0 {
		r = append(r, rds.DialPassword(cfg.DialPassword))
	}
	if cfg.DialDatabase > 0 {
		r = append(r, rds.DialDatabase(cfg.DialDatabase))
	}
	if cfg.DialConnectTimeout > 0 {
		r = append(r, rds.DialConnectTimeout(cfg.DialConnectTimeout*time.Millisecond))
	} else {
		r = append(r, rds.DialConnectTimeout(1000*time.Millisecond))
	}
	if cfg.DialReadTimeout > 0 {
		r = append(r, rds.DialReadTimeout(cfg.DialReadTimeout*time.Millisecond))
	} else {
		r = append(r, rds.DialReadTimeout(500*time.Millisecond))
	}
	if cfg.DialWriteTimeout > 0 {
		r = append(r, rds.DialWriteTimeout(cfg.DialWriteTimeout*time.Millisecond))
	} else {
		r = append(r, rds.DialWriteTimeout(500*time.Millisecond))
	}
	return
}

func HasError(err error) bool {
	if err != nil && err.Error() != rds.ErrNil.Error() {
		return true
	}
	return false
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
