# Redis

基于redigo封装的Redis库

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ant-libs-go/redis?status.png)](http://godoc.org/github.com/ant-libs-go/redis)
[![Go Report Card](https://goreportcard.com/badge/github.com/ant-libs-go/redis)](https://goreportcard.com/report/github.com/ant-libs-go/redis)

# 特性

* 简化Redis实例初始化流程，基于配置自动对Redis进行初始化
* 支持连接池、多实例等场景
* 基于Redis实现的分布式锁lock，支持加锁、等待锁释放、等待锁释放并加锁、基于锁实现幂等 等能力
* 基于Redis实现延迟定时器timer，支持业务中触发延迟回调
* 基于Redis和timer实现Storer，维护Mem、Redis、MySQL三者之间的数据同步关系(暂无法保证特殊情况下数据一致性)，有助于基于对象开发时维护对象的同步

# 快速开始

* toml 配置文件
    ```
    [redis.default]
        addr = "127.0.0.1:6379"
        pawd = "123456"
    [redis.stats]
        addr = "127.0.0.1:6379"
        pawd = "123456"
    ```

* 使用方法

	```golang
    // 初始化config包，参考config模块
    code...

    // 验证Redis实例的配置正确性与连通性。非必须
    if err = redisp.Valid(); err != nil {
        fmt.Printf("redis error: %s\n", err)
        os.Exit(-1)
    }

    // 如下方式可以直接使用Redis实例
    conn := redis.Client("default")
    defer conn.Close()

    conn.Do("GET", "key")

    // lock、timer、storer 的使用方法见代码
    ```

# 高级用法

* 分布式锁、延迟定时器、Storer 使用方法参考代码
