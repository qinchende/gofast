# [Preview] 当前持续完善中，不可生产中使用 [Preview]

# GoFast Micro-Service Framework

GoFast是一个用Go语言实现的微服务开发框架。他的产生源于目前流行的gin、go-zero、fastify等众多开源框架；同时结合了作者多年的开发实践经验，很多模块的实现方式都是作者首创；当然也免不了不断借鉴社区中优秀的设计理念。我们的目标是简洁高效易上手，在封装大量特性的同时又不失灵活性。希望你能喜欢GoFast。

更多了解：[GoFast的实现细节](https://chende.ren/tags/gofast-intr/)

GoFast的微服务：虽然也提供现成的微服务治理能力，但我们想强调的是，GoFast更专注于帮助框架使用者清晰的开发业务逻辑。大型项目上框架应该弱化微服务治理，将这一部分特性交给istio处理。

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
	"fmt"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fstx"
	"log"
	"net/http"
	"time"
)

var handler = func(str string) func(c *fst.Context) {
	return func(c *fst.Context) {
		log.Println(str)
	}
}

var handlerRender = func(str string) func(c *fst.Context) {
	return func(c *fst.Context) {
		log.Println(str)
		c.Fai(200, fst.KV{"data": str})
	}
}

func main() {
	app := fst.CreateServer(&fst.AppConfig{
		PrintRouteTrees:        true,
		HandleMethodNotAllowed: true,
		RunMode:                "debug",
		//FitMaxReqCount:         1,
		//FitMaxReqContentLen:    10 * 1024,
	})

	// 拦截器，微服务治理 ++++++++++++++++++++++++++++++++++++++
	app.Fits(fstx.AddDefaultFits)
	app.Fit(func(w *fst.GFResponse, r *http.Request) {
		log.Println("app fit before 1")
		w.NextFit(r)
		log.Println("app fit after 1")
	})

	// 根路由
	app.NoRoute(func(ctx *fst.Context) {
		ctx.JSON(http.StatusNotFound, "404-Can't find the path.")
	})
	app.NoMethod(func(ctx *fst.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, "405-Method not allowed.")
	})

	// ++ (用这种方法可以模拟中间件需要上下文变量的场景)
	app.Before(func(ctx *fst.Context) {
		ctx.Set("nowTime", time.Now())
		time.Sleep(3 * time.Second)
	})
	app.After(func(ctx *fst.Context) {
		// 处理后获取消耗时间
		val, exist := ctx.Get("nowTime")
		if exist {
			costTime := time.Since(val.(time.Time)) / time.Millisecond
			fmt.Printf("The request cost %dms", costTime)
		}
	})
	// ++ end

	// curl -H "Content-Type: application/json" -X POST  --data '{"data":"bmc","nick":"yes"}' http://127.0.0.1:8099/root?first=yang\&last=lmx
	// curl -H "Content-Type: application/x-www-form-urlencoded" -X POST  --data '{"data":"bmc","nick":"yes"}' http://127.0.0.1:8099/root?first=yang\&last=lmx
	// curl -H "Content-Type: application/x-www-form-urlencoded" -X POST  --data "data=bmc&nick=yes" http://127.0.0.1:8099/root?ids[a]=1234\&ids[b]=hello\&first=yang\&last=lmx
	type MyData struct {
		Data string `json:"data"`
		Nick string `json:"nick"`
	}
	app.Post("/root", func(ctx *fst.Context) {
		myData := MyData{}
		_ = ctx.ShouldBindBodyWith(&myData, binding.JSON)
		log.Printf("%v %+v %#v\n", myData, myData, myData)

		myDataT := MyData{}
		_ = ctx.ShouldBindBodyWith(&myDataT, binding.JSON)
		log.Printf("%v %+v %#v\n", myDataT, myDataT, myDataT)

		ids := ctx.QueryMap("ids")
		firstname := ctx.DefaultQuery("first", "Guest")
		lastname := ctx.Query("last") // shortcut for ctx.Request.URL.Query().Get("lastname")

		message := ctx.PostForm("data")
		nick := ctx.DefaultPostForm("nick", "anonymous")

		//names := ctx.PostFormMap("names")
		ctx.Suc(fst.KV{
			"message": message,
			"nick":    nick,
			"first":   firstname,
			"last":    lastname,
			"ids":     ids,
			//"data":    myData,
		})
		//ctx.String(http.StatusOK, fmt.Sprintf("file uploaded!"))
		//ctx.JSON(http.StatusOK, myData)
	})

	app.Post("/root", handler("root"))
	app.Before(handler("before root")).After(handler("after root"))

	// 分组路由1
	adm := app.AddGroup("/admin")
	adm.After(handler("after group admin")).Before(handler("before group admin"))

	tst := adm.Get("/chende", handlerRender("handle chende"))
	// 添加路由处理事件
	tst.Before(handler("before tst_url"))
	tst.After(handler("after tst_url"))
	tst.PreSend(handler("preSend tst_url"))
	tst.AfterSend(handler("afterSend tst_url"))

	// 分组路由2
	adm2 := app.AddGroup("/admin2").Before(handler("before admin2"))
	adm2.Get("/zht", handler("zht")).After(handler("after zht"))

	adm22 := adm2.AddGroup("/group2").Before(handler("before group2"))
	adm22.Get("/lmx", handler("lmx")).Before(handler("before lmx"))

	// 应用级事件
	app.OnReady(func(fast *fst.GoFast) {
		log.Println("App OnReady Call.")
		log.Printf("Listening and serving HTTP on %s\n", "127.0.0.1:8099")
	})
	app.OnClose(func(fast *fst.GoFast) {
		log.Println("App OnClose Call.")
	})
	// 开始监听接收请求
	_ = app.Listen("127.0.0.1:8099")
}
```

```sh
# run example.go and visit website 127.0.0.1:8099 on browser
$ go run example.go
```

在控制台启动后台Web服务器之后，你会看到底层的路由树构造结果：

```
[GoFast-debug] POST   /root                     --> main.glob..func1.1 (1 handlers)
[GoFast-debug] GET    /admin/chende             --> main.glob..func2.1 (1 handlers)
[GoFast-debug] GET    /admin2/zht               --> main.glob..func1.1 (1 handlers)
[GoFast-debug] GET    /admin2/group2/lmx        --> main.glob..func1.1 (1 handlers)
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
2021/02/20 02:56:00 App OnReady Call.
2021/02/20 02:56:00 Listening and serving HTTP on 127.0.0.1:8099
```

浏览器输入网址访问地址：`127.0.0.1:8099/admin/chende`，日志会输出：

```
2021/02/20 02:57:17 app fit before 1
2021/02/20 02:57:20 before root
2021/02/20 02:57:20 before group admin
2021/02/20 02:57:20 before tst_url
2021/02/20 02:57:20 handle chende
2021/02/20 02:57:20 preSend tst_url
2021/02/20 02:57:20 afterSend tst_url
2021/02/20 02:57:20 after tst_url
2021/02/20 02:57:20 after group admin
2021/02/20 02:57:20 after root
2021/02/20 02:57:20 app fit after 1
The request cost 3000ms
[GET] /admin/chende (127.0.0.1/02-20 02:57:20) 200/24 [3000]
  B:  C: 
  P: 
  R: 
  E: 
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


## License

This project is licensed under the terms of the MIT license.


（其它介绍陆续补充...）