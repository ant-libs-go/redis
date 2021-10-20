/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-11-07 20:15:51
# File Name: storer.go
# Description:
####################################################################### */

package storer

import (
	"fmt"

	"github.com/ant-libs-go/redis/lock"
	"github.com/ant-libs-go/redis/timer"
)

type Operator interface {
	SetData(inp interface{})
	GetData() (r interface{})
	LoadFromRedis() (r interface{}, err error)
	LoadFromMysql() (r interface{}, err error)
	SaveToRedis() (err error)
	SaveToMysql() (err error)
}

type Storer struct {
	token string
	op    Operator
	lk    *lock.Lock
	timer *timer.Timer
}

// 同类数据的Storer的timerId唯一，如 PASSPORT
// token 与 lock 中的 lockId 一致
func NewStorer(token string, op Operator, lk *lock.Lock, timer *timer.Timer) (r *Storer, err error) {
	o := &Storer{token: token, lk: lk, op: op, timer: timer}

	if err = o.lk.WaitAndLock(60); err != nil {
		return
	}
	if err = o.load(); err != nil {
		o.Release()
		return
	}
	if timer != nil {
		if err = timer.Add(token, 0); err != nil {
			o.Release()
			return
		}
	}
	return o, nil
}

func (this *Storer) load() (err error) {
	var data interface{}
	if data, err = this.op.LoadFromRedis(); err != nil {
		return
	}
	if data != nil {
		this.op.SetData(data)
		return
	}
	if data, err = this.op.LoadFromMysql(); err != nil {
		return
	}
	if data != nil {
		this.op.SetData(data)
		if err = this.op.SaveToRedis(); err != nil {
			return
		}
	}
	return
}

func (this *Storer) Reload() (err error) {
	return this.load()
}

func (this *Storer) GetOperator() (r Operator) {
	return this.op
}

func (this *Storer) GetToken() (r string) {
	return this.token
}

func (this *Storer) GetStorerId() (r string) {
	return fmt.Sprintf("STORER.%s", this.token)
}

func (this *Storer) Release() (err error) {
	return this.lk.Release()
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
