package httpx

import (
	"github.com/qinchende/gofast/cst"
	"net/http"
)

type Client struct {
	http.Client
}

type RequestPet struct {
	ProxyUrl   string    // 代理服务器地址
	Method     string    // GET or POST
	Url        string    // http(s)地址
	Headers    cst.WebKV // 请求头
	QueryArgs  cst.KV    // url上的参数
	BodyArgs   cst.KV    // body带的参数
	BodyFormat int8      // body数据的格式，比如 json|url-encoding|xml
}

var myClient = &Client{}                            // 框架默认client对象
var myTransports = make(map[string]*http.Transport) // 当前所有的代理实例

func PetClient(pet *RequestPet) *Client {
	if pet.ProxyUrl == "" {
		return myClient
	}
	t := &Client{}
	t.Transport = getTransport(pet.ProxyUrl)
	return t
}

func (cli *Client) Do(req *http.Request) (*http.Response, error) {
	return cli.Do(req)
}

func (cli *Client) DoGetKV(req *http.Request) (cst.KV, error) {
	return ParseJsonResponse(cli.Do(req))
}
