package middleware

import (
	"context"
	"net/http"
	"strings"

	"lms-be/configs"
	"lms-be/shared/failure"
	"lms-be/shared/logger"
	"lms-be/transport/http/response"

	"github.com/dgrijalva/jwt-go"
)

// JWT is struct to using a JWT Verifier
// Usage on Handlers :
// ctx:= r.Context().Value(middleware.ValueKeyContextue)
// log.Info().Interface("IDOpd", ctx["IDOpd"]).Interface("UserID", ctx["userId"]).Msg("Context")
// ctx key could refer on user.NewUserLoginClaims
type JWT struct {
	Config *configs.Config
}

// JwtKeyContext is list of key for value in context
type JwtKeyContext string

const (
	ValueKeyContext JwtKeyContext = "value"
)

// ProvideJWTMiddleware is the middleware for JWT
func ProvideJWTMiddleware(config *configs.Config) *JWT {
	return &JWT{
		Config: config,
	}
}

// VerifyToken is function to verify the token is valid or not
func (j *JWT) VerifyToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Split(r.Header.Get("Authorization"), " ")

		if len(authorization) <= 1 {
			response.WithError(w, failure.Unauthorized("Token invalid."))
			return
		}

		token, err := jwt.Parse(authorization[1], func(token *jwt.Token) (interface{}, error) {
			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				response.WithError(w, failure.Unauthorized("Token invalid."))
				return nil, nil
			}

			if method != jwt.SigningMethodHS256 {
				response.WithError(w, failure.Unauthorized("Token invalid."))
				return nil, nil
			}

			return []byte(j.Config.Token.JWT.AccessToken), nil
		})

		if err != nil {
			logger.ErrorWithStack(err)
			response.WithError(w, failure.Unauthorized("Token invalid."))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			response.WithError(w, failure.Unauthorized("Token invalid."))
			return
		}

		ctx := context.WithValue(r.Context(), ValueKeyContext, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// GetClaimsValue is function to retrieve value from Context
func GetClaimsValue(ctx context.Context, key string) interface{} {
	return ctx.Value(ValueKeyContext).(jwt.MapClaims)[key]
}
