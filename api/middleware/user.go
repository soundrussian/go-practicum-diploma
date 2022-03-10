package middleware

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/soundrussian/go-practicum-diploma/pkg/curruser"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"net/http"
)

// CurrentUser is middleware that gets user_id from JWT token (if it is present and signed)
// and saves it in request context
func CurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, logger := logging.CtxLogger(r.Context())

		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			logger.Err(err).Msg("failed to fetch token from request context")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			logger.Error().Msg("jwt token is nil or invalid")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var userIDAsFloat float64
		var userID uint64
		var v interface{}
		var ok bool

		if v, ok = claims["user_id"]; !ok {
			logger.Error().Msg("no user_id in jwt claims")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// user_id claim is float64, so we should cast it to float64 first,
		// and then convert it to uint64
		if userIDAsFloat, ok = v.(float64); !ok {
			logger.Error().Msgf("could not convert %v to float64", v)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID = uint64(userIDAsFloat)

		logger = logger.With().Uint64(logging.CurrentUserKey, userID).Logger()
		ctx = logging.SetCtxLogger(ctx, logger)
		ctx = curruser.SetCurrentUser(ctx, userID)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
