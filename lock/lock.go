/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-08-15 11:34:17
# File Name: lock.go
# Description:
####################################################################### */

package lock

import (
	"errors"
	"fmt"
	"time"

	"github.com/ant-libs-go/util"
	rds "github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

/*
	lk, _ := New("hahahalock", redis.GetPool("default"))
	lk.Transaction().WaitAndLock(86400)
	fmt.Println(lk.Commit())
*/

type LockStatus int

var (
	ErrLock          = errors.New("lock error")
	ErrTransFinished = errors.New("transaction finished")
)

const (
	LockStatusNormal   LockStatus = 0
	LockStatusDoing    LockStatus = 1
	LockStatusFinished LockStatus = 2
)

type Lock struct {
	key         string
	transaction bool
	status      LockStatus
	client      *rds.Pool
}

func New(key interface{}, client *rds.Pool) (r *Lock, err error) {
	if key == nil {
		err = fmt.Errorf("key is empty")
		return
	}
	var k []byte
	if k, err = util.GobSerialize(key); err != nil {
		return
	}
	r = &Lock{client: client}
	if r.key, _ = key.(string); len(r.key) == 0 {
		r.key = uuid.NewV3(uuid.Nil, string(k)).String()
	}
	r.key = fmt.Sprintf("LOCK.%s", r.key)
	r.refreshLockStatus()
	return
}

func (this *Lock) refreshLockStatus() (err error) {
	conn := this.client.Get()
	defer conn.Close()

	var v int
	v, err = rds.Int(conn.Do("GET", this.key))
	if err != nil {
		return
	}
	this.status = LockStatus(v)
	return
}

func (this *Lock) GetLockStatus() (r LockStatus) {
	return this.status
}

func (this *Lock) Transaction() *Lock {
	this.transaction = true
	return this
}

func (this *Lock) Wait() (err error) {
	for {
		if err = this.refreshLockStatus(); err != nil {
			return
		}
		if this.transaction == true &&
			this.GetLockStatus() == LockStatusFinished {
			err = ErrTransFinished
			return
		}
		if this.GetLockStatus() == LockStatusNormal {
			break
		}
		time.Sleep(5 * time.Millisecond) // 5毫秒
	}
	return
}

func (this *Lock) Lock(aliveSeconds int64) (err error) {
	conn := this.client.Get()
	defer conn.Close()

	status := LockStatusFinished
	if this.transaction {
		status = LockStatusDoing
	}
	defer this.refreshLockStatus()

	var suc bool
	suc, err = rds.Bool(conn.Do("SET", this.key, int(status), "EX", aliveSeconds, "NX"))
	if err != nil {
		return
	}
	if suc == false {
		err = ErrLock
	}
	return
}

func (this *Lock) WaitAndLock(aliveSeconds int64) (err error) {
	for {
		if err = this.Wait(); err != nil {
			return
		}
		err = this.Lock(aliveSeconds)
		if err == nil || err.Error() != ErrLock.Error() {
			break
		}
	}
	return
}

func (this *Lock) Commit() (err error) {
	if this.transaction == false {
		return fmt.Errorf("no transaction")
	}
	defer this.refreshLockStatus()

	conn := this.client.Get()
	defer conn.Close()

	var aliveSeconds int64
	aliveSeconds, err = rds.Int64(conn.Do("TTL", this.key))
	if err != nil {
		return
	}
	var suc bool
	suc, err = rds.Bool(conn.Do("SET", this.key, int(LockStatusFinished), "EX", aliveSeconds, "XX"))
	if err != nil {
		return
	}
	if suc == false {
		err = fmt.Errorf("commit error")
	}
	return
}

func (this *Lock) Release() (err error) {
	conn := this.client.Get()
	defer this.refreshLockStatus()
	defer conn.Close()
	_, err = conn.Do("DEL", this.key)
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
