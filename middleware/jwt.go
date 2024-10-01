package middleware

import (
	"context"
	"hypruntime/config"
	"hypruntime/logger"
	"hypruntime/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaimsKey string

const JWTClaims JWTClaimsKey = "jwt_claims"

func HandleJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context = r.Context()
		tokenStr := r.Header.Get("Authorization")

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
		if err != nil {
			if config.IsDevEnvironment() {
				logger.Debug(r.Context()).Err(err).Msg("JWT parse error")
				next.ServeHTTP(w, r)
				return
			}
			logger.Error(r.Context()).Err(err).Msg("JWT parse error")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			claimsJson, err := utils.JsonSerialize(claims)
			if err != nil {
				logger.Error(r.Context()).Err(err).Msg("JWT claims serialization error")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, JWTClaims, string(claimsJson))
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetJWTClaims(ctx context.Context) string {
	if claims, ok := ctx.Value(JWTClaims).(string); ok {
		return claims
	}
	return ""
}
