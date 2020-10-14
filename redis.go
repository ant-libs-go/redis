/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-13 11:18:19
# File Name: redis.go
# Description:
####################################################################### */

package redis

import (
	"sync"
	"time"

	rds "github.com/gomodule/redigo/redis"
)

var (
	lock  sync.RWMutex
	pools map[string]*rds.Pool
)

func init() {
	pools = map[string]*rds.Pool{}
}

func New(name string, addr string, opts ...Option) {
	dialOptions := []rds.DialOption{
		rds.DialConnectTimeout(time.Duration(1000) * time.Millisecond),
		rds.DialReadTimeout(time.Duration(1000) * time.Millisecond),
		rds.DialWriteTimeout(time.Duration(1000) * time.Millisecond),
	}
	for _, opt := range opts {
		o := opt(nil)
		if v, ok := o.(rds.DialOption); ok {
			dialOptions = append(dialOptions, v)
		}
	}

	pool := &rds.Pool{
		MaxIdle:     20,
		IdleTimeout: 60 * time.Second,
		TestOnBorrow: func(c rds.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (r rds.Conn, err error) {
			return rds.Dial("tcp", addr, dialOptions...)
		},
	}
	for _, opt := range opts {
		opt(pool)
	}

	lock.Lock()
	pools[name] = pool
	lock.Unlock()
}

func GetPool(name string) *rds.Pool {
	lock.RLock()
	pool, ok := pools[name]
	lock.RUnlock()
	if !ok {
		return nil
	}
	return pool
}

func Get(name string) rds.Conn {
	pool := GetPool(name)
	if pool == nil {
		return nil
	}
	return pool.Get()
}

func HasError(err error) bool {
	if err != nil && err.Error() != rds.ErrNil.Error() {
		return true
	}
	return false
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
