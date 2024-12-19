package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamvkosarev/go-shared-utils/middleware/auth"
	"github.com/iamvkosarev/go-shared-utils/middleware/logger"
	"github.com/iamvkosarev/wishlist/back/internal/http-server/handlers/wishlist/create"
	"github.com/iamvkosarev/wishlist/back/internal/http-server/handlers/wishlist/get"
	"github.com/iamvkosarev/wishlist/back/internal/http-server/handlers/wishlist/get_all"
	"github.com/iamvkosarev/wishlist/back/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, storage *sqlite.Storage, ssoURL string) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.URLFormat)
	mux.Use(logger.NewLogger(log))

	mux.Group(
		func(r chi.Router) {
			r.Use(auth.Auth(log, ssoURL))
			r.Post("/api/wishlist", create.NewWishlistHandler(log, storage))
			r.Get("/api/wishlist", get_all.NewWishlistHandler(log, storage))
			r.Get(fmt.Sprintf("/api/wishlist/{%s}", get.WishlistIDParam), get.NewWishlistHandler(log, storage))

		},
	)

	return mux
}
