package mid

//var midTimeoutBody = "<html><head><title>Timeout</title></head><body><h1>Timeout</h1></body></html>"
const (
	midTimeoutBody  = "<html>Timeout!</html>"      // 超时
	midFusingBody   = "<html>Fusing!</html>"       // 熔断
	midSheddingBody = "<html>LoadShedding!</html>" // 降载
)
