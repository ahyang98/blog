package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	cache "github.com/patrickmn/go-cache"
	"strings"
	"time"
)

type JWTHandler struct {
	c             *cache.Cache
	signingMethod jwt.SigningMethod
	refreshKey    []byte
	JWTKey        []byte
	rcExpiration  time.Duration
	tkExpiration  time.Duration
	SsidKeyFmt    string
}

func NewJWTHandler(c *cache.Cache) *JWTHandler {
	return &JWTHandler{
		c:             c,
		signingMethod: jwt.SigningMethodHS512,
		refreshKey:    []byte(`sbUZPISeSMJIwJ4pfc1AdkkpHCFPUJPJ`),
		JWTKey:        []byte(`sbUZPISeSMJIwJ4pfc1AdkkpHCFPUJPJ`),
		rcExpiration:  0,
		tkExpiration:  0,
		SsidKeyFmt:    "users:ssid:%s",
	}
}

func (j *JWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header(JWTHttpHeaderKey, "")
	ctx.Header(RefreshTokenHeaderKey, "")
	userClaims := ctx.MustGet("user").(UserClaims)
	j.c.Set(fmt.Sprintf(j.SsidKeyFmt, userClaims.Ssid), "", j.rcExpiration)
	return nil
}

func (j *JWTHandler) SetLoginToken(ctx *gin.Context, uid uint) error {
	ssid := uuid.New().String()
	err := j.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return j.SetJWTToken(ctx, uid, ssid)
}

func (j *JWTHandler) SetJWTToken(ctx *gin.Context, uid uint, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tkExpiration)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(j.signingMethod, claims)
	tokenStr, err := token.SignedString(j.JWTKey)
	if err != nil {
		return err
	}
	ctx.Header(JWTHttpHeaderKey, tokenStr)
	return nil
}

func (j *JWTHandler) SetRefreshToken(ctx *gin.Context, uid uint, ssid string) error {
	claims := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.rcExpiration)),
		},
	}
	token := jwt.NewWithClaims(j.signingMethod, claims)
	tokenStr, err := token.SignedString(j.refreshKey)
	if err != nil {
		return err
	}
	ctx.Header(RefreshTokenHeaderKey, tokenStr)
	return nil
}

func (j *JWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, ok := j.c.Get(fmt.Sprintf(j.SsidKeyFmt, ssid))
	if !ok {
		return fmt.Errorf("invalid ssid")
	}
	return nil
}

func (j *JWTHandler) GetUserClaim(ctx *gin.Context) (UserClaims, error) {
	tokenStr := j.ExtractToken(ctx)
	var uc UserClaims
	token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (any, error) {
		return j.JWTKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return uc, fmt.Errorf("invalid token")
	}
	return uc, nil
}

func (j *JWTHandler) GetRefreshClaim(ctx *gin.Context) (RefreshClaims, error) {
	tokenStr := j.ExtractToken(ctx)
	var rc RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (any, error) {
		return j.refreshKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return rc, fmt.Errorf("invalid token")
	}
	return rc, nil
}

func (j *JWTHandler) ExtractToken(ctx *gin.Context) string {
	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	splits := strings.Split(auth, " ")
	if len(splits) != 2 {
		return ""
	}
	return splits[1]
}
