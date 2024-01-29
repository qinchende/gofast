package sdx

import (
	"encoding/base64"
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/randx"
	"github.com/qinchende/gofast/skill/timex"
	"github.com/qinchende/gofast/store/jde"
	"strings"
	"time"
)

func JwtSessBuilder(c *fst.Context) {
	// 不可重复执行 token 检查，Sess构造的过程
	if c.Sess != nil {
		return
	}

	// 每个请求对应的SESSION对象都是新创建的，线程安全。
	ss := new(JwtSession)
	c.Sess = ss
	ss.raw, _ = c.GetString(PmsToken)

	// 请求没有tok，赋予当前请求新的token，同时走后面的逻辑
	if ss.raw == "" {
		ss.createNewToken()
		return
	}

	// 有 tok ，解析出 [payload、hmac]，实际上 token = [payload].[hmac]
	reqPayload, reqHmac := parseJwt(ss.raw)
	if reqPayload == "" {
		ss.createNewToken()
		return
	}

	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确，需要分配新的Token。
	isValid := checkJwt(reqPayload, MySessDB.Secret, reqHmac)
	if !isValid {
		ss.createNewToken()
		return
	}

	ss.payload = reqPayload
	if err := ss.parsePayloadValues(); err != nil {
		c.CarryMsg(err.Error())
		c.AbortFai(110, "Load jwt data error.", nil)
	}

	// 如果过期，利用当前sid重建Session记录。
	if err := ss.checkExpire(); err != nil {
		c.CarryMsg(err.Error())
		c.AbortFai(110, "Load jwt data error.", nil)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//	type Token struct {
//		Raw       string                 // The raw token.  Populated when you Parse a token
//		Method    SigningMethod          // The signing method used or to be used
//		Header    map[string]interface{} // The first segment of the token
//		Claims    Claims                 // The second segment of the token
//		Signature string                 // The third segment of the token.  Populated when you Parse a token
//		Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
//	}
//
//	type StandardClaims struct {
//		Audience  string `json:"aud,omitempty"`
//		ExpiresAt int64  `json:"exp,omitempty"`
//		Id        string `json:"jti,omitempty"`
//		IssuedAt  int64  `json:"iat,omitempty"`
//		Issuer    string `json:"iss,omitempty"`
//		NotBefore int64  `json:"nbf,omitempty"`
//		Subject   string `json:"sub,omitempty"`
//	}
const (
	jwtAudience  = "aud" // 接受者标识
	jwtExpire    = "exp" // 到期时间搓
	jwtId        = "jti" // token id
	jwtIssueAt   = "iat" // 发布时间
	jwtIssuer    = "iss" // 发布者标识
	jwtNotBefore = "nbf" // 开始生效时间
	jwtSubject   = "sub" // 内容主题
)

type JwtSession struct {
	raw     string        // raw token string
	payload string        // content values string
	values  cst.WebKV     // map[string]string
	guid    string        // unique key
	expAt   time.Duration // 在什么时间点过期
	changed bool          // 值是否改变
}

// JwtSession 需要实现 sessionKeeper 所有接口
var _ fst.SessionKeeper = &JwtSession{}
var _JwtSessionInitializer JwtSession

func (ss *JwtSession) Get(key string) (v string, ok bool) {
	v, ok = ss.values[key]
	return
}

func (ss *JwtSession) GetValues() cst.WebKV {
	return ss.values
}

func (ss *JwtSession) Set(key string, val string) {
	if ss.values == nil {
		ss.values = make(cst.WebKV)
	}
	ss.changed = true
	ss.values[key] = val
}

func (ss *JwtSession) SetValues(kvs cst.WebKV) {
	if ss.values == nil {
		ss.values = make(cst.WebKV)
	}
	ss.changed = true
	for k, v := range kvs {
		ss.values[k] = v
	}
}

func (ss *JwtSession) SetUid(uid string) {
	ss.Set(MySessDB.UidField, uid)
}

func (ss *JwtSession) GetUid() (uid string) {
	uid, _ = ss.Get(MySessDB.UidField)
	return
}

func (ss *JwtSession) Save() {
}

func (ss *JwtSession) Del(key string) {
	delete(ss.values, key)
	ss.changed = true
}

// 从当前开始，过多少秒后过期
func (ss *JwtSession) ExpireS(exp uint32) {
	ss.expAt = timex.NowAddSDur(int(exp))
	ss.changed = true
}

func (ss *JwtSession) TokenIsNew() bool {
	return ss.changed
}

func (ss *JwtSession) Token() string {
	if ss.changed {
		ss.raw = ss.buildToken()
	}
	return ss.raw
}

func (ss *JwtSession) Destroy() {
	*ss = _JwtSessionInitializer
}

func (ss *JwtSession) Recreate() {
	ss.createNewToken()
}

// 新生成一个SDX Session对象，生成新的tok
func (ss *JwtSession) createNewToken() {
	ss.guid = genJwtGuid()
	ss.expAt = timex.NowDur() + MySessDB.TTL
	ss.changed = true
}

// crypto
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func parseJwt(tok string) (string, string) {
	dot := strings.IndexByte(tok, '.')
	// 格式明显不对，直接返回空
	if dot <= 0 {
		return "", ""
	}
	return tok[:dot], tok[(dot + 1):]
}

func checkJwt(data, secret, sHmac string) bool {
	md5Val := md5Base64(lang.STB(data), lang.STB(secret))
	return sHmac == md5Val
}

func (ss *JwtSession) parsePayloadValues() error {
	bs, err := base64.RawURLEncoding.DecodeString(ss.payload)
	if err != nil {
		return err
	}
	if ss.values == nil {
		ss.values = make(cst.WebKV)
	}
	if err = jde.DecodeBytes(&ss.values, bs); err != nil {
		return err
	}

	if val, ok := ss.Get(jwtId); !ok || val == "" {
		return errors.New("token must include jwt id")
	} else {
		ss.guid = val
	}

	if val, ok := ss.Get(jwtExpire); !ok {
		return errors.New("token must include jwt expire time")
	} else {
		ss.expAt = time.Duration(lang.ParseIntFast(val)) * time.Second
	}
	return nil
}

func (ss *JwtSession) checkExpire() error {
	now := timex.NowDur()
	diffDur := ss.expAt - now

	// 令牌时间搓无效
	if diffDur <= 0 || diffDur > MySessDB.TTL {
		return errors.New("token expiration")
	}
	// 还剩余不到一半的有效期，需要自动延迟token有效期
	if diffDur < MySessDB.TTL/2 {
		ss.expAt = timex.NowDur() + MySessDB.TTL
		ss.changed = true
	}
	return nil
}

func (ss *JwtSession) buildToken() string {
	ss.Set(jwtId, ss.guid)
	ss.Set(jwtExpire, lang.ToString(timex.ToS(ss.expAt)))
	jsonVal, _ := jde.EncodeToBytes(&ss.values)

	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(jsonVal)))
	base64.RawURLEncoding.Encode(buf, jsonVal)

	//payload := base64.RawURLEncoding.EncodeToString(jsonVal)

	md5Val := md5Base64(buf, lang.STB(MySessDB.Secret))

	ss.changed = false
	return lang.BTS(buf) + "." + md5Val
}

// YXRJT0l5ckpYNldBTjYzNHZw
// 闪电侠Guid
func genJwtGuid() string {
	size := int((MySessDB.SidSize*3 + 3) / 4)
	// TODO：要想办法保证sid的唯一性
	sid := randx.RandomBytes(size)
	guid := base64.RawURLEncoding.EncodeToString(sid)
	return guid[:MySessDB.SidSize] // 要确保guid的长度一致性
}
