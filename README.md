# redis
基于redigo封装的Redis库，支持连接池、多实例等，且对分布式锁、storer等常用场景进行封装

# 功能
 - 简化Redis实例初始化流程，基于配置自动对Redis进行初始化
 - 支持连接池、多实例等场景
 - 基于Redis实现的分布式锁lock，支持加锁、等待锁释放、等待锁释放并加锁、基于锁实现幂等 等能力
 - 基于Redis实现延迟定时器timer，支持业务中触发延迟回调
 - 基于Redis和timer实现Storer，维护Mem、Redis、MySQL三者之间的数据同步关系(暂无法保证特殊情况下数据一致性)，有助于基于对象开发时维护对象的同步


# 基本使用
 - toml 配置文件
    ```
    [redis.default]
        addr = "127.0.0.1:6379"
        pawd = "123456"
    [redis.stats]
        addr = "127.0.0.1:6379"
        pawd = "123456"
    ```

 - 使用方法
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
