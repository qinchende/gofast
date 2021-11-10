package mid

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/jwtx"
	"github.com/qinchende/gofast/logx"
	"net/http"
	"net/http/httputil"
)

const (
	jwtAudience    = "aud"
	jwtExpire      = "exp"
	jwtId          = "jti"
	jwtIssueAt     = "iat"
	jwtIssuer      = "iss"
	jwtNotBefore   = "nbf"
	jwtSubject     = "sub"
	noDetailReason = "no detail reason"
)

var (
	errInvalidToken = errors.New("invalid auth token")
	errNoClaims     = errors.New("no auth params")
)

func JwtAuthHandler(secret string) fst.CtxHandler {
	jwtParser := jwtx.NewTokenParser()

	return func(ctx *fst.Context) {
		//w := ctx.GFResponse
		r := ctx.ReqRaw

		tok, err := jwtParser.ParseToken(r, secret, secret)
		if err != nil {
			unauthorizedPanic(r, err)
			return
		}

		if !tok.Valid {
			unauthorizedPanic(r, errInvalidToken)
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			unauthorizedPanic(r, errNoClaims)
			return
		}

		// 上下文中加入一些token
		rc := r.Context()
		for k, v := range claims {
			switch k {
			case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
				// ignore the standard claims
			default:
				rc = context.WithValue(rc, k, v)
			}
		}
	}
}

//
//func JwtAuthorize(secret string) http.HandlerFunc {
//	jwtParser := jwtx.NewTokenParser()
//
//	return func(w *fst.GFResponse, r *http.Request) {
//		tok, err := jwtParser.ParseToken(r, secret, secret)
//		if err != nil {
//			unauthorized(w, r, err)
//			return
//		}
//
//		if !tok.Valid {
//			unauthorized(w, r, errInvalidToken)
//			return
//		}
//
//		claims, ok := tok.Claims.(jwt.MapClaims)
//		if !ok {
//			unauthorized(w, r, errNoClaims)
//			return
//		}
//
//		// 上下文中加入一些token
//		ctx := r.Context()
//		for k, v := range claims {
//			switch k {
//			case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
//				// ignore the standard claims
//			default:
//				ctx = context.WithValue(ctx, k, v)
//			}
//		}
//	}
//}

func detailAuthLog(r *http.Request, reason string) {
	// discard dump error, only for debug purpose
	details, _ := httputil.DumpRequest(r, true)
	logx.Errorf("authorize failed: %s\n=> %+v", reason, string(details))
}

//func unauthorized(w *fst.GFResponse, r *http.Request, err error) {
//	if err != nil {
//		detailAuthLog(r, err.Error())
//	} else {
//		detailAuthLog(r, noDetailReason)
//	}
//
//	w.ErrorF("Authorize failed, rejected with code %d", http.StatusUnauthorized)
//	w.ResWrap.WriteHeader(http.StatusUnauthorized)
//	w.AbortFit()
//}

func unauthorizedPanic(r *http.Request, err error) {
	detailAuthLog(r, err.Error())
	fst.RaisePanicErr(err)
}
