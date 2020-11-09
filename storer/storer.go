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
	key   string
	op    Operator
	lk    *lock.Lock
	timer *timer.Timer
}

func NewStorer(storerId string, op Operator, lk *lock.Lock, timer *timer.Timer) (r *Storer, err error) {
	o := &Storer{lk: lk, op: op, timer: timer}
	o.key = fmt.Sprintf("STORER.%s", storerId)

	if err = o.lk.WaitAndLock(60); err != nil {
		return
	}
	if err = o.load(); err != nil {
		return
	}
	if timer == nil {
		return
	}
	if err = timer.Add(storerId, 0); err != nil {
		return
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

func (this *Storer) Release() (err error) {
	return this.lk.Release()
}

func (this *Storer) GetOperator() (r Operator) {
	return this.op
}

func (this *Storer) GetData() (r interface{}) {
	return this.op.GetData()
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
