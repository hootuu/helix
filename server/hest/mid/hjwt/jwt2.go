package hjwt

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/herr"
	"time"
)

const (
	dfIdentityKey   = "helix_identity"
	dfTokenLookup   = "header: helix_authorization, query: helix_token, cookie: helix_token"
	dfTokenHeadName = "helix_authorization"
)

type JwtMid struct {
	Code      string
	ginJwtMid *jwt.GinJWTMiddleware

	auth    func(dict dict.Dict) (interface{}, *herr.Error)
	refresh func(dict dict.Dict) (interface{}, *herr.Error)
}

func newJwtMid(code string) (*JwtMid, error) {
	mid := &JwtMid{Code: code}
	ginJwtMid := &jwt.GinJWTMiddleware{
		Realm:            mid.Code,
		SigningAlgorithm: hcfg.GetString(mid.ck("signing.algorithm"), ""),
		Key:              []byte(hcfg.GetString(mid.ck("key"), "")),
		KeyFunc:          nil,
		Timeout:          hcfg.GetDuration(mid.ck("timeout"), 3*time.Hour),
		MaxRefresh:       hcfg.GetDuration(mid.ck("max.refresh"), 24*time.Hour),
		Authenticator: func(c *gin.Context) (interface{}, error) {
			c.Request.Body
			mid.auth()
		},
		Authorizator:    nil,
		PayloadFunc:     nil,
		Unauthorized:    nil,
		LoginResponse:   nil,
		LogoutResponse:  nil,
		RefreshResponse: nil,
		IdentityHandler: nil,
		IdentityKey:     hcfg.GetString(mid.ck("identity.key"), dfIdentityKey),
		TokenLookup:     hcfg.GetString(mid.ck("token.lookup"), dfTokenLookup),
		TokenHeadName:   hcfg.GetString(mid.ck("token.head.name"), dfTokenHeadName),
	}

	return mid, nil
}

func (mid *JwtMid) ck(item string) string {
	return "hjwt." + mid.Code + "." + item
}

func (mid *JwtMid) IdentityHandler() func(ctx *gin.Context) interface{} {
	return func(ctx *gin.Context) interface{} {
		claims := jwt.ExtractClaims(ctx)
		return mid.identityHandler(dict.New(claims))
	}
}

func (mid *JwtMid) Authenticator() func(ctx *gin.Context) (interface{}, *herr.Error) {
	return func(ctx *gin.Context) (interface{}, *herr.Error) {
		return mid.authenticator(dict.New(nil)) //todo
	}
}

func init() {
	jwt.New()
	authMiddleware := jwt.GinJWTMiddleware{
		Realm:                 "",
		SigningAlgorithm:      "",
		Key:                   nil,
		KeyFunc:               nil,
		Timeout:               0,
		TimeoutFunc:           nil,
		MaxRefresh:            0,
		Authenticator:         nil,
		Authorizator:          nil,
		PayloadFunc:           nil,
		Unauthorized:          nil,
		LoginResponse:         nil,
		LogoutResponse:        nil,
		RefreshResponse:       nil,
		IdentityHandler:       nil,
		IdentityKey:           "",
		TokenLookup:           "",
		TokenHeadName:         "",
		TimeFunc:              nil,
		HTTPStatusMessageFunc: nil,
		PrivKeyFile:           "",
		PrivKeyBytes:          nil,
		PubKeyFile:            "",
		PrivateKeyPassphrase:  "",
		PubKeyBytes:           nil,
		SendCookie:            false,
		CookieMaxAge:          0,
		SecureCookie:          false,
		CookieHTTPOnly:        false,
		CookieDomain:          "",
		SendAuthorization:     false,
		DisabledAbort:         false,
		CookieName:            "",
		CookieSameSite:        0,
		ParseOptions:          nil,
		ExpField:              "",
	}
}
