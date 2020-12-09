/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-28 13:11:55
# File Name: timer.go
# Description:
####################################################################### */

package timer

import (
	"fmt"
	"time"

	"github.com/ant-libs-go/redis"
	"github.com/ant-libs-go/safe_stop"
	"github.com/ant-libs-go/util"
	rds "github.com/gomodule/redigo/redis"
)

type Timer struct {
	key          string
	client       *rds.Pool
	ticker       *time.Ticker
	delaySeconds int64
	callback     func(token string, tm time.Duration)
}

func NewTimer(timerId string, client *rds.Pool, fn func(token string, tm time.Duration), delaySeconds int64) *Timer {
	o := &Timer{client: client, callback: fn, delaySeconds: delaySeconds}
	o.key = fmt.Sprintf("TIMER.%s", timerId)
	return o
}

func (this *Timer) Start(interval time.Duration) {
	this.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-this.ticker.C:
				this.run()
			}
		}
	}()
}

func (this *Timer) run() {
	safe_stop.Lock(1)
	defer safe_stop.Unlock()
	conn := this.client.Get()
	defer conn.Close()

	values, err := rds.StringMap(conn.Do("ZRANGEBYSCORE", this.key, 0, time.Now().Unix(), "WITHSCORES"))
	if redis.HasError(err) {
		return
	}
	for token, ts := range values {
		suc, _ := rds.Bool(conn.Do("ZREM", this.key, token))
		if suc == false {
			continue
		}
		go func(token string, ts string) {
			safe_stop.Lock(1)
			defer safe_stop.Unlock()
			this.callback(token, time.Second*time.Duration(util.StrToInt32(ts, 0)))
		}(token, ts)
	}
}

func (this *Timer) Add(token string, delaySeconds int64) (err error) {
	conn := this.client.Get()
	defer conn.Close()

	if delaySeconds == 0 {
		delaySeconds = this.delaySeconds
	}
	if delaySeconds == 0 {
		delaySeconds = 60
	}

	_, err = conn.Do("ZADD", this.key, "NX", time.Now().Unix()+delaySeconds, token)
	return
}

func (this *Timer) Stop() {
	this.ticker.Stop()
}

func (this *Timer) Close() (err error) {
	conn := this.client.Get()
	defer conn.Close()

	this.Stop()
	_, err = conn.Do("DEL", this.key)
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
