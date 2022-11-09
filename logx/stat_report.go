// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/sysx/host"
)

var (
	clusterName = proc.Env("CLUSTER_NAME")
)

func InfoReport(kv cst.KV) {
	if len(clusterName) > 0 {
		kv["cluster"] = clusterName
	}
	kv["host"] = host.Hostname()
	InfoKV(kv)
}
