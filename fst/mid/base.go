package mid

import (
	"net/http"
)

type funcServeHTTP func(w http.ResponseWriter, r *http.Request)

type FitHelper struct {
	nextHandler funcServeHTTP
}

func (fh *FitHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.nextHandler(w, r)
}
