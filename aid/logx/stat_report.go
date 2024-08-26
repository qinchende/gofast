// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/aid/proc"
	"github.com/qinchende/gofast/aid/sysx/host"
	"github.com/qinchende/gofast/core/cst"
)

var (
	clusterName = proc.Env("CLUSTER_NAME")
)

func InfoReport(kv cst.KV) {
	if len(clusterName) > 0 {
		kv["cluster"] = clusterName
	}
	kv["host"] = host.Hostname()
	StatKV(kv)
}
