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

	"github.com/ant-libs-go/redis"
	"github.com/ant-libs-go/util"
	rds "github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

/*
	lk, _ := New("hahahalock", redis.DefaultPool())
	lk.Transaction().WaitAndLock(86400)
	fmt.Println(lk.Commit())

	1、lk.Lock(10)，加锁，当锁存在时抛出错误
	2、lk.Wait()，等待状态为未加锁
	3、lk.WaitAndLock(10)，加锁，当锁存在时等待，直到加锁成功
	4、lk.Transaction().WaitAndLock(10)、lk.Commit()，开启 transaction 后加 doing 锁，Commit将锁切换为 finished
	5、lk.Release()，锁主动释放
*/

type LockStatus int

var (
	ErrLock          = errors.New("lock error")
	ErrTransFinished = errors.New("transaction finished")
	ErrNoTransaction = errors.New("no transaction")
	ErrCommit        = errors.New("commit error")
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

func New(lockId interface{}, client *rds.Pool) (r *Lock, err error) {
	var key string

	if v, ok := lockId.(string); ok {
		if len(v) == 0 {
			err = fmt.Errorf("lockId is empty string")
			return
		}
		key = v
	} else {
		if lockId == nil {
			err = fmt.Errorf("lockId is empty object")
			return
		}
		var k []byte
		if k, err = util.GobEncode(lockId); err != nil {
			return
		}
		key = uuid.NewV3(uuid.Nil, string(k)).String()
	}

	r = &Lock{key: fmt.Sprintf("LOCK.%s", key), client: client}
	r.refreshLockStatus()
	return
}

func (this *Lock) refreshLockStatus() (err error) {
	conn := this.client.Get()
	defer conn.Close()

	var v int
	v, err = rds.Int(conn.Do("GET", this.key))
	if redis.HasError(err) {
		return
	}
	err = nil
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

// 阻塞业务进程并直到锁状态为 normal 时释放阻塞
// 当开启 transaction 时，锁状态为 Finished 时抛出错误
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

// 加锁，需指定锁有效时间
// 当开启 transaction 时，将锁状态设置为 Doing
func (this *Lock) Lock(aliveSeconds int64) (err error) {
	conn := this.client.Get()
	defer conn.Close()

	status := LockStatusFinished
	if this.transaction {
		status = LockStatusDoing
	}
	defer this.refreshLockStatus()

	var res interface{}
	res, err = conn.Do("SET", this.key, int(status), "EX", aliveSeconds, "NX")
	if redis.HasError(err) {
		return
	}
	if res == nil {
		err = ErrLock
	}
	return
}

// 阻塞业务进程并直到加锁成功，需指定锁有效时间
func (this *Lock) WaitAndLock(aliveSeconds int64) (err error) {
	for {
		if err = this.Wait(); err != nil {
			return
		}
		err = this.Lock(aliveSeconds)
		if err == nil || err != ErrLock {
			break
		}
	}
	return
}

// 仅开启 transaction 时有效，将锁状态切换为 Finished
func (this *Lock) Commit() (err error) {
	if this.transaction == false {
		return ErrNoTransaction
	}
	defer this.refreshLockStatus()

	conn := this.client.Get()
	defer conn.Close()

	var aliveSeconds int64
	aliveSeconds, err = rds.Int64(conn.Do("TTL", this.key))
	if redis.HasError(err) {
		return
	}
	var res interface{}
	res, err = conn.Do("SET", this.key, int(LockStatusFinished), "EX", aliveSeconds, "XX")
	if redis.HasError(err) {
		return
	}
	if res == nil {
		err = ErrCommit
	}
	return
}

// 主动释放锁
func (this *Lock) Release() (err error) {
	conn := this.client.Get()
	defer this.refreshLockStatus()
	defer conn.Close()
	_, err = conn.Do("DEL", this.key)
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
