package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaimsKey string

const JWTClaims JWTClaimsKey = "jwt_claims"

func HandleJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context = r.Context()
		tokenStr := r.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		privKeysStr := os.Getenv("MODUS_PRIV_KEYS")
		if privKeysStr == "" {
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}
			if config.IsDevEnvironment() {
				token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
				if err != nil {
					logger.Debug(r.Context()).Err(err).Msg("JWT parse error")
					next.ServeHTTP(w, r)
					return
				}
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					ctx = AddClaimsToContext(ctx, claims)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		var privKeysUnmarshalled map[string]string
		err := json.Unmarshal([]byte(privKeysStr), &privKeysUnmarshalled)
		if err != nil {
			logger.Error(r.Context()).Err(err).Msg("JWT private keys unmarshalling error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var token *jwt.Token

		for _, privKey := range privKeysUnmarshalled {
			token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return jwt.ParseRSAPublicKeyFromPEM([]byte(privKey))
			})
			if err != nil {
				if config.IsDevEnvironment() {
					logger.Debug(r.Context()).Err(err).Msg("JWT parse error")
					next.ServeHTTP(w, r)
					return
				}
				logger.Error(r.Context()).Err(err).Msg("JWT parse error")
				continue
			} else {
				break
			}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx = AddClaimsToContext(ctx, claims)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AddClaimsToContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	claimsJson, err := utils.JsonSerialize(claims)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("JWT claims serialization error")
		return ctx
	}
	return context.WithValue(ctx, JWTClaims, string(claimsJson))
}

func GetJWTClaims(ctx context.Context) string {
	if claims, ok := ctx.Value(JWTClaims).(string); ok {
		return claims
	}
	return ""
}
