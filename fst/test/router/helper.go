package router

import (
	"gofast/fst"
)

func init() {
	initGoFastServer()
}

var gftApp *fst.GoFast

func initGoFastServer() {
	// 新建Server
	gftApp = fst.CreateServer(&fst.AppConfig{
		RunMode: fst.ProductMode,
	})

	addMiddleware()
	addRouters()
	gftApp.ReadyToListen()
}

func addMiddleware() {

}

func addRouters() {

}

func getNewServer() *fst.GoFast {
	app := fst.CreateServer(&fst.AppConfig{
		RunMode: fst.ProductMode,
	})
	return app
}
