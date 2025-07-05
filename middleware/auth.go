package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gitlab.com/durfy/durpify/handlers"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

func InitAuthMiddleware(allowedGroups []string, jwks string) *AuthConfig {
	return &AuthConfig{
		allowedGroups: allowedGroups,
		jwks:          jwks,
	}
}

type AuthConfig struct {
	allowedGroups []string
	jwks          string
}

type StandardMessage struct {
	Message string `json:"message" example:"message"`
}

func (cfg *AuthConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var groups []string

		tokenString, err := getToken(w)
		if err != nil {
			resp := handlers.NewFailureResponse(
				err.Error(),
				http.StatusUnauthorized,
				[]string{},
			)
			resp.SendReponse(w)
		}

		token, err := cfg.validateToken(tokenString)
		if err != nil {
			resp := handlers.NewFailureResponse(
				"Failed to Validate Token",
				http.StatusUnauthorized,
				[]string{err.Error()},
			)
			resp.SendReponse(w)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			resp := handlers.NewFailureResponse(
				"Invalid Authorization token claim",
				http.StatusUnauthorized,
				[]string{},
			)
			resp.SendReponse(w)
			return
		}

		groupsClaim, ok := claims["groups"].([]interface{})
		if !ok {
			resp := handlers.NewFailureResponse(
				"Missing or invalid groups in the token",
				http.StatusUnauthorized,
				[]string{},
			)
			resp.SendReponse(w)
			return
		}

		for _, group := range groupsClaim {
			if groupName, ok := group.(string); ok {
				groups = append(groups, groupName)
			}
		}

		isAllowed := false
		for _, allowedGroup := range cfg.allowedGroups {
			for _, group := range groups {
				if group == allowedGroup {
					isAllowed = true
					break
				}
			}
			if isAllowed {
				break
			}
		}

		if !isAllowed {
			resp := handlers.NewFailureResponse(
				"Unauthorized to use this endpoint",
				http.StatusUnauthorized,
				[]string{},
			)
			resp.SendReponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *AuthConfig) validateToken(tokenString string) (*jwt.Token, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {
			log.Error("There was an error with the jwt.Keyfunc" + err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	jwks, err := keyfunc.Get(cfg.jwks, options)
	defer jwks.EndBackground()
	if err != nil {
		return nil, errors.New("Failed to get JWKS")
	}

	token, err := jwt.Parse(tokenString, jwks.Keyfunc)
	if err != nil {
		return nil, errors.New("Failed to Parse JWT")
	}

	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	return token, nil
}

func getToken(w http.ResponseWriter) (string, error) {
	tokenString := w.Header().Get("Authorization")

	if tokenString == "" {
		return "", errors.New("No Token Detected")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	return tokenString, nil
}
