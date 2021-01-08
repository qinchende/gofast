# GoFast Web Framework

GoFast是一个用Go语言实现的高效Web开发框架。他的产生源于目前流行的Gin框架，同时也不断借鉴社区中优秀的设计理念，目标是打造一个更高效的Web开发框架，在提供更丰富封装特性的同时又不失灵活性。希望你能喜欢GoFast。

更多了解：[GoFast的实现细节](https://chende.ren/tags/gofast-intr/)

## Installation

To install GoFast package, you need to install Go and set your Go workspace first.

The first need [Go](https://golang.org/) installed (**version 1.15+ is required**), then you can use the below Go command to install GoFast.

```sh
$ go get -u github.com/qinchende/gofast
```

## Quick start

```sh
# assume the following codes in example.go file
$ cat example.go
```

```go
package main

import (
	"gofast/fst"
	"log"
	"net/http"
)

var handler = func(str string) func(c *fst.Context) {
	return func(c *fst.Context) {
		log.Println(str)
	}
}

var handlerRender = func(str string) func(c *fst.Context) {
	return func(c *fst.Context) {
		log.Println(str)
		c.JSON(200, fst.KV{"data": str})
	}
}

func main() {
	gft := fst.CreateServer(&fst.AppConfig{
		PrintRouteTrees: true,
		RunMode:         "debug",
	})

	// 根路由
	gft.NoRoute(func(ctx *fst.Context) {
		ctx.JSON(http.StatusNotFound, "404-Can't find the path.")
	})
	gft.NoMethod(func(ctx *fst.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, "405-Method not allowed.")
	})

	gft.Post("/root", handler("root"))
	gft.Before(handler("before root")).After(handler("after root"))

	// 分组路由1
	adm := gft.AddGroup("/admin")
	adm.After(handler("after group admin")).Before(handler("before group admin"))

	tst := adm.Get("/chende", handlerRender("handle chende"))
	// 添加路由处理事件
	tst.Before(handler("before tst_url"))
	tst.After(handler("after tst_url"))
	tst.PreSend(handler("preSend tst_url"))
	tst.AfterSend(handler("afterSend tst_url"))

	// 分组路由2
	adm2 := gft.AddGroup("/admin2").Before(handler("before admin2"))
	adm2.Get("/zht", handler("zht")).After(handler("after zht"))

	adm22 := adm2.AddGroup("/group2").Before(handler("before group2"))
	adm22.Get("/lmx", handler("lmx")).Before(handler("before lmx"))

	// 应用级事件
	gft.OnReady(func(fast *fst.GoFast) {
		log.Println("App OnReady Call.")
		log.Printf("Listening and serving HTTP on %s\n", "127.0.0.1:8099")
	})
	gft.OnClose(func(fast *fst.GoFast) {
		log.Println("App OnClose Call.")
	})
	// 开始监听接收请求
	_ = gft.Listen("127.0.0.1:8099")
}

```

```sh
# run example.go and visit 127.0.0.1:8099/admin/cd/user_list on browser
$ go run example.go
```

在控制台启动后台Web服务器之后，你会看到底层的路由树构造结果：

```
++++++++++The route tree:

(GET)
└── /admin                                                       [false-/2]
    ├── /chende                                                  [true-]
    └── 2/                                                       [false-zg]
        ├── zht                                                  [true-]
        └── group2/lmx                                           [true-]
(POST)
└── /root                                                        [true-]

++++++++++THE END.
2021/01/04 01:18:24 Listening and serving HTTP on 127.0.0.1:8099
```

浏览器输入网址访问地址：`127.0.0.1:8099/admin/chende`，日志会输出：

```
2021/01/06 17:35:40 before root
2021/01/06 17:35:40 before group admin
2021/01/06 17:35:40 before tst_url
2021/01/06 17:35:40 handle chende
2021/01/06 17:35:40 preSend tst_url
2021/01/06 17:35:40 afterSend tst_url
2021/01/06 17:35:40 after tst_url
2021/01/06 17:35:40 after group admin
2021/01/06 17:35:40 after root
```

## Core feature

#### Like gin feature
GoFast目前复用了Gin的很多特性，除特别说明之外，使用方式一样。

#### Server Handlers

应用启动之后，在开始监听端口之后调用OnReady事件，应用关闭退出之前调用OnClose事件
```go
app.OnReady(func(fast *fst.GoFast) {
	log.Println("App OnReady Call.")
	log.Printf("Listening and serving HTTP on %s\n", "127.0.0.1:8099")
})

app.OnClose(func(fast *fst.GoFast) {
	log.Println("App OnClose Call.")
})

```

#### Router Handlers

分组或路由项事件是一样的，现在支持下面四个，以后慢慢扩展和调整
```go
tst.Before(handler("before tst_url"))
tst.After(handler("after tst_url"))
tst.PreSend(handler("preSend tst_url"))
tst.AfterSend(handler("afterSend tst_url"))
```

## benchmark
> Gin是非常优秀的框架，GoFast主要与Gin做性能上的比较，力争进一步提高处理能力。

`gofast/fst/test/performance` 目录中有一套简单的基准测试代码，大家可以在自己的机器上试试看。

在该目录下运行命令：`go test -bench=. -benchtime=10s`

我笔记本的运行结果如下：
```go
goos: windows
goarch: amd64
pkg: gofast/fst/test/performance
BenchmarkGinWebRouter-2         70346450        165 ns/op       0 B/op      0 allocs/op
BenchmarkGoFastWebRouter-2      133092091       89.9 ns/op      0 B/op      0 allocs/op
PASS
ok      gofast/fst/test/performance     33.007s

```
经过多次测试发现相比Gin，GoFast处理能力提升60%以上。

测试过程中有些参数可以调整，对实际测试结果影响比较大的是`middlewareNum`参数，说明Gin中嵌套执行中间件函数的方式性能不好。类似递归调用，层级太多影响性能。
```go
// router_helper.go
var routersLevel = 1                 // 路由数量的基数，实际值=routersSum
var routersSum = 1000 * routersLevel // 1000*routersNum
var middlewareNum = 10               // 中间件函数的数量
var reqPoolSize = routersSum         // 内置请求对象，用于模拟发起的不同Router请求
var differentReqNum = 1              // 用多少个不同路由的请求来测试
```

（其它介绍陆续补充...）