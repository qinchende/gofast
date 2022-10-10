package httpx

import (
	"github.com/qinchende/gofast/logx"
	"net/http"
	"net/url"
)

func getTransport(proxyUrl string) *http.Transport {
	if proxyUrl == "" {
		return nil
	}

	if myTransports[proxyUrl] != nil {
		return myTransports[proxyUrl]
	}

	netURL, err := url.Parse(proxyUrl)
	if err != nil {
		logx.Debugs(err)
		return nil
	}

	trans := &http.Transport{
		Proxy: http.ProxyURL(netURL),
	}
	myTransports[proxyUrl] = trans
	return trans
}
