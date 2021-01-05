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

var handler = func(str string) func(ctx *fst.Context) {
	return func(ctx *fst.Context) {
		log.Println(str)
	}
}

func main() {
	app, home := fst.CreateServer(&fst.AppConfig{
		PrintRouteTrees: true,
		RunMode:         "debug",
	})

	// 应用级事件
	app.OnReady(func(fast *fst.GoFast) {
		log.Println("App OnReady Call.")
	})
	app.OnClose(func(fast *fst.GoFast) {
		log.Println("App OnClose Call.")
	})

	// 根路由
	home.NoRoute(func(ctx *fst.Context) {
		ctx.JSON(http.StatusNotFound, "404-Can't find the path.")
	})
	home.NoMethod(func(ctx *fst.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, "405-Method not allowed.")
	})

	home.Post("/root", handler("root"))
	home.Before(handler("before root")).After(handler("after root"))

	// 分组路由1
	adm := home.AddGroup("/admin").After(handler("after admin"))
	adm.Get("/chende", handler("chende")).After(handler("after chende"))

	// 分组路由2
	adm2 := home.AddGroup("/admin2").Before(handler("before admin2"))
	adm2.Get("/zht", handler("zht")).After(handler("after zht"))

	adm22 := adm2.AddGroup("/group2").Before(handler("before group2"))
	adm22.Get("/lmx", handler("lmx")).Before(handler("before lmx"))

	// 开始监听接收请求
	log.Printf("Listening and serving HTTP on %s\n", "127.0.0.1:8099")
	log.Fatal(app.Listen("127.0.0.1:8099"))
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
2021/01/04 01:20:38 before root
2021/01/04 01:20:38 chende
2021/01/04 01:20:38 after chende
2021/01/04 01:20:38 after admin
2021/01/04 01:20:38 after root
```

（其它介绍陆续补充...）