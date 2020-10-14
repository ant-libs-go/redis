/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-13 11:27:36
# File Name: redis2/options.go
# Description:
####################################################################### */

package redis

import (
	"time"

	rds "github.com/gomodule/redigo/redis"
)

type Option func(o interface{}) interface{}

func WithPoolMaxIdle(inp int) Option {
	return func(o interface{}) (r interface{}) {
		switch obj := o.(type) {
		case *rds.Pool:
			obj.MaxIdle = inp
		}
		return
	}
}

func WithPoolMaxActive(inp int) Option {
	return func(o interface{}) (r interface{}) {
		switch obj := o.(type) {
		case *rds.Pool:
			obj.MaxActive = inp
		}
		return
	}
}

func WithPoolIdleTimeout(inp time.Duration) Option {
	return func(o interface{}) (r interface{}) {
		switch obj := o.(type) {
		case *rds.Pool:
			obj.IdleTimeout = inp
		}
		return
	}
}

func WithPoolWait(inp bool) Option {
	return func(o interface{}) (r interface{}) {
		switch obj := o.(type) {
		case *rds.Pool:
			obj.Wait = inp
		}
		return
	}
}

func WithPoolMaxConnLifetime(inp time.Duration) Option {
	return func(o interface{}) (r interface{}) {
		switch obj := o.(type) {
		case *rds.Pool:
			obj.MaxConnLifetime = inp
		}
		return
	}
}

func WithDialUsername(inp string) Option {
	return func(o interface{}) (r interface{}) {
		switch o.(type) {
		case *rds.Pool:
		default:
			return rds.DialUsername(inp)
		}
		return
	}
}

func WithDialPassword(inp string) Option {
	return func(o interface{}) (r interface{}) {
		switch o.(type) {
		case *rds.Pool:
		default:
			return rds.DialPassword(inp)
		}
		return
	}
}

func WithDialDatabase(inp int) Option {
	return func(o interface{}) (r interface{}) {
		switch o.(type) {
		case *rds.Pool:
		default:
			return rds.DialDatabase(inp)
		}
		return
	}
}

func WithDialConnectTimeout(inp time.Duration) Option {
	return func(o interface{}) (r interface{}) {
		switch o.(type) {
		case *rds.Pool:
		default:
			return rds.DialConnectTimeout(inp)
		}
		return
	}
}

func WithDialReadTimeout(inp time.Duration) Option {
	return func(o interface{}) (r interface{}) {
		switch o.(type) {
		case *rds.Pool:
		default:
			return rds.DialReadTimeout(inp)
		}
		return
	}
}

func WithDialWriteTimeout(inp time.Duration) Option {
	return func(o interface{}) (r interface{}) {
		switch o.(type) {
		case *rds.Pool:
		default:
			return rds.DialWriteTimeout(inp)
		}
		return
	}
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
