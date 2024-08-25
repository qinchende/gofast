# [Preview] 当前持续完善中，不可生产中使用 [Preview]

# GoFast Micro-Service Framework

GoFast是一个用Go语言实现的微服务开发框架。他的产生源于目前流行的gin、go-zero、fastify等众多开源框架；同时结合了作者多年的开发实践经验，很多模块的实现方式都是作者首创；当然也免不了不断借鉴社区中优秀的设计理念。我们的目标是简洁高效易上手，在封装大量特性的同时又不失灵活性。希望你能喜欢GoFast。

GoFast的微服务：虽然也提供现成的微服务治理能力，但我们想强调的是，GoFast更专注于帮助框架使用者清晰的开发业务逻辑。大型项目上框架应该弱化微服务治理，将这一部分特性交给istio处理。

## Installation

To install GoFast package, you need to install Go and set your Go workspace first.

The first need [Go](https://golang.org/) installed (**version 1.20+ is required**), then you can use the below Go command to install GoFast.

```sh
$ go get -u github.com/qinchende/gofast
```

## Quick start

```sh
# assume the following codes in example.go file
$ cat main.go
```

```go
import (
    "fmt"
    "github.com/qinchende/gofast/aid/conf"
    "github.com/qinchende/gofast/core/cst"
    "github.com/qinchende/gofast/core/logx"
    "github.com/qinchende/gofast/fst"
    "github.com/qinchende/gofast/sdx"
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
        c.FaiData(cst.KV{"data": str})
    }
}

func main() {
    appCfg := &fst.GfConfig{
        RunMode: "debug",
    }
    _ = conf.LoadConfigFromJsonBytes(&appCfg.WebConfig, []byte("{}"))
    _ = conf.LoadConfigFromJsonBytes(&appCfg.LogConfig, []byte("{}"))
    _ = conf.LoadConfigFromJsonBytes(&appCfg.SdxConfig, []byte("{}"))
    appCfg.WebConfig.PrintRouteTrees = true
    appCfg.LogConfig.LogLevel = "debug"
    appCfg.LogConfig.FileFolder = "_logs_"

	app := fst.CreateServer(appCfg)
	logx.MustSetup(&appCfg.LogConfig)

	// 拦截器，微服务治理 ++++++++++++++++++++++++++++++++++++++
	app.UseHttpHandler(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("app enter before 1")
			next(w, r)
			log.Println("app enter after 1")
		}
	})
	//app.UseGlobal(sdx.SuperHandlers)
	app.Before(sdx.PmsParser) // 解析请求参数，构造 ctx.Pms

	// 根路由
	app.Reg404(func(c *fst.Context) {
		c.Json(http.StatusNotFound, "404-Can't find the path.")
	})
	app.Reg405(func(c *fst.Context) {
		c.Json(http.StatusMethodNotAllowed, "405-Method not allowed.")
	})

	// ++ (用这种方法可以模拟中间件需要上下文变量的场景)
	app.Before(func(c *fst.Context) {
		c.Set("nowTime", time.Now())
		//time.Sleep(3 * time.Second)
	})
	app.After(func(c *fst.Context) {
		// 处理后获取消耗时间
		val, exist := c.Get("nowTime")
		if exist {
			costTime := time.Since(val.(time.Time)) / time.Millisecond
			fmt.Printf("The request cost %dms", costTime)
		}
	})
	// ++ end

	// curl -H "Content-Type: application/json" -X POST  --data '{"data":"bmc","nick":"yes"}'
	//      http://127.0.0.1:8099/root?first=yang\&last=lmx
	// curl -H "Content-Type: application/x-www-form-urlencoded" -X POST  --data "data=bmc&nick=yes"
	//      http://127.0.0.1:8099/root?ids[a]=1234\&ids[b]=hello\&first=yang\&last=lmx
	type MyData struct {
		Data string `v:"must,len=[1:16]"`
		Nick string `v:"must,len=[1:32]"`
	}
	app.Post("/root", func(c *fst.Context) {
		var myData MyData
		c.PanicIfErr(c.BindAndValid(&myData), "数据解析错误")
		log.Printf("%v %+v %#v\n", myData, myData, myData)

		ids, _ := c.GetString("ids")
		firstname := c.GetStringDef("first", "Guest")
		lastname := c.GetStringMust("last")

		message := c.GetStringMust("data")
		nick := c.GetStringDef("nick", "anonymous")

		//names := c.PostFormMap("names")
		c.SucData(cst.KV{
			"message": message,
			"nick":    nick,
			"first":   firstname,
			"last":    lastname,
			"ids":     ids,
			//"data":    myData,
		})
		//c.String(http.StatusOK, fmt.Sprintf("file uploaded!"))
		//c.JSON(http.StatusOK, myData)
	})

	//app.Post("/root", handler("root"))
	app.Before(handler("before root")).After(handler("after root"))

	// 分组路由1
	adm := app.Group("/admin")
	adm.After(handler("after group admin")).Before(handler("before group admin"))

	tst := adm.Get("/sdx", handlerRender("handle sdx"))
	// 添加路由处理事件
	tst.Before(handler("before tst_url"))
	tst.After(handler("after tst_url"))
	tst.BeforeSend(handler("beforeSend tst_url"))
	tst.AfterSend(handler("afterSend tst_url"))

	// 分组路由2
	adm2 := app.Group("/admin2").Before(handler("before admin2"))
	adm2.Get("/zht", handler("zht")).After(handler("after zht"))

	adm22 := adm2.Group("/group2").Before(handler("before group2"))
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
	app.Listen("127.0.0.1:8099")
}

```

```sh
# run main.go and visit website 127.0.0.1:8099 on browser
$ go run main.go
```

在控制台启动后台Web服务器之后，你会看到底层的路由树构造结果：

```
API server listening at: 127.0.0.1:59863
2024/08/22 15:22:21 Current GoFast version: v0.5.0.20240306.2
2024/08/22 15:22:21 App OnReady Call.
2024/08/22 15:22:21 Listening and serving HTTP on 127.0.0.1:8099
[08-22 15:22:21][debug]: POST   /root                     main.main.func6 (1 hds)
[08-22 15:22:21][debug]: GET    /admin/sdx                main.init.func2.1 (1 hds)
[08-22 15:22:21][debug]: GET    /admin2/zht               main.init.func1.1 (1 hds)
[08-22 15:22:21][debug]: GET    /admin2/group2/lmx        main.init.func1.1 (1 hds)

+++++++++++++++The route tree:
(GET)
└── /admin                                                       [0-/2]
    ├── /sdx                                                     [1]
    └── 2/                                                       [0-zg]
        ├── zht                                                  [1]
        └── group2/lmx                                           [1]
(POST)
└── /root                                                        [1]
++++++++++++++++++++++++++++++
```

浏览器输入网址访问地址：`127.0.0.1:8099/admin/sdx`，日志会输出：

```
2024/08/22 15:41:07 app enter before 1
2024/08/22 15:41:07 before root
2024/08/22 15:41:07 before group admin
2024/08/22 15:41:07 before tst_url
2024/08/22 15:41:07 handle sdx
2024/08/22 15:41:07 beforeSend tst_url
2024/08/22 15:41:07 afterSend tst_url
2024/08/22 15:41:07 after tst_url
2024/08/22 15:41:07 after group admin
2024/08/22 15:41:07 after root
2024/08/22 15:41:07 app enter after 1
The request cost 0ms
[GET] /admin/sdx (127.0.0.1/08-23 15:41:07) [200/63/0]
  B: {}
  P: {"nowTime":"2024-08-22T15:41:07+08:00"}
  R: {"status":"fai","code":0,"msg":"","data":{"data":"handle sdx"}}
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


#### Sample ORM

本框架集成了自研ORM框架，简单够用性能好，你用了就会爱上她。

```go
func QueryUserCache(c *fst.Context) {
    userId := c.GetIntMust("user_id")

	var ccUser acc.SysUser
	ct := cf.DDemo.QueryPrimaryCache(&ccUser, userId)

	userId += 1
	var ccUser2 acc.SysUser
	ct = cf.DDemo.QueryPrimaryCache(&ccUser2, userId)

	//kvs := make(cst.KV)
	//ct := cf.DDemo.QuerySqlRow(&kvs, "select * from sys_user where id=?;", userId)

	if ct > 0 {
		c.SucData(cst.KV{"name1": ccUser.Name, "name2": ccUser2.Name})
		// c.SucData(kvs)
	} else {
		c.FaiMsg("can't find the record")
	}
}
```


## benchmark


## License

This project is licensed under the terms of the MIT license.


（其它介绍陆续补充...）