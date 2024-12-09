package main

import (
	"github.com/iamvkosarev/go-shared-utils/cors"
	"github.com/iamvkosarev/go-shared-utils/jwts"
	"github.com/iamvkosarev/go-shared-utils/logs"
	"github.com/iamvkosarev/wishlist/back/api"
	"github.com/iamvkosarev/wishlist/back/internal/storage/storage"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

const JWT_SECRET_KEY = "JWT_SECRET"

type API interface {
	CreateWishlistHandle() http.Handler
}

func main() {
	godotenv.Load()
	var newAPI API = api.NewAPI(
		api.Deps{
			Storage:        storage.NewLocalStorage(),
			TokenValidator: jwts.NewJWTValidator(os.Getenv(JWT_SECRET_KEY)),
			CORS:           cors.NewCORS([]string{"https://kosarev.app", "http://localhost:63342"}),
			Logger:         logs.NewHttpLogger(true, true),
		},
	)
	http.Handle("/wishlist/create", newAPI.CreateWishlistHandle())
	http.ListenAndServe(":8081", nil)
}
