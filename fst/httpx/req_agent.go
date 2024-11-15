package httpx

import (
	"github.com/qinchende/gofast/aid/logx"
	"net/http"
	"net/url"
)

var myTransports = make(map[string]*http.Transport) // 当前所有的代理实例

func getTransport(proxyUrl string) *http.Transport {
	if proxyUrl == "" {
		return nil
	}

	if myTransports[proxyUrl] != nil {
		return myTransports[proxyUrl]
	}

	netURL, err := url.Parse(proxyUrl)
	if err != nil {
		logx.Debug().Msg(err.Error())
		return nil
	}

	trans := &http.Transport{
		Proxy: http.ProxyURL(netURL),
	}
	myTransports[proxyUrl] = trans
	return trans
}
