package get_all

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/iamvkosarev/go-shared-utils/api/response"
	"github.com/iamvkosarev/go-shared-utils/api/user"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/wishlist/back/internal/model"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Wishlists []model.Wishlist `json:"wishlists"`
}

type UserWishlistsProvider interface {
	GetWishlists(userID int64) ([]model.Wishlist, error)
}

func NewWishlistHandler(log *slog.Logger, wishlistsProvider UserWishlistsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get_all.NewWishlistHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, err := user.GetUserId(r)
		if err != nil {
			log.Error("failed to get user id", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get user id"))
			return
		}

		wishlists, err := wishlistsProvider.GetWishlists(userID)
		if err != nil {
			log.Error("failed to get wishlists", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get wishlists"))
			return
		}
		responseOK(w, r, wishlists)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, wishlists []model.Wishlist) {
	render.JSON(
		w, r, Response{
			Response:  resp.Ok(),
			Wishlists: wishlists,
		},
	)
}
