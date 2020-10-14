/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-13 13:36:06
# File Name: redis_test.go
# Description:
####################################################################### */

package redis

import (
	"fmt"
	"testing"
	//. "github.com/smartystreets/goconvey/convey"
	//rds "github.com/gomodule/redigo/redis"
)

func TestEncode(t *testing.T) {
	New("qa01", "192.168.0.230:6379")

	conn := Get("qa01")
	a, _ := conn.Do("SET", "aaa", 10, "EX", 10, "NX")
	if a == nil {
		fmt.Println("aaa")
	}
	fmt.Println("bbb")

	/*
		Convey("TestEncode", t, func() {
			Convey("TestDecode err should return nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("TestDecode result should equal rawStr", func() {
				So(decStr, ShouldEqual, rawStr)
			})
		})
	*/
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
