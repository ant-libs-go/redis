/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-13 13:36:06
# File Name: redis_test.go
# Description:
####################################################################### */

package redis

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ant-libs-go/config"
	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/config/parser"
	//. "github.com/smartystreets/goconvey/convey"
	//rds "github.com/gomodule/redigo/redis"
)

var globalCfg *config.Config

func TestMain(m *testing.M) {
	config.New(parser.NewTomlParser(),
		options.WithCfgSource("/tmp/app.toml"),
		options.WithCheckInterval(1))
	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {

	//fmt.Println(Valid())
	for {
		cli := DefaultClient()
		defer cli.Close()
		fmt.Println(cli.Do("SET", "testkey", "12345"))
		time.Sleep(time.Second)
	}

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
