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
)

func main() {
	app, home := fst.CreateServer(&fst.FConfig{
		PrintRouteTrees: true,
		CurrMode:        "debug",
	})
	var handler = func(str string) func(c *fst.Context) {
		return func(c *fst.Context) {
			log.Println(str)
		}
	}

	// 根路由
	home.Method("GET", "/user/:name", handler("home handler 1"))
	home.Before(handler("home before 1")).After(handler("home after 1"))
	home.Get("/user/:name/age", handler("home handler 2"))

	// 分组路由
	adm := home.AddGroup("/admin/cd")
	adm.Before(handler("before admin 1")).After(handler("after admin 1"), handler("after admin 2"))
	ul := adm.Get("/user_list", handler("user_list handler"))
	ul.After(handler("user_list after 1"))

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
2020/12/31 13:50:27 Listening and serving HTTP on 127.0.0.1:8099
++++++++++The route tree:

(GET)
└── /                                                            [false-ua]
    ├── user/                                                    [false-]
    │   └── :name                                                [true-/]
    │       └── /age                                             [true-]
    └── admin/cd/user_list                                       [true-]

++++++++++THE END.
```

浏览器输入网址访问地址：127.0.0.1:8099/admin/cd/user_list，日志会输出：

```
2020/12/31 13:50:31 home before 1
2020/12/31 13:50:31 before admin 1
2020/12/31 13:50:31 user_list handler
2020/12/31 13:50:31 user_list after 1
2020/12/31 13:50:31 after admin 1
2020/12/31 13:50:31 after admin 2
2020/12/31 13:50:31 home after 1
```

（其它介绍陆续补充...）