package service

import (
	"app/internal/config"
	"app/internal/http/resp"
	"app/pkg/utils"
	"app/pkg/utils/helper"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const AudAccess = "access"
const AudRefresh = "refresh"
const KeyContextTokenClaim = "tokenClaim"

type AuthInfo struct {
	ID                     string `json:"id"`
	Username               string `json:"username"`
	AccessToken            string `json:"accessToken"`
	AccessTokenExpireAt    int64  `json:"accessTokenExpireAt"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshToken           string `json:"refreshToken"`
	RefreshTokenExpireAt   int64  `json:"refreshTokenExpireAt"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}

func (authInfo *AuthInfo) ToMap() map[string]any {
	return map[string]any{
		"id":                     authInfo.ID,
		"username":               authInfo.Username,
		"accessToken":            authInfo.AccessToken,
		"accessTokenExpireAt":    authInfo.AccessTokenExpireAt,
		"accessTokenExpireTime":  authInfo.AccessTokenExpireTime,
		"refreshToken":           authInfo.RefreshToken,
		"refreshTokenExpireAt":   authInfo.RefreshTokenExpireAt,
		"refreshTokenExpireTime": authInfo.RefreshTokenExpireTime,
	}
}

var AuthService = &authService{}

type authService struct{}

func (*authService) GetJwtKey() string {
	return config.GetConfig().ServerApi.JWTKey
}

func (*authService) GetJwtAccessTokenExpireTime() time.Duration {
	return time.Duration(config.GetConfig().ServerApi.JWTAccessTokenExpireTime) * time.Second
}

func (*authService) GetJwtRefreshTokenExpireTime() time.Duration {
	return time.Duration(config.GetConfig().ServerApi.JWTRefreshTokenExpireTime) * time.Second
}

func (*authService) JWTID(uid string) string {
	return fmt.Sprintf("%s-%s", utils.StringUtil.UniqueID(), uid)
}

func (*authService) AuthIDFromJWTID(jwtID string) string {
	a := strings.Split(jwtID, "-")
	id := a[len(a)-1]
	return id
}

func (svc *authService) GetAuthIDByTokenClaims(claims *jwt.RegisteredClaims) string {
	uid := svc.AuthIDFromJWTID(claims.ID)
	return uid
}

func (svc *authService) GetTokenClaims(c *gin.Context, aud ...string) (*jwt.RegisteredClaims, error) {
	v, ok := c.Get(KeyContextTokenClaim)
	if ok {
		claims, ok := v.(*jwt.RegisteredClaims)
		if ok {
			return claims, nil
		}
	}
	token := RequestService.GetAccessToken(c)
	if token == "" {
		return nil, resp.ErrInvalidToken
	}
	claims, err := utils.JwtUtil.ParseToken(token, svc.GetJwtKey())
	if err != nil {
		return nil, resp.ErrInvalidToken
	}
	if len(aud) == 0 {
		aud = []string{AudAccess}
	}
	for _, v := range aud {
		if !helper.InArray(claims.Audience, v) {
			return nil, resp.ErrInvalidToken
		}
	}
	c.Set(KeyContextTokenClaim, claims)
	return claims, nil
}

func (svc *authService) CurrentAuthID(c *gin.Context) string {
	claims, err := svc.GetTokenClaims(c)
	if err != nil {
		return ""
	}
	uid := svc.GetAuthIDByTokenClaims(claims)
	return uid
}

func (svc *authService) GenerateAccessToken(ctx context.Context, id string, username string, expireTime time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        svc.JWTID(id),
		Subject:   username,
		Issuer:    "wallet",
		Audience:  []string{AudAccess},
		ExpiresAt: jwt.NewNumericDate(expireTime),
	}
	token, err := utils.JwtUtil.GenerateToken(claims, svc.GetJwtKey())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (svc *authService) GenerateRefreshToken(ctx context.Context, id string, username string, expireTime time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        svc.JWTID(id),
		Subject:   username,
		Audience:  []string{AudRefresh},
		Issuer:    "wallet",
		ExpiresAt: jwt.NewNumericDate(expireTime),
	}
	token, err := utils.JwtUtil.GenerateToken(claims, svc.GetJwtKey())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (svc *authService) Login(ctx context.Context, id string, username string) (AuthInfo, error) {
	var loginInfo AuthInfo
	loginInfo.ID = id
	loginInfo.Username = username
	var err error
	now := time.Now()
	expireTime := now.Add(svc.GetJwtAccessTokenExpireTime())
	loginInfo.AccessTokenExpireAt = expireTime.UnixMilli()
	loginInfo.AccessTokenExpireTime = svc.GetJwtAccessTokenExpireTime().Milliseconds()
	loginInfo.AccessToken, err = svc.GenerateAccessToken(ctx, id, username, expireTime)
	if err != nil {
		return loginInfo, err
	}
	expireTime = now.Add(svc.GetJwtRefreshTokenExpireTime())
	loginInfo.RefreshTokenExpireTime = svc.GetJwtRefreshTokenExpireTime().Milliseconds()
	loginInfo.RefreshToken, err = svc.GenerateRefreshToken(ctx, id, username, expireTime)
	if err != nil {
		return loginInfo, err
	}
	return loginInfo, nil
}
