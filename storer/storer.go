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
	LoadFromRedis(key string) (r interface{}, err error)
	LoadFromMysql() (r interface{}, err error)
	SaveToRedis(key string) (err error)
	SaveToMysql() (err error)
}

type Storer struct {
	token string
	key   string
	op    Operator
	lk    *lock.Lock
	timer *timer.Timer
}

func NewStorer(token string, op Operator, lk *lock.Lock, timer *timer.Timer) (r *Storer, err error) {
	o := &Storer{token: token, lk: lk, op: op, timer: timer}
	o.key = fmt.Sprintf("STORER.%s", token)

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
	if data, err = this.LoadFromRedis(); err != nil {
		return
	}
	if data != nil {
		this.SetData(data)
		return
	}
	if data, err = this.LoadFromMysql(); err != nil {
		return
	}
	if data != nil {
		this.SetData(data)
		if err = this.SaveToRedis(); err != nil {
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

func (this *Storer) SetData(inp interface{}) {
	this.op.SetData(inp)
}

func (this *Storer) GetData() (r interface{}) {
	return this.op.GetData()
}

func (this *Storer) LoadFromRedis() (r interface{}, err error) {
	return this.op.LoadFromRedis(this.key)
}

func (this *Storer) LoadFromMysql() (r interface{}, err error) {
	return this.op.LoadFromMysql()
}

func (this *Storer) SaveToRedis() (err error) {
	return this.op.SaveToRedis(this.key)
}

func (this *Storer) SaveToMysql() (err error) {
	return this.op.SaveToMysql()
}

func (this *Storer) GetToken() (r string) {
	return this.token
}

func (this *Storer) Release() (err error) {
	return this.lk.Release()
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
