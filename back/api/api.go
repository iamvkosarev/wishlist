package api

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/iamvkosarev/wishlist/back/internal/model"
	"net/http"
	"strings"
)

type Storage interface {
	GetUser(email string) (*model.User, error)
}

type CORS interface {
	EnableCORS(next http.Handler, methods []string) http.Handler
}

type TokenValidator interface {
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type HTTPLogger interface {
	Error(writer http.ResponseWriter, message string, statusCode int)
	InternalError(writer http.ResponseWriter, message string, statusCode int, err error)
	Success(writer http.ResponseWriter, message string, statusCode int)
}

type Deps struct {
	Storage        Storage
	TokenValidator TokenValidator
	CORS           CORS
	Logger         HTTPLogger
}

type API struct {
	deps Deps
}

func NewAPI(deps Deps) *API {
	return &API{deps: deps}
}

func (api *API) CreateWishlistHandle() http.Handler {
	return api.deps.CORS.EnableCORS(http.HandlerFunc(api.createWishlistHandler), []string{"POST"})
}

func (api *API) createWishlistHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		api.deps.Logger.Error(writer, "Not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		api.deps.Logger.Error(writer, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		api.deps.Logger.Error(writer, "Invalid Authorization Header format", http.StatusUnauthorized)
	}

	token := authHeader[len(bearerPrefix):]

	claims, err := api.deps.TokenValidator.ValidateToken(token)
	if err != nil {
		api.deps.Logger.Error(writer, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	fmt.Println(claims)
}
