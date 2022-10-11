package cst

type SdxConfig struct {
	// sdx 实现模块的配置参数
	NeedSysCheck     bool   `v:"def=true"`                        // 是否启动CPU使用情况的定时检查工作
	NeedSysPrint     bool   `v:"def=true"`                        // 定时打印系统检查日志
	EnableTimeout    bool   `v:"def=true"`                        // 默认启动超时拦截
	DefTimeoutMS     int64  `v:"def=3000"`                        // 每次请求的超时时间（单位：毫秒）
	MaxContentLength int64  `v:"def=33554432"`                    // 最大请求字节数，32MB（33554432），传0不限制
	MaxConnections   int32  `v:"def=1000000,range=[0:100000000]"` // 最大同时请求数，默认100万同时进入，传0不限制
	JwtSecret        string `v:""`                                // JWT认证的秘钥
}
