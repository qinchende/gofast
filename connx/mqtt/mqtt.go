package mqtt

type (
	ClientConfig struct {
		ConnStr    string   `v:"must"`
		ConnStrR   string   `v:"must=false"`
		MaxOpen    int      `v:"def=100,range=[1:1000]"`
		MaxIdle    int      `v:"def=100"`
		RedisNodes []string `v:"must=false,len=[10:300]"`
	}
)
